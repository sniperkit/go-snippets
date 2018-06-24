import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

import analytics from '../analytics';
import { newLocalModel, updateLocalModel, editModel, removeLocalModel } from './locals';

const blankUncrawlable = {
  meta: {
    title: "",
    description: "",
  },
};

const UNCRAWLABLE_NEW = "UNCRAWLABLE_NEW";
export function newUncrawlable(keyId) {
  return newLocalModel(Schemas.UNCRAWLABLE, UNCRAWLABLE_NEW, Object.assign({}, blankUncrawlable, { keyId }));
}

const UNCRAWLABLE_EDIT = "UNCRAWLABLE_EDIT";
export function editUncrawlable(collection) {
  return editModel(Schemas.UNCRAWLABLE, UNCRAWLABLE_EDIT, Object.assign({}, blankUncrawlable, collection));
}

const UNCRAWLABLE_CANCEL_EDIT = "UNCRAWLABLE_CANCEL_EDIT";
export function cancelUncrawlableEdit(collection) {
  return removeLocalModel(Schemas.UNCRAWLABLE, UNCRAWLABLE_CANCEL_EDIT, collection);
}

const UNCRAWLABLE_UPDATE = "UNCRAWLABLE_UPDATE";
export function updateUncrawlable(collection) {
  return updateLocalModel(Schemas.UNCRAWLABLE, UNCRAWLABLE_UPDATE, collection);
}

export const UNCRAWLABLES_FETCH_REQUEST = "UNCRAWLABLES_FETCH_REQUEST";
export const UNCRAWLABLES_FETCH_SUCCESS = "UNCRAWLABLES_FETCH_SUCCESS";
export const UNCRAWLABLES_FETCH_FAILURE = "UNCRAWLABLES_FETCH_FAILURE";

export function fetchUncrawlables(page = 1, pageSize = 25) {
  return {
    [CALL_API]: {
      types: [UNCRAWLABLES_FETCH_REQUEST, UNCRAWLABLES_FETCH_SUCCESS, UNCRAWLABLES_FETCH_FAILURE],
      schema: Schemas.UNCRAWLABLE,
      endpoint: '/collections',
      data: { page, pageSize },
    },
  };
}

export function loadUncrawlables(page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchUncrawlables(page, pageSize));
  };
}


export const UNCRAWLABLE_FETCH_REQUEST = "UNCRAWLABLE_FETCH_REQUEST";
export const UNCRAWLABLE_FETCH_SUCCESS = "UNCRAWLABLE_FETCH_SUCCESS";
export const UNCRAWLABLE_FETCH_FAILURE = "UNCRAWLABLE_FETCH_FAILURE";

export function fetchUncrawlable(id = "") {
  return {
    [CALL_API]: {
      types: [UNCRAWLABLE_FETCH_REQUEST, UNCRAWLABLE_FETCH_SUCCESS, UNCRAWLABLE_FETCH_FAILURE],
      schema: Schemas.UNCRAWLABLE_ARRAY,
      endpoint: `/collections/${id}`,
      data: { id },
    },
  };
}

export function loadUncrawlable(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchUncrawlable(id));
  };
}

export const UNCRAWLABLE_SAVE_REQUEST = "UNCRAWLABLE_SAVE_REQUEST";
export const UNCRAWLABLE_SAVE_SUCCESS = "UNCRAWLABLE_SAVE_SUCCESS";
export const UNCRAWLABLE_SAVE_FAILURE = "UNCRAWLABLE_SAVE_FAILURE";

export function saveUncrawlable(collection = {}) {
  return (dispatch) => {
    analytics.track(collection.id ? "Created Uncrawlable" : "Saved Uncrawlable", collection);
    return dispatch({
      [CALL_API]: {
        types: [UNCRAWLABLE_SAVE_REQUEST, UNCRAWLABLE_SAVE_SUCCESS, UNCRAWLABLE_SAVE_FAILURE],
        schema: Schemas.UNCRAWLABLE,
        method: (collection.id == "new") ? "POST" : "PUT",
        endpoint: "/collections",
        data: collection,
      },
    }).then((action) => {
      if (action.type == UNCRAWLABLE_SAVE_SUCCESS) {
        cancelUncrawlableEdit(collection.keyId, collection.subject);
      }
    });
  };
}

export const UNCRAWLABLE_DELETE_REQUEST = "UNCRAWLABLE_DELETE_REQUEST";
export const UNCRAWLABLE_DELETE_SUCCESS = "UNCRAWLABLE_DELETE_SUCCESS";
export const UNCRAWLABLE_DELETE_FAILURE = "UNCRAWLABLE_DELETE_FAILURE";

export function deleteUncrawlable(collection = {}) {
  analytics.track("Deleted Uncrawlable", collection);
  return (dispatch) => {
    return dispatch({
      [CALL_API]: {
        types: [UNCRAWLABLE_DELETE_REQUEST, UNCRAWLABLE_DELETE_SUCCESS, UNCRAWLABLE_DELETE_FAILURE],
        schema: Schemas.UNCRAWLABLE,
        endpoint: "/collections",
        method: "DELETE",
        data: collection,
      },
    });
  };
}
