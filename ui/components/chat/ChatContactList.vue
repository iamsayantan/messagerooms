<template>
  <div class="chat-contact">
    <v-toolbar flat dense class="chat-contact--toolbar">
      <v-text-field flat solo full-width prepend-icon="search" label="Search"></v-text-field>
    </v-toolbar>
    <vue-perfect-scrollbar class="chat-history--scrollbar">
      <v-divider></v-divider>
      <v-list two-line class="chat-contact--list">
        <v-subheader>Contacts</v-subheader>
        <template v-for="(room, index) in rooms">
          <v-divider :key="index"></v-divider>
          <v-list-tile avatar :key="room.room_name + index" :to="contactRoute(room.id)">
            <v-list-tile-avatar color="primary">
              <span class="white--text headline">{{ firstLetter(room.room_name)}}</span>
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title>
                {{room.room_name}}
              </v-list-tile-title>
              <v-list-tile-sub-title>This is an awesome room</v-list-tile-sub-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-circle dot small :color="userStatusColor(room)"></v-circle>
            </v-list-tile-action>
          </v-list-tile>
        </template>
      </v-list>
    </vue-perfect-scrollbar>
  </div>
</template>

<script>
  import {getUser} from '@/api/user';
  import VCircle from '@/components/circle/VCircle';
  import VuePerfectScrollbar from 'vue-perfect-scrollbar';

  import {mapGetters} from 'vuex'

  export default {
    components: {
      VuePerfectScrollbar,
      VCircle
    },
    mounted() {
      this.fetchMessageRooms()
    },
    data: () => ({}),
    computed: {
      ...mapGetters(['rooms']),
      users() {
        return getUser();
      }
    },
    methods: {
      async fetchMessageRooms() {
        try {
          const { data } = this.$axios.get('/api/rooms/v1')
          console.log(data)
        } catch (e) {
          console.error(e)
        }
      },
      contactRoute(id) {
        return '/chat/contact/' + id;
      },
      firstLetter(name) {
        return name.charAt(0);
      },
      userStatusColor(item) {
        return ((item.room_name % 2) === 0) ? 'green' : 'grey';
      }
    }
  };
</script>

<style>

</style>
