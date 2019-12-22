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
                            <input type="text" name="username" id="username" v-model="username" placeholder="Enter your username" class="form-control" required>
                          </div>
                          <div class="form-wrapper">
                            <label>Password</label>
                            <input type="password" name="password" id="username" v-model="password" placeholder="Enter your password" class="form-control" required>
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
import { postRequest } from '../utils/request'
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
      this.showSpinner = true
      const loginPayload = {
        nickname: this.username,
        password: this.password
      }

      try {
        const response = await postRequest('http://localhost:9050/user/v1/login', loginPayload)
        this.$store.commit('authenticate', response)
        
        storeAccessToken(response.access_token)
        storeUser(response.user)
        
        this.$router.push({
          name: 'chat'
        })
        console.log(response)
      } catch (e) {
        console.error(e)
      }
      this.showSpinner = false
    },
    // authLoginUser() {
    //   var apiKey = process.env.VUE_APP_COMMETCHAT_API_KEY;
    //   this.showSpinner = true;

    //   CometChat.login(this.username, apiKey).then(
    //     () => {
    //       this.showSpinner = false;
    //       this.$router.push({
    //         name: "chat"
    //       });
    //     },
    //     error => {
    //       this.showSpinner = false;
    //       alert("Whops. Something went wrong. This commonly happens when you enter a username that doesn't exist. Check the console for more information")
    //       console.log("Login failed with error:", error.code);
    //     }
    //   );
    // }
  }
};
</script>
