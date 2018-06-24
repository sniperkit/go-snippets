;(() => {
    const TARGET = 'http://127.0.0.1:12345'

    chrome.commands.onCommand.addListener((command) => {
        chrome.storage.local.get({'playlist': false}, ({playlist}) => {
            chrome.tabs.query({url: 'https://www.youtube.com/watch*'}, (tabs) => {
                for (let tab of tabs) {
                    if (tab.audible) {
                        let xhr = new XMLHttpRequest();
                        let encodedURL = encodeURIComponent(tab.url.substr(24));
                        xhr.open('GET', `${TARGET}/${encodedURL}?playlist=${playlist}`);
                        xhr.send();

                        break;
                    }
                }
            })
        });
    });
})();
