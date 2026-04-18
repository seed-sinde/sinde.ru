export const useTraitsMobileHeader = () => {
  const route = useRoute()
  const {load, t} = useI18nSection("traits")
  onServerPrefetch(load)
  if (import.meta.client) void load()
  const {user} = useAuth()

  const normalizePath = (rawPath: string) => {
    const normalizedPath = String(rawPath || "/").replace(/\/{2,}/g, "/")
    return normalizedPath.length <= 1 ? "/" : normalizedPath.replace(/\/+$/, "")
  }

  const getRouteUuid = () => {
    const rawUuid = route.params.uuid
    return typeof rawUuid === "string" ? rawUuid.trim() : Array.isArray(rawUuid) ? String(rawUuid[0] || "").trim() : ""
  }

  const routeState = computed(() => {
    const normalizedPath = normalizePath(route.path)
    const isTraitDetailPage = normalizedPath.startsWith("/traits/trait/")
    const isTraitSetPage = normalizedPath.startsWith("/traits/") && !isTraitDetailPage
    const isTraitsIndexPage = normalizedPath === "/traits"
    const activeTab = normalizeTabRouteValue(route.query.tab, ["traits", "saved"], "traits") as "traits" | "saved"
    const currentUuid = isTraitDetailPage || isTraitSetPage ? getRouteUuid() : ""

    return {
      isTraitDetailPage,
      isTraitSetPage,
      isTraitsIndexPage,
      activeTab,
      traitUuid: isTraitDetailPage ? currentUuid : "",
      setUuid: isTraitSetPage ? currentUuid : ""
    }
  })

  const mobileBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
    const state = routeState.value

    if (state.isTraitDetailPage) return [{label: t("detail.breadcrumb"), to: "/traits"}]

    if (state.isTraitsIndexPage || state.isTraitSetPage) {
      if (state.activeTab === "saved") {
        return [
          {label: t("page.title"), to: "/traits"},
          {label: t("saved.breadcrumb"), current: true, kind: "tab"}
        ]
      }

      const currentSetUuid = state.setUuid
      const primaryTraitUuid = String(user.value?.primary_trait_uuid || "").trim()

      return currentSetUuid && currentSetUuid === primaryTraitUuid
        ? [
            {label: t("page.title"), to: "/traits"},
            {label: "", current: true, kind: "tab", badge: t("primary_badge")}
          ]
        : [{label: t("page.title"), ...(currentSetUuid ? {to: "/traits"} : {current: true})}]
    }

    return []
  })

  const mobileHeaderActions = computed<MobileHeaderAction[]>(() => {
    const state = routeState.value
    return state.isTraitDetailPage && state.traitUuid
      ? [
          {kind: "traits-copy-uuid", uuid: state.traitUuid},
          {kind: "traits-paste-uuid", mode: "trait"}
        ]
      : state.isTraitsIndexPage || state.isTraitSetPage
        ? [
            ...(state.setUuid ? [{kind: "traits-copy-uuid" as const, uuid: state.setUuid}] : []),
            {kind: "traits-paste-uuid", mode: "set"}
          ]
        : []
  })

  usePageMobileHeaderMeta(mobileBreadcrumbItems, mobileHeaderActions)
}
