<template>
  <v-app id="login" class="primary">
    <v-content>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4 lg4>
            <v-card class="elevation-1 pa-3">
              <v-card-text>
                <div class="layout column align-center">
                  <img src="../static/m.png" alt="Vue Material Admin" width="120" height="120">
                  <h1 class="flex my-4 primary--text">Messagerooms</h1>
                </div>
                <v-form>
                  <v-text-field append-icon="person" name="nickname" label="Nickname" type="text"
                                v-model="login.nickname"></v-text-field>
                  <v-text-field append-icon="lock" name="password" label="Password" id="password" type="password"
                                v-model="login.password"></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn block color="primary" @click="doLogin" :loading="loading">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  export default {
    layout: 'default',
    data: () => ({
      loading: false,
      login: {
        nickname: null,
        password: null
      }
    }),

    methods: {
      doLogin() {
        this.loading = true;

        this.$auth.loginWith('local', {
          data: {
            nickname: this.login.nickname,
            password: this.login.password
          }
        })
          .catch(e => {
            console.error(e)
            this.loading = false
          })
      }
    }

  };
</script>
<style scoped lang="css">
  #login {
    height: 50%;
    width: 100%;
    position: absolute;
    top: 0;
    left: 0;
    content: "";
    z-index: 0;
  }
</style>
