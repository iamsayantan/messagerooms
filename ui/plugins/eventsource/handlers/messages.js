const message = {
  eventType: 'MessageEvent',
  handle: function (event, store) {
    const eventData = JSON.parse(event.data)

    // match with the personal topics and take appropriate action.
    // topics are appended with user's id with it
    if (eventData.topic.startsWith('NewMessage')) {
      const { message, room } = eventData.payload

      if (store.getters.selected_room === room.id) {
        store.commit('appendMessage', message)
      }
    }
  }
}

export default message
