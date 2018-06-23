var Groupify = angular.module('Groupify', []);

Groupify.controller('MainCtrl', function($scope, $http, $timeout) {
  var spotifyApi = new SpotifyWebApi();
  $scope.queue = [];
  $scope.query = "";
  $scope.trackResults = [];
  backup_avatar = "http://www.140proof.com/images/logos/140-proof-logo-shadow-500px.png";

  (function tick() {
    $http.get('/api/v1/queue/list')
    .then(function(res){

      // sum time to play up to each track in collection
      var sum = 0;

      track = res.data.now_playing.track;
      if (track && res.data.now_playing.time_remaining) {
        if( ! track.queued_by_avatar )
          track.queued_by_avatar = backup_avatar;
        console.log( "currently playing queued by avatar: " + track.queued_by_avatar );
        sum = track.time_remaining = parseInt(res.data.now_playing.time_remaining);
        $scope.current_track = track;
      }

      $scope.queue = res.data.queue;
      for(var i = 0; i < $scope.queue.length; i++) {
        track = $scope.queue[i];
        track.time_to_play = sum;
        if( ! track.queued_by_avatar )
          track.queued_by_avatar = backup_avatar;
        sum += parseInt(track.time);
        console.log( track );
      }

      $timeout(tick, 1000);
    });
  })();

  $scope.next = function(){
    $http.post('/api/v1/queue/next')
    .then(function(res){
      console.log("skipping to next track");
    });
  };

  $scope.search = function(){
    spotifyApi.searchTracks($scope.query, {limit: 10, offset: 0}, function(err, data) {
      $scope.trackResults = data.tracks.items;
    });
  };

  $scope.enqueue = function(track){
    $http.post('/api/v1/queue/add', {
      track_id: track.id
    })
    .then(function(res){
      console.log("Enqueued track " + track.name);
    });
  };

  $scope.dequeue = function(track){
    $http.post('/api/v1/queue/delete', {
      track_id: track.id
    })
    .then(function(res){
      // FIXME: handle errors
      console.log("De-queued track " + track.name);
    });
  };

});

Groupify.filter('secondsToTime', function() {
  // shameless copy/paste
  // http://codeaid.net/javascript/convert-seconds-to-hours-minutes-and-seconds-(javascript)
  return function(secs) {
    var hr  = Math.floor(secs / 3600);
    var min = Math.floor((secs - (hr * 3600))/60);
    var sec = secs - (hr * 3600) - (min * 60);

    if (hr  < 10) { hr  = "0" + hr; }
    if (min < 10) { min = "0" + min;}
    if (sec < 10) { sec = "0" + sec;}
    if (hr)       { hr  = "00"; }
    return hr + ':' + min + ':' + sec;
  };
});

Groupify.filter('secondsToMinutes', function() {
  return function(secs) {
    var min = Math.floor(secs / 60);
    var sec = secs - (min * 60);

    if (min < 10) { min = "0" + min;}
    if (sec < 10) { sec = "0" + sec;}
    return min + ':' + sec;
  };
});
