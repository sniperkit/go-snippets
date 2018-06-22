var tvcontrollers = angular.module('tvControllers', ['youtube-embed'] );

tvcontrollers.controller('HomeCtrl', [ '$scope', '$http', '$location',
								function($scope, $http, $location) {
									$scope.categories = [];
									//$scope.categories = [ {"CategoryName": "One", "videos": [{"Title": "Video Title"}]}];	

									$http.get('playlists').success(function(data) {
                    $scope.categories = data;
									});

									$scope.navigateToHome = function(id) {
										$location.path('/');
									};
									$scope.navigateToVideo = function(id) {
										$location.path('/video/' + id);
									};
									$scope.navToCategory = function(id) {
										$location.path('/category/' + id);
									};
								}]);

tvcontrollers.controller('CategoryDetailCtrl', [ '$scope', '$http', '$location',
								function($scope, $http, $location) {
									var tag_id = $location.path().split('/');
									tag_id = tag_id[tag_id.length - 1];
									$scope.categoryName = tag_id;

									$http.get('videos?tag=' + tag_id).success(function(data){
										$scope.videos = data;
									});

									$scope.navigateToVideo = function(id) {
										$location.path('/video/' + id);
									};
								}]);

tvcontrollers.controller('VideoCtrl', [ '$scope', '$http', '$location',
								function($scope, $http, $location) {
								  var vid_id = $location.path().split('/');
									vid_id = vid_id[vid_id.length - 1];

									$scope.youtubevideo = vid_id;

									$http.get('videos/' + vid_id).success(function(data) {
										$scope.video = data;
										
									//	$scope.setupYoutubePlayer();
									});
									
									$scope.setupYoutubePlayer = function() {
										// Load the IFrame Player API code asynchronously.
										var tag = document.createElement('script');
										tag.src = "https://www.youtube.com/player_api";
										var firstScriptTag = document.getElementsByTagName('script')[0];
										firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
									};
								}]);

// Replace the 'ytplayer' element with an <iframe> and
										// YouTube player after the API code downloads.
										var player;
										function onYouTubePlayerAPIReady() {
											var id = document.getElementById('video_id').value;
											player = new YT.Player('ytplayer', {
												height: '390',
												width: '640',
												videoId: id
											});
										}	
