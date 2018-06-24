import { fileSizeString, dateString, fileTypeString } from './format';

export function selectContent(state, hash) {
  return state.entities.content[hash];
}

export function selectContentConsensus(state, hash) {
  return state.entities.consensus[hash];
}

export function selectContentUrls(state, hash) {
  const { urls } = state.entities;

  return Object.keys(urls)
    .filter(u => urls[u].hash == hash)
    .map(url => urls[url]);
}

export function selectContentMetadata(state, hash) {
  const { metadata } = state.entities;
  return Object.keys(metadata)
    .filter(key => metadata[key].subject == hash)
    .map(key => metadata[key]);
}

export function selectRecentContentUrls(state) {
  const pages = state.pagination.contentRecentUrls;
  if (!pages.created) {
    return [];
  }
  return pages.created.ids.map(urlId => state.entities.urls[urlId]);
}

export function contentStats(url = {}) {
  return {
    size: fileSizeString(url.contentLength),
    "content type": fileTypeString(url.contentSniff),
    "last get": url.lastGet ? dateString(url.lastGet) : "never",
    status: url.status,
  };
}
