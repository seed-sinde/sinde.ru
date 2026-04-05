<script setup lang="ts">
  const mobileBreadcrumbItems = useState<BreadcrumbItem[]>('mobile-header-breadcrumb-items', () => [])
  const mobileBreadcrumbOwner = useState<string>('mobile-header-breadcrumb-owner', () => '')
  const isCompactHeader = ref(false)
  const headerInstanceId = `lab-nav-header:${Math.random().toString(36).slice(2)}`
  let removeCompactHeaderListener: (() => void) | null = null
  const props = withDefaults(
    defineProps<{
      title: string
      breadcrumbItems?: BreadcrumbItem[]
      breadcrumbContainerClass?: string
    }>(),
    {
      breadcrumbItems: () => [],
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
  watch(
    resolvedBreadcrumbItems,
    next => {
      mobileBreadcrumbOwner.value = headerInstanceId
      mobileBreadcrumbItems.value = next
    },
    { immediate: true }
  )
  onMounted(() => {
    if (!import.meta.client) return
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
  })
</script>
<template>
  <div>
    <ClientOnly>
      <Teleport v-if="$slots.actions && isCompactHeader" to="#mobile-header-actions">
        <div class="flex items-center gap-1">
          <slot name="actions" :compact="true"></slot>
        </div>
      </Teleport>
    </ClientOnly>
    <div
      v-if="showBreadcrumbRow"
      class="border-(--lab-border) bg-(--lab-bg-overlay) sticky top-0 z-30 hidden border-b px-3 py-2 lg:block lg:px-6">
      <div class="flex min-h-9 flex-wrap items-center gap-2">
        <LabNavBreadcrumb
          v-if="resolvedBreadcrumbItems.length"
          :items="resolvedBreadcrumbItems"
          :container-class="breadcrumbContainerClass" />
        <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
          <slot name="actions" :compact="false"></slot>
        </div>
      </div>
    </div>
  </div>
</template>
