import { walk, filter } from '../utils/tree';

const completionColors = ["#ee005c", "#e01f54", "#d43a4e", "#c85647", "#bc7040", "#b08a3a", "#a4a533", "#98c02c", "#8adf24"];
export function completionColor(node) {
  if (!node.numLeaves) {
    return node.coverage && node.coverage.length ? completionColors[completionColors.length - 1] : completionColors[0];
  }
  const compl = ((node.numLeavesArchived || 0) / node.numLeaves);
  if (compl == 1) return completionColors[completionColors.length - 1];
  return completionColors[Math.floor(compl * completionColors.length)];
}

export function searchTree(tree, query) {
  return filter(tree, (n) => ~~n.name.indexOf(query) >= 0);
}

export function selectNode(state,id) {
  const { tree } = state.coverage;
  let node;
  walk(tree, (n) => {
    if (n && n.id == id) {
      node = n;
    }
  });
  return node;
}

export function flattenTree(tree) {
  let links = [], nodes = [];

  if (tree == undefined) {
    return {
      links,
      nodes
    };
  }

  walk(tree,(n) => {
    nodes.push({
      id : n.id,
      name: n.name,
      radius: 5,
      numLeaves: n.numLeaves,
      numLeavesArchived: n.numLeavesArchived,
      numChildren: n.numChildren,
      coverage: n.coverage,
      collapsed: n.collapsed,
    });

    if (n.children) {
      n.children.forEach((c) => {
        links.push({
          source: n.id,
          target: c.id,
        });
      })
    }

    return true;
  });

  return {
    links,
    nodes,
  }  
}
