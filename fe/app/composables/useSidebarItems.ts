import { type InterfaceMessageKey } from '~/data/interfacePreferences'
const SIDEBAR_TRANSLATION_KEYS = {
  '/lab': 'sidebar.lab',
  '/traits': 'sidebar.traits',
  '/kitchen': 'sidebar.kitchen',
  '/edu': 'sidebar.wiki',
  '/impact': 'sidebar.impact',
  '/perception': 'sidebar.perception',
  '/experience': 'sidebar.experience',
  '/careers': 'sidebar.careers',
  '/family': 'sidebar.family'
} as const satisfies Record<string, InterfaceMessageKey>
export const useSidebarItems = () => {
  const { t } = useInterfacePreferences()
  return computed<MenuItem[]>(() =>
    sidebarItems.map(item => ({
      ...item,
      label: t(SIDEBAR_TRANSLATION_KEYS[item.to as keyof typeof SIDEBAR_TRANSLATION_KEYS] || 'sidebar.wiki')
    }))
  )
}
