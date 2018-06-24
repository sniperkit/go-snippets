import { push } from 'react-router-redux';
import { USERS_API } from '../middleware/users';

import analytics from '../analytics';
import Schemas from '../schemas';
import { setMessage, resetMessage } from './app';
import { updateLocalModel, editModel } from './locals';
import { selectSessionUser } from '../selectors/session';

export const GROUPS_FETCH_REQUEST = "GROUPS_FETCH_REQUEST";
export const GROUPS_FETCH_SUCCESS = "GROUPS_FETCH_SUCCESS";
export const GROUPS_FETCH_FAILURE = "GROUPS_FETCH_FAILURE";

export function fetchGroups(page = 1, pageSize = 25) {
  return {
    [USERS_API]: {
      types: [GROUPS_FETCH_REQUEST, GROUPS_FETCH_SUCCESS, GROUPS_FETCH_FAILURE],
      schema: Schemas.GROUP,
      endpoint: '/groups',
      data: { page, pageSize },
    },
  };
}

export function loadGroups(page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchGroups(page, pageSize));
  };
}

export const GROUP_FETCH_REQUEST = "GROUP_FETCH_REQUEST";
export const GROUP_FETCH_SUCCESS = "GROUP_FETCH_SUCCESS";
export const GROUP_FETCH_FAILURE = "GROUP_FETCH_FAILURE";

export function fetchGroup(id = "") {
  return {
    [USERS_API]: {
      types: [GROUP_FETCH_REQUEST, GROUP_FETCH_SUCCESS, GROUP_FETCH_FAILURE],
      schema: Schemas.GROUP,
      endpoint: `/groups/${id}`,
      data: { id },
    },
  };
}

export function loadGroup(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchGroup(id));
  };
}

export const GROUP_SAVE_REQUEST = "GROUP_SAVE_REQUEST";
export const GROUP_SAVE_SUCCESS = "GROUP_SAVE_SUCCESS";
export const GROUP_SAVE_FAILURE = "GROUP_SAVE_FAILURE";

export function saveGroup(group) {
  return {
    [USERS_API]: {
      types : [GROUP_SAVE_REQUEST, GROUP_SAVE_SUCCESS, GROUP_SAVE_FAILURE],
      schema: Schemas.GROUP,
      method: group.id ? "PUT" : "POST",
      endpoint: `/groups`,
      data: group,
    },
  };
}

export const GROUP_DELETE_REQUEST = "GROUP_DELETE_REQUEST";
export const GROUP_DELETE_SUCCESS = "GROUP_DELETE_SUCCESS";
export const GROUP_DELETE_FAILURE = "GROUP_DELETE_FAILURE";

export function deleteGroup(group) {
  return {
    [USERS_API]: {
      types : [GROUP_SAVE_REQUEST, GROUP_SAVE_SUCCESS, GROUP_SAVE_FAILURE],
      schema: Schemas.GROUP,
      method: group.id ? "PUT" : "POST",
      endpoint: `/groups`,
      data: group,
    },
  };
}