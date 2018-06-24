/**
 * 图片处理相关
 */
import * as store from './localstore.js'


const KEY_BACKGROUND = 'background'

export const storeBase64 = url => {
    
    localStorage.clear(KEY_BACKGROUND)

    return new Promise((resolve,reject)=>{
        let xhr = new XMLHttpRequest()
        let fr = new FileReader()
        xhr.open("GET", url, true)
        xhr.responseType = "blob"        
        xhr.send(null) 
        xhr.onload = () => { 
            fr.onload = e => {       
                localStorage.setItem(KEY_BACKGROUND, e.target.result)   
                console.log('Store OK')
                resolve(e.target.result)             
            }
            fr.readAsDataURL(xhr.response)            
        }
    })
}

export const fetchBase64 = () => {
    return localStorage.getItem(KEY_BACKGROUND)
}

// localStorage.clear()