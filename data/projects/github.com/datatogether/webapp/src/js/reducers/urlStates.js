import { URL_SET_LOADING, URL_SET_SUCCESS, URL_SET_ERROR } from '../actions/url';

const initialState = {};

export default function urlStates(state = initialState, action) {
  switch (action.type) {
    case URL_SET_LOADING:
      return Object.assign({}, state, {
        [action.data.url]: Object.assign({}, state[action.data.url], { loading: action.data.loading }),
      });
    case URL_SET_SUCCESS:
      return Object.assign({}, state, {
        [action.data.url]: Object.assign({}, state[action.data.url], { success: action.data.success, loading: false }),
      });
    case URL_SET_ERROR:
      return Object.assign({}, state, {
        [action.data.url]: Object.assign({}, state[action.data.url], { error: action.data.error, loading: false }),
      });
    default:
      return state;
  }
}
