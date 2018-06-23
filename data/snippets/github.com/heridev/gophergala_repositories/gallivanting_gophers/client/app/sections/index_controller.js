(function() {
  'use strict';

  angular.module('gg.index', [])
    .controller('indexCtrl', ['$log', '$scope', IndexCtrl]);

  function IndexCtrl($log, $scope) {
    $log.log('Index Controller Started!');

    $scope.updateTrigger = function(data) {
      $scope.boardData = data;
    };
  }
})();
