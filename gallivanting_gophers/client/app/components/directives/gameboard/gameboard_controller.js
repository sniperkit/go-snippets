(function() {
  'use strict';

  angular.module('gg.components.gameboard', [])
  .controller('gameboardCtrl', ['$log', '$scope', 'gameService', gameboardCtrl]);

  function gameboardCtrl($log, $scope, gameService) {
    $log.log('Controller Started - Gameboard!');

    gameService.getBoard().then(function(data){
      $scope.updateBoard(data);
    });

    $scope.gopherSelected = function(id) {
      gameService.selectedGopher = id;
    };
  }
})();
