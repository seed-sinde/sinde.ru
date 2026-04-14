<template>
  <div class="lab-tabs-root" :class="containerClass">
    <div class="relative min-w-0">
      <div ref="tabsScrollerRef" class="lab-scroll-hidden min-w-0 overflow-x-auto overflow-y-hidden">
        <div role="tablist" :class="resolvedListClass">
          <component
            :is="tabComponent(item)"
            v-for="(item, index) in items"
            :key="String(item.value)"
            :ref="(el: Element | { $el?: Element } | null) => setTabRef(index, el)"
            v-bind="tabComponentProps(item)"
            role="tab"
            :aria-selected="isActive(item.value) ? 'true' : 'false'"
            :tabindex="isActive(item.value) ? 0 : -1"
            :aria-disabled="item.disabled ? 'true' : undefined"
            :class="[
              resolvedButtonClass,
              noSelect ? 'select-none' : '',
              isActive(item.value) ? activeClass : inactiveClass
            ]"
            @click="onTabClick($event, item)"
            @keydown="onTabKeydown($event, index)"
          >
            <slot name="tab" :item="item" :active="isActive(item.value)">
              <span>{{ item.label }}</span>
              <span
                v-if="item.badge !== undefined && item.badge !== null"
                class="ml-1.5 text-xs text-(--lab-text-soft)"
                :class="isActive(item.value) ? 'text-(--lab-accent)' : ''"
              >
                {{ item.badge }}
              </span>
            </slot>
          </component>
        </div>
      </div>
      <div
        class="lab-scroll-fade lab-scroll-fade-x-left"
        :class="{ 'lab-scroll-fade-visible': tabsScrollEdges.left }"
        aria-hidden="true"
      />
      <div
        class="lab-scroll-fade lab-scroll-fade-x-right"
        :class="{ 'lab-scroll-fade-visible': tabsScrollEdges.right }"
        aria-hidden="true"
      />
    </div>
    <div v-if="renderPanels" :class="panelClass">
      <div v-for="item in items" v-show="isActive(item.value)" :key="`panel:${String(item.value)}`">
        <slot :name="`panel-${String(item.value)}`" :value="item.value" :item="item" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { twMerge } from 'tailwind-merge'
