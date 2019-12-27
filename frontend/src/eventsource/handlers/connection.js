const connection = {
  eventType: 'ClientConnection',
  handle: function (event) {
    // const parsed = JSON.parse(event.data)
    // console.log("[ClientConnection] Parsed: ", parsed)
    console.log('[ClientConnection] Data', event.data)
  }
}

export default connection