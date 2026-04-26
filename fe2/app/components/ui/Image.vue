<script setup lang="ts">
import type {CSSProperties} from "vue"
import IcRoundClose from "~icons/ic/round-close"
import IcRoundNavigateBefore from "~icons/ic/round-navigate-before"
import IcRoundNavigateNext from "~icons/ic/round-navigate-next"
import IcRoundOpenInFull from "~icons/ic/round-open-in-full"
import IcRoundRotateLeft from "~icons/ic/round-rotate-left"
import IcRoundRotateRight from "~icons/ic/round-rotate-right"
import IcRoundSwapHoriz from "~icons/ic/round-swap-horiz"
import IcRoundSwapVert from "~icons/ic/round-swap-vert"
defineOptions({inheritAttrs: false})
type ViewerMode = "fit" | "original"
type ImageItem = {
  src: string
  title?: string | null
  author?: string | null
  attribution?: string | null
  sourceUrl?: string | null
  licenseUrl?: string | null
  license?: string | null
  alt?: string | null
  thumbnailSrc?: string | null
}
type Item = {
  src: string
  title: string | null
  alt: string
  author: string | null
  attribution: string | null
  sourceUrl: string | null
  licenseUrl: string | null
  license: string | null
  thumbnailSrc: string | null
}
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    src?: string
    title?: string
    alt?: string
    items?: ImageItem[]
    initialIndex?: number
    routeQueryParam?: string
    routeQuerySync?: boolean
  }>(),
  {
    modelValue: false,
    src: "",
    title: "",
    alt: "",
    items: () => [],
    initialIndex: 0,
    routeQueryParam: "image",
    routeQuerySync: true
  }
)
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  opened: []
  closed: []
  "mode-change": [mode: ViewerMode]
  "active-index-change": [index: number]
}>()
const attrs = useAttrs()
const triggerClass = computed(() => attrs.class)
const triggerAttrs = computed(() =>
  Object.fromEntries(Object.entries(attrs).filter(([key]) => key !== "class"))
)
const {locale, key, load, t} = useI18nSection("ui")
const fallback = {
  "viewer.flip_x": "Отразить по горизонтали",
  "viewer.flip_y": "Отразить по вертикали",
  "viewer.rotate_left": "Повернуть влево",
  "viewer.rotate_right": "Повернуть вправо",
  "viewer.scale": "Масштаб",
  "viewer.close": "Закрыть",
  "viewer.previous": "Предыдущее",
  "viewer.next": "Следующее",
  "viewer.open": "Посмотреть"
} as const
const label = (k: keyof typeof fallback) => {
  const v = t(k)
  return v === k ? fallback[k] : v
}
const text = (v: unknown) => String(v ?? "").trim() || null
const mod = (v: number, n: number) =>
  n ? ((Math.trunc(Number.isFinite(v) ? v : 0) % n) + n) % n : 0
