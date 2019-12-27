import Handlers from './handlers'
import store from '../store'

import { EventSourcePolyfill } from 'event-source-polyfill'

export default function configureEventSources (eventsourceRoute) {
  const eventSource = new EventSourcePolyfill(eventsourceRoute, {
    headers: {
      'Authorization': store.getters.auth.accessToken
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
      handler.handle(event)
    })
  }
}