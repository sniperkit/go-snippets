import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';
import { selectUrl } from '../selectors/url';

export const URL_ARCHIVE_REQUEST = "URL_ARCHIVE_REQUEST";
export const URL_ARCHIVE_SUCCESS = "URL_ARCHIVE_SUCCESS";
export const URL_ARCHIVE_FAILURE = "URL_ARCHIVE_FAILURE";

export function archiveUrl(url = "", callback) {
  return (dispatch) => {
    return dispatch({
      [CALL_API]: {
        types: [URL_ARCHIVE_REQUEST, URL_ARCHIVE_SUCCESS, URL_ARCHIVE_FAILURE],
        schema: Schemas.URL,
        endpoint: '/archive',
        data: {
          url,
        },
      },
    }).then((action) => {
      if (typeof callback == "function") {
        callback(action);
      }
    });
  };
}

export const URL_SET_LOADING = "URL_SET_LOADING";
export function setUrlLoading(urlString = "", loading = false) {
  return {
    type: URL_SET_LOADING,
    data: {
      url: urlString,
      loading,
    },
  };
}

export const URL_SET_SUCCESS = "URL_SET_SUCCESS";
export function setUrlSuccess(urlString = "", success = false) {
  return {
    type: URL_SET_SUCCESS,
    data: {
      url: urlString,
      success,
    },
  };
}

export const URL_SET_ERROR = "URL_SET_ERROR";
export function setUrlError(urlString = "", error = "") {
  return {
    type: URL_SET_ERROR,
    data: {
      url: urlString,
      error,
    },
  };
}

export const URL_FETCH_REQUEST = "URL_FETCH_REQUEST";
export const URL_FETCH_SUCCESS = "URL_FETCH_SUCCESS";
export const URL_FETCH_FAILURE = "URL_FETCH_FAILURE";

export function fetchUrl(url = "", silentError = false) {
  return {
    [CALL_API]: {
      types: [URL_FETCH_REQUEST, URL_FETCH_SUCCESS, URL_FETCH_FAILURE],
      schema: Schemas.URL,
      endpoint: '/url',
      silentError,
      data: {
        url,
      },
    },
  };
}

export function loadUrl(url = "", silentError = false) {
  return (dispatch, getState) => {
    if (selectUrl(getState(), url)) {
      return null;
    }
    return dispatch(fetchUrl(url, silentError));
  };
}

export const URL_FETCH_OUTBOUND_LINKS_REQUEST = "URL_FETCH_OUTBOUND_LINKS_REQUEST";
export const URL_FETCH_OUTBOUND_LINKS_SUCCESS = "URL_FETCH_OUTBOUND_LINKS_SUCCESS";
export const URL_FETCH_OUTBOUND_LINKS_FAILURE = "URL_FETCH_OUTBOUND_LINKS_FAILURE";

export function fetchOutboundLinks(url = "") {
  return {
    [CALL_API]: {
      types: [URL_FETCH_OUTBOUND_LINKS_REQUEST, URL_FETCH_OUTBOUND_LINKS_SUCCESS, URL_FETCH_OUTBOUND_LINKS_FAILURE],
      schema: Schemas.LINK_ARRAY,
      endpoint: '/url/links/outbound',
      data: { url },
    },
  };
}

export function loadOutboundLinks(url = "") {
  return (dispatch) => {
    // TODO - check for local url copy
    return dispatch(fetchOutboundLinks(url));
  };
}

export const URL_FETCH_INBOUND_LINKS_REQUEST = "URL_FETCH_INBOUND_LINKS_REQUEST";
export const URL_FETCH_INBOUND_LINKS_SUCCESS = "URL_FETCH_INBOUND_LINKS_SUCCESS";
export const URL_FETCH_INBOUND_LINKS_FAILURE = "URL_FETCH_INBOUND_LINKS_FAILURE";

export function fetchInboundLinks(url = "") {
  return {
    [CALL_API]: {
      types: [URL_FETCH_INBOUND_LINKS_REQUEST, URL_FETCH_INBOUND_LINKS_SUCCESS, URL_FETCH_INBOUND_LINKS_FAILURE],
      schema: Schemas.LINK_ARRAY,
      endpoint: '/url/links/inbound',
      data: { url },
    },
  };
}

export function loadInboundLinks(url = "") {
  return (dispatch) => {
    // TODO - check for local url copy
    return dispatch(fetchInboundLinks(url));
  };
}