const route = useRoute()
const router = useRouter()
const ownOpen = ref(false)
const activeIndex = ref(0)
const imageRef = ref<HTMLImageElement | null>(null)
const viewportRef = ref<HTMLElement | null>(null)
const toolbarRef = ref<HTMLElement | null>(null)
const thumbsRef = ref<HTMLElement | null>(null)
const natural = reactive({width: 0, height: 0})
const viewport = reactive({width: 0, height: 0})
const mode = ref<ViewerMode>("fit")
const loaded = ref(false)
const nav = ref<-1 | 1>(-1)
const flip = reactive({x: 1, y: 1})
const rotation = ref(0)
const touch = reactive({
  startX: 0,
  startY: 0,
  deltaX: 0,
  deltaY: 0,
  tracking: false,
  swiped: false,
  animating: false
})
const items = computed<Item[]>(() => {
  const list = props.items
    .map((item, i) => {
      const src = item.src.trim()
      const title = text(item.title)
      const alt = text(item.alt)
      return src
        ? {
            src,
            title,
            alt: alt || title || `Изображение ${i + 1}`,
            author: text(item.author),
            attribution: text(item.attribution),
            sourceUrl: text(item.sourceUrl),
            licenseUrl: text(item.licenseUrl),
            license: text(item.license),
            thumbnailSrc: text(item.thumbnailSrc)
          }
        : null
    })
    .filter((item): item is Item => item !== null)
  const title = text(props.title)
  return list.length || !props.src.trim()
    ? list
    : [
        {
          src: props.src.trim(),
          title,
          alt: text(props.alt) || title || "Изображение",
          author: null,
          attribution: null,
          sourceUrl: null,
          licenseUrl: null,
          license: null,
          thumbnailSrc: null
        }
      ]
})
const routeIndex = computed(() => {
  if (!props.routeQuerySync || !items.value.length) return null
  const k = props.routeQueryParam.trim()
  const raw = Array.isArray(route.query[k]) ? route.query[k]?.[0] : route.query[k]
  const parsed = Number.parseInt(String(raw || "").trim(), 10)
  return k && Number.isFinite(parsed) && parsed > 0 ? mod(parsed - 1, items.value.length) : null
})
const activeItem = computed(() => items.value[activeIndex.value] || null)
const credits = computed(() =>
  Boolean(
    activeItem.value?.title ||
    activeItem.value?.author ||
    activeItem.value?.attribution ||
    activeItem.value?.sourceUrl ||
    activeItem.value?.license
  )
)
const multiple = computed(() => items.value.length > 1)
const visibleModel = computed(() => props.modelValue || ownOpen.value)
const isVisible = computed(() =>
  Boolean((visibleModel.value || routeIndex.value !== null) && activeItem.value?.src)
)
const displayTitle = computed(() => activeItem.value?.title || text(props.title))
const overlayNav = computed(() => multiple.value && viewport.width >= 768)
const quarter = computed(() => Math.abs(rotation.value / 90) % 2 === 1)
const oriented = computed(() => ({
  width: quarter.value ? natural.height : natural.width,
  height: quarter.value ? natural.width : natural.height
}))
const available = computed(() => ({
  width: Math.max(
    0,
    viewport.width - (viewport.width < 640 ? 24 : 160) - (overlayNav.value ? 96 : 0)
  ),
  height: Math.max(
    0,
    viewport.height -
      (viewport.width < 640 ? 24 : 48) -
      (multiple.value || credits.value || displayTitle.value ? 12 : 0)
  )
}))
const fits = computed(() =>
  Boolean(
    oriented.value.width &&
    oriented.value.height &&
    available.value.width &&
    available.value.height &&
    oriented.value.width <= available.value.width &&
    oriented.value.height <= available.value.height
  )
)
const fitScale = computed(() =>
  oriented.value.width && oriented.value.height && available.value.width && available.value.height
    ? Math.min(
        1,
        available.value.width / oriented.value.width,
        available.value.height / oriented.value.height
      )
    : 1
)
const effectiveMode = computed<ViewerMode>(() => (fits.value ? "original" : mode.value))
const overflow = computed(() =>
  Boolean(
    oriented.value.width &&
    oriented.value.height &&
    available.value.width &&
    available.value.height &&
    (oriented.value.width > available.value.width || oriented.value.height > available.value.height)
  )
)
const canPan = computed(() => effectiveMode.value === "original" && overflow.value)
const canSwipe = computed(() => multiple.value && !canPan.value)
const scaleLabel = computed(() =>
  loaded.value
    ? `${Math.round((effectiveMode.value === "original" ? 1 : fitScale.value) * 100)}%`
    : "..."
)
const frameStyle = computed<CSSProperties>(() => {
  const scale = effectiveMode.value === "fit" ? fitScale.value : 1
  return natural.width && natural.height
    ? {
        width: `${Math.max(1, Math.round(oriented.value.width * scale))}px`,
        height: `${Math.max(1, Math.round(oriented.value.height * scale))}px`,
        flex: "0 0 auto"
      }
    : {}
})
const imageStyle = computed<CSSProperties>(() => {
  const scale = effectiveMode.value === "fit" ? fitScale.value : 1
  return natural.width && natural.height
    ? {
        position: "absolute",
        left: "50%",
        top: "50%",
        width: `${Math.max(1, Math.round(natural.width * scale))}px`,
        height: `${Math.max(1, Math.round(natural.height * scale))}px`,
        maxWidth: "none",
        maxHeight: "none",
        transform: `translate(-50%, -50%) rotate(${rotation.value}deg) scale(${flip.x}, ${flip.y})`,
        transformOrigin: "center center"
      }
    : {}
})
const contentStyle = computed<CSSProperties>(() => ({
  transform: `translate3d(${touch.deltaX}px, 0, 0)`,
  opacity: canSwipe.value
    ? String(Math.max(0.58, 1 - Math.abs(touch.deltaX) / Math.max(1, viewport.width * 1.15)))
    : "1",
  transition: `transform ${touch.animating ? 300 : 0}ms ease, opacity ${touch.animating ? 300 : 0}ms ease`
}))
const toolButtons = computed(() => [
  {icon: IcRoundSwapHoriz, title: label("viewer.flip_x"), click: () => flipAxis("x")},
  {icon: IcRoundSwapVert, title: label("viewer.flip_y"), click: () => flipAxis("y")},
  {icon: IcRoundRotateLeft, title: label("viewer.rotate_left"), click: () => rotate(-90)},
  {icon: IcRoundRotateRight, title: label("viewer.rotate_right"), click: () => rotate(90)}
])
const scroll = useScrollableEdges({
  toolbar: {target: toolbarRef, axis: "x"},
  thumbs: {target: thumbsRef, axis: "x"},
  viewport: {target: viewportRef, axis: "both"}
})
const toolbarEdges = scroll.items.toolbar.edges
const thumbEdges = scroll.items.thumbs.edges
const viewportEdges = scroll.items.viewport.edges
let resizeObserver: ResizeObserver | null = null
let revealFrame: number | null = null
let syncFrame: number | null = null
let swipeTimer: number | null = null
const syncQuery = async (i: number | null) => {
  if (!props.routeQuerySync || !import.meta.client) return
  const k = props.routeQueryParam.trim()
  const cur = Array.isArray(route.query[k]) ? route.query[k]?.[0] : route.query[k]
  const next = i === null ? null : String(i + 1)
  if (!k || String(cur || "").trim() === String(next || "")) return
  const query = {...route.query}
  if (next === null) Reflect.deleteProperty(query, k)
  else query[k] = next
  await router.replace({path: route.path, query})
}
const resetMode = () => {
  const next: ViewerMode = natural.width && natural.height && fits.value ? "original" : "fit"
  if (mode.value === next) return
  mode.value = next
  emit("mode-change", next)
}
const resetTouch = () => {
  if (swipeTimer !== null && import.meta.client) window.clearTimeout(swipeTimer)
  swipeTimer = null
  touch.startX = 0
  touch.startY = 0
  touch.deltaX = 0
  touch.deltaY = 0
  touch.tracking = false
  touch.swiped = false
  touch.animating = false
}
const cancelFrames = () => {
  if (revealFrame !== null && import.meta.client) window.cancelAnimationFrame(revealFrame)
  if (syncFrame !== null && import.meta.client) window.cancelAnimationFrame(syncFrame)
  revealFrame = null
  syncFrame = null
}
const updateViewport = () => {
  const el = viewportRef.value
  if (!el) return
  viewport.width = el.clientWidth
  viewport.height = el.clientHeight
  scroll.sync()
}
const observeViewport = () => {
  resizeObserver?.disconnect()
  if (viewportRef.value) resizeObserver?.observe(viewportRef.value)
  updateViewport()
}
const alignOriginal = () => {
  const el = viewportRef.value
  if (!import.meta.client || !el || effectiveMode.value !== "original") return
  el.scrollLeft = Math.max(0, Math.round((el.scrollWidth - el.clientWidth) / 2))
  el.scrollTop = Math.max(0, Math.round((el.scrollHeight - el.clientHeight) / 2))
  scroll.sync()
}
const centerThumb = () => {
  const el = thumbsRef.value
  if (!import.meta.client || !el || !multiple.value) return
  const thumb = el.querySelector<HTMLElement>(`[data-thumb-index="${activeIndex.value}"]`)
  if (!thumb) return
  el.scrollTo({
    left: Math.min(
      Math.max(0, el.scrollWidth - el.clientWidth),
      Math.max(0, thumb.offsetLeft - Math.max(0, (el.clientWidth - thumb.offsetWidth) / 2))
    ),
    behavior: "smooth"
  })
  requestAnimationFrame(scroll.sync)
}
const syncLoaded = () => {
  const img = imageRef.value
  if (!img?.complete || (!img.naturalWidth && !img.naturalHeight)) return
  natural.width = img.naturalWidth
  natural.height = img.naturalHeight
  resetMode()
  alignOriginal()
  loaded.value = true
}
const queueLoadedSync = () =>
  nextTick(() => {
    if (!import.meta.client) return syncLoaded()
    if (syncFrame !== null) window.cancelAnimationFrame(syncFrame)
    syncFrame = window.requestAnimationFrame(() => {
      syncFrame = null
      syncLoaded()
    })
  })
