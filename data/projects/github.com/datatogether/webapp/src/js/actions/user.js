import { USERS_API } from '../middleware/users';
import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';
import { updateLocalModel, editModel } from './locals';
import { selectUserByUsername, selectUserById } from '../selectors/user';

export const EDIT_USER = 'EDIT_USER';
export function editUser(id) {
  return (dispatch, getState) => {
    const user = selectUserById(getState(), id);
    if (!user) {
      return null;
    }

    return dispatch(editModel(Schemas.USER, EDIT_USER, user));
  };
}

const USER_UPDATE = 'USER_UPDATE';
export function updateUser(user) {
  return updateLocalModel(Schemas.USER, USER_UPDATE, user);
}

export const SAVE_USER_REQUEST = 'SAVE_USER_REQUEST';
export const SAVE_USER_SUCCESS = 'SAVE_USER_SUCCESS';
export const SAVE_USER_FAILURE = 'SAVE_USER_FAILURE';

export function saveUser(user) {
  return (dispatch) => {
    analytics.track("User Updated", user);
    return dispatch({
      [USERS_API]: {
        types: [SAVE_USER_REQUEST, SAVE_USER_SUCCESS, SAVE_USER_FAILURE],
        endpoint: `/users/${user.id}`,
        method: 'PUT',
        schema: Schemas.USER,
        data: user,
      },
    }).then((action) => {
      if (action.type == SAVE_USER_SUCCESS) {
        dispatch(setMessage("settings successfully saved"));
        setTimeout(() => {
          dispatch(resetMessage());
        }, 3500);
        return dispatch(push(`/users/${action.response.entities.users[user.id].username}`));
      }

      return null;
    });
  };
}

export const USER_REQUEST = 'USER_REQUEST';
export const USER_SUCCESS = 'USER_SUCCESS';
export const USER_FAILURE = 'USER_FAILURE';

// Fetches a single user from Github API.
// Relies on the custom API middleware defined in ../middleware/api.js.
function fetchUser(login) {
  return {
    [USERS_API]: {
      types: [USER_REQUEST, USER_SUCCESS, USER_FAILURE],
      endpoint: `/users/${login}`,
      schema: Schemas.USER,
    },
  };
}

// Fetches a single user from Github API unless it is cached.
// Relies on Redux Thunk middleware.
export function loadUser(login, requiredFields = []) {
  return (dispatch, getState) => {
    const user = getState().entities.users[login];
    if (user && requiredFields.every(key => Object.prototype.hasOwnProperty.call(user, key))) {
      return null;
    }

    return dispatch(fetchUser(login));
  };
}

export function fetchUserByUsername(username) {
  return {
    [USERS_API]: {
      types: [USER_REQUEST, USER_SUCCESS, USER_FAILURE],
      endpoint: `/users?username=${username}`,
      schema: Schemas.USER,
    },
  };
}

export function loadUserByUsername(username, requiredFields = []) {
  return (dispatch, getState) => {
    const user = selectUserByUsername(getState(), username);
    if (user && requiredFields.every(key => Object.prototype.hasOwnProperty.call(user, key))) {
      return null;
    }

    return dispatch(fetchUserByUsername(username));
  };
}

export const USERS_SEARCH_REQUEST = 'USERS_SEARCH_REQUEST';
export const USERS_SEARCH_SUCCESS = 'USERS_SEARCH_SUCCESS';
export const USERS_SEARCH_FAILURE = 'USERS_SEARCH_FAILURE';

export function usersSearch(query) {
  return {
    [USERS_API]: {
      types: [USERS_SEARCH_REQUEST, USERS_SEARCH_SUCCESS, USERS_SEARCH_FAILURE],
      endpoint: `/search?q=${query}`,
      schema: Schemas.USER_ARRAY,
    }
  }
}

export const USERS_REQUEST = 'USERS_REQUEST';
export const USERS_SUCCESS = 'USERS_SUCCESS';
export const USERS_FAILURE = 'USERS_FAILURE';

export function fetchUsers(userType, page = 1, pageSize = 50) {
  return {
    [USERS_API]: {
      types: [USERS_REQUEST, USERS_SUCCESS, USERS_FAILURE],
      endpoint: `/users?type=${userType}&page=${page}&page_size=${pageSize}`,
      schema: Schemas.USER_ARRAY,
    }
  }
}

export function loadUsers(userType, page = 1, pageSize = 50) {
  return (dispatch, getState) => {
    // TODO - prevent overfetching
    return dispatch(fetchUsers(userType, page, pageSize));
  }
}

export const USER_COMMUNITIES_REQUEST = 'USER_COMMUNITIES_REQUEST';
export const USER_COMMUNITIES_SUCCESS = 'USER_COMMUNITIES_SUCCESS';
export const USER_COMMUNITIES_FAILURE = 'USER_COMMUNITIES_FAILURE';

export function fetchUserCommunities(id, page = 1, pageSize = 50) {
  return {
    [USERS_API]: {
      types: [USER_COMMUNITIES_REQUEST, USER_COMMUNITIES_SUCCESS, USER_COMMUNITIES_FAILURE],
      endpoint: `/users/communities`,
      data: { id, page, pageSize },
      schema: Schemas.USER_ARRAY,
    },
    id, page, pageSize
  }
}

export function loadUserCommunities(id, page = 1, pageSize = 50) {
  return (dispatch, getState) => {
    // TODO - prevent overfetching
    return dispatch(fetchUserCommunities(id, page, pageSize));
  }
}

export const USER_COLLECTIONS_REQUEST = "USER_COLLECTIONS_REQUEST";
export const USER_COLLECTIONS_SUCCESS = "USER_COLLECTIONS_SUCCESS";
export const USER_COLLECTIONS_FAILURE = "USER_COLLECTIONS_FAILURE";

export function fetchUserCollections(creator, page = 1, pageSize = 25) {
  return ({
    [CALL_API]: {
      types: [USER_COLLECTIONS_REQUEST, USER_COLLECTIONS_SUCCESS, USER_COLLECTIONS_FAILURE],
      schema: Schemas.COLLECTION,
      endpoint: `/collections?creator=${creator}`,
      data: { creator, page, pageSize },
    },
    creator, page, pageSize,
  });
}

export function loadUserCollections(creator, page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchUserCollections(creator, page, pageSize));
  };
}
