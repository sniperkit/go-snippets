import { createStore, applyMiddleware } from 'redux';
import { browserHistory } from 'react-router';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import api from '../middleware/api';
import users from '../middleware/users';
import coverage from '../middleware/coverage';
import locals from '../middleware/locals';
import rootReducer from '../reducers';


export default function configureStore(preloadedState) {
  return createStore(
    rootReducer,
    preloadedState,
    applyMiddleware(thunk, api, locals, users, coverage, routerMiddleware(browserHistory)),
  );
}
