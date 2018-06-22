renderGitSenseIframe();

function renderGitSenseIframe() {
    var renderTo = document.getElementById("gitsense-iframe-body");

    if ( renderTo === null )
        throw("No gitsense body found");

    var iframeSrc     = IFRAME_SRC,
        targetOrigin  = TARGET_ORIGIN,
        lastHeight    = 0,
        overlayWindow = null,
        resize        = null,
        peekWidth     = 350;

    renderTo.innerHTML = "";
    renderTo.parentNode.style.padding = 0;
    renderTo.style.height = "300px" 

    var iframe = document.createElement("iframe");

    iframe.style.width  = "100%";
    iframe.style.height = "100%";
    iframe.style.border = 0;
    iframe.setAttribute("src", iframeSrc+window.location.hash);
    renderTo.appendChild(iframe);

    setTimeout(
        function() {
            renderTo.style.height = null;
            getFrameSize();
        },
        1000
    );

    window.addEventListener("message", receiveMessage, false);
    
    window.onhashchange = function() {
        iframe.contentWindow.postMessage("hash:"+window.location.hash, targetOrigin);
    }

    function getFrameSize() { 
        iframe.contentWindow.postMessage("height", targetOrigin);
        setTimeout(getFrameSize,250);
    }

    function receiveMessage(event) {
        if ( event.origin !== targetOrigin )
            return;

        var type = typeof(event.data);

        switch(type) {
            case "number":
                setHeight(event.data);

                return;
            case "string":
                var meta = event.data.split(":")[0];
    
                if ( meta === "hash" )
                    setHash(event.data.replace(/^hash:/,""));
                else if ( meta === "commit" )
                    viewCommit(event.data.replace(/^commit:/, ""));
                else if ( meta === "file" )
                    viewFile(event.data.replace(/^file:/,""));
                else if ( meta === "gswin" )
                    openGitSenseWindow(event.data.replace(/^gswin:/,""));
                else
                    throw("Unrecognized meta data '"+meta+"'") 

                return;
            default:
                throw("Unrecognized type '"+type+"'");
        }
    }

    function setHeight(height) {
        if ( height === lastHeight )
            return;

        iframe.style.height = height+"px";
        lastHeight = height;
    }

    function setHash(hash) {
        if ( window.location.hash !== hash )
            window.location.hash = hash;
    }

    function openGitSenseWindow(href, file) {
        var regexp = new RegExp("^"+targetOrigin);

        if ( ! href.match(regexp) )
            return;

        var iframe = openOverlayWindow(href);

            setTimeout(
                function() {
                    iframe.contentWindow.postMessage("init", targetOrigin);
                },
                500
            );
    }

    function openOverlayWindow(href) {
        if ( overlayWindow !== null ) 
            overlayWindow.parentNode.removeChild(overlayWindow);

        if ( resize !== null )
            resize.parentNode.removeChild(resize);

        var width  = window.innerWidth - 30,
            height = window.innerHeight - 25;

        overlayWindow = document.createElement("body");
        overlayWindow.style.width  = (width - peekWidth)+"px";
        overlayWindow.style.height = height+"px";
        overlayWindow.style.backgroundColor = "white";
        overlayWindow.style.zIndex = 1000000;
        overlayWindow.style.position = "fixed";
        overlayWindow.style.top = 10+"px";
        overlayWindow.style.left = peekWidth+"px";
        overlayWindow.style.border = "0px";
        overlayWindow.style.boxShadow = "0px 0px 26px 0px rgba(48,48,48,1)";
        overlayWindow.style.overflow = "hidden";

        var title = document.createElement("div");
        title.style.backgroundColor = "black";
        title.style.height = "30px";
        title.style.display = "table";
        title.style.fontSize = "12px";
        title.style.color = "white";

        var sameLink = document.createElement("a");
        sameLink.setAttribute("class", "aui-icon aui-icon-small aui-iconfont-link");
        sameLink.setAttribute("title", "Open in current window");
        sameLink.href = href;
        sameLink.style.fontWeight = "bold";
        sameLink.style.display = "block";
        sameLink.style.overflow = "hidden";
        sameLink.style.textOverflow = "ellipsis";
        sameLink.style.whiteSpace = "nowrap";
        sameLink.style.color = "white";

        var sameLinkCell = document.createElement("div");
        sameLinkCell.style.display = "table-cell";
        sameLinkCell.style.verticalAlign = "middle";
        sameLinkCell.style.paddingRight = "15px";
        sameLinkCell.style.paddingTop = "1px";
        sameLinkCell.appendChild(sameLink);

        var newLinkText = document.createTextNode(href);
        var newLink = document.createElement("a");
        newLink.setAttribute("title", "Open in new window");
        newLink.setAttribute("target", "_blank");
        newLink.href = href;
        newLink.style.fontWeight = "bold";
        newLink.style.display = "block";
        newLink.style.width = "500px";
        newLink.style.overflow = "hidden";
        newLink.style.textOverflow = "ellipsis";
        newLink.style.whiteSpace = "nowrap";
        newLink.style.color = "white";
        newLink.appendChild(newLinkText);

        var newLinkCell = document.createElement("div");
        newLinkCell.style.display = "table-cell";
        newLinkCell.style.width = "100%";
        newLinkCell.style.verticalAlign = "middle";
        newLinkCell.style.paddingLeft = "15px";
        newLinkCell.style.paddingTop = "1px";
        newLinkCell.style.fontFamily = "monospace";
        newLinkCell.appendChild(newLink);

        var close = document.createElement("span");
        close.setAttribute("class", "aui-icon aui-icon-small aui-iconfont-close-dialog");
        close.style.marginRight = "10px";
        close.style.cursor = "pointer";
        close.style.fontSize = "14px";
        close.style.fontWeight = "bold";

        var closeCell = document.createElement("div");
        closeCell.style.display = "table-cell";
        closeCell.style.verticalAlign = "middle";
        closeCell.appendChild(close);

        var size = document.createElement("span");
        size.setAttribute("class", "aui-icon aui-icon-small aui-iconfont-focus");
        size.setAttribute("title", "Increase window width");
        size.style.marginRight = "15px";
        size.style.cursor = "pointer";
        size.style.fontSize = "14px";
        size.style.fontWeight = "bold";
        size.style.transform = "rotate(45deg)";

        var sizeCell = document.createElement("div");
        sizeCell.style.display = "table-cell";
        sizeCell.style.verticalAlign = "middle";
        sizeCell.appendChild(size);

        title.appendChild(newLinkCell);
        title.appendChild(sameLinkCell);
        title.appendChild(sizeCell);
        title.appendChild(closeCell);

        var iframe = document.createElement("iframe");
        iframe.style.border = "0px";
        iframe.style.width  = "100%";
        iframe.style.height = (height - 30)+"px";
        iframe.src = href;

        overlayWindow.appendChild(title);
        overlayWindow.appendChild(iframe);

        document.body.appendChild(overlayWindow);

        resize = document.createElement("div");
        resize.style.position = "fixed";
        resize.style.left     = (parseInt(overlayWindow.style.left) - 2)+"px";
        resize.style.top      = overlayWindow.style.top;
        resize.style.height   = overlayWindow.style.height;
        resize.style.width    = "5px";
        resize.style.cursor   = "col-resize";
        resize.style.zIndex   = overlayWindow.style.zIndex+1;
        resize.style.backgroundColor = "transparent";

        var screen = null;

        resize.onmousedown = function(e) {
            screen = document.createElement("div");
            screen.style.width  = window.innerWidth+"px";
            screen.style.height = window.innerHeight+"px"; 
            screen.style.position = "fixed";
            screen.style.top    = 0;
            screen.style.left   = 0;
            screen.style.zIndex = overlayWindow.style.zIndex+1;
            screen.style.backgroundColor = "white";
            screen.style.cursor  = "col-resize";
            screen.style.opacity = .3;
            document.body.appendChild(screen);

            var box = document.createElement("div");
            box.style.position = "fixed";
            box.style.left     = (parseInt(overlayWindow.style.left) - 2)+"px";
            box.style.top      = overlayWindow.style.top;
            box.style.height   = (parseInt(overlayWindow.style.height) - 6)+"px";
            box.style.width    = (parseInt(overlayWindow.style.width) - 6)+"px";
            box.style.cursor   = "col-resize";
            box.style.border   = "2px solid #333";

            screen.appendChild(box);

            screen.onmousemove = function(e) {
                var diff = parseInt(resize.style.left) - e.clientX;
                box.style.width = (parseInt(overlayWindow.style.width)+diff)+"px";
                box.style.left  = e.clientX+"px";
            }

            screen.onmouseup = function(e) {
                screen.parentNode.removeChild(screen);
                screen = null;

                var diff = parseInt(resize.style.left) - e.clientX;
                overlayWindow.style.width = (parseInt(overlayWindow.style.width)+diff)+"px";
                overlayWindow.style.left  = e.clientX+"px";
                resize.style.left = e.clientX+"px";
                iframe.contentWindow.postMessage("resize", targetOrigin);
                peekWidth = e.clientX;
            }
        }

        document.body.appendChild(resize);

        close.onclick = function() {
            overlayWindow.parentNode.removeChild(overlayWindow);
            overlayWindow = null;
            resize.parentNode.removeChild(resize);
            resize = null;
        }

        size.onclick = function() {
            var expand = size.className.match(/unfocus/) ? false : true;

            if ( expand ) {
                resize.style.left         = "15px";
                overlayWindow.style.left  = "15px";
                overlayWindow.style.width = (width - 15)+"px";
                size.setAttribute("class", "aui-icon aui-icon-small aui-iconfont-unfocus");
                size.setAttribute("title", "Shrink window width");
            } else {
                resize.style.left         = peekWidth+"px";
                overlayWindow.style.left  = peekWidth+"px";
                overlayWindow.style.width = (width - peekWidth)+"px";
                size.setAttribute("class", "aui-icon aui-icon-small aui-iconfont-focus");
                size.setAttribute("title", "Increase window width");
            }

            iframe.contentWindow.postMessage("resize", targetOrigin);
        }

        return iframe;
    }
}
