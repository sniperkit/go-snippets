import { fileSizeString, dateString } from './format';

export function concatUrlString(u, max = 80) {
  return (u.length < max) ? u : `${u.slice(0, max)}...`;
}

export function containsContent(url = {}) {
  return (url.contentSniff && url.contentSniff != "text/html; charset=utf-8");
}

export function selectUrl(state, urlString = "") {
  const { urls } = state.entities;
  const id = Object.keys(urls).find(u => (u == urlString));
  if (!id) { return undefined; }

  let url = urls[id];
  url.state = state.urlStates[urlString] || {};
  return url;
}

export function selectUrlByHash(state, hash = "") {
  const { urls } = state.entities;
  const id = Object.keys(urls).find(u => (urls[u].hash == hash));
  if (!id) { return undefined; }

  let url = urls[id];
  url.state = state.urlStates[url.url] || {};
  return url;
}

export function urlStats(url = {}) {
  return {
    "last get": url.lastGet ? dateString(url.lastGet) : "never",
    status: url.status,
    size: fileSizeString(url.contentLength),
    "content type": url.contentSniff,
  };
}

export function selectOutboundLinks(state, urlString = "") {
  const { urls, links } = state.entities;
  return Object.keys(links)
    .map(key => links[key])
    .filter(link => link.src.url == urlString)
    .map((link) => {
      let url = urls[link.dst.url] || link.dst;
      url.state = state.urlStates[link.dst.url] || {};
      return url;
    });
}

export function selectInboundLinks(state, urlString = "") {
  const { urls, links } = state.entities;
  return Object.keys(links)
    .map(key => links[key])
    .filter(link => link.dst.url == urlString)
    .map((link) => {
      let url = urls[link.src.url] || link.src;
      url.state = state.urlStates[link.dst.url] || {};
      return url;
    });
}
