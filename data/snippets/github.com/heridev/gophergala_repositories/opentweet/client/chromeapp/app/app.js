angular
	.module('openTweet', ['ionic', 'header', 'tweetPage', 'followPage', 'settingsPage']);

angular
	.module('backend', [])
	.factory('Settings', ['$q', '$http', function($q, $http) {
		function get() {
			var dfd = $q.defer();
			chrome.storage.local.get({
				settings: {
					server: '',
					username: '',
					password: ''
				}
			}, function(items) {
				dfd.resolve(items.settings);
			});
			return dfd.promise;
		}
		return {
			save: function(settings) {
				var dfd = $q.defer();
				chrome.storage.local.set({
					settings: settings
				}, dfd.resolve);
				return dfd.promise;
			},
			fetch: get,
		    sendTweet: function(text) {
			console.log(text)
				get().then(function(settings) {
					// TODO - Server should allow CORS
				    // return $http.post(settings.server + '/tweet', {
				    // 	headers: {
				    // 	    'Authorization': 'Basic ' + window.btoa(settings.username + ':' + settings.password)
				    // 	},
				    // 	data: {
				    // 	    "tweet": text
				    // 	}
				    // })
				    return $http({
					method: 'POST',
					url: settings.server + '/tweet',
						headers: {
							'Authorization': 'Basic ' + window.btoa(settings.username + ':' + settings.password)
						},
						data: {
							"tweet": text
						}
					})
				});
			},
			register: function(server, username, password) {
				// TODO - Server should allow CORS
				return $http.post(server + '/users', {
					"user": username,
					"password": password
				});
			}
		}
	}])
	.factory('User', ['$q', function($q) {
		function getWhoIFollow() {
			var dfd = $q.defer();
			chrome.storage.local.get({
				whoIFollow: []
			}, function(items) {
				dfd.resolve(items.whoIFollow);
			});
			return dfd.promise;
		}

		function setWhoIFollow(people) {
			var dfd = $q.defer();
			chrome.storage.local.set({
				whoIFollow: people
			}, function() {
				dfd.resolve();
			});
			return dfd.promise;
		}

		return {
			whoIFollow: getWhoIFollow,
			followPerson: function(user) {
				return getWhoIFollow().then(function(people) {
					people.push(user);
					return setWhoIFollow(people);
				});
			},
			unFollowPerson: function(user) {
				return getWhoIFollow().then(function(people) {
					if (people.indexOf(user) !== -1) {
						people.splice(people.indexOf(user), 1);
					}
					return setWhoIFollow(people);
				});
			},
			parse: function(data) {
				var username = data.split(/@/);
				var server = username[1].split(':');
				return {
					username: username[0],
					server: server[0],
					port: server[1] || 12315
				}
			}
		}
	}])
	.factory('Tweets', ['$q', 'User', function($q, User) {
		return {
			get: function(personStr, from, to, callbacks) {
				var user = User.parse(personStr);
				var defer = $q.defer();
				var message = ['OT v1', user.username, from, to].join('\r\n') + '\r\n';
				// TODO Fetch from actual servers
				var tcpClient = new TcpClient(user.server, user.port);
				tcpClient.connect(function() {
					var data = [];
					tcpClient.addResponseListener(function(response) {
						data.push(response);
					});
					tcpClient.addResponseErrorListener(function() {
						var len = data.map(function(d) {
							return d.byteLength;
						}).reduce(function(prev, curr) {
							return prev + curr;
						}, 0);

						var tmp = new Uint8Array(len),
							offset = 0;
						data.forEach(function(d) {
							tmp.set(new Uint8Array(d), offset);
							offset += d.byteLength;
						});
						tcpClient._arrayBufferToString(tmp.buffer, function(str) {
							var data = str.split('\r\n');
							var result = [];
							for (var i = 0; i < data.length - 1; i += 2) {
								result.push({
									user: personStr,
									timestamp: new Date(parseInt(data[i], 10) * 1000),
									tweet: data[i + 1]
								});
							}
							defer.resolve(result);
						});
					});
					tcpClient.sendMessage(message);
				});
				return defer.promise;
			}
		}
	}]);

angular
	.module('tweetPage', ['backend'])
	.controller('TweetListCtrl', ['Tweets', 'User', '$scope', function(tweets, user, $scope) {
		$scope.tweetsFromWhoIFollow = [];
		$scope.lastFetchTime = Math.round(moment().subtract(1, 'day').toDate().getTime()/1000.0);
		function update() {
			user.whoIFollow().then(function(people) {
				var userCount = people.length;
                        	var now = Math.round(new Date().getTime()/1000.0);
				angular.forEach(people, function(person) {
					return tweets.get(person, $scope.lastFetchTime, now).then(function(tweets) {
						Array.prototype.push.apply($scope.tweetsFromWhoIFollow, tweets);
					}).finally(function() {
						if (--userCount === 0) {
							$scope.lastFetchTime = now;
							$scope.$broadcast('scroll.refreshComplete');
						}
					})
				});
			});
		}
		$scope.doRefresh = update;
		update();
	}])
	.filter('fromNow', function() {
		return function(input) {
			return moment(input).fromNow();
		}
	}).filter('userName', ['User', function(user) {
		return function(input) {
			return user.parse(input).username;
		}
	}]);

angular
	.module('followPage', ['backend', 'ionic'])
	.controller('FollowCtrl', ['User', '$ionicPopup', '$scope',
		function(User, $ionicPopup, $scope) {
			$scope.follow = function(user) {
				User.followPerson(user).then(function() {
					$scope.whoIFollow.push(user);
				});
			};
			$scope.unfollow = function(user) {
				User.unFollowPerson(user).then(refreshWhoIFollow);
			};

			function refreshWhoIFollow() {
				User.whoIFollow().then(function(people) {
					$scope.whoIFollow = people;
				});
			}
			refreshWhoIFollow();
		}
	]);

angular
	.module('settingsPage', ['backend', 'ionic'])
	.controller('SettingsCtrl', ['Settings', '$ionicPopup', '$scope', function(Settings, $ionicPopup, $scope) {
		Settings.fetch().then(function(settings) {
			$scope.server = settings.server;
			$scope.username = settings.username;
			$scope.password = settings.password;
		});
		$scope.saveSettings = function(server, username, password) {
			Settings.save({
				username: username,
				password: password,
				server: server
			}).then(function() {
				$ionicPopup.alert({
					title: 'Settings saved successfully',
				});
			});
		}

		$scope.registerUser = function(server, username, password) {
			Settings.register(server, username, password).then(function() {
				$ionicPopup.alert({
					title: 'Settings saved successfully',
				});
			}, function() {
				$ionicPopup.alert({
					title: 'Could not register new user on server',
				});
			})
		};
	}]);

angular
	.module('header', ['ionic', 'backend'])
	.controller('HeaderCtrl', ['$scope', '$ionicModal', 'Settings', function($scope, $ionicModal, Settings) {
		$ionicModal.fromTemplateUrl('compose.html', function(modal) {
			$scope.compose = modal;
		}, {
			scope: $scope,
			animation: 'slide-in-down'
		});

		$scope.sendTweet = function(text) {
			$scope.compose.hide();
		    console.log(text)
		    Settings.sendTweet(text);
		};
	}]);
