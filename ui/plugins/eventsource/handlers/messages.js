const message = {
  eventType: 'MessageEvent',
  handle: function (event) {
    const parsed = JSON.parse(event.data)
    console.log('[MessageEvent] Data', parsed)
  }
}

export default message