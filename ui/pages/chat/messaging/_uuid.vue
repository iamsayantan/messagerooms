<template>
  <v-container class="fill-height pa-0 ma-0 messaging fluid" id="messaging">
    <template v-if="!$vuetify.breakpoint.smAndDown">
      <v-layout row>
        <v-flex lg3 class="white">
          <chat-history></chat-history>
        </v-flex>
        <v-flex lg9>
          <chat-window v-if="$route.params.uuid" @roomJoin="selectAndFetchRoom"></chat-window>
        </v-flex>
      </v-layout>
    </template>
    <template v-else>
      <v-layout column>
        <v-flex sm12 class="white" v-if="showSidebar">
          <chat-history>
          </chat-history>
        </v-flex>
        <v-flex sm12 v-if="showWindow">
          <chat-window v-if="$route.params.uuid"></chat-window>
        </v-flex>
      </v-layout>
    </template>
  </v-container>
</template>
<script>
  import {mapGetters} from 'vuex'
  import ChatHistory from '../../../components/chat/ChatHistory';
  import ChatWindow from '../../../components/chat/ChatWindow';
  export default {
    components: {
      ChatHistory,
      ChatWindow
    },
    mounted() {
      this.fetchMessageRooms()
      if (this.$route.params.uuid) {
        this.selectAndFetchRoom(this.$route.params.uuid)
      }
    },
    data () {
      return {
      };
    },
    computed: {
      ...mapGetters(['selected_room']),
      showSidebar () {
        return this.$route.params.uuid === undefined;
      },
      showWindow () {
        return this.$route.params.uuid !== undefined;
      },
    },
    watch: {
      '$route.params.uuid': function (uuid) {
        if (!uuid) return
        this.selectAndFetchRoom(uuid)
      }
    },
    methods: {
      async fetchMessageRooms() {
        try {
          const { data } = await this.$axios.get('/rooms/v1')
          this.$store.commit('storeRooms', data.rooms)
        } catch (e) {
          console.error(e)
        }
      },

      async selectAndFetchRoom(roomID) {
        this.$store.commit('selectRoom', roomID)
        try {
          const { data } = await this.$axios.get(`/rooms/v1/${this.selected_room}`)
          this.$store.commit('storeRoomDetails', {room: data.room_details, is_member: data.is_member})
          console.log(data.room_details)
        } catch (e) {
          console.error(e)
        }
      }
    }
  };
</script>
