import { 
  COVERAGE_NODE_TOGGLE,
  COVERAGE_NODE_SUCCESS,
  // COVERAGE_USER_FAILURE,
  // COVERAGE_LOGIN_SUCCESS,
  // COVERAGE_ADD_HISTORY_ENTRY,
} from '../actions/coverage';

// initial state is an empty tree
const initialState = {
  tree : undefined,
};

export default function coverageReducer(state = initialState, action) {
  switch (action.type) {
    case COVERAGE_NODE_TOGGLE:
      return {
        tree: copyWalk(state.tree, (n) => {
          if (n && n.id != action.id) {
            return n;
          } else {
            if (n.children) {
              return Object.assign({}, n, { _children: n.children, children: undefined });
            } else if (n._children) {
              return Object.assign({}, n, { children: n._children, _children: undefined });
            } else {
              return n;
            }
          }
        })
      }
    case COVERAGE_NODE_SUCCESS:
      if (action.tree.id == "root") {
        return {
          tree: action.tree,
        }
      }
      return Object.assign({}, state, { 
        tree: copyWalk(state.tree, (n) => (n && n.id == action.tree.id) ? action.tree : n)
      });
    default:
      return state;
  }
}

// walk a provide tree, calling fn on each node with a new copy
function copyWalk(tree, fn) {
  let copy = fn(copyNode(tree))

  if (copy.children) {
    for (let i = 0; i < copy.children.length; i++) {
      copy.children[i] = copyWalk(copy.children[i], fn);
    }
  }

  return copy;
}

function copyNode(node) {
  let copy = Object.assign({}, node);
  if (copy.children) {
    var a = [];
    // TODO - make faster
    copy.children.forEach(n => {
      a.push(n);
    });
    copy.children = a;
  }

  return copy;
}


