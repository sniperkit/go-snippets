/* globals window, WebSocket, __BUILD__ */
import { normalize } from 'normalizr';
import Schemas from '../schemas';
// import { setErrorMessage, resetErrorMessage } from '../actions/app';


// "expectations" allows us to have duplexed communication behave
// in a request/response style by storing requests we've sent by id
// with resolve & reject functions that correspond to expected outcomes.
// expecting is an object in the form:
// {
//    requestId : {
//      successActionType: resolve(action),
//      failureActionType: reject(action),
//    },
//    ...
//  }
let expectations = {};

// TODO - we should assign an ID to each request & track
// expectations as objects that track against that ID. Server will need
// to return that request ID every time, and ids will need to be unique
// system-wide? meh, only kinda.

// Generate an id for this request
function newRequestId() {
  /* eslint-disable no-bitwise, no-mixed-operators */
  return 'xxxxxxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
  /* eslint-enable no-bitwise, no-mixed-operators */
}

// add an expected response to the expecting list.
// @param {string} - requestId
// @param {object} - outcomes is an object with action types as keys &
//                   promise resolve / reject as values
function addExpectation(requestId, outcomes) {
  expectations[requestId] = outcomes;
}

// received events should call meetExpectation to check the list of
// currently expected events, resolve their promises, and remove
// them from the list
function meetExpectation(action) {
  const id = Object.keys(expectations).find((requestId) => {
    if (action.requestId == requestId) {
      const func = expectations[requestId][action.type];
      if (typeof func == 'function') {
        func(action);
      }
      return true;
    }

    return false;
  });

  if (id) {
    // remove any found requestId from the list of expectations
    // regardless of weather function was called or not
    delete expectations[id];
  }
}

console.log(__BUILD__);

// the URL to dial to
const WEBSOCKET_URL = (__BUILD__.DEVELOP) ? 
`ws://${window.location.hostname}:${__BUILD__.WEBSOCKET_PORT}/ws`:
__BUILD__.WEBSOCKET_URL;

// Variable to hold the connection
let conn;

// Connect the websocket using a passed-in redux dispatch
// @return {promise}
export function connect(dispatch, reconnectTimeout = 6500) {
  return new Promise((resolve, reject) => {
    if (window.WebSocket) {
      conn = new WebSocket(WEBSOCKET_URL);
      conn.onclose = (evt) => {
        console.log("WebSocket conn closed", evt);
        if (reconnectTimeout) {
          setTimeout(() => {
            connect(dispatch, reconnectTimeout);
          }, reconnectTimeout);
        }
      };

      conn.onmessage = (evt) => {
        let res;
        try {
          res = JSON.parse(evt.data);
        } catch (e) {
          console.log(evt.data);
          console.log(e);
          return;
        }
        console.log(res);
        if (res.schema) {
          try {
            res.response = normalize(res.data, Schemas[res.schema]);
            // console.log(JSON.stringify(res.response));
          } catch (e) {
            console.log(e);
            console.log(res);
            return;
          }
        }
        meetExpectation(res);
        dispatch(res);
      };

      conn.onopen = (evt) => {
        console.log("WebSocket conn opened", evt);
        resolve(evt);
      };

      conn.onerror = (evt) => {
        console.log("WebSocket conn error", evt);
        conn.close();
        if (reject) {
          reject(evt);
        }
      };

      return;
    }

    reject("websocket connections not supported");
  });
}

// check weather the websocket is connected or not
// @return {boolean}
export function connected() {
  return !!conn;
}

// disconnect closes the connection.
// @return {promise}
export function disconnect() {
  return new Promise((resolve) => {
    if (conn) {
      conn.close();
    }
    return resolve();
  });
}

// sendCallApiAction is meant to work with api middleware
// as a seamless upgrade to traditional JSON API calls
// "actions" submitted here are expected to be the CALL_API
// object referenced in actions
export function callApiAction(store, next, action) {
  let { endpoint } = action;
  const { schema, types, data, timeout = 15000, silentError } = action;

  // if we're in "connecting" state, defer until connection
  // occurs
  if (conn.readyState == 0) {
    const prev = conn.onopen;
    conn.onopen = (evt) => {
      callApiAction(store, next, action);
      prev(evt);
    };
    return null;
  }

  // websockets doesn't use endpoint but it's worth warning
  // when endpoint isn't set to maintain cross-compatibility
  // with the json API
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

  const requestId = newRequestId();
  const [requestType, successType, failureType] = types;

  function actionWith(d) {
    const finalAction = Object.assign({ requestId }, action, d);
    // delete finalAction[CALL_API]
    return finalAction;
  }

  // fire an action indicating a request will be made
  next(actionWith({ type: requestType }));
  
  // make the request
  conn.send(JSON.stringify({
    type: requestType,
    requestId,
    silentError,
    data,
  }));

  // return a promise with it's resolve/reject wired
  // to an expectation
  return new Promise((resolve, reject) => {
    // TODO - should we create a reject func if one isn't defined?
    // this would necessitate the need for a "silentError" flag.

    addExpectation(requestId, {
      [successType]: resolve,
      [failureType]: (act) => {
        if (reject) {
          reject(`${act.type} ${act.error}`);
        }
      },
    });

    if (timeout) {
      // If the server takes too long, meet our own expectation
      // with a request timeout error
      setTimeout(() => {
        if (expectations[requestId]) {
          if (typeof expectations[requestId][failureType] == "function") {
            expectations[requestId][failureType]({
              requestId,
              type: failureType,
              error: `request ${requestId} timed out`,
            });
          }
          delete expectations[requestId];
        }
      }, timeout);
    }
  });
}
