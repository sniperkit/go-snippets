import test from 'ava';
import { handler } from '../lib';

function request(content) {
  return JSON.stringify({ content });
}

function responseFor(content) {
  return JSON.parse(handler(request(content)));
}

test('returns an error response if the code cannot be parsed', t => {
  let resp = responseFor("a +% b");

  t.is(resp.status, "error");
  t.true(resp.errors.length > 0);
});


test('returns an ok response with an ast', t => {
  let resp = responseFor("var a = function() {};");

  t.is(resp.status, "ok");
  t.true("ast" in resp);
});

test('accepts module syntax', t => {
  let resp = responseFor("export default 42;");

  t.is(resp.status, "ok");
  t.true("ast" in resp);
});
