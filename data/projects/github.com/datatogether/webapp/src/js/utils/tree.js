// walk a provide tree, calling fn on each node
// if fn returns a falsy value walk will halt
// and return the current node
export function walk(tree, fn, depth = 0) {
  if (!tree) { return; }
  fn(tree, depth)
  if (tree.children) {
    tree.children.forEach((c) => {
      walk(c,fn, depth+1)
    })
  }
}

export function widestDepthCount(tree) {
  let depths = [0];
  walk(tree, (node, depth) => {
    if (depths.length < depth + 1) {
      depths.push(0);
    }
    depths[depth]++;
  });

  return depths.reduce((w, a) => {
    return (a > w) ? a : w;
  }, 0);
}

export function maxDepthCount(tree) {
  let max = [0];
  walk(tree, (node, depth) => {
    if (max < depth) {
      max = depth;
    }
  });
  return max;
}

// copy a tree, keeping a branch if the node
// or any children pass the filter test
export function filter(tree, fn, depth = 0) {
  if (!tree) { return; }

  // copy the node
  let node = Object.assign({}, tree);
  let include = fn(node, depth);

  if (node.children) {
    node.children = node.children.reduce((children, node) => {
      node = filter(node, fn, depth + 1)
      if (node) {
        include || (include = true);
        children.push(node);
      }
      return children;
    }, []);
    if (node.children.length === 0) {
      node.children = undefined;
    }
  }

  return include ? node : undefined;
}

// walk a provide tree, calling fn on each node with a new copy
function copyWalk(tree, fn) {
  let copy = fn(copyNode(tree))

  if (copy.children) {
    for (let i = 0; i < copy.children.length; i++) {
      copy.children[i] = copyWalk(copy.children[i], fn);
    }
  }

  return copy;
}

function copyNode(node) {
  let copy = Object.assign({}, node);
  if (copy.children) {
    var a = [];
    // TODO - make faster
    copy.children.forEach(n => {
      a.push(n);
    });
    copy.children = a;
  }

  return copy;
}