sdes.gitlab.data.repo = function(rule) {
    "use strict";

    var apiUrl      = rule.host.api,
        accessToken = rule.host.secret === "" ? "" : "?private_token="+rule.host.secret;

    this.get = function(projectId, callback) {
        $.ajax({
            url: apiUrl+"/projects/"+projectId+accessToken,
            success: function (data) {
                callback(data);
            },
            error: function(e) {
                callback(null, e);
            }
        });
    }
}
