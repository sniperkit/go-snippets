import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

import { selectContent, selectContentUrls, selectContentMetadata } from '../selectors/content';

export const CONTENT_REQUEST = 'CONTENT_REQUEST';
export const CONTENT_SUCCESS = 'CONTENT_SUCCESS';
export const CONTENT_FAILURE = 'CONTENT_FAILURE';

export function fetchContent(hash) {
  return {
    [CALL_API]: {
      types: [CONTENT_REQUEST, CONTENT_SUCCESS, CONTENT_FAILURE],
      endpoint: `/content/${hash}`,
      schema: Schemas.CONTENT,
      data: { hash },
    },
  };
}

export function loadContent(hash) {
  return (dispatch, getState) => {
    if (selectContent(getState(), hash)) {
      return null;
    }

    return dispatch(fetchContent(hash));
  };
}

// export const CONTENT_CONSENSUS_REQUEST = 'CONTENT_CONSENSUS_REQUEST';
// export const CONTENT_CONSENSUS_SUCCESS = 'CONTENT_CONSENSUS_SUCCESS';
// export const CONTENT_CONSENSUS_FAILURE = 'CONTENT_CONSENSUS_FAILURE';

// export function fetchContentConsensus(hash) {
//   return {
//     [CALL_API]: {
//       types: [CONTENT_CONSENSUS_REQUEST, CONTENT_CONSENSUS_SUCCESS, CONTENT_CONSENSUS_FAILURE],
//       endpoint: `/consensus/${hash}`,
//       schema: Schemas.CONSENSUS,
//       data: { hash },
//     },
//   };
// }

// export function loadContentConsensus(hash) {
//   return (dispatch, getState) => {
//     if (selectContentConsensus(getState(), hash)) {
//       return null;
//     }

//     return dispatch(fetchContentConsensus(hash));
//   };
// }


export const CONTENT_METADATA_REQUEST = 'CONTENT_METADATA_REQUEST';
export const CONTENT_METADATA_SUCCESS = 'CONTENT_METADATA_SUCCESS';
export const CONTENT_METADATA_FAILURE = 'CONTENT_METADATA_FAILURE';

export function fetchContentMetadata(hash) {
  return {
    [CALL_API]: {
      types: [CONTENT_METADATA_REQUEST, CONTENT_METADATA_SUCCESS, CONTENT_METADATA_FAILURE],
      endpoint: `/metadata/${hash}`,
      schema: Schemas.METADATA_ARRAY,
      data: { hash },
    },
  };
}

export function loadContentMetadata(hash) {
  return (dispatch, getState) => {
    if (selectContentMetadata(getState(), hash).length) {
      return null;
    }

    return dispatch(fetchContentMetadata(hash));
  };
}


export const CONTENT_URLS_REQUEST = 'CONTENT_URLS_REQUEST';
export const CONTENT_URLS_SUCCESS = 'CONTENT_URLS_SUCCESS';
export const CONTENT_URLS_FAILURE = 'CONTENT_URLS_FAILURE';

export function fetchContentUrls(hash) {
  return {
    [CALL_API]: {
      types: [CONTENT_URLS_REQUEST, CONTENT_URLS_SUCCESS, CONTENT_URLS_FAILURE],
      endpoint: `/content/${hash}/urls`,
      schema: Schemas.URL_ARRAY,
      data: { hash },
    },
  };
}

export function loadContentUrls(hash) {
  return (dispatch, getState) => {
    if (selectContentUrls(getState(), hash).length) {
      return null;
    }

    return dispatch(fetchContentUrls(hash));
  };
}


export const CONTENT_RECENT_URLS_REQUEST = 'CONTENT_RECENT_URLS_REQUEST';
export const CONTENT_RECENT_URLS_SUCCESS = 'CONTENT_RECENT_URLS_SUCCESS';
export const CONTENT_RECENT_URLS_FAILURE = 'CONTENT_RECENT_URLS_FAILURE';

export function fetchRecentContentUrls(page, pageSize) {
  return {
    [CALL_API]: {
      types: [CONTENT_RECENT_URLS_REQUEST, CONTENT_RECENT_URLS_SUCCESS, CONTENT_RECENT_URLS_FAILURE],
      endpoint: `/content`,
      schema: Schemas.URL_ARRAY,
      data: { page, pageSize },
      page,
      pageSize,
    },
  };
}

export function loadRecentContentUrls(page = 1, pageSize = 30) {
  return (dispatch) => {
    // if (selectContentUrls(getState(), "").length) {
    //   return null;
    // }

    return dispatch(fetchRecentContentUrls(page, pageSize));
  };
}
