'use strict';

/**
 * @ngdoc function
 * @name foodPlanOrganizerApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the foodPlanOrganizerApp
 */

function getMonday(d) {
  d = new Date(d);
  var day = d.getDay(),
    diff = d.getDate() - day + (day === 0 ? -6 : 1);
  return new Date(d.setDate(diff));
}
var getWeek = function(weekDiff) {
  var startDate = new Date();
  startDate = startDate.setDate(startDate.getDate() + weekDiff * 7);
  startDate = getMonday(startDate);

  var dates = [];
  for (var d = 0; d < 7; d++) {
    var date = new Date(startDate.getTime());
    date.setDate(date.getDate() + d);
    dates.push(date);
  }
  return dates;
}

angular.module('foodPlanOrganizerApp')
.controller('PlanningCtrl', function($scope) {
  $scope.weeks = [getWeek(0), getWeek(1), getWeek(2)];
});