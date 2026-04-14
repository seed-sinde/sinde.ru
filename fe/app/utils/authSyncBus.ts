type AuthSyncEventType = 'summary-refresh'
type AuthSyncEvent = {
  type: AuthSyncEventType
}
const AUTH_SYNC_CHANNEL_NAME = 'sinde-auth-sync'
const AUTH_SYNC_WINDOW_EVENT = 'sinde:auth-sync'
let authSyncChannel: BroadcastChannel | null | undefined
function getAuthSyncChannel() {
  if (!import.meta.client) return null
  if (authSyncChannel !== undefined) return authSyncChannel
  authSyncChannel = typeof BroadcastChannel === 'function' ? new BroadcastChannel(AUTH_SYNC_CHANNEL_NAME) : null
  return authSyncChannel
}
export function emitAuthSyncEvent(type: AuthSyncEventType) {
  if (!import.meta.client) return
  const payload: AuthSyncEvent = { type }
  getAuthSyncChannel()?.postMessage(payload)
  window.dispatchEvent(new CustomEvent<AuthSyncEvent>(AUTH_SYNC_WINDOW_EVENT, { detail: payload }))
}
export function subscribeAuthSyncEvents(handler: (event: AuthSyncEvent) => void) {
  if (!import.meta.client) {
    return () => {}
  }
  const channel = getAuthSyncChannel()
  const onMessage = (event: MessageEvent<AuthSyncEvent>) => {
    if (event.data?.type) {
      handler(event.data)
    }
  }
  const onWindowEvent = (event: Event) => {
    const customEvent = event as CustomEvent<AuthSyncEvent>
    if (customEvent.detail?.type) {
      handler(customEvent.detail)
    }
  }
  channel?.addEventListener('message', onMessage)
  window.addEventListener(AUTH_SYNC_WINDOW_EVENT, onWindowEvent as EventListener)
  return () => {
    channel?.removeEventListener('message', onMessage)
    window.removeEventListener(AUTH_SYNC_WINDOW_EVENT, onWindowEvent as EventListener)
  }
}
