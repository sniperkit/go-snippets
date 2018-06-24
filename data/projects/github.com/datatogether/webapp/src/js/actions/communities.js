import { USERS_API } from '../middleware/users';
import Schemas from '../schemas';

import analytics from '../analytics';
import { selectMetadata } from '../selectors/metadata';

export const COMMUNITY_USERS_REQUEST = "COMMUNITY_USERS_REQUEST";
export const COMMUNITY_USERS_SUCCESS = "COMMUNITY_USERS_SUCCESS";
export const COMMUNITY_USERS_FAILURE = "COMMUNITY_USERS_FAILURE";

export function fetchCommunityUsers(id, page = 1, pageSize = 25) {
  return ({
    [USERS_API]: {
      types: [COMMUNITY_USERS_REQUEST, COMMUNITY_USERS_SUCCESS, COMMUNITY_USERS_FAILURE],
      schema: Schemas.USER_ARRAY,
      endpoint: `/communities/users`,
      data: { id, page, pageSize },
    },
    id, page, pageSize,
  });
}

export function loadCommunityUsers(id, page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchCommunityUsers(id, page, pageSize));
  };
}