sdes.gitsense.data.repo = function(rule) {
    "use strict";

    var varUtil = new sdes.utils.variable(),
        cache   = sdes.cache.gitsense.data.repo,
        baseUrl = rule.gitsense.baseUrl+"/insight";

    if ( varUtil.isNoU(cache) )
        cache = {};

    this.queryHost = function(host, id, callback) {
        $.ajax({
            url: baseUrl+"/query/"+host+"/repos/"+id,
            success: function(repo) {
                callback(repo);
            },
            error: function(e) {
                callback(null, e);
            }
        });
    }
}
