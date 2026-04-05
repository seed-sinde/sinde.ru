export const useInterfacePreferencesSync = () => {
  const { user, updateProfile } = useAuth()
  const { localeCode, themePreference, setInterfaceLocale, setThemePreference, readThemePreferenceFromSettings } =
    useInterfacePreferences()
  const accountLocale = computed<InterfaceLocaleCode>(() =>
    normalizeInterfaceLocale(user.value?.locale || localeCode.value)
  )
  const accountTheme = computed<ThemePreference>(() =>
    user.value ? readThemePreferenceFromSettings(user.value.settings) : themePreference.value
  )
  const saveInterfacePreferences = async (input: { locale?: unknown; theme?: unknown }) => {
    const previousLocale = localeCode.value
    const previousTheme = themePreference.value
    const nextLocale = normalizeInterfaceLocale(input.locale ?? accountLocale.value)
    const nextTheme = normalizeThemePreference(input.theme ?? accountTheme.value)
    setInterfaceLocale(nextLocale)
    setThemePreference(nextTheme)
    if (!user.value) {
      return {
        locale: nextLocale,
        theme: nextTheme
      }
    }
    try {
      await updateProfile({
        locale: interfaceLocaleToTag(nextLocale),
        settings: {
          appearance: {
            theme_preference: nextTheme
          }
        }
      })
      return {
        locale: nextLocale,
        theme: nextTheme
      }
    } catch (err) {
      setInterfaceLocale(previousLocale)
      setThemePreference(previousTheme)
      throw err
    }
  }
  return {
    accountLocale,
    accountTheme,
    saveInterfacePreferences
  }
}
