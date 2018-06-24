import Vue from 'vue'
import store from './store'
import App from './app.vue'


import './styles/style.sass'

const root = document.createElement('div')
document.getElementsByTagName('body')[0].appendChild(root)

new Vue({
    store,
    render: (h) => h(App)
}).$mount(root)