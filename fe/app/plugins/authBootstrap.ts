export default defineNuxtPlugin(async () => {
  const interfacePreferences = useInterfacePreferences()
  const uiPreferences = useUiPreferencesStore()
  const {
    user,
    loaded,
    sharedUserSummary,
    sharedAdminSummary,
    ensureLoaded,
    isAuthenticated,
    loadSharedUserSummary
  } = useAuth()
  if (import.meta.server && import.meta.prerender) {
    user.value = null
    loaded.value = true
    sharedUserSummary.value = null
    sharedAdminSummary.value = null
    return
  }
  if (import.meta.client) {
    uiPreferences.restorePersisted()
  }
  await ensureLoaded()
  if (user.value) {
    interfacePreferences.applyAccountPreferences(user.value)
  }
  watch(
    () => user.value,
    nextUser => {
      if (!nextUser) return
      interfacePreferences.applyAccountPreferences(nextUser)
    },
    { deep: true }
  )
  if (!isAuthenticated.value) return
  if (import.meta.server && !sharedUserSummary.value) {
    await loadSharedUserSummary()
  }
})
