/*
 * 数据存储结构 
 * {
 *  key : value
 * }
 */
 

/**
 * 存储数据 
 * @param {Object} data 
 */
export const storeData = data => {
     chrome.storage.local.set(data)   
}

/**
 * 获取数据
 * @param {string} key 
 * @return {Promise}
 */
export const getData = key => {

    return new Promise((resolve, reject) => {
        chrome.storage.local.get(key, value => {        
            resolve(value)
        })
    })
}

/**
 * 清空本地存储
 */
export const clearAll = () => {
    chrome.storage.local.clear(() => {
        console.log('Clear all localStore.')
    })
}

/**
 * 获取所有本地存储
 * @return {Promise}
 */
export const getAll = () => {
    return getData(null)
}

/**
 * 判断是否是空数据
 * @param {Object} obj 
 * @return {Boolean} 
 */
export const isEmptyData = obj => {

    if (Object.prototype.toString.call(obj) === '[object Object]') {
        for (let key in obj) {
            return false
        }
        return true
    }
    return false
}


//clearAll()

// getAll().then( v=>{
//     console.log(v)
// })