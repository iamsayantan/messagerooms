import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import  {getAccessToken, getAuthenticatedUser } from './utils/auth'

Vue.config.productionTip = false

const accessToken = getAccessToken()
const user = getAuthenticatedUser()

console.log('AccessToken', accessToken)
console.log('AuthUser', user)

if (accessToken && user) {
  store.commit('authenticate', { user, access_token: accessToken })
}

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
