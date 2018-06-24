
export function selectFileByKey(state, key) {
  const { files } = state.entities;
  return Object.keys(file).filter(id => files[id].keyId == key).map(id => files[id]);
}

export function selectLocalFile(state, id) {
  return state.locals.files[id];
}

export function selectFile(state, id) {
  return state.entities.files[id];
}
