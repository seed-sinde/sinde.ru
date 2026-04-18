const SIDEBAR_TRANSLATION_KEYS = {
  "/lab": "sidebar.lab",
  "/traits": "sidebar.traits",
  "/kitchen/recipes": "sidebar.kitchen",
  "/edu": "sidebar.wiki"
} as const
export const useSidebarItems = () => {
  const {load, t} = useI18nSection("ui")
  onServerPrefetch(load)
  if (import.meta.client) void load()
  return computed<MenuItem[]>(() =>
    sidebarItems.map(item => ({
      ...item,
      label: t(SIDEBAR_TRANSLATION_KEYS[item.to as keyof typeof SIDEBAR_TRANSLATION_KEYS] || "sidebar.wiki")
    }))
  )
}
