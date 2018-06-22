angular.module "feedBag"
  .controller "MainCtrl", ($scope, socket) ->
    $scope.Activities = []

    socket.syncUpdates "activity", $scope.Activities, -> console.log arguments
