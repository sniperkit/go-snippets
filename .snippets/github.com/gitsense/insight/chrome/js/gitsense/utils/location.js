sdes.gitsense.utils.location = function() {
    "use strict";

    var _this   = this,
        varUtil = new sdes.utils.variable();

    this.parse = function(_hash) {
        var search = _this.parseSearch(),
            hash   = _this.parseHash(_hash);
        
        return {
            search: search,
            hash: hash
        };
    }

    this.parseSearch = function() {
        var search  = decodeURI(window.location.search),
            params  = search.match(/^\s*$/) ? [] : search.replace(/^\?/, "").split("&"),
            map     = {};

        for ( var i = 0; i < params.length; i++ ) {
            var temp = params[i].split("=");

            if ( temp[0] === "r" )
                temp[0] = "repos";

            map[temp[0]] = temp[1];
        }

        return map;
    }

    this.parseHash = function(hash, returnString) {
        var winhash = 
                varUtil.isNoU(hash) ? 
                    decodeURIComponent(window.location.hash.replace(/^#/,"")) :
                    decodeURIComponent(hash.replace(/^#/,""));

        var params        = winhash.match(/^\s*$/) ? [] : winhash.replace(/^#/, "").split("&"),
            host          = "",
            admin         = false,
            hostUser      = false,
            quickAnalysis = false,
            diffs         = null,
            branches      = null,
            repos         = null,
            query         = null,
            cdo           = null,
            cmo           = null,
            tab           = null,
            range         = null,
            pill          = null,
            repoOwner     = null,
            repoName      = null,
            repoFilter    = null,
            viewMode      = null,
            _query        = null,
            _tab          = null;

        for ( var i = 0; i < params.length; i++ ) {
            var temp = params[i].split("=");

            if ( temp[0] === "" )
                continue;

            switch(temp[0]) {
                case "a":
                    continue;
                case "admin":
                    if ( temp[1] === "true" )
                        admin = true;

                    continue;
                case "host-user":
                    if ( temp[1] === "true" )
                        hostUser = true;

                    continue;
                case "b":
                    branches = temp[1].split("::").sort();

                    //for ( var j = 0; j < branches.length; j++ )
                    //    branches[j] = decodeURI(branches[j]);

                    continue;    
                case "h":
                    host = temp[1];
                    continue;    
                case "r":
                    repos = temp[1].split("::").sort();
                    continue;    
                case "p":
                    pill = temp[1];
                    continue;
                case "ro":
                    repoOwner = temp[1];
                    continue;
                case "rn":
                    repoName = temp[1];
                    continue;
                case "rf":
                    repoFilter = temp[1];
                    continue;
                case "q":
                    query = _this.parseQuery(temp[1]);
                    continue;
                case "_q":
                    _query = _this.parseQuery(temp[1]);
                    continue;
                case "cdo":
                    cdo = temp[1] === "true" ? true : false;
                    continue;
                case "cmo":
                    cmo = temp[1] === "true" ? true : false;
                    continue;
                case "rg":
                    range = temp[1];
                    continue;
                case "t":
                    tab = temp[1];
                    continue;
                case "_t":
                    _tab = temp[1];
                    continue;
                case "qa":
                    quickAnalysis = temp[1] === "true" ? true : false;
                    continue;
                case "vm":
                    viewMode = temp[1];
                    continue;
                case "dc":
                    if ( diffs === null )
                        diffs = {};

                    diffs.diff = temp[1] === "true" ? true : false;
                    continue;
                case "drb":
                    if ( diffs === null )
                        diffs = {};

                    diffs.relatedBranches = temp[1] === "" ? [] : temp[1].split(",");
                    continue;
                case "dcl":
                    diffs.callLines = temp[1] !== "" ? temp[1].split(",") : [];
                    continue; 
                case "df":
                    if ( temp[1] === "" )
                        continue;

                    if ( diffs === null )
                        diffs = {};

                    diffs.filters = JSON.parse(temp[1]);
                    continue;
                case "dl":
                    if ( diffs === null )
                        diffs = {};

                    diffs.lang = temp[1];

                    continue;
                case "dp":
                    if ( diffs === null )
                        diffs = {};

                    diffs.path= temp[1];

                    continue;
                case "drg":
                    if ( diffs === null )
                        diffs = {};

                    diffs.range = temp[1];

                    continue;
                case "dvm":
                    if ( diffs === null )
                        diffs = {};

                    diffs.viewMode = temp[1];

                    continue;
                default:
                    console.warn("Ignoring unrecognized hash paramater '"+temp[0]+"'");
            }
        }

        if ( query === null )
            query = _this.parseQuery();

        var results = {
            host: host,
            branches: branches,
            repos: repos,
            query: query,
            diffs: diffs,
            admin: admin,
            hostUser: hostUser,
            cdo: cdo,
            cmo: cmo,
            tab: tab,
            range: range,
            quickAnalysis: quickAnalysis,
            viewMode: viewMode,
            pill: pill,
            repoFilter: repoFilter,
            repoOwner: repoOwner,
            repoName: repoName,
            _query: _query,
            _tab: _tab
        };

        if ( returnString )
            return _this.getHashString(results);

        return results;
    }

    this.parseQuery = function(q) {
        var oargs = 
                varUtil.isNoU(q) ? 
                    [] : 
                    typeof(q) === "string" ?
                        q.split("^_^") :
                    q,
             args = [];

        var filters = {
            actions: [],
            authors: [],
            chgs: "",
            churn: "",
            schurn: "",
            commits: [],
            diff: [],
            mergedby: [],
            langs: [],
            paths: [],
            days: [],
            weeks: [],
            months: [],
            follow: false,
            bh: false,
            unique: false,
            onlatest: false,
            nomerges: false,
            noempty: false,
            casesensitive: false,
            mergeconflicts: false,
            starred: "none",
            lasthours: -1,
            lastdays: -1,
            tmzoffset: 0,
            likecm: "",
            likepom: 0,
            likemr: 0,
            likepr: 0,
            head: "",
            window: ""
        };

        for ( var i = 0; i < oargs.length; i++ ) {
            var arg   = oargs[i],
                temp  = arg.split(":"),
                type  = temp.length === 2 ? temp.shift() : "",
                value = temp.join(":");

            if ( arg === "" ) 
                continue;

            if ( 
                type === "diff" && 
                ! value.match(/\.\.\./) &&
                ! value.match(/\.\./) &&
                ! value.match(/->/)
            ) {
                console.warn("Invalid diff query \""+arg+"\" ... ignoring");
                continue;
            }

            if (
                type !== "window" &&
                type !== "chgs" &&
                type !== "churn" &&
                type !== "schurn" &&
                type !== "head" &&
                type !== "diff" && 
                type !== "mergedby" && 
                type !== "tmzoffset" &&
                type !== "starred" &&
                type !== "likecm" &&
                type !== "likepom" &&
                type !== "likepr" &&
                type !== "likemr" &&
                value !== "false" && 
                value !== "true" &&
                ! type.match(/s$/)
            ) {
                type += "s";
            }

            if ( value !== "" && ! varUtil.isNoU(filters[type]) ) {
                var typeValue = typeof(filters[type]);

                if ( 
                    typeValue === "string" || 
                    typeValue === "number" ||
                    filters[type].length === undefined 
                ) {
                    filters[type] = 
                        type === "tmzoffset" || 
                        type === "likepr" || 
                        type === "likemr" || 
                        type === "likepom" ?
                            parseInt(value) : 
                            value;
                } else {
                    filters[type].push(value);
                }
            } else {
                args.push(arg);
            }
        }

        return {
            q: q,
            oargs: oargs,
            args: args,
            filters: filters
        };
    }

    this.clearHash = function()  {
        window.location.hash = "";
    }

    this.getHashString = function(hash, q)  {
        if ( ! varUtil.isNoU(q) )
            throw("ERROR: Creating a hash based on a query string is no longer supported");

        var values = [];

        if ( ! varUtil.isNoU(hash.host) && hash.host !== "" )
            values.push("h="+hash.host);

        if ( ! varUtil.isNoU(hash.branches) ) {
            var b = hash.branches.join("::");

            values.push("b="+b);
        }

        if ( ! varUtil.isNoU(hash.repos) ) {
            var r = hash.repos.join("::");

            values.push("r="+r);
        }

        if ( 
            ! varUtil.isNoU(hash.query) &&
            ! varUtil.isNoU(hash.query.oargs) &&
            hash.query.oargs.length !== 0  
        ) {
            values.push("q="+hash.query.oargs.join("^_^"));
        }

        if ( ! varUtil.isNoU(hash.cdo) )
            values.push("cdo="+hash.cdo);

        if ( ! varUtil.isNoU(hash.cmo) )
            values.push("cmo="+hash.cmo);

        if ( ! varUtil.isNoU(hash.tab) )
            values.push("t="+hash.tab);

        if ( ! varUtil.isNoU(hash.range) )
            values.push("rg="+hash.range);

        if ( ! varUtil.isNoU(hash.viewMode) )
            values.push("vm="+hash.viewMode);

        if ( ! varUtil.isNoU(hash.quickAnalysis) && hash.quickAnalysis ) 
            values.push("qa=true");

        if ( ! varUtil.isNoU(hash.pill) )
            values.push("p="+hash.pill);
            
        if ( ! varUtil.isNoU(hash.repoOwner) )
            values.push("ro="+hash.repoOwner);
            
        if ( ! varUtil.isNoU(hash.repoName) )
            values.push("rn="+hash.repoName);
            
        if ( ! varUtil.isNoU(hash.repoFilter) )
            values.push("rf="+hash.repoFilter);
            
        if ( ! varUtil.isNoU(hash.diffs) ) {
            values.push("dc="+(hash.diffs.diff ? "true" : "false"));
            values.push("drg="+(varUtil.isNoU(hash.diffs.range) ? "" : hash.diffs.range));
            values.push("dl="+(varUtil.isNoU(hash.diffs.lang) ? "" : hash.diffs.lang));
            values.push("dp="+(varUtil.isNoU(hash.diffs.path) ? "" : hash.diffs.path));
            values.push("dvm="+(varUtil.isNoU(hash.diffs.viewMode) ? "" : hash.diffs.viewMode));

            values.push(
                "drb="+(
                    varUtil.isNoU(hash.diffs.relatedBranches) ? 
                        "" : 
                        hash.diffs.relatedBranches.join(",")
                )
            );

            values.push(
                "dcl="+
                (
                    varUtil.isNoU(hash.diffs.callLines) ? 
                        "" : 
                        hash.diffs.callLines.join(",")
                )
            );

            values.push(
                "df="+
                (
                    varUtil.isNoU(hash.diffs.filters) ? 
                        "" : 
                        JSON.stringify(hash.diffs.filters)
                )
            );
        }

        if ( hash.admin )
            values.push("admin=true");

        if ( hash.hostUser )
            values.push("host-user=true");

        return encodeURIComponent(values.join("&"));
    }

    this.setHash = function(hash, q) {
        window.location.hash = _this.getHashString(hash, q);
    }
}   
