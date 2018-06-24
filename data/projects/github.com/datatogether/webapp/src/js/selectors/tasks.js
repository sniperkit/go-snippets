

export function selectTasks(state) {
  const { tasks } = state.entities;
  return Object.keys(tasks).map(k => tasks[k]);
}

export function selectTask(state, id = "") {
  if (!state.entities.tasks[id]) { return undefined; }
  const p = Object.assign({}, state.entities.tasks[id]);
  p.stats || (p.stats = {});
  if (p.subTasks) {
    p.subTasks = p.subTasks.map(subId => state.entities.tasks[subId]);
  }
  return p;
}

export function selectCollectionTasks(state, collectionId = "") {
  const { tasks } = state.entities;
  return Object.keys(tasks).filter(id => tasks[id].params.collectionId == collectionId).map(id => tasks[id]);
}

export function selectCollectionActiveTasks(state, collectionId = "") {
  const { tasks } = state.entities;
  return Object.keys(tasks).filter(id => (tasks[id].params.collectionId == collectionId && tasks[id].progress)).map(id => tasks[id]);
}

export function selectUrlActiveTasks(state) {
  const { tasks } = state.entities;
  return Object.keys(tasks).filter(id => (tasks[id].progress)).map(id => tasks[id]);
}