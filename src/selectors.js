import { createSelector } from 'reselect'

const rawState = state => state
const selector = createSelector(rawState, (state) => {...state})
export default selector
