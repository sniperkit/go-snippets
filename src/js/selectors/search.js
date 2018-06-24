
export function selectSearchQuery(state) {
  return state.app.query;
}

export function selectSearchResults(state) {
  const query = selectSearchQuery(state);
  if (!query) { return []; }

  const { searchResults } = state.entities;
  return Object.keys(searchResults).map(key => searchResults[key]).filter((res) => {
    return (res.url.indexOf(query) >= 0);
  });
}
