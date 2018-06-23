'use strict';
angular.module('foodPlanOrganizerApp')
.factory('Food', function($resource, Settings) {
  return $resource(Settings.host() + '/food/', {
    id: '@id'
  });
});