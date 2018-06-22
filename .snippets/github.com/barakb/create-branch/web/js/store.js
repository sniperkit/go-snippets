import { createBranchRequest, branchAdded, setUser } from "./actions"
import { CreateBranch, CreateBranchContainer } from "./createBranch";
import { FilterableBranchesTableComponent } from "./branchesView";
import { rootReducer } from "./reducers";
import { thunkMiddleware as thunk} from './redux-thunk';
const { createStore, applyMiddleware} = Redux;
const { Provider } = ReactRedux;


export const store = createStore(rootReducer, applyMiddleware(thunk));


ReactDOM.render(
   <Provider store={store}>
        <FilterableBranchesTableComponent />
   </Provider>, document.getElementById('listBranches')
);

ReactDOM.render(
   <Provider store={store}>
      <CreateBranchContainer />
   </Provider>, document.getElementById('createBranch')
);


const isReadOnly = (branches, login) => {
    let ret = true;
    let mine = Object.keys(branches).filter(function(value){
        return branches[value] == login;
    });
    return mine.length < Object.keys(branches).length;
}

const loadBranches = () => {
    fetch('api/get_branches/', {
    	method: 'get',
    	credentials: 'same-origin',
    	cache: 'no-cache'
    }).then(function(response) {
        return response.json()
    }).then(function(res){
        store.dispatch(setUser(res.login));
        let branches = res.branches;
        let repos = res.repos;
        let isXAPRepo = function(repo){return !repo.startsWith("InsightEdge/")};
        let xapRepos = repos.filter(isXAPRepo);
        for (var key of Object.keys(branches)) {
                if(Object.keys(branches[key]).filter(isXAPRepo).length ==  xapRepos.length){
                    store.dispatch(branchAdded(key,  branches[key], isReadOnly(branches[key], document.currentLoginName)));
                }
        }
    }).catch(function(err) {
        console.info("err is ", err)
    });
}


loadBranches();