<template>
  <nav ref="containerRef" :aria-label="ariaLabel" :class="resolvedContainerClass">
    <ol :class="listClass">
      <li
        v-for="(item, index) in normalizedItems"
        :key="`${item.label}:${index}`"
        class="inline-flex min-w-0 items-center gap-2">
        <template v-if="index === 0">
          <NuxtLink
            v-if="leadingSectionIcon && canLinkToSection"
            :to="leadingSectionTo"
            class="inline-flex shrink-0 items-center"
            :aria-label="t('nav.open_section', { label: matchedSidebarItem?.label || item.label })">
            <Icon :name="leadingSectionIcon" :class="leadingSectionIconClass" aria-hidden="true" />
          </NuxtLink>
          <Icon
            v-else-if="leadingSectionIcon"
            :name="leadingSectionIcon"
            :class="leadingSectionIconClass"
            aria-hidden="true" />
        </template>
        <template v-else>
          <Icon :name="separatorIcon" :class="separatorClass" aria-hidden="true" />
        </template>
        <NuxtLink v-if="item.to && !item.current" :to="item.to" :class="itemClasses(item, false)">
          <span class="inline-flex min-w-0 items-center gap-2">
            <span v-if="item.label">{{ item.label }}</span>
            <LabBaseBadge v-if="item.badge" variant="info" size="xs">{{ item.badge }}</LabBaseBadge>
          </span>
        </NuxtLink>
        <span v-else :class="itemClasses(item, true)" :aria-current="item.current ? 'page' : undefined">
          <span class="inline-flex min-w-0 items-center gap-2">
            <span v-if="item.label">{{ item.label }}</span>
            <LabBaseBadge v-if="item.badge" variant="info" size="xs">{{ item.badge }}</LabBaseBadge>
          </span>
        </span>
      </li>
    </ol>
  </nav>
</template>
<script setup lang="ts">
  const route = useRoute()
  const { t } = useInterfacePreferences()
  const translatedSidebarItems = useSidebarItems()
  const props = withDefaults(
    defineProps<{
      items: Array<{
        label: string
        to?: string | Record<string, any>
        current?: boolean
        badge?: string
      }>
      ariaLabel?: string
      containerClass?: string
      listClass?: string
      separatorIcon?: string
      separatorClass?: string
      leadingIcon?: string
      leadingIconClass?: string
    }>(),
    {
      ariaLabel: '',
      containerClass: '',
      listClass: 'flex min-w-max items-center gap-x-2 gap-y-1 text-sm whitespace-nowrap',
      separatorIcon: 'ic:round-chevron-right',
      separatorClass: 'text-(--lab-text-soft) h-4 w-4 shrink-0',
      leadingIcon: '',
      leadingIconClass: 'text-(--lab-text-muted) h-4 w-4 shrink-0'
    }
  )
  const containerRef = ref<HTMLElement | null>(null)
  const normalizedItems = computed(() => props.items || [])
  const ariaLabel = computed(() => props.ariaLabel || t('nav.breadcrumbs'))
  const resolvedContainerClass = computed(() =>
    [props.containerClass, 'min-w-0 overflow-x-auto overscroll-x-contain'].filter(Boolean)
  )
  const normalizeToPath = (to: unknown) => {
    if (!to) return ''
    if (typeof to === 'string') return to
    if (typeof to === 'object' && to && 'path' in to) {
      return String((to as { path?: string }).path || '')
    }
    return ''
  }
  const firstItemPath = computed(() => {
    const first = normalizedItems.value[0]
    return normalizeToPath(first?.to)
  })
  const routeSectionItem = computed(() => {
    const path = String(route.path || '')
    if (!path || path === '/') return null
    return translatedSidebarItems.value
      .filter(item => path === item.to || path.startsWith(`${item.to}/`))
      .sort((left, right) => right.to.length - left.to.length)[0] || null
  })
  const matchedSidebarItem = computed(() => {
    const firstPath = firstItemPath.value
    if (!firstPath) return routeSectionItem.value || null
    return translatedSidebarItems.value.find(item => item.to === firstPath) || routeSectionItem.value || null
  })
  const leadingSectionTo = computed(() => matchedSidebarItem.value?.to || '')
  const canLinkToSection = computed(() => {
    if (!leadingSectionTo.value) return false
    return normalizedItems.value.length > 1 || route.path !== leadingSectionTo.value
  })
  const leadingSectionIcon = computed(() => {
    if (props.leadingIcon) return props.leadingIcon
    return matchedSidebarItem.value?.icon || 'ic:round-folder'
  })
  const leadingSectionIconClass = computed(() => {
    if (props.leadingIconClass !== 'text-(--lab-text-muted) h-4 w-4 shrink-0') {
      return props.leadingIconClass
    }
    return ['h-4 w-4 shrink-0', matchedSidebarItem.value?.iconColor || 'text-(--lab-text-muted)']
  })
  const itemClasses = (item: { current?: boolean }, isCurrent: boolean) => {
    return [
      'inline-flex min-w-0 items-center wrap-break-word',
      isCurrent ? 'text-(--lab-text-primary)' : 'text-(--lab-text-muted) hover:text-(--lab-text-primary)'
    ]
  }
  const scrollCurrentItemIntoView = async () => {
    await nextTick()
    const container = containerRef.value
    if (!container) return
    const currentItem = container.querySelector('[aria-current="page"]')
    if (currentItem instanceof HTMLElement) {
      currentItem.scrollIntoView({
        block: 'nearest',
        inline: 'end'
      })
      return
    }
    container.scrollTo({
      left: container.scrollWidth,
      top: 0
    })
  }
  watch(
    () => [route.fullPath, normalizedItems.value.length],
    () => {
      scrollCurrentItemIntoView()
    },
    { immediate: true }
  )
</script>
