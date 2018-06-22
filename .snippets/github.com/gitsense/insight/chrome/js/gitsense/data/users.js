sdes.gitsense.data.users = function() {
    "use strict";

    this.getAvatars = function(emails, callback) {
        chrome.storage.local.get(emails, function(avatars){
            callback(avatars);
        });
    }

    this.storeAvatars = function(avatars, callback) {
        chrome.storage.local.set(avatars, function() {
            if ( callback !== undefined )
                callback(true);
        });
    }
}
