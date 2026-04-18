let removeThemeMediaListener: (() => void) | null = null
const themeColorFor = (theme: "dark" | "light") => (theme === "light" ? "#f4f4f5" : "#0c0c0f")
const toPublicAssetUrl = (baseURL: string, path: string) => {
  const normalizedBaseURL = String(baseURL || "")
    .trim()
    .replace(/\/+$/, "")
  const normalizedPath = String(path || "")
    .trim()
    .replace(/^\/+/, "")
  if (!normalizedPath) return normalizedBaseURL || "/"
  return normalizedBaseURL ? `${normalizedBaseURL}/${normalizedPath}` : `/${normalizedPath}`
}
const isRecord = (value: unknown): value is Record<string, unknown> =>
  Boolean(value) && typeof value === "object" && !Array.isArray(value)
const readThemePreferenceFromSettings = (settings: unknown): ThemePreference => {
  if (!isRecord(settings)) return "system"
  const appearance = isRecord(settings.appearance) ? settings.appearance : null
  return normalizeThemePreference(appearance?.theme_preference)
}
export const useInterfacePreferences = () => {
  const runtimeConfig = useRuntimeConfig()
  const store = useUiPreferencesStore()
  const showTranslationKeys = useState('interface-show-translation-keys', () => false)
  const systemTheme = useState<"dark" | "light">("interface-system-theme", () => "dark")
  const syncSystemTheme = () => {
    if (!import.meta.client) return
    const media = window.matchMedia("(prefers-color-scheme: light)")
    systemTheme.value = media.matches ? "light" : "dark"
    if (removeThemeMediaListener) return
    const onChange = (event: MediaQueryListEvent) => {
      systemTheme.value = event.matches ? "light" : "dark"
    }
    media.addEventListener("change", onChange)
    removeThemeMediaListener = () => {
      media.removeEventListener("change", onChange)
      removeThemeMediaListener = null
    }
  }
  if (import.meta.client) {
    syncSystemTheme()
  }
  const localeCode = computed<InterfaceLocaleCode>(() => normalizeInterfaceLocale(store.interfaceLocale))
  const localeTag = computed(() => interfaceLocaleToTag(localeCode.value))
  const themePreference = computed<ThemePreference>(() => normalizeThemePreference(store.themePreference))
  const effectiveTheme = computed<"dark" | "light">(() =>
    themePreference.value === "system" ? systemTheme.value : themePreference.value
  )
  const publicAsset = (path: string) => toPublicAssetUrl(String(runtimeConfig.public.baseURL || ""), path)
  const faviconLightSrc = computed(() => publicAsset("/favicon-light.svg"))
  const faviconDarkSrc = computed(() => publicAsset("/favicon-dark.svg"))
  const themeColor = computed(() => themeColorFor(effectiveTheme.value))
  const faviconSrc = computed(() => (effectiveTheme.value === "light" ? faviconLightSrc.value : faviconDarkSrc.value))
  const localeOptions = computed<SelectOptionInput[]>(() =>
    INTERFACE_LOCALE_OPTIONS.map(option => ({
      value: option.code,
      label: option.nativeLabel
    }))
  )
  const themeOptions = computed<SelectOptionInput[]>(() =>
    INTERFACE_THEME_OPTIONS.map(value => ({
      value,
      label: value
    }))
  )
  const setInterfaceLocale = (value: unknown) => {
    store.setInterfaceLocale(normalizeInterfaceLocale(value))
  }
  const setThemePreference = (value: unknown) => {
    store.setThemePreference(normalizeThemePreference(value))
  }
  const applyAccountPreferences = (user?: AuthUser | null) => {
    if (!user) return
    setInterfaceLocale(user.locale)
    setThemePreference(readThemePreferenceFromSettings(user.settings))
  }
  return {
    localeCode,
    localeTag,
    themePreference,
    effectiveTheme,
    themeColor,
    publicAsset,
    faviconSrc,
    faviconLightSrc,
    faviconDarkSrc,
    localeOptions,
    themeOptions,
    showTranslationKeys,
    toggleTranslationKeys: () => {
      showTranslationKeys.value = !showTranslationKeys.value
    },
    setInterfaceLocale,
    setThemePreference,
    applyAccountPreferences,
    readThemePreferenceFromSettings
  }
}
