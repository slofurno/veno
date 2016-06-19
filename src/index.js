import React, { Component, PropTypes } from 'react'
import { Provider } from 'react-redux'
import { Router, Route, browserHistory } from 'react-router'
import { render } from 'react-dom'
import { createStore, applyMiddleware } from 'redux'
import thunk from 'redux-thunk'

import rootReducer from './reducers'
const store = createStore(rootReducer, applyMiddleware(thunk))

import App from './containers/app'

class Root extends Component {
  render() {
    const { store, history } = this.props
    return (
      <Provider store={store}>
        <Router history={history}>
          <Route path="/" component={App}>
          </Route>
        </Router>
      </Provider>
    )
  }
}

render(<Root store={store} history={browserHistory}/>, document.getElementById('root'))
