// import * as ActionTypes from '../actions/app';

const initialState = {
  sshKeys: {},
  groups: {},
  files: {},
  metadata: {},
  collections: {},
  uncrawlables: {},
  tasks: {},
  users: {},
};

// updates an entity cache in response to any actuion with response.local.
// see local middleware
export default function locals(state = initialState, action) {
  if (action.locals) {
    if (action.locals.remove) {
      const newState = Object.assign({}, state, action.locals.entities);
      Object.keys(action.locals.remove.entities).forEach((model) => {
        Object.keys(action.locals.remove.entities[model]).forEach((id) => {
          delete newState[model][id];
        });
      });
      return newState;
    } else if (action.locals.entities) {
      return Object.assign({}, state, action.locals.entities);
    }
  }
  return state;
}
