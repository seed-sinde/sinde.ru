const AUTH_SYNC_COOLDOWN_MS = 3000
const AUTH_SYNC_RUNTIME_KEY = '__traitsAuthSyncRuntime__'
type AuthSyncRuntime = {
  cleanup: (() => void) | null
}
const authSyncRuntime = (() => {
  const scope = globalThis as typeof globalThis & {
    [AUTH_SYNC_RUNTIME_KEY]?: AuthSyncRuntime
  }
  if (!scope[AUTH_SYNC_RUNTIME_KEY]) {
    scope[AUTH_SYNC_RUNTIME_KEY] = { cleanup: null }
  }
  return scope[AUTH_SYNC_RUNTIME_KEY]!
})()
export default defineNuxtPlugin(nuxtApp => {
  if (!import.meta.client) return
  authSyncRuntime.cleanup?.()
  const { loadMe, loaded, refreshSharedSummaries, ensureSummaryRealtime } =
    useAuth()
  let inFlight = false
  let syncReady = loaded.value
  let lastSyncAt = Date.now()
  const markSyncReady = () => {
    syncReady = true
    lastSyncAt = Date.now()
  }
  const stopLoadedWatch = watch(
    loaded,
    next => {
      if (!next) return
      markSyncReady()
    },
    { immediate: true }
  )
  ensureSummaryRealtime()
  const syncAuthState = async (force = false) => {
    if (!syncReady && !force) return
    await nuxtApp.runWithContext(async () => {
      const now = Date.now()
      if (!force && now - lastSyncAt < AUTH_SYNC_COOLDOWN_MS) return
      if (inFlight) return
      const { canAttemptSessionRestore, isAuthenticated, setAnonymousState } = useAuth()
      if (!canAttemptSessionRestore.value && !isAuthenticated.value) {
        setAnonymousState()
        markSyncReady()
        return
      }
      inFlight = true
      try {
        await loadMe()
        await refreshSharedSummaries()
      } finally {
        inFlight = false
        markSyncReady()
      }
    })
  }
  const onVisibilityChange = () => {
    if (document.visibilityState !== 'visible') return
    void syncAuthState()
  }
  const onAuthStateStale = () => {
    void syncAuthState(true)
  }
  const onAuthUserRefreshed = (event: Event) => {
    const customEvent = event as CustomEvent<{ user?: AuthUser | null }>
    void nuxtApp.runWithContext(async () => {
      const payloadUser = customEvent.detail?.user ?? null
      const { user, loaded } = useAuth()
      user.value = payloadUser
      loaded.value = true
      markSyncReady()
      await refreshSharedSummaries()
    })
  }
  const onFocus = () => {
    void syncAuthState()
  }
  const onPageFinish = () => {
    void syncAuthState()
  }
  window.addEventListener('focus', onFocus)
  window.addEventListener('auth-state-stale', onAuthStateStale)
  window.addEventListener('auth-user-refreshed', onAuthUserRefreshed as EventListener)
  document.addEventListener('visibilitychange', onVisibilityChange)
  nuxtApp.hook('page:finish', onPageFinish)
  const cleanup = () => {
    stopLoadedWatch()
    window.removeEventListener('focus', onFocus)
    window.removeEventListener('auth-state-stale', onAuthStateStale)
    window.removeEventListener('auth-user-refreshed', onAuthUserRefreshed as EventListener)
    document.removeEventListener('visibilitychange', onVisibilityChange)
    nuxtApp.hooks.removeHook('page:finish', onPageFinish)
    if (authSyncRuntime.cleanup === cleanup) {
      authSyncRuntime.cleanup = null
    }
  }
  authSyncRuntime.cleanup = cleanup
  if (import.meta.hot) {
    import.meta.hot.dispose(() => {
      cleanup()
    })
  }
})