const resetImage = () => {
  cancelFrames()
  loaded.value = false
  natural.width = 0
  natural.height = 0
  flip.x = 1
  flip.y = 1
  rotation.value = 0
  nextTick(updateViewport)
}
const setIndex = (i: number, direction?: -1 | 1) => {
  const next = mod(i, items.value.length)
  if (activeIndex.value === next) return
  nav.value = direction || (next > activeIndex.value ? -1 : 1)
  activeIndex.value = next
  emit("active-index-change", next)
  void syncQuery(next)
}
const previous = () => multiple.value && setIndex(activeIndex.value - 1, 1)
const next = () => multiple.value && setIndex(activeIndex.value + 1, -1)
const open = (i = props.initialIndex) => {
  activeIndex.value = mod(i, items.value.length)
  ownOpen.value = true
  emit("update:modelValue", true)
  void syncQuery(activeIndex.value)
}
const close = () => {
  ownOpen.value = false
  void syncQuery(null)
  emit("update:modelValue", false)
}
const rotate = (deg: number) => {
  rotation.value = mod(rotation.value + deg, 360)
  nextTick(() => {
    resetMode()
    alignOriginal()
  })
}
const flipAxis = (axis: "x" | "y") => (flip[axis] *= -1)
const toggleZoom = () => {
  if (fits.value) return
  mode.value = mode.value === "fit" ? "original" : "fit"
  emit("mode-change", mode.value)
  nextTick(alignOriginal)
}
const onImageLoad = (e: Event) => {
  const img = e.target as HTMLImageElement | null
  if (!img) return
  natural.width = img.naturalWidth
  natural.height = img.naturalHeight
  resetMode()
  nextTick(alignOriginal)
  if (!import.meta.client) return
  if (revealFrame !== null) window.cancelAnimationFrame(revealFrame)
  revealFrame = window.requestAnimationFrame(() => {
    loaded.value = true
    revealFrame = null
  })
}
const onBackdropClick = () => (touch.swiped ? (touch.swiped = false) : close())
const onTouchStart = (e: TouchEvent) => {
  const v = e.touches[0]
  if (!v || !canSwipe.value) return resetTouch()
  touch.startX = v.clientX
  touch.startY = v.clientY
  touch.deltaX = 0
  touch.deltaY = 0
  touch.tracking = true
  touch.swiped = false
  touch.animating = false
}
const onTouchMove = (e: TouchEvent) => {
  const v = e.touches[0]
  if (!v || !touch.tracking || !canSwipe.value) return
  touch.deltaX = v.clientX - touch.startX
  touch.deltaY = v.clientY - touch.startY
}
const onTouchEnd = () => {
  if (!touch.tracking || !canSwipe.value) return resetTouch()
  const direction = touch.deltaX < 0 ? -1 : 1
  touch.tracking = false
  touch.animating = true
  if (Math.abs(touch.deltaX) >= 56 && Math.abs(touch.deltaX) > Math.abs(touch.deltaY)) {
    touch.swiped = true
    nav.value = direction
    touch.deltaX = direction * Math.max(Math.round(viewport.width * 1.35), 420)
    window.setTimeout(() => {
      if (direction < 0) next()
      else previous()
      touch.deltaX = 0
      touch.deltaY = 0
      touch.animating = false
    }, 300)
    swipeTimer = window.setTimeout(() => {
      touch.swiped = false
      swipeTimer = null
    }, 420)
    return
  }
  touch.deltaX = 0
  touch.deltaY = 0
  window.setTimeout(() => (touch.animating = false), 300)
}
const onKeydown = (e: KeyboardEvent) => {
  if (!isVisible.value) return
  if (e.key === "Escape") close()
  else if (e.key === "ArrowLeft") previous()
  else if (e.key === "ArrowRight") next()
}
const syncInitialIndex = () =>
  (activeIndex.value = routeIndex.value ?? mod(props.initialIndex, items.value.length))
