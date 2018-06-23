var application = angular
.module('abbita', ['angularFileUpload', 'ui.router'])
.controller('NewSessionController', ['$scope', 'FileUploader','$state', function($scope, FileUploader, $state) {
	var uploader = $scope.uploader = new FileUploader({url: '/api/sessions/new', alias: 'fileUpload'});
	uploader.onSuccessItem = function(fileItem, response, status, headers) {
		console.info('onSuccessItem', response);
		$state.go('session',{id:response.id})
	};
	uploader.onErrorItem = function(fileItem, response, status, headers) {
		console.info('onErrorItem', fileItem, response, status, headers);
	};
	
	$scope.upload = function(item) {
		item.formData.push({name: item.file.name})
		item.upload();
	}
}]).controller('SessionController', ['$scope','$state','$http','$stateParams', function($scope, $state, $http, $stateParams) {
	var promise = $http.get('/api/sessions/'+$stateParams.id);
	promise.then(function(response) {
		$scope.session = response.data;
	})
	
}]);

application.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider){
	$urlRouterProvider.otherwise("/");
	$stateProvider
	.state('home', {
		url: "/",
		templateUrl: 'templates/home.html',
		controller: 'NewSessionController'
	})
	.state('session', {
		url: "/session/:id",
		templateUrl: 'templates/session.html',
		controller: 'SessionController'
	})
}])	

application.filter('bytes', function() {
	return function(bytes, precision) {
		if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
		if (typeof precision === 'undefined') precision = 1;
		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
		number = Math.floor(Math.log(bytes) / Math.log(1024));
		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
	}
});