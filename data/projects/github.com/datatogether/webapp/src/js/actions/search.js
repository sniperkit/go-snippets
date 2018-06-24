import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

export const SEARCH_SET = "SEARCH_SET";
export function setSearch(query = "") {
  return {
    type: SEARCH_SET,
    query,
  };
}

export const SEARCH_CLEAR = "SEARCH_CLEAR";
export function clearSearch() {
  return { type: SEARCH_CLEAR };
}

export const SEARCH_REQUEST = "SEARCH_REQUEST";
export const SEARCH_SUCCESS = "SEARCH_SUCCESS";
export const SEARCH_FAILURE = "SEARCH_FAILURE";

export function search(query = "") {
  if (!query || query.length < 3) {
    return setSearch(query);
  }

  return (dispatch) => {
    dispatch(setSearch(query));

    return dispatch({
      [CALL_API]: {
        types: [SEARCH_REQUEST, SEARCH_SUCCESS, SEARCH_FAILURE],
        schema: Schemas.SEARCH_RESULT_ARRAY,
        endpoint: "/search",
        data: {
          query,
          page: 1,
        },
      },
    });
  };
}

