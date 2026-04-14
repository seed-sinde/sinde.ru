import { computed, type ComputedRef, type Ref } from 'vue'
import type { InterfaceLocaleCode } from '../types/ui'
export const resolveLocalizedPage = <T>(
  pages: Record<InterfaceLocaleCode, T>,
  localeCode: Ref<InterfaceLocaleCode>
): ComputedRef<T> => {
  return computed(() => pages[localeCode.value] ?? pages.ru)
}
