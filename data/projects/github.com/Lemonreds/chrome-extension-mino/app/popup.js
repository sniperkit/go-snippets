import Vue from 'vue'
import App from './popup/app.vue'

const root = document.createElement('div')
document.getElementsByTagName('body')[0].appendChild(root)

new Vue({
    render: (h) => h(App)
}).$mount(root)