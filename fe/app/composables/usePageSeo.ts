import type {ComputedRef, Ref} from "vue"

type SeoValue = string | null | undefined
type SeoInput = SeoValue | Ref<SeoValue> | ComputedRef<SeoValue> | (() => SeoValue)
type UsePageSeoInput = {
  title: SeoInput
  description: SeoInput
  ogTitle?: SeoInput
  ogDescription?: SeoInput
}
export const usePageSeo = (input: UsePageSeoInput) => {
  const toGetter = (value: SeoInput) => () => toValue(value) ?? ""
  useSeoMeta({
    title: toGetter(input.title),
    description: toGetter(input.description),
    ogTitle: toGetter(input.ogTitle ?? input.title),
    ogDescription: toGetter(input.ogDescription ?? input.description)
  })
}
