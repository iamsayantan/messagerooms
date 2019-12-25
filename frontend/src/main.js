import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

import App from './App.vue'
import router from './router'
import store from './store'

import  {getAccessToken, getAuthenticatedUser } from './utils/auth'

Vue.config.productionTip = false

const accessToken = getAccessToken()
const user = getAuthenticatedUser()

const axiosInstance = axios.create({
  baseURL: 'http://localhost:9050/',
  headers: {
    'Content-Type': 'application/json'
  }
})

Vue.use(VueAxios, axiosInstance)

if (accessToken && user) {
  store.commit('authenticate', { user, access_token: accessToken })
}

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
