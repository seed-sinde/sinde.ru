const PRELOAD_RECOVER_TS_KEY = 'nuxt:preload-recover:last-reload-at'
const PRELOAD_RECOVER_COOLDOWN_MS = 10_000
export default defineNuxtPlugin(() => {
  if (!import.meta.client) return
  const onPreloadError = (event: Event) => {
    const custom = event as CustomEvent
    const message = String((custom as any)?.payload?.message || (custom as any)?.error?.message || '')
    if (!message && custom.type !== 'vite:preloadError') return
    if (typeof custom.preventDefault === 'function') custom.preventDefault()
    const now = Date.now()
    let lastReloadAt = 0
    try {
      lastReloadAt = Number(window.sessionStorage.getItem(PRELOAD_RECOVER_TS_KEY) || 0)
    } catch {
      lastReloadAt = 0
    }
    if (now - lastReloadAt < PRELOAD_RECOVER_COOLDOWN_MS) return
    try {
      window.sessionStorage.setItem(PRELOAD_RECOVER_TS_KEY, String(now))
    } catch {
      // ignore sessionStorage access issues
    }
    window.location.reload()
  }
  window.addEventListener('vite:preloadError', onPreloadError)
})
