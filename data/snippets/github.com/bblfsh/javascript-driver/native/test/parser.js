import test from 'ava';
import { parse, parseExpression } from '../lib';

test('exports a parse function', t => {
  t.is(typeof parse, 'function');
});

test('exports a parseExpression function', t => {
  t.is(typeof parseExpression, 'function');
});
