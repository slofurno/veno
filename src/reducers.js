import { combineReducers } from 'redux'

function todo(state = [], action) {
  switch(action.type) {
  default:
    return state
  }
}

const rootReducer = combineReducers({
  todo,
})

export default rootReducer
