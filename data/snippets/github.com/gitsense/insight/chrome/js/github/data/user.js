sdes.github.data.user = function(rule) {
    "use strict";

    var apiUrl = rule.host.api;

    this.get = function(user, callback) {
        $.ajax({
            beforeSend: getBeforeSend(),
            url: apiUrl+"/users/"+user,
            success: function (data) {
                callback(data);
            },
            error: function(e) {
                callback(null, e);
            }
        });
    }

    function getBeforeSend() {
        if ( rule.host.secret === undefined || rule.host.secret === "" ) 
            return null;
    
        return function (xhr){ 
            xhr.setRequestHeader(
                "Authorization", 
                "Basic "+btoa("username:"+rule.host.secret)
            );
        };
    }
}
