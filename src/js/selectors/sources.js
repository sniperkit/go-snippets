

export function selectSources(state) {
  const { sources } = state.entities;
  return Object.keys(sources).map(k => sources[k]);
}

export function selectSource(state, id = "") {
  return state.entities.sources[id];
}

export function selectSourceUndescribedUrls(state, id = "") {
  const pages = state.pagination.sourceUndescribedUrls;
  if (!pages[id]) {
    return [];
  }

  return pages[id].ids.map(urlId => state.entities.urls[urlId]);
}

export function selectRecentSources(state) {
  const pages = state.pagination.sources;
  if (!pages.created) {
    return [];
  }
  return pages.created.ids.map(id => state.entities.sources[id]);
}

export function selectSourceAttributedUrls(state, id = "") {
  const pages = state.pagination.sourceAttributedUrls;
  if (!pages[id]) {
    return [];
  }
  return pages[id].ids.map(urlId => state.entities.urls[urlId]);
}
