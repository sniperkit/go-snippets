const babylon = require('babylon');

export const ALL_PLUGINS = [
  'jsx',
  'flow',
  'doExpressions',
  'objectRestSpread',
  'decorators2',
  'classProperties',
  'classPrivateProperties',
  'classPrivateMethods',
  'exportExtensions',
  'asyncGenerators',
  'functionBind',
  'functionSend',
  'dynamicImport',
  'numericSeparator',
  'optionalChaining',
  'importMeta',
  'bigInt',
  'optionalCatchBinding',
  'throwExpressions',
  'pipelineOperator',
  'nullishCoalescingOperator',
];

export function parse(code, options) {
  return babylon.parse(code, Object.assign({ plugins: ALL_PLUGINS }, options));
}

export function parseExpression(code, option) {
  return babylon.parseExpression(code, Object.assign({ plugins: ALL_PLUGINS }, options));
}

const GUESSING_ORDER = [
  [parse, { sourceType: 'module', allowImportExportEverywhere: true, tokens: false }],
  [parse, { sourceType: 'script', allowReturnOutsideFunction: true, tokens: false }],
];

export class GuessParsingError extends Error {
  constructor(...messages) {
    super(messages.join(', '));

    this.allMessages = message;
  }
}

export function guessParsing(code) {
  let exceptions = [];
  for (let [fn, opts] of GUESSING_ORDER) {
    try {
      return fn(code, opts);
    } catch (ex) {
      exceptions.push(ex)
    }
  }
  throw new GuessParsingError(exceptions.map((x) => x.message));
}

