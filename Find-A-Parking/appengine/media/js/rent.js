$(document).on("ready",init);
var myMap = {};

var callback =function  () {
	
	google.maps.event.addListener( myMap.getMap(),'click' , addMarker);
	console.log("Added")
}


function init () {
	$("#rent").on("click",sendParkingInfo);
	myMap = new DMaps("map",25.670708,-100.308172,callback);
}

function addMarker (e) {
	myMap.addMarker(e.latLng,"parking");
	$("#lat").val(e.latLng.lat());
	$("#lng").val(e.latLng.lng());

}

function sendParkingInfo () {
	$.ajax({
	  type: "POST",
	  url: "createPark",
	  dataType: "json",
	  data: { name: $("#name").val(), email: $("#email").val(), price: $("#price").val(), lat: $("#lat").val(), lng: $("#lng").val() },
	success: function  (result) {
		console.log(result);
	},
	}).done(function( msg ) {
	    alert("Parking saved");
	  });
}

