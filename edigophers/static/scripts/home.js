

$(function() {
  function initialize() {

    var mapOptions = {
      center: user.location,
      zoom: 10
    };

    var map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
    var marker = new google.maps.Marker({
      position: user.location,
      map: map,
      title:user.name,
      icon: 'images/markers/darkgreen_MarkerA.png'
    });

    // Add the circle for this city to the map.
    var radius = new google.maps.Circle({
      strokeColor: '#00FF00',
      strokeOpacity: 0.2,
      strokeWeight: 2,
      fillColor: '#00AA00',
      fillOpacity: 0.35,
      map: map,
      center: user.location,
      radius: 20000
    });

    var removeAnimation = $.noop
    $.each(users, function(index, user){
      var marker = new google.maps.Marker({
        position: user.location,
        map: map,
        title: user.name,
        animation: google.maps.Animation.DROP
      });

      $("#usr-zoom-" + user.name).click(function(e){
        removeAnimation()
        $("h2")[0].scrollIntoView();
        map.setCenter(user.location);
        marker.setAnimation(google.maps.Animation.BOUNCE);
        removeAnimation = function(){ marker.setAnimation(null); }
        e.preventDefault();
      })
    });
  };

  initialize();
})
