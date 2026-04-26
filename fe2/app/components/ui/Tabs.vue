<template>
  <div v-if="!isInvalidTabPath">
    <div class="relative">
      <div ref="scrollerRef" class="ui-tabs-scroll min-w-0 overflow-x-auto overflow-y-hidden p-1">
        <div role="tablist" class="flex min-w-max items-center gap-2">
          <UiButton
            v-for="(tab, index) in tabs"
            :key="tabKey(tab)"
            :ref="bindTabRef(index)"
            variant="ghost"
            role="tab"
            :to="tab.path"
            :disabled="tab.disabled"
            :tabindex="tab.disabled ? -1 : 0"
            :aria-selected="tab.active ? 'true' : 'false'"
            :aria-current="tab.active ? 'page' : undefined"
            @click="openTab(tab)"
            @keydown.tab="switchTab(index, $event)"
          >
            {{ tab.label }}
          </UiButton>
        </div>
      </div>
      <div
        v-if="scrollEdges.left"
        class="pointer-events-none absolute top-1 bottom-1 left-0 w-0.5 bg-(--accent)"
      />
      <div
        v-if="scrollEdges.right"
        class="pointer-events-none absolute top-1 right-0 bottom-1 w-0.5 bg-(--accent)"
      />
    </div>
    <div v-if="activeTab" :key="tabKey(activeTab)" class="pt-4">
      <slot :name="activeSlot" :item="activeTab" :to="activeTab.to" :path="activeTab.path">
        <slot :item="activeTab" :to="activeTab.to" :path="activeTab.path" />
      </slot>
    </div>
  </div>
</template>
<script setup lang="ts">
type TabItem = {
  to: string
  label: string
  key?: string
  disabled?: boolean
}
type ResolvedTab = TabItem & {
  active: boolean
  path: string
}
interface Props {
  items?: TabItem[]
  basePath?: string
  modelValue?: string
}
const props = withDefaults(defineProps<Props>(), {
  items: () => [],
  basePath: ""
})
const emit = defineEmits<{
  "update:modelValue": [value: string]
  change: [value: string, item: TabItem]
}>()
const route = useRoute()
const scrollerRef = ref<HTMLElement | null>(null)
const tabRefs = ref<Array<HTMLElement | null>>([])
const currentPath = ref(route.path)
const {edges: scrollEdges, sync: syncScrollEdges} = useScrollableEdges(scrollerRef, {axis: "x"})
const trimSlashes = (v: string) => v.replace(/^\/+|\/+$/g, "")
const normalizePath = (v: string) => `/${trimSlashes(v)}`.replace(/\/$/, "") || "/"
const joinPath = (base: string, to: string) =>
  normalizePath(`${normalizePath(base)}/${trimSlashes(to)}`)
const tabKey = (tab: TabItem) => (tab.key ?? trimSlashes(tab.to)) || tab.to
const itemByValue = (value?: string) =>
  props.items.find(tab => tab.to === value || tab.key === value)
const baseFromPath = (path: string) => {
  const hit = props.items.find(
    tab => path === joinPath(path.slice(0, -normalizePath(tab.to).length), tab.to)
  )
  return hit ? path.slice(0, -normalizePath(hit.to).length) : path
}
const matchedBase = computed(() =>
  normalizePath(
    (route.matched[route.matched.length - 1]?.path ?? route.path).split("/:")[0] ?? route.path
  )
)
const strippedBase = computed(() => {
  if (props.basePath) return normalizePath(props.basePath)
  const base = baseFromPath(currentPath.value)
  return normalizePath(base === currentPath.value ? matchedBase.value : base)
})
const tabPath = (tab: TabItem) => joinPath(strippedBase.value, tab.to)
const matchedTab = computed(() => props.items.find(tab => tabPath(tab) === currentPath.value))
const isInvalidTabPath = computed(
  () => props.items.length > 0 && currentPath.value !== strippedBase.value && !matchedTab.value
)
const activeSource = computed(
  () => matchedTab.value ?? itemByValue(props.modelValue) ?? props.items[0]
)
const tabs = computed<ResolvedTab[]>(() =>
  props.items.map(tab => {
    const path = tabPath(tab)
    return {
      ...tab,
      path,
      active: activeSource.value ? tabKey(tab) === tabKey(activeSource.value) : false
    }
  })
)
const activeTab = computed(() => tabs.value.find(tab => tab.active))
const activeSlot = computed(() =>
  activeTab.value ? `panel-${tabKey(activeTab.value)}` : "default"
)
const setTabRef = (i: number, el: unknown) => {
  tabRefs.value[i] =
    el instanceof HTMLElement
      ? el
      : el && typeof el === "object" && "$el" in el && el.$el instanceof HTMLElement
        ? el.$el
        : null
}
const bindTabRef = (i: number) => (el: unknown) => setTabRef(i, el)
const nextTabIndex = (from: number, step: number) => {
  for (let index = from + step; index >= 0 && index < tabs.value.length; index += step) {
    if (!tabs.value[index]?.disabled) return index
  }
  return -1
}
const showTabNotFound = () =>
  showError(createError({statusCode: 404, statusMessage: "Tab not found"}))
const openTab = (tab: ResolvedTab) => {
  if (tab.disabled || tab.active) return
  emit("change", tab.to, tab)
}
const switchTab = async (from: number, event: KeyboardEvent) => {
  const step = event.shiftKey ? -1 : 1
  const index = nextTabIndex(from, step)
  const tab = tabs.value[index]
  if (!tab) return
  event.preventDefault()
  await nextTick()
  tabRefs.value[index]?.focus()
}
if (isInvalidTabPath.value) showTabNotFound()
watch(
  () => route.path,
  path => {
    currentPath.value = path
  }
)
watch(isInvalidTabPath, invalid => {
  if (invalid) showTabNotFound()
})
watch(activeTab, (tab, oldTab) => {
  if (tab && tabKey(tab) !== (oldTab && tabKey(oldTab))) emit("update:modelValue", tab.to)
})
watch(tabs, async () => {
  await nextTick()
  syncScrollEdges()
})
</script>

<style scoped>
.ui-tabs-scroll {
  scrollbar-width: none;
}
.ui-tabs-scroll::-webkit-scrollbar {
  display: none;
}
</style>
