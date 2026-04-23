type InterfaceLocaleCode = "ru" | "en" | "zh" | "ja"

const formatI18nMessage = (template: string, params?: Record<string, string | number>) =>
  !params ? template : template.replace(/\{([a-z0-9_]+)\}/gi, (_, key: string) => String(params[key] ?? ""))

export const useI18nSection = (namespace: string) => {
  const i18n = useI18nStore()
  const locale = useState<InterfaceLocaleCode>("i18n-locale", () => "ru")
  const showTranslationKeys = useState("i18n-show-translation-keys", () => false)
  const key = computed(() => `i18n:${locale.value}:${namespace}`)
  const load = () => i18n.loadSection(locale.value, namespace)
  const t = (messageKey: string, params?: Record<string, string | number>) =>
    showTranslationKeys.value ? messageKey : formatI18nMessage(i18n.t(locale.value, namespace, messageKey), params)

  return {locale, key, load, t}
}
