import test from 'ava';
import {
  parse, 
  processRequest,
  collectComments,
  makeLineMap
} from '../src/parse';

test('parse correct code', t => {
  let { ast, diagnostics } = parse(`
    // Fancy comment

    function sum(a: number, b: number): number {
      return a + b;
    }
  `);

  t.is(diagnostics.length, 0);
  t.is(ast.kind, 'SourceFile');
  t.is(ast.statements.length, 1);
  t.is(ast.statements[0].kind, 'FunctionDeclaration');
  t.is(ast.comments.length, 1);
});

test('parse incorrect code', t => {
  let { ast, diagnostics } = parse(`
    def sum(a, b):
      return a + b
  `);

  t.true(diagnostics.length > 0);
  t.is(ast.kind, 'SourceFile');
});

test('parse includes jsDoc comments', t => {
  let { ast, diagnostics } = parse(`
  /**
   * Returns the sum of two numbers.
   * @function sum
   * @param {number} a - First number
   * @param {number} b - Second number
   * @returns {number} the sum
   */
  function sum(a: number, b: number): number {
    return a + b;
  }
  `);

  t.is(diagnostics.length, 0);
  t.is(ast.kind, 'SourceFile');
  t.is(ast.statements[0].jsDoc.length, 1);
  t.is(ast.statements[0].jsDoc[0].kind, 'JSDocComment');
});

test('collectComments only collects non-jsdoc', t => {
  let source = `// foo bar baz

  /* this is a multi
   * line comment
   */

  /**
   * this is jsdoc
   */

  // fooooo`;
  let lineMap = makeLineMap(source);
  let comments = collectComments(lineMap, source);
  
  t.is(comments.length, 3);
  let [ singleLine, multiLine, singleLineEof ] = comments;

  let assertComment = (cmt, text, line, col) => {
    t.is(cmt.text, text);
    t.is(cmt.line, line);
    t.is(cmt.col, col);
    t.is(cmt.kind, 'Comment');
  };
  assertComment(singleLine, 'foo bar baz', 1, 1);
  assertComment(multiLine, `this is a multi
   * line comment`, 3, 3);
  assertComment(singleLineEof, 'fooooo', 11, 3);
})

test('processRequest with non json input', t => {
  let result;
  try {
    result = JSON.parse(processRequest(`hello world`));
  } catch (e) {
    t.fail();
  }

  t.is(result.status, 'fatal');
  t.true(result.ast === null);
  t.true(result.errors.length > 0);
});

test('processRequest with invalid code', t => {
  let result;
  try {
    result = JSON.parse(processRequest(`{"content":"a = lambda x: x"}`));
  } catch (e) {
    t.fail();
  }

  t.is(result.status, 'error');
  t.true(result.ast === null);
  t.true(result.errors.length > 0);
});

test('processRequest with valid code', t => {
  let result;
  try {
    result = JSON.parse(processRequest(`{"content":"const a: number = 1;"}`));
  } catch (e) {
    t.fail();
  }

  t.is(result.status, 'ok');
  t.is(result.errors.length, 0);
  t.is(result.ast.kind, 'SourceFile');
  t.is(result.ast.statements.length, 1);
  t.is(result.ast.statements[0].kind, 'VariableStatement');
});

