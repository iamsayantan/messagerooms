import {is} from "vee-validate/dist/rules.esm";

export const state = () => ({
  drawer: true,
  eventsource: {
    connection_id: null
  },
  rooms: [],
  selected_room: null,
  selected_room_details: {
    room: {},
    is_member: false
  },
  room_messages: []
})

export const mutations = {
  toggleDrawer(state) {
    state.drawer = !state.drawer
  },

  drawer(state, val) {
    state.drawer = val
  },

  storeEventsourceConnection(state, connID) {
    state.eventsource.connection_id = connID
  },

  storeRooms(state, rooms) {
    state.rooms = rooms
  },

  selectRoom(state, roomId) {
    state.selected_room = roomId
  },

  storeRoomDetails(state, {room, is_member}) {
    state.selected_room_details.room = room
    state.selected_room_details.is_member = is_member
  },

  appendRoom(state, room) {
    state.rooms.push(room)
  },

  storeMessages(state, messages) {
    state.room_messages = messages
  },

  appendMessage(state, message) {
    state.room_messages = [...state.room_messages, message]
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
  },
  selected_room_details(state) {
    return state.selected_room_details
  },
}
