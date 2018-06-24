import * as store from './localstore.js'
import KEYS from './storekeys.js'

/**
 * 默认设置
 */
export const defaluts = {
    application: {
        name: 'application',
        text: 'Application',
        status: true,
        storeKey: KEYS.SETTING_APPLICATION
    },
    bookmark: {
        name: 'bookmark',
        text: 'Bookmarks',
        status: true,
        storeKey: KEYS.SETTING_BOOKMARK
    },
    commonsites: {
        name: 'commonsites',
        text: 'CommonSites',
        status: false,
        storeKey: KEYS.SETTING_COMMONSITES
    }
}

/**
 * bing设置
 */
export const bing = {
    name: 'bing',
    status: true,
    inital : true 
}

/**
 * 获取一个配置项
 * 当配置项在本地存储不存在时 
 * 需要指定一个默认值 写入本地存储
 * @param {string} key 
 * @param {Object} value 
 * @return {Promise}
 */
function fetch(key, value) {
    return new Promise((resolve, reject) => {
        store.getData(key).then(data => {
            if (store.isEmptyData(data)) { //maybe frist store this data
                store.storeData({
                    [key]: value
                })
                resolve(value)
            }
            resolve(data[key])
        })
    })
}

export const fetchApplication = fetch(KEYS.SETTING_APPLICATION, defaluts.application)
export const fetchBookmark = fetch(KEYS.SETTING_BOOKMARK, defaluts.bookmark)
export const fetchCommonsites = fetch(KEYS.SETTING_COMMONSITES, defaluts.commonsites)
export const fetchBing = fetch(KEYS.SETTING_BING, bing)


// 获取上面3个的数据数组
export const fetchFeatures = Promise.all([fetchApplication, fetchBookmark, fetchCommonsites]).then(datas => {

    let data = {
        application: datas[0],
        bookmark: datas[1],
        commonsites: datas[2]
    }
    return Promise.resolve(data)
})


/**
 * 修改设置项
 * @param {string} type 
 * @param {boolean} value 
 */
export const modify = (type, value) => {
    switch (type) {
        case 'application':
            set(KEYS.SETTING_APPLICATION, value, defaluts.application)
            break;
        case 'bookmark':
            set(KEYS.SETTING_BOOKMARK, value, defaluts.bookmark)
            break;
        case 'commonsites':
            set(KEYS.SETTING_COMMONSITES, value, defaluts.commonsites)
            break;
        case 'bing':
            set(KEYS.SETTING_BING, value, bing)
            break;
        default:
            break;
    }
}

/**
 * 
 * @param {String} key 
 * @param {Boolean} value 
 * @param {Object} source 
 */
function set(key, status, source) {
    if (typeof status === 'boolean') {

        Object.assign(source, {
            status: status
        })
        store.storeData({
            [key]: source
        })
    }
}


/**
 * 默认颜色
 */
export const colors = ['#3b5998','#5f7ec1','#2b6b73','#337ab7','#f0ad4e',
                        '#3b987a','#cdbfe3','#d9534f', '#983b59','#987a3b','#88983b','#674a34','#67758a']
                        
/**
 * 存储颜色到本地存储
 * @param {String} value 
 */
export const setColor = value => {
    store.storeData({
        [KEYS.SETTING_COLOR] : value
    })
}
/**
 * 获取本地存储的颜色
 * @return {String}
 */
export const fetchColor = fetch(KEYS.SETTING_COLOR,colors[0])


// store.getAll().then(v=>{
//     console.log(v)
// })

//store.clearAll()