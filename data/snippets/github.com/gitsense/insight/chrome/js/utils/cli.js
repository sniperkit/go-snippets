sdes.utils.cli = function(options) {
    var optToValue = undefined;

    if ( options !== undefined ) 
    {
        optToValue = {};

        for ( var i = 0; i < options.length; i++ ) 
        {
            var opt     = options[i].opt;
            var longOpt = options[i].longOpt;
            var numArgs = options[i].numberOfArgs;
            var value   = numArgs == -1 ? 1 : options[i]["values"][0];

            if ( opt != "" )
                optToValue[opt] = value;

            if ( longOpt != "" )
                optToValue[longOpt] = value;
        }
    }

    this.getArgs = function(input) 
    {
        input = input.replace(/^\s+/,"").replace(/\s+$/,"");
        
        var chars = input.split("");
        var args  = [];
        var arg   = "";
        var stack = [];
        var quote = null;
        var c     = null;
    
        for ( var i = 0; i < chars.length; i++ ) {
            c = chars[i];
    
            if ( quote === null && c === " " ) {
                if ( stack.length !== 0 ) {
                    arg = stack.join("");
                    args.push(arg);
                    stack = [];
                }
    
                continue;
            }
    
            if ( c == "'" || c == '"' ) {
                if ( quote == c ) {
                    arg = stack.join("");
                    stack = [];
                    args.push(arg);
                    quote = null;
                    continue;
                } else if ( quote === null ) {
                    quote = c;
                    continue;
                }
            }
     
            stack.push(c);
        } 
      
        if ( quote !== null )
            return "";
    
        if ( stack.length != 0 ) {
            arg = stack.join("");
            args.push(arg);
        }
    
        return args;
    }
    
    this.getOptionValue = function(opt) {
        if ( optToValue === undefined )
            throw("optToValue hash is undefined");
            
        return optToValue[opt];
    }
}
