

export function selectPrimers(state) {
  const { primers } = state.entities;
  return Object.keys(primers).map(k => primers[k]);
}

export function selectPrimer(state, id = "") {
  if (!state.entities.primers[id]) { return undefined; }
  const p = Object.assign({}, state.entities.primers[id]);
  p.stats || (p.stats = {});
  if (p.subPrimers) {
    p.subPrimers = p.subPrimers.map(subId => state.entities.primers[subId]);
  }
  return p;
}