watch(isVisible, nextVisible => {
  if (!nextVisible) return emit("closed")
  syncInitialIndex()
  emit("opened")
  nextTick(observeViewport)
  nextTick(alignOriginal)
  nextTick(centerThumb)
  queueLoadedSync()
})
watch(
  routeIndex,
  i => {
    if (i === null) return
    activeIndex.value = i
    ownOpen.value = true
    emit("update:modelValue", true)
  },
  {immediate: true}
)
watch([() => props.initialIndex, items], syncInitialIndex)
watch(
  () => activeItem.value?.src,
  () => {
    resetTouch()
    resetImage()
    queueLoadedSync()
    nextTick(centerThumb)
  }
)
watch([fits, () => viewport.width, () => viewport.height], () => {
  if (!isVisible.value || !natural.width || !natural.height) return
  if (fits.value && mode.value !== "original") {
    mode.value = "original"
    emit("mode-change", mode.value)
  }
  nextTick(alignOriginal)
})
watch(
  [() => activeIndex.value, () => viewport.width],
  () => isVisible.value && nextTick(centerThumb),
  {flush: "post"}
)
onMounted(() => {
  window.addEventListener("keydown", onKeydown)
  resizeObserver = new ResizeObserver(updateViewport)
  observeViewport()
})
onBeforeUnmount(() => {
  resetTouch()
  cancelFrames()
  window.removeEventListener("keydown", onKeydown)
  resizeObserver?.disconnect()
})
defineExpose({open, close})
await useAsyncData(key.value, load, {watch: [locale]})
</script>

