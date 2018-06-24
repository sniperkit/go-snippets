import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

// import { selectConsensus } from '../selectors/consensus';

export const CONSENSUS_REQUEST = 'CONSENSUS_REQUEST';
export const CONSENSUS_SUCCESS = 'CONSENSUS_SUCCESS';
export const CONSENSUS_FAILURE = 'CONSENSUS_FAILURE';

export function fetchConsensus(subject) {
  return {
    [CALL_API]: {
      types: [CONSENSUS_REQUEST, CONSENSUS_SUCCESS, CONSENSUS_FAILURE],
      endpoint: `/consensus?subject=${subject}`,
      schema: Schemas.CONSENSUS,
      data: { subject },
    },
  };
}

export function loadConsensus(subject) {
  return fetchConsensus(subject);
  // return (dispatch) => {
  //   return dispatch(fetchConsensus(subject));
  // };
}
