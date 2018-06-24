import { CALL_API } from '../middleware/api';
import Schemas from '../schemas';

import analytics from '../analytics';
import { newLocalModel, updateLocalModel, editModel, removeLocalModel } from './locals';

const blankTask = {
  meta: {
    title: "",
    description: "",
  },
};

const TASK_NEW = "TASK_NEW";
export function newTask(keyId) {
  return newLocalModel(Schemas.TASK, TASK_NEW, Object.assign({}, blankTask, { keyId }));
}

const TASK_EDIT = "TASK_EDIT";
export function editTask(task) {
  return editModel(Schemas.TASK, TASK_EDIT, Object.assign({}, blankTask, task));
}

const TASK_CANCEL_EDIT = "TASK_CANCEL_EDIT";
export function cancelTaskEdit(task) {
  return removeLocalModel(Schemas.TASK, TASK_CANCEL_EDIT, task);
}

const TASK_UPDATE = "TASK_UPDATE";
export function updateTask(task) {
  return updateLocalModel(Schemas.TASK, TASK_UPDATE, task);
}

export const TASKS_FETCH_REQUEST = "TASKS_FETCH_REQUEST";
export const TASKS_FETCH_SUCCESS = "TASKS_FETCH_SUCCESS";
export const TASKS_FETCH_FAILURE = "TASKS_FETCH_FAILURE";

export function fetchTasks(page = 1, pageSize = 25) {
  return {
    [CALL_API]: {
      types: [TASKS_FETCH_REQUEST, TASKS_FETCH_SUCCESS, TASKS_FETCH_FAILURE],
      schema: Schemas.TASK,
      endpoint: '/tasks',
      data: { page, pageSize },
    },
  };
}

export function loadTasks(page = 1, pageSize = 25) {
  return (dispatch) => {
    // TODO - add pagination & pagination check
    return dispatch(fetchTasks(page, pageSize));
  };
}


export const TASK_FETCH_REQUEST = "TASK_FETCH_REQUEST";
export const TASK_FETCH_SUCCESS = "TASK_FETCH_SUCCESS";
export const TASK_FETCH_FAILURE = "TASK_FETCH_FAILURE";

export function fetchTask(id = "") {
  return {
    [CALL_API]: {
      types: [TASK_FETCH_REQUEST, TASK_FETCH_SUCCESS, TASK_FETCH_FAILURE],
      schema: Schemas.TASK_ARRAY,
      endpoint: `/tasks/${id}`,
      data: { id },
    },
  };
}

export function loadTask(id = "") {
  return (dispatch) => {
    // TODO - check for local url copy via getState
    return dispatch(fetchTask(id));
  };
}

export const TASK_ENQUEUE_REQUEST = "TASK_ENQUEUE_REQUEST";
export const TASK_ENQUEUE_SUCCESS = "TASK_ENQUEUE_SUCCESS";
export const TASK_ENQUEUE_FAILURE = "TASK_ENQUEUE_FAILURE";

export function enqueueTask(task = {}) {
  return (dispatch) => {
    analytics.track("Enqueued Task", task);
    return dispatch({
      [CALL_API]: {
        types: [TASK_ENQUEUE_REQUEST, TASK_ENQUEUE_SUCCESS, TASK_ENQUEUE_FAILURE],
        schema: Schemas.TASK,
        method: "POST",
        endpoint: "/tasks",
        data: task,
      },
    }).then((action) => {
      if (action.type == TASK_ENQUEUE_SUCCESS) {
        cancelTaskEdit(task.keyId, task.subject);
      }
    });
  };
}

export const TASK_CANCEL_REQUEST = "TASK_CANCEL_REQUEST";
export const TASK_CANCEL_SUCCESS = "TASK_CANCEL_SUCCESS";
export const TASK_CANCEL_FAILURE = "TASK_CANCEL_FAILURE";

export function cancelTask(task = {}) {
  return (dispatch) => {
    analytics.track("Cancelled Task", task);
    return dispatch({
      [CALL_API]: {
        types: [TASK_CANCEL_REQUEST, TASK_CANCEL_SUCCESS, TASK_CANCEL_FAILURE],
        schema: Schemas.TASK,
        endpoint: `/tasks/${task.id}`,
        method: "DELETE",
        data: task,
      },
    });
  };
}

export function archiveCollection(collection) {
  return enqueueTask({
    title: `Archive collection to IPFS: ${collection.title}`,
    taskType: "ipfs.addcollection",
    params: {
      collectionId : collection.id,
    },
  });
}

export function archiveUrl(url) {
  return enqueueTask({
    title: `Archive url ${url}`,
    taskType: "ipfs.addurl",
    params: {
      url,
    },
  });
}