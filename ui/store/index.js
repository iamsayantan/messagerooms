export const state = () => ({
  drawer: true,
  rooms: [],
  selected_room: null,
  room_messages: []
})

export const mutations = {
  toggleDrawer(state) {
    state.drawer = !state.drawer
  },
  drawer(state, val) {
    state.drawer = val
  },
  storeRoms(state, rooms) {
    state.rooms = [...state.rooms, ...rooms]
  },
  appendRoom(state, room) {
    state.rooms.push(room)
  }
}

export const getters = {
  rooms(state) {
    return state.rooms
  },
  selected_room(state) {
    return state.selected_room
  },
  room_messages(state) {
    return state.room_messages
  }
}
