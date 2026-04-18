const DEFAULT_LIMIT = 30
const STORAGE_KEY = "ui.preferences.v1"
const createDefaultMineralsFilters = (): MineralsFiltersSnapshot => ({
  q: "",
  sort: "name_asc",
  limit: DEFAULT_LIMIT,
  offset: 0,
  imageFilter: "any",
  crystalSystems: [],
  crystalSystemMode: "any",
  chemistryAll: [],
  chemistryAny: [],
  chemistryNone: []
})
export const useUiPreferencesStore = defineStore("uiPreferences", () => {
  const interfaceLocale = ref<InterfaceLocaleCode>("ru")
  const themePreference = ref<ThemePreference>("system")
  const mineralsFilters = reactive<MineralsFiltersSnapshot>(createDefaultMineralsFilters())
  const restored = ref(false)
  const hasActiveMineralsFilters = () =>
    Boolean(
      mineralsFilters.q.trim() ||
      mineralsFilters.sort !== "name_asc" ||
      mineralsFilters.limit !== DEFAULT_LIMIT ||
      mineralsFilters.offset > 0 ||
      mineralsFilters.imageFilter !== "any" ||
      mineralsFilters.crystalSystems.length ||
      mineralsFilters.crystalSystemMode !== "any" ||
      mineralsFilters.chemistryAll.length ||
      mineralsFilters.chemistryAny.length ||
      mineralsFilters.chemistryNone.length
    )
  const mineralsFiltersSnapshot = (): MineralsFiltersSnapshot => ({
    q: mineralsFilters.q,
    sort: mineralsFilters.sort,
    limit: mineralsFilters.limit,
    offset: mineralsFilters.offset,
    imageFilter: mineralsFilters.imageFilter,
    crystalSystems: mineralsFilters.crystalSystems.slice(),
    crystalSystemMode: mineralsFilters.crystalSystemMode,
    chemistryAll: mineralsFilters.chemistryAll.slice(),
    chemistryAny: mineralsFilters.chemistryAny.slice(),
    chemistryNone: mineralsFilters.chemistryNone.slice()
  })
  const replaceMineralsFilters = (next: Partial<MineralsFiltersSnapshot>) => {
    mineralsFilters.q = String(next.q || "").trim()
    mineralsFilters.sort = next.sort === "name_desc" ? "name_desc" : "name_asc"
    mineralsFilters.limit =
      typeof next.limit === "number" && Number.isFinite(next.limit) && next.limit > 0 ? next.limit : DEFAULT_LIMIT
    mineralsFilters.offset =
      typeof next.offset === "number" && Number.isFinite(next.offset) && next.offset >= 0 ? next.offset : 0
    mineralsFilters.imageFilter =
      next.imageFilter === "with" || next.imageFilter === "without" ? next.imageFilter : "any"
    mineralsFilters.crystalSystems = Array.isArray(next.crystalSystems) ? next.crystalSystems.slice() : []
    mineralsFilters.crystalSystemMode = next.crystalSystemMode === "all" ? "all" : "any"
    mineralsFilters.chemistryAll = Array.isArray(next.chemistryAll) ? next.chemistryAll.slice() : []
    mineralsFilters.chemistryAny = Array.isArray(next.chemistryAny) ? next.chemistryAny.slice() : []
    mineralsFilters.chemistryNone = Array.isArray(next.chemistryNone) ? next.chemistryNone.slice() : []
  }
  const clearMineralsFilters = () => {
    replaceMineralsFilters(createDefaultMineralsFilters())
  }
  const snapshot = (): UiPreferencesSnapshot => ({
    interfaceLocale: interfaceLocale.value,
    themePreference: themePreference.value,
    mineralsFilters: mineralsFiltersSnapshot()
  })
  const restorePersisted = () => {
    if (!import.meta.client || restored.value) return
    try {
      const raw = window.localStorage.getItem(STORAGE_KEY)
      if (!raw) {
        restored.value = true
        return
      }
      const parsed = JSON.parse(raw) as Partial<UiPreferencesSnapshot>
      interfaceLocale.value = normalizeInterfaceLocale(parsed.interfaceLocale)
      themePreference.value = normalizeThemePreference(parsed.themePreference)
      replaceMineralsFilters(parsed.mineralsFilters || {})
    } catch {
      window.localStorage.removeItem(STORAGE_KEY)
      interfaceLocale.value = "ru"
      themePreference.value = "system"
      clearMineralsFilters()
    } finally {
      restored.value = true
    }
  }
  const persistState = () => {
    if (!import.meta.client || !restored.value) return
    try {
      const hasMeaningfulState =
        interfaceLocale.value !== "ru" || themePreference.value !== "system" || hasActiveMineralsFilters()
      if (!hasMeaningfulState) {
        window.localStorage.removeItem(STORAGE_KEY)
        return
      }
      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(snapshot()))
    } catch {
      // Ignore storage issues.
    }
  }
  if (import.meta.client) {
    onNuxtReady(() => {
      restorePersisted()
    })
    watch(
      () => [
        interfaceLocale.value,
        themePreference.value,
        mineralsFilters.q,
        mineralsFilters.sort,
        mineralsFilters.limit,
        mineralsFilters.offset,
        mineralsFilters.imageFilter,
        mineralsFilters.chemistryAll,
        mineralsFilters.chemistryAny,
        mineralsFilters.chemistryNone
      ],
      persistState,
      {deep: true}
    )
  }
  return {
    restored,
    interfaceLocale,
    themePreference,
    mineralsFilters,
    snapshot,
    restorePersisted,
    hasActiveMineralsFilters,
    mineralsFiltersSnapshot,
    replaceMineralsFilters,
    clearMineralsFilters,
    setInterfaceLocale(value: InterfaceLocaleCode) {
      interfaceLocale.value = normalizeInterfaceLocale(value)
    },
    setThemePreference(value: ThemePreference) {
      themePreference.value = normalizeThemePreference(value)
    }
  }
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUiPreferencesStore, import.meta.hot))
}
