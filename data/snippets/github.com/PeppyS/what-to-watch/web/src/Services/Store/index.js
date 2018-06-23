import reducers from './reducers'
import { createStore, applyMiddleware, compose } from 'redux'
import { createLogger } from 'redux-logger'
import thunk from 'redux-thunk'

const initialState = {}
const middleware = [thunk, createLogger({
    collapsed: true,
    diff: true,
    duration: true
})]

const store = createStore(reducers, initialState, compose(applyMiddleware(...middleware)))

export default store
