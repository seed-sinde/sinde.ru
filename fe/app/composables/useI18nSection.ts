import {useI18nStore} from "~/stores/i18n"

const formatI18nMessage = (template: string, params?: Record<string, string | number>) =>
  !params ? template : template.replace(/\{([a-z0-9_]+)\}/gi, (_, key: string) => String(params[key] ?? ""))

export const useI18nSection = (namespace: I18nNamespace | string) => {
  const i18n = useI18nStore()
  const {localeCode, showTranslationKeys} = useInterfacePreferences()
  const locale = computed(() => localeCode.value)
  const key = computed(() => `i18n:${locale.value}:${namespace}`)
  const load = () => i18n.loadSection(locale.value, namespace)
  const t = (messageKey: string, params?: Record<string, string | number>) => {
    if (showTranslationKeys.value) return messageKey
    return formatI18nMessage(i18n.t(locale.value, namespace, messageKey), params)
  }
  return {locale, key, load, t}
}
