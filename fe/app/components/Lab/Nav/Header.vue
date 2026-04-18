<script setup lang="ts">
const mobileBreadcrumbItems = useState<BreadcrumbItem[]>('mobile-header-breadcrumb-items', () => [])
const mobileBreadcrumbOwner = useState<string>('mobile-header-breadcrumb-owner', () => '')
const mobileHeaderActions = useMobileHeaderActions()
const mobileHeaderActionsOwner = useMobileHeaderActionsOwner()
const isCompactHeader = ref(false)
const headerInstanceId = `lab-nav-header:${Math.random().toString(36).slice(2)}`
let removeCompactHeaderListener: (() => void) | null = null
const props = withDefaults(
  defineProps<{
    title: string
    breadcrumbItems?: BreadcrumbItem[]
    mobileActions?: MobileHeaderAction[]
    breadcrumbContainerClass?: string
  }>(),
  {
    breadcrumbItems: () => [],
    mobileActions: () => [],
    breadcrumbContainerClass: ''
  }
)
const slots = useSlots()
const resolvedBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
  if (Array.isArray(props.breadcrumbItems) && props.breadcrumbItems.length) {
    return props.breadcrumbItems
  }
  return props.title ? [{ label: props.title, current: true }] : []
})
const showBreadcrumbRow = computed(() => resolvedBreadcrumbItems.value.length > 0 || Boolean(slots.actions))
const syncMobileBreadcrumbItems = (items: BreadcrumbItem[]) => {
  mobileBreadcrumbOwner.value = headerInstanceId
  mobileBreadcrumbItems.value = items
}
const syncMobileHeaderActions = (actions: MobileHeaderAction[]) => {
  mobileHeaderActionsOwner.value = headerInstanceId
  mobileHeaderActions.value = [...actions]
}
onMounted(() => {
  if (!import.meta.client) return
  syncMobileBreadcrumbItems(resolvedBreadcrumbItems.value)
  syncMobileHeaderActions(props.mobileActions)
  watch(resolvedBreadcrumbItems, next => {
    syncMobileBreadcrumbItems(next)
  })
  watch(
    () => props.mobileActions,
    next => {
      syncMobileHeaderActions(next)
    }
  )
  const mediaQuery = window.matchMedia('(min-width: 64rem)')
  const syncCompactHeader = () => {
    isCompactHeader.value = !mediaQuery.matches
  }
  syncCompactHeader()
  mediaQuery.addEventListener('change', syncCompactHeader)
  removeCompactHeaderListener = () => mediaQuery.removeEventListener('change', syncCompactHeader)
})
onBeforeUnmount(() => {
  removeCompactHeaderListener?.()
  removeCompactHeaderListener = null
  if (mobileBreadcrumbOwner.value === headerInstanceId) {
    mobileBreadcrumbItems.value = []
    mobileBreadcrumbOwner.value = ''
  }
  if (mobileHeaderActionsOwner.value === headerInstanceId) {
    mobileHeaderActions.value = []
    mobileHeaderActionsOwner.value = ''
  }
})
</script>
<template>
  <div class="sticky top-0 z-30 h-full">
    <ClientOnly>
      <Teleport v-if="$slots.actions && isCompactHeader && !props.mobileActions.length" to="#mobile-header-actions">
        <div class="flex items-center gap-1">
          <slot name="actions" :compact="true" />
        </div>
      </Teleport>
    </ClientOnly>
    <div
      v-if="showBreadcrumbRow"
      class="hidden h-10 border-b border-(--lab-border) bg-(--lab-bg-overlay) px-4 lg:block"
    >
      <div class="flex h-10 items-center gap-2">
        <LabNavBreadcrumb
          v-if="resolvedBreadcrumbItems.length"
          :items="resolvedBreadcrumbItems"
          :container-class="breadcrumbContainerClass"
        />
        <div v-if="$slots.actions" class="flex items-center gap-2">
          <slot name="actions" :compact="false" />
        </div>
      </div>
    </div>
  </div>
</template>
