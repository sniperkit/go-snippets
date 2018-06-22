(function() {
  'use strict';

  angular.module('gg.components.controls', [])
  .controller('controlsCtrl', ['$log', '$scope', 'gameService', controlsCtrl]);

  function controlsCtrl($log, $scope, gameService) {
    $log.log('Controller Started - Controls!');

    $scope.direction = 1;

    $scope.selectedGopher = function() {
      return gameService.selectedGopher;
    }

    $scope.createBoard = function() {
      gameService.createBoard().then(function(data) {
        $scope.updateTrigger(data);
      });
    };

    $scope.moveGopher = function() {
      gameService.moveGopher($scope.direction).then(function(data){
        $scope.updateTrigger(data);
      });
    };

  }
})();
