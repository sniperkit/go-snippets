import * as store from './localstore'
import KEYS from './storekeys'


/**
 * 所有的搜索引擎
 */
export const allEngines = [{
    name: "Google",
    icon: "./static/engines/google.png",
    action: "https://www.google.com/search",
    inputName: "q"
}, {
    name: "Bing",
    icon: "./static/engines/bing.png",
    action: "https://www.bing.com/search",
    inputName: "q"
}, {
    name: "Yahoo",
    icon: "./static/engines/yahoo.png",
    action: "https://search.yahoo.com/search?",
    inputName: "p"
}, {
    name: "Stackoverflow",
    icon: "./static/engines/stackoverflow.png",
    action: "https://stackoverflow.com/search?",
    inputName: "q"
}, {
    name: "Baidu",
    icon: "./static/engines/baidu.png",
    action: "https://www.baidu.com/s?",
    inputName: "wd"
}]


/**
 * 获取默认的搜索引擎 index 
 * @return {number}
 */
export const getEngineIndx = store.getData(KEYS.ENGINIE).then(data => {


    // 没有默认的搜索引擎 
    // 设置 谷歌
    if (store.isEmptyData(data)) {
        let defaults = {
            [KEYS.ENGINIE]: 0
        }
        store.storeData(defaults)
        return defaults[KEYS.ENGINIE]
    }
    return data[KEYS.ENGINIE]
})

/**
 *  设置新的搜索引擎
 * @param {number} index 
 */
export const setEngineIndx = index => {
    store.storeData({
        [KEYS.ENGINIE]: index
    })
}