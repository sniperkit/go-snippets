(function () {
  'use strict';

  angular.module('gg.components.game')
  .factory('gameService', ['$log', '$q', '$http', GameService]);

  function GameService($log, $q, $http) {
    var gameService = {};

    $log.log('Game Service - Started');

    gameService.selectedGopher = -1;

    gameService.getBoard = function() {
      return _getBoard($log, $q, $http);
    };

    gameService.createBoard = function() {
      return _createBoard($log, $q, $http);
    };

    gameService.moveGopher = function(direction) {
      return _moveGopher(gameService.selectedGopher, direction, $log, $q, $http);
    };

    return gameService;
  }

  function _getBoard($log, $q, $http) {
    var deferred = $q.defer();

    $http.get('http://localhost:8888/board').success(function (response) {
      deferred.resolve(response);
    }).error(function (response) {
      deferred.reject(response);
    });

    return deferred.promise;
  }

  function _createBoard($log, $q, $http) {
    var deferred = $q.defer();

    $http.post('http://localhost:8888/board').success(function (response) {
      deferred.resolve(response);
    }).error(function (response) {
      deferred.reject(response);
    });

    return deferred.promise;
  }

  function _moveGopher(gopher, direction, $log, $q, $http) {
    var deferred = $q.defer();

    $http.post('http://localhost:8888/move', { "gopher": parseInt(gopher, 10), "direction": parseInt(direction, 10) }).success(function (response) {
      deferred.resolve(response);
    }).error(function (response) {
      deferred.reject(response);
    });

    return deferred.promise;
  }

})();
