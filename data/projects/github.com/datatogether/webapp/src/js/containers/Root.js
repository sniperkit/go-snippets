/* global __BUILD__ */
/* eslint-disable global-require */
if (__BUILD__.PRODUCTION) {
  module.exports = require('./Root.prod');
} else if (__BUILD__.STAGING) {
  module.exports = require('./Root.prod');
} else if (__BUILD__.DEVELOP) {
  module.exports = require('./Root.dev');
}
/* eslint-enable global-require */
