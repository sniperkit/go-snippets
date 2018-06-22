import Ember from 'ember';

export function formatDate(input, options) {
  var format = options.hash['format'] || 'llll';
  return moment(input).format(format);
}

export default Ember.Handlebars.makeBoundHelper(formatDate);