<template>
  <slot
    :open="open"
    :close="close"
    :active-index="activeIndex"
    :active-item="activeItem"
    :items="items"
  >
    <button
      v-if="activeItem"
      type="button"
      :class="[
        'group ui-focus relative block max-w-full bg-transparent p-0 text-left',
        triggerClass
      ]"
      v-bind="triggerAttrs"
      @click="open()"
    >
      <span class="relative block aspect-4/3 w-full overflow-hidden">
        <img
          :src="activeItem.thumbnailSrc || activeItem.src"
          :alt="activeItem.alt"
          class="h-full w-full object-cover"
        />
        <span
          class="pointer-events-none absolute inset-x-0 bottom-0 flex items-center justify-center gap-1.5 bg-zinc-950/70 p-2 text-xs font-medium text-white opacity-0 transition-opacity group-hover:opacity-100 group-focus-visible:opacity-100"
        >
          <IcRoundOpenInFull class="h-4 w-4 shrink-0" />
          <span class="truncate">{{ label("viewer.open") }}</span>
        </span>
      </span>
    </button>
  </slot>

  <Teleport to="body">
    <div v-if="isVisible" class="fixed inset-0 z-60 bg-(--bg) text-(--text)">
      <div class="flex h-full w-full flex-col">
        <div class="z-10" @click.stop>
          <div class="relative">
            <div
              ref="toolbarRef"
              class="ui-image-scroll overflow-x-auto overflow-y-hidden bg-(--bg) px-2.5 py-1.5 ring-1 ring-(--border-color) sm:px-3"
            >
              <div class="flex min-w-max items-center justify-end gap-2.5">
                <div class="flex shrink-0 items-center justify-end gap-1.5">
                  <button
                    v-for="button in toolButtons"
                    :key="button.title"
                    type="button"
                    :title="button.title"
                    :aria-label="button.title"
                    class="ui-focus inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-(--elevated) text-sm transition hover:ring-2 hover:ring-(--accent)"
                    @click="button.click"
                  >
                    <component :is="button.icon" class="h-5 w-5" />
                  </button>
                </div>

                <div
                  class="flex shrink-0 items-center gap-1.5 text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  {{ label("viewer.scale") }}
                  <p class="w-[4ch] text-right tabular-nums">
                    {{ scaleLabel }}
                  </p>
                </div>

                <slot
                  name="toolbar-extra"
                  :active-index="activeIndex"
                  :active-item="activeItem"
                  :items="items"
                />

                <button
                  type="button"
                  :title="label('viewer.close')"
                  :aria-label="label('viewer.close')"
                  class="ui-focus inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-(--elevated) transition hover:ring-2 hover:ring-(--accent)"
                  @click="close"
                >
                  <IcRoundClose class="h-5 w-5" />
                </button>
              </div>
            </div>
            <span
              :class="[
                'pointer-events-none absolute inset-y-0 left-0 w-6 bg-(--bg) opacity-0 transition-opacity',
                toolbarEdges.left && 'opacity-90'
              ]"
              aria-hidden="true"
            />
            <span
              :class="[
                'pointer-events-none absolute inset-y-0 right-0 w-6 bg-(--bg) opacity-0 transition-opacity',
                toolbarEdges.right && 'opacity-90'
              ]"
              aria-hidden="true"
            />
          </div>
        </div>

        <div class="relative min-h-0 flex-1 overflow-hidden">
          <button
            v-if="overlayNav"
            type="button"
            :title="label('viewer.previous')"
            :aria-label="label('viewer.previous')"
            class="ui-focus absolute inset-y-0 left-0 z-10 flex w-20 items-center justify-center bg-(--elevated) transition hover:ring-2 hover:ring-(--accent) sm:w-32"
            @click.stop="previous"
          >
            <IcRoundNavigateBefore class="h-8 w-8" />
          </button>
          <button
            v-if="overlayNav"
            type="button"
            :title="label('viewer.next')"
            :aria-label="label('viewer.next')"
            class="ui-focus absolute inset-y-0 right-0 z-10 flex w-20 items-center justify-center bg-(--elevated) transition hover:ring-2 hover:ring-(--accent) sm:w-32"
            @click.stop="next"
          >
            <IcRoundNavigateNext class="h-8 w-8" />
          </button>

          <div
            ref="viewportRef"
            :class="[
              'relative h-full w-full bg-(--bg)',
              canPan ? 'overflow-auto' : 'overflow-hidden'
            ]"
            @click="onBackdropClick"
            @touchstart.passive="onTouchStart"
            @touchmove.passive="onTouchMove"
            @touchend="onTouchEnd"
            @touchcancel="resetTouch"
          >
            <div
              :class="[
                'box-border flex min-h-full min-w-full px-3 py-3 sm:px-20 sm:py-6',
                multiple && 'sm:px-24',
                (multiple || credits) && 'py-2 sm:py-3',
                effectiveMode === 'original' && overflow
                  ? 'h-max w-max items-center justify-center'
                  : 'h-full w-full items-center justify-center'
              ]"
              :style="contentStyle"
            >
              <div
                :class="[
                  effectiveMode === 'original' && overflow
                    ? 'relative shrink-0'
                    : 'relative flex h-full w-full items-center justify-center overflow-hidden'
                ]"
              >
                <Transition :name="nav < 0 ? 'ui-image-next' : 'ui-image-prev'">
                  <div
                    v-if="activeItem"
                    :key="activeItem.src"
                    :class="[
                      effectiveMode === 'original' && overflow
                        ? 'relative shrink-0'
                        : 'flex h-full w-full items-center justify-center'
                    ]"
                  >
                    <div :style="frameStyle" class="relative shrink-0 overflow-visible">
                      <img
                        ref="imageRef"
                        :src="activeItem.src"
                        :alt="activeItem.alt"
                        :class="[
                          'pointer-events-auto block transition-opacity duration-300 ease-out select-none',
                          effectiveMode === 'fit' || !natural.width || !natural.height
                            ? 'cursor-zoom-in object-contain'
                            : fits
                              ? 'max-w-full cursor-default'
                              : 'max-w-none cursor-zoom-out',
                          loaded ? 'opacity-100' : 'opacity-0'
                        ]"
                        :style="imageStyle"
                        decoding="async"
                        @load="onImageLoad"
                        @click.stop="toggleZoom"
                      />
                    </div>
                  </div>
                </Transition>
                <div
                  v-if="!loaded"
                  class="pointer-events-none absolute inset-0 flex items-center justify-center"
                >
                  <UiLoader text="" />
                </div>
              </div>
            </div>
          </div>

          <span
            :class="[
              'pointer-events-none absolute inset-y-0 left-0 z-20 w-6 bg-(--bg) opacity-0 transition-opacity',
              viewportEdges.left && !overlayNav && 'opacity-90'
            ]"
            aria-hidden="true"
          />
          <span
            :class="[
              'pointer-events-none absolute inset-y-0 right-0 z-20 w-6 bg-(--bg) opacity-0 transition-opacity',
              viewportEdges.right && !overlayNav && 'opacity-90'
            ]"
            aria-hidden="true"
          />
          <span
            :class="[
              'pointer-events-none absolute inset-x-0 top-0 z-20 h-6 bg-(--bg) opacity-0 transition-opacity',
              viewportEdges.top && 'opacity-90'
            ]"
            aria-hidden="true"
          />
          <span
            :class="[
              'pointer-events-none absolute inset-x-0 bottom-0 z-20 h-6 bg-(--bg) opacity-0 transition-opacity',
              viewportEdges.bottom && 'opacity-90'
            ]"
            aria-hidden="true"
          />
        </div>

        <div
          v-if="multiple || credits || displayTitle"
          class="relative bg-(--bg) px-3 pt-0 pb-3 ring-1 ring-(--border-color) sm:px-4"
          @click.stop
        >
          <div v-if="multiple" ref="thumbsRef" class="ui-image-scroll min-w-0 overflow-x-auto pb-2">
            <div class="flex min-w-max gap-2">
              <button
                v-for="(item, i) in items"
                :key="`${item.src}:${i}`"
                type="button"
                :data-thumb-index="i"
                :aria-label="`Открыть изображение ${i + 1}`"
                :aria-current="i === activeIndex ? 'true' : undefined"
                :class="[
                  'ui-focus relative h-16 w-16 shrink-0 overflow-visible transition sm:h-18 sm:w-18',
                  i === activeIndex
                    ? 'opacity-100 ring-2 ring-(--accent)'
                    : 'opacity-70 hover:ring-2 hover:ring-(--accent)'
                ]"
                @click="setIndex(i)"
              >
                <span
                  :class="[
                    'absolute inset-x-0 top-0 h-1',
                    i === activeIndex ? 'bg-(--accent)' : 'bg-transparent'
                  ]"
                  aria-hidden="true"
                />
                <img
                  :src="item.thumbnailSrc || item.src"
                  :alt="item.alt"
                  class="h-full w-full object-cover pt-1.5"
                  loading="lazy"
                  decoding="async"
                />
              </button>
            </div>
          </div>
          <span
            :class="[
              'pointer-events-none absolute top-0 bottom-2 left-0 w-6 bg-(--bg) opacity-0 transition-opacity',
              multiple && thumbEdges.left && 'opacity-90'
            ]"
            aria-hidden="true"
          />
          <span
            :class="[
              'pointer-events-none absolute top-0 right-0 bottom-2 w-6 bg-(--bg) opacity-0 transition-opacity',
              multiple && thumbEdges.right && 'opacity-90'
            ]"
            aria-hidden="true"
          />
          <UiCredits
            v-if="credits"
            :title="displayTitle"
            :author="activeItem?.author"
            :attribution="activeItem?.attribution"
            :source-url="activeItem?.sourceUrl"
            :license-url="activeItem?.licenseUrl"
            :license="activeItem?.license"
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.ui-image-scroll {
  scrollbar-width: none;
}
.ui-image-scroll::-webkit-scrollbar {
  display: none;
}
.ui-image-next-enter-active,
.ui-image-next-leave-active,
.ui-image-prev-enter-active,
.ui-image-prev-leave-active {
  transition:
    transform 300ms ease,
    opacity 300ms ease;
}
.ui-image-next-enter-from,
.ui-image-prev-enter-from,
.ui-image-next-leave-to,
.ui-image-prev-leave-to {
  opacity: 0;
}
.ui-image-next-leave-active,
.ui-image-prev-leave-active {
  position: absolute;
  inset: 0;
}
.ui-image-next-enter-from,
.ui-image-prev-leave-to {
  transform: translate3d(96px, 0, 0);
}
.ui-image-prev-enter-from,
.ui-image-next-leave-to {
  transform: translate3d(-96px, 0, 0);
}
</style>
