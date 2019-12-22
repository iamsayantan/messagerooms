import vuex from 'vuex'
import Vue from 'vue'

Vue.use(vuex)
const store = new vuex.Store({
    state: {
        auth: {
            user: null,
            accessToken: null
        }
    },
    mutations: {
        authenticate(state, {user, access_token}) {
            state.auth = Object.assign({}, {user, accessToken: access_token})
        }
    },
    getters: {
        auth: (state) => {
            return state.auth
        }
    }
})

export default store