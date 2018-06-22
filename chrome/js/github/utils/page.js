var lastGitSenseShow = null;

sdes.github.utils.page = function(rule) {
    var varUtil = new sdes.utils.variable();

    this.parse = function(callback) {
        var names = window.location.pathname.replace(/^\/|\/$/, "").split(/\//);

        if ( names[names.length - 1] === "" )
            names.pop();

        if ( names.length === 1 && names[0] === "login" ) {
            callback(null);
            return;
        }

        if ( names.length !== 1 && names[0] !== "orgs" ) {
            if ( names.length > 1 )
                parseRepo(names, callback);
            else 
                callback(null);

            return;
        }

        parseOrg(names.length === 1 ? names[0] : names[1], callback);
    }

    function parseOrg(name, callback) {
        var show   = window.location.search.match(/gitsense=insight/) ? true : false,
            stopAt = new Date().getTime()+2000;

        if ( lastGitSenseShow && show )
            return;

        lastGitSenseShow = show;

        parse();

        function parse() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Given up on looking for elements");
                return;
            }

            var orgElems   = document.getElementsByClassName("orgnav"),
                containers = document.getElementsByClassName("container");

            if ( 
                orgElems === null || 
                containers == null ||
                orgElems.length === 0 ||
                containers.length === 0 ||
                containers.length < 3
            ) {
                setTimeout(parse, 50);
                return;
            }

            var container = null;

            for ( var i = 0; i < containers.length; i++ )  {
                if ( ! containers[i].className.match(/org-profile/) )
                    continue;

                container = containers[i];
                break;
            }

            callback({
                type: "org",
                org: name,
                tabs: orgElems[0],
                container: container,
                show: show
            });
        }
    }

    function parseRepo(names, callback) {
        var owner     = names.shift(),
            repo      = names.shift(),
            home      = names.length === 0 ? true : false,
            search    = names.length !== 0 && names[0] === "search" ? true : false,
            pulls     = names.length !== 0 && names[0] === "pulls" ? true : false,
            pull      = names.length !== 0 && names[0] === "pull" ? true : false,
            show      = window.location.search.match(/gitsense=insight/) ? true : false,
            stopAt    = new Date().getTime() + 1000,
            trackerId = "gitsense-content-tracker",
            iframeId  = "gitsense-content-iframe",
            container = document.getElementById(trackerId),
            iframe    = document.getElementById(iframeId);

        if ( container !== null ) {
            stopAt += 2000;
            waitForNewContainers();
            return;
        }

        findContainers();
        lastGitSenseShow = show; 

        function waitForNewContainers() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Giving up in waiting for new containers");
                return;
            }

            var container = document.getElementById(trackerId);

            if ( container !== null ) {
                setTimeout(waitForNewContainers, 50);
                return;
            }

            stopAt = new Date().getTime() + 2000;
            findContainers();
        }

        function findContainers() {
            if ( new Date().getTime() > stopAt ) {
                console.error("Given up on looking for containers");
                return;
            }

            var containers = document.getElementsByClassName("container");

            if ( containers === null || containers.length < 3 ) {
                setTimeout(findContainers, 50);
                return;
            }

            var tabs    = null,
                content = null,
                fileNav = null,
                branch  = null;

            for ( var i = 0; i < containers.length; i++ ) {
                var container = containers[i];

                if ( container.className.match(/js-repo-nav/) && ! search ) {
                    if ( container.childNodes === undefined || container.childNodes[1] === undefined )
                        continue;

                    tabs = container;
                } else if ( 
                    (container.className === "container" && search) ||
                    (container.className.match(/experiment-repo-nav/) && ! search)
                ) {
                    content    = container;
                    content.id = trackerId;

                    if ( show )
                        $(content).hide();
                }
            }

            var prtabs   = pull ? document.getElementsByClassName("tabnav-pr") : null,
                pullsNav = pulls ? document.getElementsByClassName("subnav") : null;

            if ( 
                (tabs === null && !search) || 
                content === null ||
                (pull && (prtabs === null || prtabs.length === 0)) ||
                (pulls && (pullsNav === null || pullsNav.length === 0))
            ) {
                setTimeout(findContainers, 50);
                return;
            }

            if ( home ) {
                var fileNavElems = document.getElementsByClassName("file-navigation"),
                    selectElems  = document.getElementsByClassName("branch-select-menu");

                if ( 
                    fileNavElems === null || 
                    fileNavElems.length === 0 ||
                    selectElems === null ||
                    selectElems.length === 0 ||
                    selectElems[0].children.length === 0 
                ) {
                    setTimeout(findContainers, 50);
                    return;
                }

                fileNav = fileNavElems[0];

                var button = selectElems[0].children[0],
                    type   = button.children[0].innerHTML.replace(/^ /, ""),
                    value  = button.children[1].innerHTML.replace(/^ /g, "");

                if ( type === "Branch:" )
                    branch = value;
            }

            callback({
                type: "repo",
                owner: owner,
                repo: repo,
                home: home,
                show: show,
                tabs: tabs,
                search: search ? getSearch() : null,
                content: content,
                fileNav: fileNav,
                branch: branch,
                pull: pull ? getPull(prtabs[0]) : null,
                pulls: pulls ? { nav: pullsNav[0] } : null
            });
        }

        function getSearch() {
            var navs         = document.getElementsByClassName("underline-nav")[0],
                temp         = window.location.search.split("&"),
                selectedType = null,
                typeToNav    = {};

            for ( var i = 0; temp.length; i++ ) {
                var query = temp[i].split("="),
                    key   = query[0].replace(/\?/,"");

                if ( key !== "type" )
                    continue;

                if ( query.length === 1 || query[1] === "" )
                    selectedType = "code";
                else
                    selectedType = query[1].toLowerCase();

                break;
            }

            if ( selectedType === null )
                return null;

            for ( var i = 0; i < navs.children.length; i++ ) {
                var nav  = navs.children[i],
                    text = nav.innerText.toLowerCase(),
                    type = null;

                if ( text.match(/^\s*code/) )
                    type = "code";
                else if ( text.match(/^\s*commits/) )
                    type = "commits";
                else if ( text.match(/^\s*issues/) )
                    type = "issues";
                else if ( text.match(/^\s*gscode/) )
                    type = "gscode";
                else if ( text.match(/^\s*gscommits/) )
                    type = "gscommits";
                else if ( text.match(/^\s*gsdiffs/) )
                    type = "gsdiffs";

                nav.selected = nav.className.match(/selected/) ? true : false;

                typeToNav[type] = nav;
            }

            return { selectedType: selectedType, typeToNav: typeToNav, navs: navs };
        }

        function getPull(prtabs) {
            var number = parseInt(names[1]),
                tabs   = null;

            for ( var i = 0; i < prtabs.children.length; i++ ) {
                var elem = prtabs.children[i];

                if ( ! elem.className.match("tabnav-tabs") )
                    continue;
                    
                tabs = elem;
                break;
            }

            if ( tabs === null )
                return null;

            var selected = "";

            for ( var i = 0; i < tabs.children.length; i++ ) {
                var tab = tabs.children[i];

                if ( ! tab.className.match(/selected/) )
                    continue;
    
                var text = tab.innerText.toLowerCase();

                if ( text.match(/conversation/) )
                    selected = "conversation";
                else if ( text.match(/commits/) )
                    selected = "commits";
                else if ( text.match(/files/) )
                    selected = "files";
            }

            return { tabs: tabs, selectedTab: selected, number: number };
        }
    }
}
