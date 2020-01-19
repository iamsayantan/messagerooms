<template>
  <div class="theme--dark py-5 darken-1">
    <div>
      <v-btn icon large flat slot="activator">
        <v-avatar>
          <img src="../../static/avatar/man_4.jpg" alt="Micahel Wang">
        </v-avatar>
      </v-btn>
    </div>
    <v-list class="mini-menu">
      <template v-for="item in items">
        <!-- Top level -->
        <v-list-tile :to="item.to" :key="item.icon" class="py-2 mini-tile my-2" avatar>
          <v-icon :color="item.iconColor" class="mini-icon" size="36">{{ item.icon }}</v-icon>
        </v-list-tile>
      </template>
      <v-list-tile :key="'new-room'" @click="dialog = true" class="py-2 mini-tile my-2" style="cursor: pointer" avatar>
        <v-icon class="mini-icon" size="36">add</v-icon>
      </v-list-tile>
    </v-list>

    <v-dialog v-model="dialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Create Room</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-text-field v-model="room_create.room_name" label="Room name" required></v-text-field>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" flat @click="dialog = false">Close</v-btn>
          <v-btn color="blue darken-1" flat @click="createRoom">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</div>
</template>

<script>
export default {
  props: {
    items: {
      type: Array,
    }
  },

  data() {
    return {
      dialog: false,
      room_create: {
        room_name: null
      }
    }
  },

  mounted() {
    console.log(this.items)
  },

  methods: {
    async createRoom() {
      if (!this.room_create.room_name) return
      try {
        const { data } = await this.$axios.post(`/rooms/v1/create`, this.room_create);
        this.$store.commit('appendRoom', data.room)
        this.dialog = false
      } catch (e) {
        console.error(e)
      }
    }
  }

};
</script>

<style lang="stylus">
  .mini-tile
    a.list__tile--active
      background:hsla(0,0%,100%,.12)

</style>
