import Vue from 'vue'
import Router from 'vue-router'
import Chat from './views/Chat.vue'
import Login from './views/Login.vue'
// import store from './store'

Vue.use(Router)

const router =  new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'homepage',
      redirect: 'login'
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/chat',
      name: 'chat',
      component: Chat
    }
  ]
})

// router.beforeEach((to, from, next) => {
//   const auth = store.getters.auth
//   const loggedIn = auth.accessToken && auth.user
//   if (!loggedIn) {
//     if (to.name != 'login')
//       next({ name: 'login' })
//   } else {
//     if (to.name == 'login')
//       next({ name: 'chat' })
//   }
// })
export default router