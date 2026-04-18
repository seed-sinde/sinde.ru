<template>
  <div>
    <div class="relative">
      <div ref="tabsScrollerRef" class="lab-scroll-hidden overflow-x-auto overflow-y-hidden">
        <div role="tablist" class="flex min-w-max items-end gap-0">
          <component
            :is="item.link ? nuxtLinkComponent : 'button'"
            v-for="(item, index) in resolvedItems"
            :key="String(item.value)"
            :ref="bindTabRef(index)"
            :to="item.routeTarget"
            :replace="item.link ? linkReplace : undefined"
            :type="item.link ? undefined : 'button'"
            role="tab"
            :tabindex="item.disabled ? -1 : item.active ? 0 : -1"
            :aria-selected="item.active ? 'true' : 'false'"
            :aria-disabled="item.disabled ? 'true' : undefined"
            :disabled="item.disabled"
            :class="[
              tabBaseClass,
              item.active ? activeStyles : inactiveStyles,
              item.disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
            ]"
            @click="onTabClick($event, item)"
          >
            <slot name="tab" :item="item" :active="item.active">
              <span>{{ item.label }}</span>
              <span
                v-if="item.badge !== undefined && item.badge !== null"
                class="ml-1.5 text-xs"
                :class="item.active ? 'text-(--lab-accent)' : 'text-(--lab-text-soft)'"
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

    <div v-if="renderPanels && activeValue !== null" class="space-y-4 p-4">
      <div :key="String(activeValue)">
        <slot :name="`panel-${String(activeValue)}`" :value="activeValue" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component, ComponentPublicInstance } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: LabTabValue
    items: LabTabItem[]
    renderPanels?: boolean
    routeQueryKey?: string
    routeParamKey?: string
    routePath?: string
    routePersistDefault?: boolean
    routeTargetMap?: TabRouteTargetMap
  }>(),
  {
    renderPanels: true,
    routeQueryKey: '',
    routeParamKey: '',
    routePath: '',
    routePersistDefault: false
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: LabTabValue]
  change: [value: LabTabValue, item: LabTabItem]
}>()

const route = useRoute()
const nuxtLinkComponent = resolveComponent('NuxtLink') as Component

const tabBaseClass = 'lab-tabs-focus relative -mb-px inline-flex h-9 shrink-0 items-center justify-center border-x-0 border-t-0 border-b-2 border-transparent bg-transparent px-2.5 text-xs font-medium whitespace-nowrap outline-none select-none'
const activeStyles = 'border-(--lab-accent) bg-[color-mix(in_srgb,var(--lab-accent)_14%,transparent)] text-(--lab-accent)'
const inactiveStyles = 'text-(--lab-text-muted) hover:border-(--lab-border-strong) hover:bg-[color-mix(in_srgb,var(--lab-bg-surface-subtle)_72%,transparent)] hover:text-(--lab-text-primary)'

const tabsScrollerRef = ref<HTMLElement | null>(null)
const tabRefs = ref<Array<HTMLElement | null>>([])
const { edges: tabsScrollEdges } = useScrollableEdges(tabsScrollerRef, { axis: 'x' })

const findItemValue = (raw: unknown) => {
  const candidate = Array.isArray(raw) ? raw[0] : raw
  if (typeof candidate !== 'string' && typeof candidate !== 'number') return null
  return props.items.find(item => stringifyTabValue(item.value) === stringifyTabValue(candidate))?.value ?? null
}

const activeValue = computed<LabTabValue | null>(() => {
  const paramValue = props.routeParamKey ? findItemValue(route.params[props.routeParamKey]) : null
  if (paramValue !== null) return paramValue
  const queryValue = props.routeQueryKey ? findItemValue(route.query[props.routeQueryKey]) : null
  if (queryValue !== null) return queryValue
  if (findItemValue(props.modelValue) !== null) return props.modelValue
  return props.items[0]?.value ?? null
})

const getTabRoute = (item: LabTabItem) =>
  buildTabRouteLocation(route, item.value, {
    ...(props.routeQueryKey ? { queryKey: props.routeQueryKey } : {}),
    defaultValue: props.items[0]?.value ?? props.modelValue,
    ...(props.routePath ? { path: props.routePath } : {}),
    persistDefault: props.routePersistDefault,
    ...(props.routeTargetMap ? { targetMap: props.routeTargetMap } : {})
  })

const linkReplace = computed(() => Boolean(props.routeQueryKey && !props.routeParamKey && !props.routeTargetMap))
const activeKey = computed(() => (activeValue.value === null ? '' : stringifyTabValue(activeValue.value)))
const resolvedItems = computed(() =>
  props.items.map(item => {
    const routeTarget = item.disabled ? undefined : getTabRoute(item)
    return {
      ...item,
      active: stringifyTabValue(item.value) === activeKey.value,
      link: Boolean(routeTarget),
      routeTarget
    }
  })
)

const setTabRef = (i: number, el: Element | ComponentPublicInstance | null) => {
  tabRefs.value[i] = el instanceof HTMLElement ? el : el && typeof el === 'object' && '$el' in el && el.$el instanceof HTMLElement ? el.$el : null
}

const bindTabRef = (i: number) => (el: Element | ComponentPublicInstance | null) => setTabRef(i, el)

const onTabClick = (event: MouseEvent, item: LabTabItem & { link?: boolean }) => {
  if (item.disabled || item.value === activeValue.value) return
  if (item.link) {
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0) return
    return
  }
  emit('update:modelValue', item.value)
  emit('change', item.value, item)
}

watch(
  () => activeValue.value,
  newVal => {
    if (newVal === null || stringifyTabValue(newVal) === stringifyTabValue(props.modelValue)) return
    emit('update:modelValue', newVal)
  }
)

watch(
  () => activeValue.value,
  async (newVal) => {
    if (!import.meta.client || newVal === null) return
    const idx = resolvedItems.value.findIndex(item => item.active)
    if (idx < 0) return
    await nextTick()
    tabRefs.value[idx]?.scrollIntoView({ block: 'nearest', inline: 'nearest' })
  },
  { immediate: true }
)
</script>
