sdes.utils.variable = function() {
    this.isNoU = function(v) {
        if ( v === undefined || v === null )
            return true;

        return false; 
    }
}
