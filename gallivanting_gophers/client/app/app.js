(function() {
  'use strict';

  angular.module('gg', [
      'ui.router',
      'ui.bootstrap',
      'templates',
      'gg.components',
      'gg.index'
    ])
    .config(appConfig);

  function appConfig($stateProvider, $urlRouterProvider, $httpProvider, $urlMatcherFactoryProvider) {
    // $httpProvider.interceptors.push('authInterceptor');
    // $httpProvider.defaults.headers.post = {
    //   'Content-Type': 'application/json;charset=utf-8'
    // };

    $urlRouterProvider.otherwise('/');
    $urlMatcherFactoryProvider.strictMode(false);
  }

})();
