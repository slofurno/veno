package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	add         chan *Conn
	remove      chan *Conn
	send        chan []byte
	connections []*Conn
}

var hub = &Hub{
	add:         make(chan *Conn, 16),
	remove:      make(chan *Conn, 16),
	send:        make(chan []byte, 64),
	connections: []*Conn{},
}

type Conn struct {
	*websocket.Conn
	send chan []byte
}

func (s *Conn) writeLoop() {
	defer func() {
		s.Close()
	}()

	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				s.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			wc, err := s.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = wc.Write(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			wc.Close()
		}
	}
}

func (s *Hub) addConnection(conn *Conn) {
	s.connections = append(s.connections, conn)
}

func (s *Hub) removeConnection(conn *Conn) {
	func() {
		for i, c := range s.connections {
			if conn == c {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				return
			}
		}
	}()

	close(conn.send)
}

func (s *Hub) broadcast(msg []byte) {
	for _, n := range s.connections {
		n.send <- msg
	}
}

func (s *Hub) Process() {
	for {
		select {
		case conn := <-s.add:
			s.addConnection(conn)
		case msg := <-s.send:
			s.broadcast(msg)
		case conn := <-s.remove:
			s.removeConnection(conn)
		}
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn := &Conn{ws, make(chan []byte, 32)}
	hub.addConnection(conn)

	go conn.writeLoop()

	defer func() {
		hub.removeConnection(conn)
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}

		hub.broadcast(msg)
	}
}

func main() {
	go hub.Process()

	http.HandleFunc("/ws", websocketHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8080", nil)
}