const props = withDefaults(
  defineProps<{
    modelValue: LabTabValue
    items: LabTabItem[]
    containerClass?: string
    listClass?: string
    buttonClass?: string
    activeClass?: string
    inactiveClass?: string
    panelClass?: string
    noSelect?: boolean
    renderPanels?: boolean
    routeQueryKey?: string
    routeDefaultValue?: LabTabValue
    routePath?: string
    routePersistDefault?: boolean
    routeActiveValue?: LabTabValue | null
    routeToMap?: Record<string, RouteLocationRaw | undefined>
  }>(),
  {
    containerClass:
      'max-sm:[&_.lab-tabs-root_.lab-tabs-root]:mt-2 max-sm:[&_.lab-tabs-root_.lab-tabs-root]:border-t max-sm:[&_.lab-tabs-root_.lab-tabs-root]:pt-2',
    listClass: '',
    buttonClass: '',
    activeClass:
      'border-b-(--lab-accent) bg-[color-mix(in_srgb,var(--lab-accent)_14%,transparent)] text-(--lab-accent)',
    inactiveClass:
      'text-(--lab-text-muted) hover:bg-[color-mix(in_srgb,var(--lab-bg-surface-subtle)_72%,transparent)] hover:text-(--lab-text-primary)',
    panelClass: 'space-y-3',
    noSelect: false,
    renderPanels: true,
    routeQueryKey: '',
    routePath: '',
    routePersistDefault: false,
    routeActiveValue: null
  }
)
const emit = defineEmits<{
  'update:modelValue': [value: LabTabValue]
  change: [value: LabTabValue, item: LabTabItem]
}>()
const route = useRoute()
const router = useRouter()
const LabBaseButton = resolveComponent('LabBaseButton')
const NuxtLink = resolveComponent('NuxtLink')
const tabsScrollerRef = ref<HTMLElement | null>(null)
const { edges: tabsScrollEdges, sync: syncTabsScrollEdges } = useScrollableEdges(tabsScrollerRef, { axis: 'x' })
const resolvedListClass = computed(() =>
  twMerge(
    'flex min-w-max items-end gap-0 border-b [border-color:color-mix(in_srgb,var(--lab-border)_62%,transparent)]',
    props.listClass
  )
)
const resolvedButtonClass = computed(() =>
  twMerge(
    'lab-tabs-focus text-(--lab-text-muted) relative inline-flex h-9 shrink-0 items-center justify-center border-0 border-b-2 border-transparent bg-transparent px-2.5 text-xs font-medium whitespace-nowrap outline-none transition-colors',
    props.buttonClass
  )
)
const fallbackValue = computed<LabTabValue>(() => {
  if (props.routeDefaultValue !== undefined) return props.routeDefaultValue
  const first = props.items[0]
  if (first) return first.value
  return props.modelValue
})
const allowedValues = computed(() => props.items.map((item) => item.value))
const routedValue = computed<LabTabValue | null>(() => {
  if (props.routeActiveValue !== null && props.routeActiveValue !== undefined) {
    return props.routeActiveValue
  }
  if (!props.routeQueryKey) return null
  return normalizeTabRouteValue(route.query[props.routeQueryKey], allowedValues.value, fallbackValue.value)
})
const activeValue = computed<LabTabValue>(() => {
  if (routedValue.value !== null) return routedValue.value
  if (props.items.some((item) => item.value === props.modelValue)) return props.modelValue
  return fallbackValue.value
})
const tabRefs = ref<Array<HTMLElement | null>>([])
const syncTabsScrollEdgesSoon = () => {
  if (import.meta.client) {
    requestAnimationFrame(syncTabsScrollEdges)
    return
  }
  syncTabsScrollEdges()
}
const hasElementRef = (value: Element | { $el?: Element } | null): value is { $el: Element } => {
  return Boolean(value && typeof value === 'object' && '$el' in value && value.$el)
}
const setTabRef = (index: number, el: Element | { $el?: Element } | null) => {
  const node = el instanceof HTMLElement ? el : hasElementRef(el) && el.$el instanceof HTMLElement ? el.$el : null
  tabRefs.value[index] = node
}
const isActive = (value: LabTabValue) => value === activeValue.value
const scrollActiveTabIntoView = async () => {
  await nextTick()
  const activeIndex = props.items.findIndex((item) => item.value === activeValue.value)
  if (activeIndex < 0) {
    syncTabsScrollEdgesSoon()
    return
  }
  const activeTab = tabRefs.value[activeIndex]
  if (activeTab) {
    activeTab.scrollIntoView({
      block: 'nearest',
      inline: 'nearest'
    })
  }
  syncTabsScrollEdgesSoon()
}
watch(
  routedValue,
  (next) => {
    if (next === null) return
    if (props.modelValue === next) return
    emit('update:modelValue', next)
  },
  { immediate: true }
)
watch(
  () => [activeValue.value, props.items.length],
  () => {
    void scrollActiveTabIntoView()
  },
  { immediate: true }
)
const tabToMapResolved = computed(() => {
  const out = new Map<LabTabValue, RouteLocationRaw | undefined>()
  for (const item of props.items) {
    const routeOptions: BuildTabRouteOptions = {
      defaultValue: fallbackValue.value,
      persistDefault: props.routePersistDefault
    }
    if (props.routeQueryKey) {
      routeOptions.queryKey = props.routeQueryKey
    }
    if (props.routePath) {
      routeOptions.path = props.routePath
    }
    if (props.routeToMap) {
      routeOptions.targetMap = props.routeToMap
    }
    out.set(item.value, buildTabRouteLocation(route, item.value, routeOptions))
  }
  return out
})
const tabTo = (item: LabTabItem) => tabToMapResolved.value.get(item.value)
const isLinkTab = (item: LabTabItem) => !item.disabled && Boolean(tabTo(item))
const tabComponent = (item: LabTabItem) => (isLinkTab(item) ? NuxtLink : LabBaseButton)
const tabComponentProps = (item: LabTabItem) => {
  if (isLinkTab(item)) {
    return {
      to: tabTo(item)
    }
  }
  return {
    type: 'button',
    disabled: item.disabled,
    focusClass: ''
  }
}
const emitTabChange = (item: LabTabItem) => {
  if (item.disabled) return
  if (item.value === activeValue.value) return
  emit('update:modelValue', item.value)
  emit('change', item.value, item)
}
const navigateToTab = async (item: LabTabItem) => {
  if (item.disabled) return
  if (item.value === activeValue.value) return
  const nextLocation = tabTo(item)
  emitTabChange(item)
  if (!nextLocation) return
  try {
    await router.replace(nextLocation)
  } catch {
    // ignore navigation duplication
  }
}
const hasModifierKey = (event: MouseEvent) =>
  event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0
const onTabClick = (event: MouseEvent, item: LabTabItem) => {
  if (item.disabled) {
    event.preventDefault()
    return
  }
  if (isLinkTab(item)) {
    if (item.value === activeValue.value && !hasModifierKey(event)) {
      event.preventDefault()
      return
    }
    emit('change', item.value, item)
    return
  }
  void navigateToTab(item)
}
const moveFocus = (startIndex: number, step: 1 | -1) => {
  if (!props.items.length) return
  let next = startIndex
  for (let i = 0; i < props.items.length; i += 1) {
    next = (next + step + props.items.length) % props.items.length
    const item = props.items[next]
    if (!item?.disabled) {
      const tab = tabRefs.value[next]
      tab?.focus()
      tab?.scrollIntoView({
        block: 'nearest',
        inline: 'nearest'
      })
      syncTabsScrollEdgesSoon()
      return
    }
  }
}
const onTabKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'ArrowRight') {
    event.preventDefault()
    moveFocus(index, 1)
    return
  }
  if (event.key === 'ArrowLeft') {
    event.preventDefault()
    moveFocus(index, -1)
    return
  }
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    const item = props.items[index]
    if (item) void navigateToTab(item)
  }
}
</script>
