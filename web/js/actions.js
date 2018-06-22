import { CREATE_BRANCH_REQUEST,
         UPDATE_BRANCHES_FILTER,
         DELETE_BRANCH_REQUEST,
         BRANCH_ADDED,
         CREATE_BRANCH_RESPONSE,
         DELETE_BRANCH_RESPONSE,
         BRANCH_DELETED,
         TOGGLED_BRANCH_ROW,
         TOGGLE_SOURCE_BRANCH,
         SET_USER,
         TOGGLE_CREATE_ONLY_XAP_BRANCH
  } from "./actionTypes"



export function updateBranchesFilterRequest(filterText) {
//    console.info("updateBranchesFilterRequest", filterText)
    return {type:UPDATE_BRANCHES_FILTER, filterText};
}

export function deleteBranchRequest(name) {
//    console.info("deleteBranchRequest", name);
    return function(dispatch){
        if (!name){
            return
        }
        dispatch({type:DELETE_BRANCH_REQUEST, name});
        fetch('/api/delete_branch/' + name, {
                    method: 'delete',
                    credentials: 'same-origin',
                    cache: 'no-cache'
                }).then(function(response) {
                    return response.json()
                }).then(function(branch){
                    dispatch(deleteBranchResponse(name, null))
                }).catch(function(err) {
                    console.info("err is ", err)
                });
    }
}

export function deleteBranchResponse(name, err){
//    console.info("deleteBranchResponse", name, err)
    return function(dispatch){
        dispatch({type:DELETE_BRANCH_RESPONSE, name, err:err});
        if(!err){
            dispatch(branchDeleted(name));
        }
    }
}

export function branchAdded(name, repositories, readOnly) {
    return {type:BRANCH_ADDED, name, repositories, readOnly};
}

export function setUser(login) {
    return {type:SET_USER, login};
}

export function branchDeleted(name) {
//    console.info("branch deleted", name)
    return {type:BRANCH_DELETED, name};
}

export function createBranchRequest(name, from, isXAPOnly) {
    from = !from ? "master" : from;
    return function(dispatch){
        if (!name){
            return
        }
        const branchName = name;
        dispatch({type:CREATE_BRANCH_REQUEST, branchName});
        fetch('/api/create_branch/' + branchName +'?from=' + encodeURIComponent(from) + '&isXAPOnly=' + encodeURIComponent(isXAPOnly), {
            method: 'get',
            credentials: 'same-origin',
            cache: 'no-cache'
        }).then(function(response) {
            return response.json()
        }).then(function(statuses){
            dispatch(createBranchResponse(branchName, statuses, null))
        }).catch(function(err) {
            console.info("err is ", err)
        });
    }
}

export function createBranchResponse(name, statuses, err) {
    return function(dispatch){
        dispatch({type:CREATE_BRANCH_RESPONSE, name, statuses, err:err});
        if(!err){
            dispatch(branchAdded(name, Object.keys(statuses).filter(k => statuses[k])));
        }
    }
}


export function toggleRow(name){
    return { "type":TOGGLED_BRANCH_ROW, "name":name };
}

export function toggleSourceBranch(name){
    return { "type":TOGGLE_SOURCE_BRANCH, "name":name };
}
export function toggleXAPRequest(){
    return { "type":TOGGLE_CREATE_ONLY_XAP_BRANCH};
}
