
/**
 * 下载图片
 * @param {string} url web路径
 * @param {function} cb 回调函数
 */
const download = (url, cb) => {
    typeof cb === 'function' ? chrome.downloads.download({
        url
    }, cb) : chrome.downloads.download({
        url
    })
}

export default download