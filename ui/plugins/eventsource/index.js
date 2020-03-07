import Handlers from './handlers'

import { EventSourcePolyfill } from 'event-source-polyfill'

export default ({app, store}) => {
  if (!app.$auth.loggedIn) return

  const eventsourceRoute = '/api/sse/connect'
  const eventSource = new EventSourcePolyfill(eventsourceRoute, {
    headers: {
      'Authorization': app.$auth.getToken('local')
    }
  });

  eventSource.addEventListener('open', (event) => {
    console.log('[Eventsource] Connection Open', event)
  });

  eventSource.addEventListener('close', (event) => {
    console.log('[Eventsource] Connection Close', event)
  });

  for (const handler of Handlers) {
    eventSource.addEventListener(handler.eventType, event => {
      handler.handle(event, store)
    })
  }
}
