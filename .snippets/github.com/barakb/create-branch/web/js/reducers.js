const { combineReducers } = Redux;
import { CREATE_BRANCH_REQUEST,
         CREATE_BRANCH_RESPONSE,
         UPDATE_BRANCHES_FILTER,
         DELETE_BRANCH_REQUEST,
         BRANCH_ADDED,
         DELETE_BRANCH_RESPONSE,
         BRANCH_DELETED,
         TOGGLED_BRANCH_ROW,
         TOGGLE_SOURCE_BRANCH,
         TOGGLE_CREATE_ONLY_XAP_BRANCH,
         SET_USER,
        }
 from "./actionTypes"

const createBranchInitialState = {name:'', isFetching:false, fromBranch : '', isXAPOnlyBranch:false};

function createBranch (state = createBranchInitialState , action) {
   switch (action.type){
       case CREATE_BRANCH_REQUEST:
           return {...state, isFetching:true};
       case CREATE_BRANCH_RESPONSE:
           return {...state, isFetching:false};
      case TOGGLE_SOURCE_BRANCH:
           let fromBranch = state.fromBranch === action.name ? '' : action.name
           return {...state, fromBranch};
      case TOGGLE_CREATE_ONLY_XAP_BRANCH:
           let isXAPOnlyBranch = !state.isXAPOnlyBranch
           console.info("TOGGLE_CREATE_ONLY_XAP_BRANCH state is ", state , " Next is ", isXAPOnlyBranch)
           return {...state, isXAPOnlyBranch};
       default:
           return state;
   }
}

const branches = [
             ];

const viewBranchesInitialState = {filterText:'', branches, filtered : [], login:""};

function viewBranch (state = viewBranchesInitialState , action) {
   switch (action.type){
      case UPDATE_BRANCHES_FILTER:
           let s = { ...state, filterText: action.filterText, filtered : state.branches.filter(b => -1 < b.name.indexOf(action.filterText))};
           return s;
      case DELETE_BRANCH_REQUEST:
           let branches = state.branches.map(b => b.name === action.name ? { ...b, isDeleting:true} : b)
           branches.sort((b1, b2) => b1.name > b2.name);
           let filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText))
          return  { ...state, branches, filtered };
      case DELETE_BRANCH_RESPONSE:
           branches = state.branches.map(b => b.name === action.name ? {...b, isDeleting:false} : b)
           branches.sort((b1, b2) => b1.name > b2.name);
           filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText))
          return  { ...state, branches, filtered };
      case BRANCH_ADDED:
           branches = [ ...state.branches, {name:action.name, repositories:action.repositories, expanded: false, readOnly:action.readOnly, isSource:false}];
           branches.sort((b1, b2) => b1.name > b2.name);
           filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText))
          return  { ...state, branches, filtered };
      case BRANCH_DELETED:
           branches = state.branches.filter(b => b.name !== action.name)
           branches.sort((b1, b2) => b1.name > b2.name);
           filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText))
          return  { ...state, branches, filtered };
       case TOGGLED_BRANCH_ROW:
           branches = state.branches.map( b => b.name === action.name ? {...b, expanded:!b.expanded} : b );
           branches.sort((b1, b2) => b1.name > b2.name);
           filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText))
           return  { ...state, branches, filtered };
       case TOGGLE_SOURCE_BRANCH:
           branches = state.branches.map( b => b.name === action.name ? {...b, isSource:!b.isSource} : {...b, isSource:false} );
           filtered = branches.filter(b => -1 < b.name.indexOf(state.filterText));
           return  { ...state, branches, filtered};
       case SET_USER:
           let login = action.login
           return  { ...state, login};
       default:
           return { ...state, filtered : state.branches.filter(b => -1 < b.name.indexOf(state.filterText))};
   }
}

export const rootReducer = combineReducers({
  createBranch,
  viewBranch
})

