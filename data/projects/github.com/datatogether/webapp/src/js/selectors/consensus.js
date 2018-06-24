

export function selectConsensus(state, subject) {
  return state.entities.consensus[subject];
}

export function selectAllConsensus(state) {
  return Object.keys(state.entities.consensus).map(hash => state.entities.consensus[hash]);
}
