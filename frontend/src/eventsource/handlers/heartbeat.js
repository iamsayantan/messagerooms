const heartbeat = {
  eventType: 'Heartbeat',
  handle: function (event) {
    // const parsed = JSON.parse(event.data)
    // console.log("[Heartbeat] Parsed: ", parsed)
    console.log('[Heartbeat] Data', event.data)
  }
}

export default heartbeat