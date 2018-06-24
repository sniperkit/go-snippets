import * as store from './localstore.js'
import KEYS from './storekeys'
import {storeBase64,fetchBase64} from './base64.js'

// 图片对象的具体信息
/**
 * id 
 * url: 具体路径
 * date: 日期 如20180422
 * copyright: 图片信息和作者
 */
class Image {
    constructor(id, url, date, copyright) {
        this.id = id
        this.url = url
        this.date = date
        this.copyright = copyright
    }
}

/**
 * format : 默认返回JSON格式
 * idx:   -1最多获取到明天,1 今天
 * n:     数量
 */
let defaults = {
    format: 'js',
    idx: -1,
    n: 1
}

/**
 * 返回 Promse.resolve(value)   
 * value: 对象数组 获取失败返回 []
 *        对象包含图片信息 
 */
const BingAPI = options => {

    let xhr = new XMLHttpRequest(),
        data, config, result, id

    // 合并 配置
    if (typeof options === 'undefined') {
        config = defaults
    } else {
        config = {}
        Object.assign(config, defaults, options)
    }
    // open ajax
    xhr.open('GET', `http://cn.bing.com/HPImageArchive.aspx?format=${config.format}&idx=${config.idx}&n=${config.n}`, true)
    xhr.send(null)
    // 回调函数
    return new Promise((resolve, reject) => {
        xhr.onreadystatechange = () => {
            if (xhr.readyState == 4) {
                if (xhr.status == 200) {

                    data = JSON.parse(xhr.responseText)
                    result = []
                    id = 0

                    data.images.forEach(image => {
                        result.push(
                            new Image(id++, `https://cn.bing.com${image.url}`, image.enddate, image.copyright))
                    })
                    // 根据日期 顺排序
                    resolve(result.reverse())
                } else {
                    // 错误 返回空数组
                    console.log('internet error / [bing]')
                    reject([])
                }
            }
        }
    })
}



/**
 * 获取上一次设置的壁纸image对象
 * @return {Image}
 */
export const getBackground = new Promise((resolve, reject) => {

    store.getData(KEYS.BACKGROUND).then(data => {
        if (store.isEmptyData(data)) {
            //  获取失败(可能是第一次启动) 
            //  设置最新的为默认壁纸
            BingAPI({
                n: 1
            }).then(data => {
                let current = data[0]
                saveBackground(current)
                current.base64 = fetchBase64()
                resolve(current)
            })
        } else {            
            let base64 = fetchBase64()
            if(base64 === null){
                // 存入新的壁纸时 （写入本地base64时)
                // 页面关闭或者刷新了
                // 重新写入base64
                setTimeout(() => {
                    storeBase64( data[KEYS.BACKGROUND].url)
                }, 300)
            }else{
                data[KEYS.BACKGROUND].base64 = fetchBase64()
            }  
            resolve(data[KEYS.BACKGROUND])
        }
    })
})



/**
 *  保存 image对象到本地存储
 * @param {Image} image 
 */
export const saveBackground = image => {

    store.storeData({
        [KEYS.BACKGROUND]: image
    })

    // TODO 2018/04/26
    // 不设置计数器会重复发送 2个http请求
    // 导致写入 base64 数据时候 阻塞太久 首页白屏的问题
    // 目前还没有解决方法 只能用个暴力的...
    setTimeout(() => {
        storeBase64(image.url)
    }, 300)
}



/**
 * 获取bing图片 8张
 * 同时从本地存储获取 壁纸 
 * 获取失败则设置最新的图片为壁纸
 * 
 * @return {Promise.resolve(Image,Array[Images])}
 */
export const getBingImages = new Promise((resolve,reject)=>{
    BingAPI({
        n: 8
    }).then(images => {
        resolve(images)
    })
})