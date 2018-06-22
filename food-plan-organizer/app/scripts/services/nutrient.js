'use strict';
angular.module('foodPlanOrganizerApp')
.factory('Nutrient', function($resource, Settings) {
  return $resource(Settings.host() + '/nutrients/');
});