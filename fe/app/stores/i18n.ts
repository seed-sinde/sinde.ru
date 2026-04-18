import {defineStore} from "pinia"

export const useI18nStore = defineStore("i18n", () => {
  const api = useAPI()
  const sections = ref<Record<string, I18nMessages>>({})
  const pending = ref<Record<string, Promise<I18nMessages>>>({})

  const sectionKey = (locale: InterfaceLocaleCode, namespace: I18nNamespace | string) => `${locale}:${namespace}`

  const getSection = (locale: InterfaceLocaleCode, namespace: I18nNamespace | string) =>
    sections.value[sectionKey(locale, namespace)] || null

  const loadSection = async (locale: InterfaceLocaleCode, namespace: I18nNamespace | string) => {
    const key = sectionKey(locale, namespace)
    if (sections.value[key]) return sections.value[key]
    if (pending.value[key]) return pending.value[key]
    const request = api
      .json<ApiResponseWithData<I18nMessages>>(`/i18n/${locale}/${namespace}`)
      .then(res => {
        const data = res?.data || {}
        sections.value[key] = data
        return data
      })
      .finally(() => {
        const {[key]: _pending, ...rest} = pending.value
        pending.value = rest
      })
    pending.value[key] = request
    return await request
  }

  const t = (locale: InterfaceLocaleCode, namespace: I18nNamespace | string, key: string) =>
    sections.value[sectionKey(locale, namespace)]?.[key] || key

  const clear = () => (sections.value = {})

  return {sections, getSection, loadSection, t, clear}
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useI18nStore, import.meta.hot))
}
