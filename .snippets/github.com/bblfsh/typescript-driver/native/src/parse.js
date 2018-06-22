'use strict';

const {
  createSourceFile,
  ScriptTarget,
  forEachChild,
  NodeFlags,
  SyntaxKind
} = require('typescript');

/**
 * TypeScript enumerations are transpiled to an object with the following form:
 * {
 *  '0': 'NameOfFirstValue',
 *  ...
 *  'N': 'NameOfLastValue',
 *  'NameOfFirstValue': 0,
 *  ...
 *  'NameOfLastValue': n,
 * }
 *
 * Because of that, we just have to skip the first N keys to get an array with
 * the names of all properties of an enum.
 * SyntaxKind has a property called `Count`, which is the last value of the
 * enum. That give us the starting point without having to iterate the array.
 *
 * @constant
 * @type {Array<string>}
 * @see {@link https://github.com/Microsoft/TypeScript/blob/a29e8cf2406dd41aea48be0fb68cde6a2e972564/src/compiler/types.ts#L52-L428}
 */
const nodeKinds = Object.keys(SyntaxKind).slice(SyntaxKind.Count+1);

/**
 * Parses the given TypeScript source code and returns a string with
 * the AST encoded as JSON.
 *
 * @function parse
 * @param {string} source - TypeScript source code to parse
 * @returns {Object} the representation of the parsed AST
 */
function parse(source) {
  let sourceFile = createSourceFile(
    'input.ts',
    source,
    ScriptTarget.ES6,
    false
  );

  const lineMap = makeLineMap(source);
  sourceFile.comments = collectComments(lineMap, source);

  const diagnostics = (sourceFile.parseDiagnostics || [])
    .filter(d => d.category === 1)
    .map(d => `at ${d.start}:${d.start+d.length}: ${d.messageText}`);

  // normalization of the AST and remove some properties that are not
  // useful in the output. Text we already had and fileName is not real
  // so there's no point in keeping it
  normalizeNode(lineMap, sourceFile);
  delete sourceFile.endOfFileToken;
  delete sourceFile.text;
  delete sourceFile.fileName;
  delete sourceFile.parseDiagnostics;

  return {
    ast: sourceFile,
    diagnostics
  };
}

/**
 * Creates a new mapping between lines and offsets from the source code.
 *
 * @function makeLineMap
 * @param {string} source - source code
 * @return {Array<Object>}
 */
function makeLineMap(source) {
  const src = Array.from(source);
  let lines = [];
  let offset = -1;
  let i, len;
  let line = 1;

  for (i = 0, len = src.length; i < len; i++) {
    if (offset < 0) {
      offset = i;
    }

    if (src[i] == '\n' || src[i] == '\r') {
      lines.push({ start: offset, end: i, line });
      line++;
      offset = -1;
    }
  }
  
  if (offset >= 0) {
    lines.push({ start: offset, end: i, line});
  }

  return lines;
}

/**
 * Finds the line and column of a node given its pos and the line map.
 *
 * @function findNodePos
 * @param {Array<Object>} lineMap - Mapping of offsets to lines.
 * @param {number} pos - Offset of the token start
 * @return {Object} line and col
 */
function findNodePos(lineMap, pos) {
  let l = -1, col = -1;
  for (let i = 0, len = lineMap.length; i < len; i++) {
    const { start, end, line } = lineMap[i];
    if (start <= pos && end >= pos) {
      l = line;
      col = pos - start + 1;
      break;
    }
  }

  return { line: l, col };
}

/**
 * Collects all non-JSDoc comments from the source file.
 *
 * @function collectComments
 * @param {string} source - Source file content
 * @returns {Array<ts.Node>}
 */
