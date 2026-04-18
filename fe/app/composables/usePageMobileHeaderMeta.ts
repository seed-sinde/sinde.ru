export const usePageMobileHeaderMeta = (
  breadcrumbs: ComputedRef<BreadcrumbItem[]> | Ref<BreadcrumbItem[]> | BreadcrumbItem[],
  actions: ComputedRef<MobileHeaderAction[]> | Ref<MobileHeaderAction[]> | MobileHeaderAction[]
) => {
  const route = useRoute()

  const read = <T>(v: ComputedRef<T> | Ref<T> | T) => (isRef(v) ? v.value : v)

  watchEffect(() => {
    route.meta.mobileBreadcrumbItems = read(breadcrumbs)
    route.meta.mobileHeaderActions = read(actions)
  })

  onBeforeUnmount(() => {
    route.meta.mobileBreadcrumbItems = []
    route.meta.mobileHeaderActions = []
  })
}
