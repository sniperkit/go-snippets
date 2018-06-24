
export function selectLocalCollection (state, id) {
  return state.locals.collections[id]
}

export function selectCollection (state, id = '') {
  return state.entities.collections[id]
}

export function selectCollections (state) {
  const { collections } = state.entities
  return Object.keys(collections).map(k => collections[k])
}

export function selectRecentCollections (state, count = 3) {
  const { collections } = state.entities
  return Object.keys(collections).map(k => collections[k]).slice(0, count)
}

export function selectCollectionItems (state, id) {
  const { collectionItems } = state.entities
  const ci = state.pagination.collectionItems[id]
  if (ci) {
    return ci.ids.map(id => collectionItems[id])
  }

  return []
}

export function selectCollectionsByKey (state, key) {
  const { collections } = state.entities
  return Object.keys(collections).filter(id => collections[id].creator == key).map(id => collections[id])
}
