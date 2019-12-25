<template>
  <div class="login-page">
        <div class="login">
            <div class="login-container">
                  <div class="login-form-column">
                      <form v-on:submit.prevent="login">
                          <h3>Hello!</h3>
                          <p>Welcome to our little Vue demo powered by CometChat. Login with the username "superhero1" or "superhero2" and test the chat out.
                             To create your own user, see <a href="https://prodocs.cometchat.com/reference#createuser">our documentation</a>   </p>
                          <div class="form-wrapper">
                            <label>Nickname</label>
                            <input type="text" name="username" v-model="username" placeholder="Enter your username" class="form-control input" required>
                          </div>
                          <div class="form-wrapper">
                            <label>Password</label>
                            <input type="password" name="password" v-model="password" placeholder="Enter your password" class="form-control input" required>
                          </div>
                          <button type="submit">LOG IN &nbsp;&nbsp;<span v-if="showSpinner" class="fa fa-spin fa-spinner"></span> </button>
                      </form>
                  </div>

                  <div class="login-image-column">
                      <div class="image-holder">
                          <img src="../assets/login-illustration.svg" alt="">
                      </div>
                  </div>
           </div>
           </div>
        </div>
</template>

<script>
import { storeAccessToken, storeUser } from '../utils/auth'

export default {
  data() {
    return {
      username: "",
      password: "",
      showSpinner: false
    };
  },
  methods: {
    async login() {
      this.showSpinner = true;
      const loginPayload = {
        nickname: this.username,
        password: this.password
      };

      try {
        const {data} = await this.$http.post('/user/v1/login', loginPayload);
        this.$store.commit('authenticate', data)
        storeAccessToken(data.access_token);
        storeUser(data.user);

        this.$router.push({
          name: 'chat'
        });
      } catch (e) {
        console.error('Error Response', e);
      }

      this.showSpinner = false;
    }
  }
};
</script>
