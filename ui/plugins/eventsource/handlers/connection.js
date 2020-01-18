const connection = {
  eventType: 'ClientConnection',
  handle: function (event, store) {
    const parsed = JSON.parse(event.data)
    store.commit('storeEventsourceConnection', parsed.connection_id)
  }
}

export default connection
