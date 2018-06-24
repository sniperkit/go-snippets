
export function selectLocalUncrawlable(state, id) {
  return state.locals.uncrawlables[id];
}

export function selectUncrawlable(state, id = "") {
  return state.entities.uncrawlables[id];
}

export function selectUncrawlables(state) {
  const { uncrawlables } = state.entities;
  return Object.keys(uncrawlables).map(k => uncrawlables[k]);
}

