import union from 'lodash/union';

// Creates a reducer managing pagination, given the action types to handle,
// and a function telling how to extract the key from an action.
const paginate = ({ types, mapActionToKey, removeIdsTypes=[], addIdsTypes=[] }) => {
  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected types to be an array of three elements.');
  }
  if (!types.every(t => typeof t === 'string')) {
    throw new Error('Expected types to be strings.');
  }
  if (typeof mapActionToKey !== 'function') {
    throw new Error('Expected mapActionToKey to be a function.');
  }

  const [requestType, successType, failureType] = types;

  const updatePagination = (state = {
    isFetching: false,
    // nextPageUrl: undefined,
    pageCount: 0,
    fetchedAll: false,
    ids: [],
  }, action) => {
    if (removeIdsTypes.includes(action.type) && action.response.result) {
      console.log(action.response.result);
      return Object.assign({}, state, {
        ids: state.ids.filter((id) => !action.response.result.includes(id)),
      });
    }

    if (addIdsTypes.includes(action.type) && action.response.result) {
      console.log(action.response.result);
      return Object.assign({}, state, {
        ids: union(state.ids, action.response.result),
      });
    }

    switch (action.type) {
      case requestType:
        return Object.assign({}, state, { isFetching: true });
      case successType:
        return Object.assign({}, state, {
          isFetching: false,
          ids: union(state.ids, action.response.result),
          // nextPageUrl: action.response.nextPageUrl,
          fetchedAll: (action.response.result.length < action.pageSize),
          pageCount: state.pageCount + 1,
        });
      case failureType:
        return Object.assign({}, state, {
          isFetching: false,
        });
      default:
        return state;
    }
  };

  function pageAction(state, action) {
    if (!action.page) { console.warn("expected page property for a pagainated request", action.type); }
    if (!action.pageSize) { console.warn("expected pageSize property for a pagainated request: ", action.type); }

    const key = mapActionToKey(action);
    if (typeof key !== 'string') {
      throw new Error('Expected key to be a string.');
    }
    return Object.assign({}, state, {
      [key]: updatePagination(state[key], action),
    });
  }

  return (state = {}, action) => {
    if (removeIdsTypes.includes(action.type) && action.response && action.response.result) {
      console.log(action.type);
      return pageAction(state, action);
    }
    if (addIdsTypes.includes(action.type) && action.response && action.response.result) {
      console.log(action.type);
      return pageAction(state, action);
    }

    // Update pagination by key
    switch (action.type) {
      case requestType:
        return pageAction(state, action);
      case successType:
        return pageAction(state, action);
      case failureType:
        return pageAction(state, action);
      default:
        return state;
    }
  };
};

export default paginate;
