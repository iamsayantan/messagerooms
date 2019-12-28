const heartbeat = {
  eventType: 'Heartbeat',
  handle: function (event) {
    const parsed = JSON.parse(event.data)
    console.log('[Heartbeat] Data', parsed)
  }
}

export default heartbeat