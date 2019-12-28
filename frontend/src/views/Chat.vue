<template>
  <div class="booker">
    <nav-bar :name="auth.user.nickname"/>
    <div class="chat">
      <div class="container">
        <div class="msg-header">
          <div class="active">
            <h5>#{{ room_details.room_name }}</h5>
          </div>
        </div>

        <div class="chat-page">
          <div class="msg-inbox">
            <div class="chats" id="chats">
              <div class="msg-page" id="msg-page">

                <div
                    v-if="loadingMessages"
                    class="loading-messages-container"
                >
                  <spinner :size="100"/>
                  <span class="loading-text">
                            Loading Messages
                          </span>
                </div>
                <div class="text-center img-fluid empty-chat" v-else-if="!groupMessages.length">
                  <div class="empty-chat-holder">
                    <img src="../assets/empty-state.svg" class="img-res" alt="empty chat image">
                  </div>

                  <div v-if="is_member">
                    <h2> No new message? </h2>
                    <h6 class="empty-chat-sub-title">
                      Send your first message below.
                    </h6>
                  </div>
                  <div v-else>
                    <h2> You are not a member. </h2>
                    <h6 class="empty-chat-sub-title">
                      <a href="#" @click.prevent="joinRoom">Click</a> here to join.
                    </h6>
                  </div>
                </div>

                <div v-else>
                  <div v-for="message in groupMessages" v-bind:key="message.id">
                    <div class="received-chats" v-if="message.created_by.id !== auth.user.id">
<!--                      <div class="received-chats-img">-->
<!--                        <img v-bind:src="message.sender.avatar" alt="" class="avatar">-->
<!--                      </div>-->

                      <div class="received-msg">
                        <div class="received-msg-inbox">
                          <p><span>{{ message.created_by.nickname }}</span><br>{{ message.message_text }}</p>
                        </div>
                      </div>
                    </div>


                    <div class="outgoing-chats" v-else>
                      <div class="outgoing-chats-msg">
                        <p>{{ message.message_text }}</p>
                      </div>

<!--                      <div class="outgoing-chats-img">-->
<!--                        <img v-bind:src="message.sender.avatar" alt="" class="avatar">-->
<!--                      </div>-->
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="msg-bottom">
            <form class="message-form" v-on:submit.prevent="postMessage">
              <div class="input-group">
                <input type="text" class="form-control message-input" placeholder="Type something" v-model="message_text"
                       required>
                <spinner
                    v-if="sendingMessage"
                    class="sending-message-spinner"
                    :size="30"
                />
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import configureEventSource from '../eventsource'

  import NavBar from "../components/NavBar.vue";
  import Spinner from "../components/Spinner.vue";
  import {mapGetters} from 'vuex'

  export default {
    name: "home",
    components: {
      NavBar,
      Spinner
    },
    data() {
      return {
        selected_room: '2ef9ad51-e369-4dcd-b239-357340633d66', // room id static for now.
        room_details: {
          id: null,
          room_name: null,
          created_by: {},
          users: []
        },
        is_member: false,
        message_text: null,
        sendingMessage: false,
        groupMessages: [],
        loadingMessages: false
      };
    },
    computed: {
      ...mapGetters(['auth'])
    },
    mounted() {
      this.getRoomDetails()

      // connect to sse endpoint when authenticated
      if (this.auth.accessToken) {
        configureEventSource('//localhost:9050/sse/connect');
      }
    },

    created() {
      if (!this.auth.accessToken || !this.auth.user) {
        this.$router.push('/login')
      }

      // this.getLoggedInUser();
    },
    methods: {
      async getRoomDetails() {
        try {
          const {data} = await this.$http.get(`rooms/v1/${this.selected_room}`, {
            headers: {
              'Authorization': this.auth.accessToken
            }
          });
          this.room_details = data.room_details
          this.is_member = data.is_member

          if (data.is_member) {
            this.fetchRoomMessages()
          }
        } catch (e) {
          console.error(e)
        }
      },

      async joinRoom() {
        if (this.is_member) return;
        try {
          const {data} = await this.$http.put(`rooms/v1/${this.selected_room}/join`, null, {
            headers: {
              'Authorization': this.auth.accessToken
            }
          });

          console.log('JOIN', data)

          this.fetchRoomMessages()
        } catch (e) {
          console.error(e)
        }
      },

      async fetchRoomMessages() {
        if (!this.room_details.id || !this.is_member) {
          return
        }

        this.loadingMessages = true;
        try {
          const {data} = await this.$http.get(`rooms/v1/${this.selected_room}/messages`, {
            headers: {
              'Authorization': this.auth.accessToken
            }
          });
          this.groupMessages = data.messages.reverse();
          this.$nextTick(() => {
            this.scrollToBottom()
          });
        } catch (e) {
          console.error(e)
        }

        this.loadingMessages = false;
      },

      async postMessage() {
        if (!this.message_text) return;

        this.sendingMessage = true;

        const messagePayload = {
          message_text: this.message_text
        };

        try {
          const {data} = await this.$http.post(`rooms/v1/${this.selected_room}/messages`, messagePayload, {
            headers: {
              'Authorization': this.auth.accessToken
            }
          });
          this.groupMessages = [...this.groupMessages, data.message];
          this.$nextTick(() => {
            this.scrollToBottom()
          });

          this.message_text = null;
        } catch (e) {
          console.error(e)
        }

        this.sendingMessage = false
      },

      scrollToBottom() {
        const chat = document.getElementById("msg-page");
        chat.scrollTo(0, chat.scrollHeight + 30);
      },
    }
  };
</script>
