"use strict";

window.addEventListener("message", receiveMessage, false);

var iframeEvent  = null,
    params       = null,
    iframe       = null,
    rendered     = false,
    parentOrigin = null,
    childOrigin  = null,
    lastHash     = null,
    lastHeight   = null;

function receiveMessage(event) {
    if ( typeof(event.data) === "number" ) {
        if ( event.data !== 100 )
            setHeight(event.data);

        return;
    }

    if ( typeof(event.data) === "string" && event.data.match(/^{/) ) {
        if ( rendered )
            return;

        rendered     = true;
        params       = JSON.parse(event.data);
        iframeEvent  = event;
        parentOrigin = event.origin;
        childOrigin  = params.targetOrigin;

        renderGitSenseIframe();
        return;
    }

    switch(event.origin) {
        case parentOrigin:
            if ( event.data.match(/resize:/) )
                window.innerWidth = parseInt(event.data.split(":")[1]);

            if ( event.data.match(/^hash:/) || event.data.match(/resize/) || event.data.match(/init/) )
                iframe.contentWindow.postMessage(event.data, "*");

            break;
        case childOrigin:
            if ( iframeEvent === null )
                return;

            var meta = event.data.split(":")[0];

            if ( 
                meta.toLowerCase() !== "gswin" && 
                meta !== "commit" && 
                meta !== "hash" && 
                meta !== "href" &&
                meta !== "page" &&
                meta !== "reload"
            ) {
                console.warn("Ignoring '"+meta+"' message");
                return;
            }

            iframeEvent.source.postMessage(
                params.id+":::::"+event.data,
                iframeEvent.origin
            );

            break;
        default:
            throw("Unknown caller");
    }
}

function renderGitSenseIframe() {
    var renderTo = document.getElementById("gitsense-iframe-body");

    if ( renderTo === null )
        throw("No gitsense body found");

    var lastHeight = 0,
        resize     = null,
        peekWidth  = 350,
        hash       = params.hash || "",
        slowAfter  = new Date().getTime();

    renderTo.innerHTML = "";
    renderTo.parentNode.style.padding = 0;
    renderTo.style.height = params.height === undefined ? null : (params.height-15)+"px";

    iframe = document.createElement("iframe");
    iframe.style.width  = "100%";
    iframe.style.height = "100%";
    iframe.style.border = 0;
    iframe.setAttribute("src", params.iframeSrc+hash);
    renderTo.appendChild(iframe);

    if ( params.maxHeight !== undefined ) {
        setTimeout(
            function() {
                iframe.contentWindow.postMessage(
                    "maxHeight:"+params.maxHeight, params.targetOrigin
                );
            },
            1000
        );
    }

    if ( params.noResize )
        return;

    setTimeout(
        function() {
            renderTo.style.height = null;
            getFrameSize();
        },
        50
    );

    function getFrameSize() {
        try { 
            iframe.contentWindow.postMessage("height", params.targetOrigin);
        } catch (e) {
            console.log("exception thrown");
        }

        setTimeout(getFrameSize, 1000);
    }
}

function setHeight(height) {
    if ( height === lastHeight )
        return;

    iframe.style.height = height+"px";
    lastHeight = height;

    if ( iframeEvent === null )
        return;

    iframeEvent.source.postMessage(
        params.id+":::::height:"+height,
        iframeEvent.origin
    );
}
