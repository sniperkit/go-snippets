
/* globals __BUILD__, fetch */
import 'isomorphic-fetch';
import { normalize } from 'normalizr';
  
export const API_ROOT = (__BUILD__.DEVELOP) ?
`http://${window.location.hostname}:${__BUILD__.USERS_API_PORT}`:
`${__BUILD__.USERS_API_URL}`;

// Fetches an API response and normalizes the result JSON according to schema.
// This makes every API response have the same shape, regardless of how nested it was.
function callApi(method, endpoint, schema, reqData) {
  let fullUrl = (endpoint.indexOf(API_ROOT) === -1) ? API_ROOT + endpoint : endpoint;
  let body;

  // add query params to GET requests listed on data
  if (reqData && method == "GET") {
    let addedFirst = false;
    Object.keys(reqData).forEach((key) => {
      const val = encodeURIComponent(reqData[key]);
      if (val != "") {
        fullUrl += (addedFirst) ? `&${key}=${val}` : `?${key}=${val}`;
        addedFirst = true;
      }
    });
  } else if (reqData) {
    body = JSON.stringify(reqData);
  }

  return fetch(fullUrl, {
    method,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body,
  })
    .then((response) => {
      return response.json().then(json => ({ json, response }));
    }).then(({ json, response }) => {
      if (!response.ok) {
        return Promise.reject(json);
      } else if (response.status == 204) {
        return {};
      }

      const { data } = json;

      return Object.assign({},
        normalize(data, schema),
        // { nextPageUrl }
      );
    });
}

// Action key that carries API call info interpreted by this Redux middleware.
export const USERS_API = Symbol('USERS_API');

// A Redux middleware that interprets actions with USERS_API info specified.
// Performs the call and promises when such actions are dispatched.
export default store => next => (action) => {
  const usersApi = action[USERS_API];
  if (typeof usersApi === 'undefined') {
    return next(action);
  }

  let { endpoint } = usersApi;
  const { schema, types, data, method = "GET", silentError = false } = usersApi;

  if (typeof endpoint === 'function') {
    endpoint = endpoint(store.getState());
  }

  if (typeof endpoint !== 'string') {
    throw new Error('Specify a string endpoint URL.');
  }
  if (!schema) {
    throw new Error('Specify one of the exported Schemas.');
  }

  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected an array of three action types.');
  }
  if (!types.every(type => typeof type === 'string')) {
    throw new Error('Expected action types to be strings.');
  }

  function actionWith(actionData) {
    const finalAction = Object.assign({}, action, actionData);
    delete finalAction[USERS_API];
    return finalAction;
  }

  const [requestType, successType, failureType] = types;

  // fire an action indicating a request has been made
  // TODO - this passed in object should only have "type" prop. IT's currently
  // shimmed to support fetching community users w pagingation. Fix.
  // next(actionWith({ type: requestType, data, id: usersApi.id, page: usersApi.page, pageSize: usersApi.pageSize }));
  next(actionWith({ type: requestType }));

  // make the request
  return callApi(method, endpoint, schema, data).then(
    response => next(actionWith({
      type: successType,
      response,
      // TODO - remove this, see above
      // id: usersApi.id, page: usersApi.page, pageSize: usersApi.pageSize
    })),
    (error) => {
      let msg = 'Something Bad Happened';
      if (error.meta && error.meta.message) {
        msg = error.meta.message;
      }
      return next(actionWith({
        type: failureType,
        error: msg,
        silentError,
      }));
    },
  );
};
