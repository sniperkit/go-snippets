import { APP_TOGGLE_MENU, APP_HIDE_MENU, APP_SHOW_MODAL, APP_HIDE_MODAL } from '../actions/app';
import { SEARCH_SET, SEARCH_CLEAR } from '../actions/search';

const initialState = {
  showMenu: false,
  modal: undefined,
  query: "",
};

export default function appReducer(state = initialState, action) {
  switch (action.type) {
    case APP_TOGGLE_MENU:
      return Object.assign({}, state, { showMenu: !state.showMenu });
    case APP_HIDE_MENU:
      return Object.assign({}, state, { showMenu: false });
    case APP_SHOW_MODAL:
      return Object.assign({}, state, { modal: action.modal });
    case APP_HIDE_MODAL:
      return Object.assign({}, state, { modal: undefined });
    case SEARCH_SET:
      return Object.assign({}, state, { query: action.query });
    case SEARCH_CLEAR:
      return Object.assign({}, state, { query: "" });
    // whenever the route changes, close the menu
    case "@@router/LOCATION_CHANGE":
      return Object.assign({}, state, { showMenu: false });
    default:
      return state;
  }
}
