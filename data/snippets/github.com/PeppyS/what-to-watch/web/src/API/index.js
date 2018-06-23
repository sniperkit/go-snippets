import axios from 'axios'
import { API_URL } from '../Config'

const request = async (method, url, payload) => {
    return axios({
        method,
        url,
        payload,
        baseURL: API_URL,
        timeout: 100000,
        headers: {
            'Content-Type': 'application/json',
        }
    })
}

export default {
    get: (url, payload) => request('get', url, payload)
}
