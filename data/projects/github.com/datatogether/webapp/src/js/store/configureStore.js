/* globals __BUILD__ */
/* eslint-disable global-require */
if (__BUILD__.PRODUCTION) {
  module.exports = require('./configureStore.prod');
} else if (__BUILD__.STAGING) {
  module.exports = require('./configureStore.prod');
} else if (__BUILD__.DEVELOP) {
  module.exports = require('./configureStore.dev');
}
/* eslint-enable global-require */
