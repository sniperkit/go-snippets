import test from 'ava';
import { error, fatal, ok } from '../lib';

test('error returns a valid response with error status and the given message', t => {
  let res = error(['not valid', 'really not valid']);

  t.is(res.status, 'error', 'status is error');
  t.true(res.errors.includes('not valid'), 'contains the first passed message');
  t.true(res.errors.includes('really not valid'), 'contains the second passed message');
});

test('fatal returns a valid response with fatal status and the given message', t => {
  let res = fatal(new Error('really not valid'));

  t.is(res.status, 'fatal', 'status is fatal');
  t.true(res.errors.includes('really not valid'), 'contains the passed message');
});

test('ok returns a valid response with ok status and the given ast', t => {
  let ast = {node: 'a', children: []};
  let res = ok(ast);

  t.is(res.status, 'ok', 'status is ok');
  t.is(res.ast, ast, 'ast is the passed ast');
});
