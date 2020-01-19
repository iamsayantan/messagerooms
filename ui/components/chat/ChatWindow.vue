<template>
  <v-card class="chat-room">
    <v-toolbar card dense flat class="white chat-room--toolbar" light>
      <v-btn icon>
        <v-icon color="text--secondary">keyboard_arrow_left</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <v-toolbar-title> <h4>{{ '#' + selected_room_details.room.room_name }}</h4></v-toolbar-title>
      <v-spacer></v-spacer>
      <v-tooltip v-if="!selected_room_details.is_member" left>
        <v-btn flat slot="activator" @click="joinRoom">
          Join Room
        </v-btn>
        <span>Click to join this room</span>
      </v-tooltip>
    </v-toolbar>
    <template v-if="selected_room_details.is_member">
      <vue-perfect-scrollbar id="chat-messages" class="chat-room--scrollbar grey lighten-5" v-bind:style="computeHeight">
        <v-card-text class="chat-room--list pa-3">
          <div v-chat-scroll="{always: true, smooth: true}">
            <template v-for="(message, index) in room_messages">
              <div v-bind:class="[ message.created_by.id === $auth.user.id ? 'reverse' : '']" class="messaging-item layout row my-4" :key="index">
                <v-avatar class="indigo mx-1" size="40">
                  <span class="white--text headline">{{message.created_by.nickname[0]}}</span>
                </v-avatar>
                <div class="messaging--body layout column mx-2">
                  <p :value="true" v-bind:class="[ message.created_by.id === $auth.user.id ? 'primary white--text' : 'white']" class="pa-2">
                    {{message.message_text}}
                  </p>
                  <div class="caption px-2 text--secondary" v-if="message.created_by.id !== $auth.user.id">{{message.created_by.nickname}} ({{new Date(message.created_at).toDateString()}})</div>
                  <div class="caption px-2 text--secondary"  v-else>You ({{new Date(message.created_at).toDateString()}})</div>
                </div>
                <v-spacer></v-spacer>
              </div>
            </template>
          </div>
        </v-card-text>
      </vue-perfect-scrollbar>
      <v-card-actions>
        <v-text-field
          v-model="message_text"
          full-width
          flat
          clearable
          solo
          :loading="sending"
          @click:append="postMessage"
          @keydown.enter="postMessage"
          append-icon="send"
          label="Type some message here">
          <v-icon slot="append-icon">send</v-icon>
          <v-icon slot="append-icon" class="mx-2">photo</v-icon>
          <v-icon slot="append-icon">face</v-icon>
        </v-text-field>
      </v-card-actions>
    </template>
  </v-card>
</template>
<script>
import {mapGetters} from 'vuex';
import { getChatById } from '@/api/chat';
import { getUserById } from '@/api/user';
import VuePerfectScrollbar from 'vue-perfect-scrollbar';
export default {
  components: {
    VuePerfectScrollbar
  },
  props: {
    uuid: {
      type: String,
      default: '',
    },
    height: {
      type: String,
      default: null,
    }
  },
  data() {
    return {
      message_text: null,
      sending: false
    }
  },
  computed: {
    ...mapGetters(['selected_room_details', 'room_messages']),
    computeHeight () {
      return {
        height: this.height || ''
      };
    }
  },

  watch: {
    selected_room_details: {
      handler(roomDetails) {
        this.fetchRoomMessages()
      },
      deep: true
    },
    room_messages: {
      handler(val) {
        // every time a new message is appended to the message list, we manually trigger
        // scrollToBottom(). Otherwise the scroll was not happening.
        this.$nextTick(() => {
          this.scrollToBottom()
        })
      }
    }
  },
  methods: {
    async fetchRoomMessages() {
      if (!this.selected_room_details.is_member) return;
      try {
        const {data} = await this.$axios.get(`/rooms/v1/${this.selected_room_details.room.id}/messages`)
        this.$store.commit('storeMessages', data.messages.reverse())
        this.$nextTick(() => {
          this.scrollToBottom()
        })
      } catch (e) {
        console.error(e)
      }
    },
    async postMessage() {
      if (!this.selected_room_details.is_member || !this.message_text) return;
      this.sending = true
      try {
        // we are not updating the message list from the message post response. the message will come back from
        // realtime connection. where it will be appended in the list.
        await this.$axios.post(`/rooms/v1/${this.selected_room_details.room.id}/messages`, {
          message_text: this.message_text
        });

        this.message_text = null
      } catch (e) {
        console.error(e)
      }
      this.sending = false
    },

    async joinRoom() {
      if (this.selected_room_details.is_member) return
      try {
        const { data } = await this.$axios.put(`/rooms/v1/${this.selected_room_details.room.id}/join`)
        this.$emit('roomJoin', this.selected_room_details.room.id)
      } catch (e) {
        console.error(e)
      }
    },

    getAvatar (uid) {
      return getUserById(uid).avatar;
    },
    scrollToBottom() {
      const chat = document.getElementById("chat-messages");
      chat.scrollTo(0, chat.scrollHeight);
    },
  }
};
</script>

