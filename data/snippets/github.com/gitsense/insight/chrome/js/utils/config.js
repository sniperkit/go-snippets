sdes.utils.config = function() {
    var _this = this;

    this.getRule = function() {
        for ( var i = 0; i < sdes.config.routingRules.length; i++ ) {
            var rule = sdes.config.routingRules[i],
                patt = new RegExp(rule.matches.replace("*", ".*"));

            if ( window.location.href.match(patt) )
                return rule;
        }

        return null;
    }
}
