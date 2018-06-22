var gophertv = angular.module('gophertv', ['ngRoute', 'tvControllers']);

gophertv.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/', {
        templateUrl: '/public/templates/list.html',
        controller: 'HomeCtrl'
      }).
			when('/video/:videoId', {
				templateUrl: 'public/templates/video.html',
				controller: 'VideoCtrl'
			}).
      when('/category/:categoryId', {
        templateUrl: 'public/templates/category-detail.html',
        controller: 'CategoryDetailCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
  }]);
