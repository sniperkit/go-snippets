sdes.github.data.repo = function(rule) {
    "use strict";

    var apiUrl = rule.host.api;

    this.get = function(owner, repo, callback) {
        $.ajax({
            beforeSend: getBeforeSend(),
            url: apiUrl+"/repos/"+owner+"/"+repo,
            success: function (data) {
                callback(data);
            },
            error: function(e) {
                callback(null, e);
            }
        });
    }

    this.getLatestBranchCommit = function(branch, callback) {
        $.ajax({
            url: apiUrl+"/repos/"+owner+"/"+name+"/commits/"+branch,
            success: function (commit) {
                callback(commit);
            },
            error: function(e) {
                callback(null, e);
            }
        });
    }

    this.getCommit = function(sha, callback) {
        $.ajax({
            beforeSend: getBeforeSend(),
            url: apiUrl+"/repos/"+owner+"/"+name+"/commits/"+sha,
            success: function (commit) {
                callback(commit);
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
                "Basic "+btoa(rule.host.username+":"+rule.host.secret)
            );
        };
    }
}
