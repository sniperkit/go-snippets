
/* globals __BUILD__, fetch */
import 'isomorphic-fetch';

export const API_ROOT = (__BUILD__.DEVELOP) ?
`http://${window.location.hostname}:${__BUILD__.COVERAGE_API_PORT}`:
`${__BUILD__.COVERAGE_API_URL}`;

// Fetches an API response from the coverage api
// This makes every API response have the same shape, regardless of how nested it was.
function callApi(method, endpoint, reqData) {
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

      // const { data } = json;
      return json.data;
    });
}

// Action key that carries API call info interpreted by this Redux middleware.
export const COVERAGE_API = Symbol('COVERAGE_API');

// A Redux middleware that interprets actions with COVERAGE_API info specified.
// Performs the call and promises when such actions are dispatched.
export default store => next => (action) => {
  const coverageApi = action[COVERAGE_API];
  if (typeof coverageApi === 'undefined') {
    return next(action);
  }

  let { endpoint } = coverageApi;
  const { types, data, method = "GET", silentError = false } = coverageApi;

  if (typeof endpoint === 'function') {
    endpoint = endpoint(store.getState());
  }

  if (typeof endpoint !== 'string') {
    throw new Error('Specify a string endpoint URL.');
  }

  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected an array of three action types.');
  }
  if (!types.every(type => typeof type === 'string')) {
    throw new Error('Expected action types to be strings.');
  }

  function actionWith(actionData) {
    const finalAction = Object.assign({}, action, actionData);
    delete finalAction[COVERAGE_API];
    return finalAction;
  }

  const [requestType, successType, failureType] = types;

  // fire an action indicating a request has been made
  next(actionWith({ type: requestType }));

  // make the request
  return callApi(method, endpoint, data).then(
    response => next(actionWith({
      type: successType,
      tree: response,
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
