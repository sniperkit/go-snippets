import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

import analytics from '../analytics';
import { selectFile } from '../selectors/metadata';
import { newLocalModel, updateLocalModel, editModel, removeLocalModel } from './locals';

const blankFile = {
  id: "new",
  data: "",
  // meta: {
  //   title: "",
  //   description: "",
  //   keyword: "",
  //   rights: "",
  //   landingPage: "",
  //   theme: "",
  //   identifier: "",
  //   accessLevel: "public",
  //   language: "english",
  //   license: "http://www.usa.gov/publicdomain/label/1.0/",
  // },
};

const FILE_NEW = "FILE_NEW";
export function newFile(attrs) {
  return newLocalModel(Schemas.FILE, FILE_NEW, Object.assign({}, blankFile, attrs));
}

const FILE_EDIT = "FILE_EDIT";
export function editFile(metadata) {
  return editModel(Schemas.FILE, FILE_EDIT, Object.assign({}, blankFile, metadata));
}

const FILE_CANCEL_EDIT = "FILE_CANCEL_EDIT";
export function cancelFileEdit(metadata) {
  return removeLocalModel(Schemas.FILE, FILE_CANCEL_EDIT, metadata);
}

const FILE_UPDATE = "FILE_UPDATE";
export function updateFile(metadata) {
  return updateLocalModel(Schemas.FILE, FILE_UPDATE, metadata);
}

export const FILE_REQUEST = 'FILE_REQUEST';
export const FILE_SUCCESS = 'FILE_SUCCESS';
export const FILE_FAILURE = 'FILE_FAILURE';

export function fetchFile(keyId, subject) {
  return {
    [CALL_API]: {
      types: [FILE_REQUEST, FILE_SUCCESS, FILE_FAILURE],
      endpoint: `/metadata/${keyId}/${subject}`,
      schema: Schemas.FILE,
      data: { keyId, subject },
    },
  };
}

export function loadFile(keyId, subject, requiredFields = []) {
  return (dispatch, getState) => {
    const metadata = selectFile(getState(), keyId, subject);
    if (metadata && requiredFields.every(key => Object.prototype.hasOwnProperty.call(metadata, key))) {
      return null;
    }

    return dispatch(fetchFile(keyId, subject, requiredFields));
  };
}

export const FILE_SAVE_REQUEST = "FILE_SAVE_REQUEST";
export const FILE_SAVE_SUCCESS = "FILE_SAVE_SUCCESS";
export const FILE_SAVE_FAILURE = "FILE_SAVE_FAILURE";

export function saveFile(metadata = {}) {
  // TODO - should this happen here?
  const md = {};
  Object.keys(metadata).forEach((key) => {
    if (metadata[key]) {
      md[key] = metadata[key];
    }
  });

  return (dispatch) => {
    analytics.track("Created File", metadata);
    return dispatch({
      [CALL_API]: {
        types: [FILE_SAVE_REQUEST, FILE_SAVE_SUCCESS, FILE_SAVE_FAILURE],
        schema: Schemas.FILE,
        endpoint: "/metadata",
        data: md,
      },
    }).then((action) => {
      // TODO - this isn't getting called
      if (action.type == FILE_SUCCESS) {
        cancelFileEdit(metadata.keyId, metadata.subject);
      }
    });
  };
}


export const FILE_BY_KEY_REQUEST = "FILE_BY_KEY_REQUEST";
export const FILE_BY_KEY_SUCCESS = "FILE_BY_KEY_SUCCESS";
export const FILE_BY_KEY_FAILURE = "FILE_BY_KEY_FAILURE";

export function fetchFileByKey(key, page = 1, pageSize = 25) {
  return ({
    [CALL_API]: {
      types: [FILE_BY_KEY_REQUEST, FILE_BY_KEY_SUCCESS, FILE_BY_KEY_FAILURE],
      schema: Schemas.FILE,
      endpoint: `/metadata?key=${key}`,
      data: { key, page, pageSize },
      key,
      page,
      pageSize,
    },
  });
}

export function loadFileByKey(key, page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchFileByKey(key, page, pageSize));
  };
}
