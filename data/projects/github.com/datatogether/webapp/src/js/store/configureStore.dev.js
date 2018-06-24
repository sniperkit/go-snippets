import { createStore, applyMiddleware, compose } from 'redux';
import { routerMiddleware } from 'react-router-redux';
import { browserHistory } from 'react-router';
import thunk from 'redux-thunk';
// import createLogger from 'redux-logger';
import api from '../middleware/api';
import users from '../middleware/users';
import locals from '../middleware/locals';
import rootReducer from '../reducers';
import DevTools from '../containers/DevTools';
import coverage from '../middleware/coverage';

export default function configureStore(preloadedState) {
  const store = createStore(
    rootReducer,
    preloadedState,
    compose(
      applyMiddleware(thunk, api, locals, users, coverage, routerMiddleware(browserHistory)),
      DevTools.instrument(),
    ),
  );

  if (module.hot) {
    // Enable Webpack hot module replacement for reducers
    module.hot.accept('../reducers', () => {
      /* eslint-disable global-require */
      const nextRootReducer = require('../reducers').default;
      /* eslint-enable global-require */
      store.replaceReducer(nextRootReducer);
    });
  }

  return store;
}
