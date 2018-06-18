;(() => {
    let checkbox = document.getElementsByTagName('input')[0];

    chrome.storage.local.get({'playlist': false}, ({playlist}) => {
        checkbox.checked = playlist;
    });

    checkbox.addEventListener('change', (e) => {
        chrome.storage.local.set({'playlist': e.target.checked});
    });
})();
