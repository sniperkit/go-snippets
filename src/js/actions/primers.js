import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

export const PRIMERS_FETCH_REQUEST = "PRIMERS_FETCH_REQUEST";
export const PRIMERS_FETCH_SUCCESS = "PRIMERS_FETCH_SUCCESS";
export const PRIMERS_FETCH_FAILURE = "PRIMERS_FETCH_FAILURE";

export function fetchPrimers(page = 1, pageSize = 25, baseOnly = false) {
  return {
    [CALL_API]: {
      types: [PRIMERS_FETCH_REQUEST, PRIMERS_FETCH_SUCCESS, PRIMERS_FETCH_FAILURE],
      schema: Schemas.PRIMER,
      endpoint: '/primers',
      data: { page, pageSize, baseOnly },
    },
  };
}

export function loadPrimers(page = 1, pageSize = 25, baseOnly = false) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchPrimers(page, pageSize, baseOnly));
  };
}


export const PRIMER_FETCH_REQUEST = "PRIMER_FETCH_REQUEST";
export const PRIMER_FETCH_SUCCESS = "PRIMER_FETCH_SUCCESS";
export const PRIMER_FETCH_FAILURE = "PRIMER_FETCH_FAILURE";

export function fetchPrimer(id = "") {
  return {
    [CALL_API]: {
      types: [PRIMER_FETCH_REQUEST, PRIMER_FETCH_SUCCESS, PRIMER_FETCH_FAILURE],
      schema: Schemas.PRIMER_ARRAY,
      endpoint: `/primers/${id}`,
      data: { id },
    },
  };
}

export function loadPrimer(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchPrimer(id));
  };
}

export const PRIMER_SUBPRIMERS_REQUEST = "PRIMER_SUBPRIMERS_REQUEST";
export const PRIMER_SUBPRIMERS_SUCCESS = "PRIMER_SUBPRIMERS_SUCCESS";
export const PRIMER_SUBPRIMERS_FAILURE = "PRIMER_SUBPRIMERS_FAILURE";

export function fetchSources(id = "") {
  return {
    [CALL_API]: {
      types: [PRIMER_SUBPRIMERS_REQUEST, PRIMER_SUBPRIMERS_SUCCESS, PRIMER_SUBPRIMERS_FAILURE],
      schema: Schemas.SUBPRIMER_ARRAY,
      endpoint: `/primers/${id}/subprimers`,
      data: { id },
    },
  };
}

export function loadSources(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchSources(id));
  };
}
