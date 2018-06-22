var lastGitSenseShow = null;

sdes.gitlab.utils.page = function(rule) {
    var varUtil = new sdes.utils.variable();

    this.parse = function(callback, force) {
        if ( force === undefined )
            force = false;

        var names = window.location.pathname.replace(/^\/|\/$/, "").split(/\//);

        if ( names.length === 1 ) {
            if ( names[0] === "" )
                return;

            if ( names[0] === "search" ) {
                parseSearch(callback);
                return;
            }
        }

        var ignoreNames = [ "profile", "dashboard", "admin", "users" ],
            firstName   = names[0];

        for ( var i = 0; i < ignoreNames.length; i++ ) {
            if ( firstName !== ignoreNames[i] ) 
                continue;

            callback(null);
            return;
        } 

        if ( names.length !== 1 && names[0] !== "groups" && names[0] !== "u" ) {
            if ( names.length > 1 )
                parseRepo(names, callback, force);
            else 
                callback(null);

            return;
        }

        parseUserOrGroup(names.length === 1 ? names[0] : names[1], callback);
    }

    function parseUserOrGroup(name, callback, force) {
        if ( force === undefined )
            force = false;

        var show   = window.location.search.match(/gitsense=insight/) ? true : false,
            stopAt = new Date().getTime()+2000;

        if ( ! force && lastGitSenseShow && show )
            return;

        lastGitSenseShow = show;

        parse();

        function parse() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Given up on looking for elements");
                return;
            }

            var navLinks     = document.getElementsByClassName("nav-links"),
                contentElems = document.getElementsByClassName("content");

            if ( 
                navLinks === null || 
                navLinks.length === 0 ||
                contentElems === null ||
                contentElems.length === 0
            ) {
                setTimeout(parse, 50);
                return;
            }

            var userPage  = false,
                groupPage = false;

            for ( var i = 0; i < navLinks[0].children.length; i++ ) {
                var li   = navLinks[0].children[i],
                    text = $(li).text();

                if ( text.toLowerCase().match(/personal projects/) )
                    userPage = true;
                else if ( text.toLowerCase().match(/contribution analytics/) )
                    groupPage = true;
                else if ( text.toLowerCase().match(/members/) )
                    groupPage = true;
            }

            callback({
                type: userPage ? "user" : groupPage ? "group" : "unknown",
                user: userPage ? name : null,
                group: groupPage ? name : null,
                content: contentElems[0],
                navLinks: navLinks,
                show: show
            });
        }
    }

    function parseRepo(names, callback, force) {
        var owner     = names.shift(),
            repo      = names.shift(),
            home      = names.length === 0 ? true : false,
            search    = names.length !== 0 && names[0] === "search" ? true : false,
            blame     = names.length !== 0 && names[0] === "blame" ? true : false,
            mr        = names.length !== 0 && names[0] === "merge_requests" ? true : false,
            mrNum     = mr && names.length > 1 ? names[1] : 0,
            mrPage    = mrNum !== 0 && names.length > 2 ? names[2] : null,
            show      = window.location.search.match(/gitsense=insight/) ? true : false,
            stopAt    = new Date().getTime() + 1000,
            trackerId = "gitsense-content-tracker",
            iframeId  = "gitsense-content-iframe",
            container = document.getElementById(trackerId),
            iframe    = document.getElementById(gitsenseFrameId);

        if ( ! force && lastGitSenseShow && show )
            return;

        if ( show && iframe !== null )
            findElems();
        else if ( container !== null )
            waitForNewContainers();
        else
            findElems();

        lastGitSenseShow = show; 

        function waitForNewContainers() {
            if ( new Date().getTime() > stopAt )  {
                stopAt = new Date().getTime() + 2000;
                findElems();
                return;
            }

            var container = document.getElementById(trackerId);

            if ( container !== null ) {
                setTimeout(waitForNewContainers, 50);
                return;
            }

            stopAt = new Date().getTime() + 2000;
            findElems();
        }

        function findElems() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Given up on looking for containers");
                return;
            }

            var contentElems = document.getElementsByClassName("content"),
                commitElems  = blame ? document.getElementsByClassName("pull-right") : [],
                navLinks     = document.getElementsByClassName("nav-links"),
                repoRef      = document.getElementById("repository_ref"),
                projectId    = document.getElementById("project_id");

            if ( 
                contentElems === null || 
                contentElems.length === 0 ||
                navLinks === null ||
                navLinks.length === 0 ||
                commitElems === null
            ) {
                setTimeout(findElems, 50);
                return;
            }

            callback({
                type: "repo",
                owner: owner,
                repo: repo,
                show: show,
                navLinks: navLinks,
                content: contentElems[0],
                mergeRequest: mr ? { number: mrNum, page: mrPage } : null,
                blame: blame ? getBlameInfo() : null
            });

            function getBlameInfo() {
                names.shift();

                var branch = "";

                while (names.length > 0) {
                    branch = branch === "" ? names.shift() : branch+"/"+names.shift();

                    if ( branch === repoRef.value )
                        break; 
                }

                if ( branch !== repoRef.value )
                    return null;

                return { commitElems: commitElems, branch: branch, path: names.join("/") };
            }
        }
    }

    function parseSearch(callback) {
        var stopAt = new Date().getTime() + 2000;

        parse();

        function parse() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Given up on looking for elements");
                return;
            }

            var navLinks     = document.getElementsByClassName("nav-links"),
                contentElems = document.getElementsByClassName("row-content-block"),
                boxElems     = document.getElementsByClassName("search_box");


            if ( 
                navLinks === null || 
                navLinks.length === 0 ||
                (contentElems === null && boxElems === null) ||
                (contentElems.length === 0 && boxElems.length === 0)
            ) {
                setTimeout(parse, 50);
                return;
            }

            var summary    = boxElems !== null && boxElems.length !== 0 ? boxElems[0] : contentElems[0],
                codeNav    = null,
                commitsNav = null;

            for ( var i = 0; i < navLinks[0].children.length; i++ ) {
                var li   = navLinks[0].children[i],
                    text = $(li).text();

                if ( text.toLowerCase().match(/code/) )
                    codeNav = li;
                else if ( text.toLowerCase().match(/commits/) )
                    commitsNav = li;
            }

            callback({
                type: "search",
                codeNav: codeNav,
                commitsNav: commitsNav,
                summary: summary
            });
        }
    }
}
