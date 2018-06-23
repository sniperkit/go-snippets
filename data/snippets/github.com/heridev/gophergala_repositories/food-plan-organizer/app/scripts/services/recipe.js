'use strict';
angular.module('foodPlanOrganizerApp')
.factory('Recipe', function($resource, Settings) {
  return $resource(Settings.host() + '/recipes/', {
    id: '@id'
  }, {
    update: {
      method: 'PUT'
    }
  });
});