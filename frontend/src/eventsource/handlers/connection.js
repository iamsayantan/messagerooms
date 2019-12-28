const connection = {
  eventType: 'ClientConnection',
  handle: function (event) {
    const parsed = JSON.parse(event.data)
    console.log('[ClientConnection] Data', parsed)
  }
}

export default connection