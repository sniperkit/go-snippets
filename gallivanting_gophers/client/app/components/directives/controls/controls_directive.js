angular.module('gg.components.controls')
.directive('ggControls', function() {
  return {
    restrict: 'E',
    scope: {
      updateTrigger: '=',
      gopher: '='
    },
    controller: 'controlsCtrl',
    templateUrl: 'components/directives/controls/controls.html',
    replace: true
  };
});