function collectComments(lineMap, source) {
  const src = Array.from(source);
  let comments = [];
  for (let i = 0, len = src.length; i < len; i++) {
    const c = src[i];
    const offset = i;
    if (c === '/' && i+1 < len) {
      let next = src[++i];
      if (next === '/') {
        while (src[i] !== '\n' && i < len) {
          i++;
        }
        let comment = source.substring(offset, i);
        comments.push(commentNode(lineMap, offset, comment));
      } else if (next === '*') {
        if (i+1 < len) {
          next = src[++i];
          // check if it's JSDoc, ignore it, then, because ts already parses it
          if (next === '*') {
            continue;
          }

          while (i+1 < len && !(src[i] === '*' && src[i+1] === '/')) {
            i++;
          }

          if (i+1 < len && src[i] === '*' && src[i+1] === '/') {
            let comment = source.substring(offset, i+2);
            comments.push(commentNode(lineMap, offset, comment));
          }
        }
      } else {
        i--;  
      }
    }
  }

  return comments;
}

/**
 * Creates a new comment node.
 *
 * @function commentNode
 * @param {number} offset - Offset since the start of the file until the first
 * character of the comment.
 * @param {string} comment - Comment text containing /* *\/ and //.
 * @returns {Object} comment node
 */
function commentNode(lineMap, offset, comment) {
  let text, pos, end;
  if (comment.indexOf('/*') >= 0) {
    pos = comment.indexOf('/*');
    end = comment.indexOf('*/')+2;
    text = comment.substring(pos+2, end-2).trim();
  } else {
    pos = comment.indexOf('//');
    end = comment.length - pos;
    text = comment.substring(pos+2).trim();
  }

  const { line, col } = findNodePos(lineMap, offset+pos);
  return { kind: 'Comment', text, line, col };
}

/**
 * Recursively normalizes the node and all its children. By normalization 
 * we mean replace the kind of the node by a string, instead of an integer
 * and also we remove the flags from the node.
 *
 * @function normalizeNode
 * @param {ts.Node} node - Node to normalize
 */
function normalizeNode(lineMap, node) {
  forEachChild(node, n => normalizeNode(lineMap, n));
  // node.kind substitution MUST be done after normalizing the children
  // otherwise forEachChild does not know how to access them
  node.kind = nodeKinds[node.kind];
  const { line, col } = findNodePos(lineMap, node.pos);
  node.line = line;
  node.col = col;
  if (node.jsDoc) {
    node.jsDoc = node.jsDoc.map(d => {
      normalizeNode(lineMap, d);
      return d;
    });
  }
  delete node.parent;
  delete node.transformFlags;

  if (node.flags) {
    node.flags = transformFlags(node.flags);
  } else {
    node.flags = [];
  }
}

/**
 * Transforms binary flags to human readable flags.
 *
 * @function transformFlags
 * @param {number} flags - Binary flags
 * @returns {Array<string>}
 */
function transformFlags(flags) {
  let result = [];
  for (let flag in NodeFlags) {
    if (NodeFlags.hasOwnProperty(flag) && (flags & NodeFlags[flag])) {
      result.push(flag);
    }
  }
  return result;
}

const OK = 'ok';
const FATAL = 'fatal';
const ERROR = 'error';

/**
 * Builds the result, either a successful or an errored one.
 *
 * @function printResult
 * @param {string} status - one of OK, FATAL or ERROR
 * @param {Object|null} ast - parsed AST, null if status is not ok
 * @param {Array<string>} diagnostics - List of diagnostic messages
 * @return {string} JSON encoded result
 */
function result(status, ast, ...diagnostics) {
  const output = {
    status,
    ast,
    errors: diagnostics.map(m => ({ message: m }))
  }

  try {
    return JSON.stringify(output, (k, v) => k === 'parent' ? null : v);
  } catch (e) {
    return result(FATAL, null, `Unable to encode response to JSON: ${e.message}`);
  }
}

/**
 * Processes one single request to parse a TypeScript source file. Will print a
 * fatal or error result if the process fails, a valid one with the AST in any
 * other case.
 *
 * @function processRequest
 * @param {string} line - Line read from stdin containing the request
 * @param {string} the JSON encoded response output
 */
function processRequest(line) {
  let req;
  try {
    req = JSON.parse(line);
  } catch (e) {
    return result(FATAL, null, `Unable to decode request from JSON: ${e.message}`);
  }

  const { ast, diagnostics } = parse(req.content || '');
  if (diagnostics.length > 0) {
    return result(ERROR, null, ...diagnostics);
  }

  return result(OK, ast);
}

module.exports = {
  parse,
  processRequest,
  collectComments,
  makeLineMap
};
