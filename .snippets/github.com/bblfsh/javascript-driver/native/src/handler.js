import { guessParsing, GuessParsingError } from './parser';
import { error, ok } from './response';

function parse(data) {
  try {
    let { content } = JSON.parse(data);
    let ast = guessParsing(content);
    return ok(ast);
  } catch (ex) {
    if (ex instanceof GuessParsingError) {
      return error(ex.allMessages);
    }
    return error([ex.message]);
  }
}

function jsonify(data) {
  return `${JSON.stringify(data)}`;
}

export function handler(input) {
  return jsonify(parse(input));
}
