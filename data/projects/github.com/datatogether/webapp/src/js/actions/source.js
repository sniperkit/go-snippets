import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

export const SOURCES_FETCH_REQUEST = "SOURCES_FETCH_REQUEST";
export const SOURCES_FETCH_SUCCESS = "SOURCES_FETCH_SUCCESS";
export const SOURCES_FETCH_FAILURE = "SOURCES_FETCH_FAILURE";

export function fetchSources(page = 1, pageSize = 25) {
  return {
    [CALL_API]: {
      types: [SOURCES_FETCH_REQUEST, SOURCES_FETCH_SUCCESS, SOURCES_FETCH_FAILURE],
      schema: Schemas.SOURCE,
      endpoint: '/sources',
      data: { page, pageSize },
      page,
      pageSize,
    },
    page,
    pageSize,
  };
}

export function loadSources(page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchSources(page, pageSize));
  };
}


export const SOURCE_FETCH_REQUEST = "SOURCE_FETCH_REQUEST";
export const SOURCE_FETCH_SUCCESS = "SOURCE_FETCH_SUCCESS";
export const SOURCE_FETCH_FAILURE = "SOURCE_FETCH_FAILURE";

export function fetchSource(id = "") {
  return {
    [CALL_API]: {
      types: [SOURCE_FETCH_REQUEST, SOURCE_FETCH_SUCCESS, SOURCE_FETCH_FAILURE],
      schema: Schemas.SOURCE_ARRAY,
      endpoint: `/sources/${id}`,
      data: { id },
    },
  };
}

export function loadSource(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchSource(id));
  };
}

export const SOURCE_URLS_REQUEST = "SOURCE_URLS_REQUEST";
export const SOURCE_URLS_SUCCESS = "SOURCE_URLS_SUCCESS";
export const SOURCE_URLS_FAILURE = "SOURCE_URLS_FAILURE";

export function fetchSourceUrls(id = "", page = 1, pageSize = 200) {
  return {
    [CALL_API]: {
      types: [SOURCE_URLS_REQUEST, SOURCE_URLS_SUCCESS, SOURCE_URLS_FAILURE],
      schema: Schemas.SOURCE_ARRAY,
      endpoint: `/primers/${id}/urls`,
      data: { id, page, pageSize },
      id,
      page,
      pageSize,
    },
  };
}

export function loadSourceUrls(id = "", page = 1, pageSize = 200) {
  return (dispatch) => {
    // TODO - check pagination
    return dispatch(fetchSourceUrls(id, page, pageSize));
  };
}

export const SOURCE_ATTRIBUTED_URLS_REQUEST = "SOURCE_ATTRIBUTED_URLS_REQUEST";
export const SOURCE_ATTRIBUTED_URLS_SUCCESS = "SOURCE_ATTRIBUTED_URLS_SUCCESS";
export const SOURCE_ATTRIBUTED_URLS_FAILURE = "SOURCE_ATTRIBUTED_URLS_FAILURE";

export function fetchSourceAttributedUrls(id = "", page = 1, pageSize = 200) {
  return {
    [CALL_API]: {
      types: [SOURCE_ATTRIBUTED_URLS_REQUEST, SOURCE_ATTRIBUTED_URLS_SUCCESS, SOURCE_ATTRIBUTED_URLS_FAILURE],
      schema: Schemas.SOURCE_ARRAY,
      endpoint: `/primers/${id}/urls/attributed`,
      data: { id, page, pageSize },
      id,
      page,
      pageSize,
    },
  };
}

export function loadSourceAttributedUrls(id = "", page = 1, pageSize = 200) {
  return (dispatch) => {
    // TODO - check pagination
    return dispatch(fetchSourceAttributedUrls(id, page, pageSize));
  };
}
