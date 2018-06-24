import { COVERAGE_API } from '../middleware/coverage';
import { selectNode } from '../selectors/coverage';

export const COVERAGE_NODE_TOGGLE = 'COVERAGE_NODE_TOGGLE';
export function toggleNode(id) {
  return {
    type: COVERAGE_NODE_TOGGLE,
    id,
  }
}

export const COVERAGE_NODE_REQUEST = 'COVERAGE_NODE_REQUEST';
export const COVERAGE_NODE_SUCCESS = 'COVERAGE_NODE_SUCCESS';
export const COVERAGE_NODE_FAILURE = 'COVERAGE_NODE_FAILURE';

export function fetchNode(id) {
  return {
    [COVERAGE_API] : {
      types : [COVERAGE_NODE_REQUEST,COVERAGE_NODE_SUCCESS,COVERAGE_NODE_FAILURE],
      endpoint : `/tree/${id}`,
      data: { id },
    }
  }
}

export function loadNode(id) {
  return ((dispatch, getState) => {
    const node = selectNode(getState(), id)
    if (node && node.children && node.children.length == node.numChildren) {
      return null;
    } else if (node && node._children) {
      return dispatch(toggleNode(id));
    }

    return dispatch(fetchNode(id));
  });
}