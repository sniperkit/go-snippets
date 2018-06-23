var requiredGitSenseJSFiles = [
        "3rdparty/js/md5.js",
        "3rdparty/js/jquery-2.1.4.min.js",
        "js/arch.js",
        "js/config.js",
        "js/events/page.js",
        "js/utils/cli.js",
        "js/utils/config.js",
        "js/utils/date.js",
        "js/utils/html.js",
        "js/utils/variable.js",
        "js/github/data/repo.js",
        "js/github/data/user.js",
        "js/github/utils/page.js",
        "js/gitlab/data/user.js",
        "js/gitlab/data/repo.js",
        "js/gitlab/utils/page.js",
        "js/gitsense/data/auth.js",
        "js/gitsense/data/insight.js",
        "js/gitsense/data/repo.js",
        "js/gitsense/utils/location.js",
        "content.js"
    ],
    requiredGitSenseCSSFiles = [
        "3rdparty/octicons/octicons.css"
    ];

chrome.tabs.onUpdated.addListener(
    function(tabId, changeInfo, tab) {
        if (changeInfo.status !== "loading") 
            return

        var code = 
            "var injected = window.injectedGitSense;"+
            "window.injectedGitSense = true;"+
            "injected;";

        chrome.tabs.executeScript(tabId, {code:code, runAt:"document_start"}, inject);
 
        function inject(res) {
            if ( chrome.runtime.lastError !== undefined ) {
                if ( chrome.runtime.lastError.message.match(/permission/i) )
                    return;
                else
                    throw(chrome.runtime.lastError);
            }
        
            // This is true if we have already injected our code
            if ( res[0] )
                return;

            for ( var i = 0; i < requiredGitSenseCSSFiles.length; i++ ) {
                var file = requiredGitSenseCSSFiles[i];

                chrome.tabs.insertCSS(
                    tabId, 
                    { file: file, runAt: "document_start" }
                );
            }

            for ( var i = 0; i < requiredGitSenseJSFiles.length; i++) {
                var file = requiredGitSenseJSFiles[i];

                chrome.tabs.executeScript(
                    tabId,
                    { file: file, runAt: "document_start" }
                );
            }
        }
    }
);
