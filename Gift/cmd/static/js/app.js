var nwt = {coords:{latitude:71.3796235, longitude:-109.8513831}};
/*
var newyork = {coords:{latitude:40.7127, longitude:-74.0059}};
var whitehouse = {coords:{latitude:38.8977, longitude:-77.0366}};
var microsoft = {coords:{latitude:47.643516,longitude:-122.1289033}};
var apple = {coords:{latitude:37.331411, longitude:-122.030315}};
*/

function getLocation(div) {
  $("#output").html("<b>Please wait while we load your location info: <i class=\"fa fa-refresh fa-spin\"></i></b>");
  if (navigator.geolocation) {
    var options = {enableHighAccuracy: true};
    navigator.geolocation.getCurrentPosition(setPosition, errorHandler, options);
  } else {
    out.innerHTML = "Geolocation is not supported by this browser.";
  }
}
function setPosition(position) {
  var out = document.getElementById("watch_results");
  $("#watch_results").html("Latitude: " + position.coords.latitude +
    "<br>Longitude: " + position.coords.longitude +
    "<br>Heading: " + position.coords.heading);

  $("#output").html('');

  var r = new XMLHttpRequest();
  r.open("POST", "/setgeo", true);
  r.onreadystatechange = function () {
    if (r.readyState != 4 || r.status != 200) return;
    window.location = "#watch";
  };

  r.send(
      ["latitude=" + position.coords.latitude,
      "longitude=" + position.coords.longitude,
      "heading=" + position.coords.heading].join("&")
      );
}
function errorHandler(err) {
  var out = document.getElementById("error");
  out.innerHTML = "Error loading your position information";
}

$(".gift").click(function(event) {
  $(this).attr('href', '#gift');
  $("#gift").html("<img src=\"/" + $(this).data("gift") + ".gif\"/>");
});


$("#custom-loc").submit(function() {
   var loc = {coords:{latitude:$("#latitude").val(), longitude:$("#longitude").val()}};

   setPosition(loc);
   return false;
});



