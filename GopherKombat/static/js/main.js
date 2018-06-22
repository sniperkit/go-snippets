/*global $*/
/*global document*/

var GK = {};

GK.init = function() {
    $(document.body).load("template/frame.html");
};

$(document).ready(function() {
    GK.init();    
});