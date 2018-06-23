angular.module "feedBag", ['ngAnimate', 'ngRoute', 'ui.bootstrap', 'btford.socket-io']
  .config ($routeProvider) ->
    $routeProvider
      .when "/",
        templateUrl: "app/main/main.html"
        controller: "MainCtrl"
      .otherwise
        redirectTo: "/"

