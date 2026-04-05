<template>
  <div v-if="isVisible" class="viewer-shell fixed inset-0 z-60">
    <div class="flex h-full w-full flex-col">
      <div class="z-10" @click.stop>
        <div class="viewer-toolbar flex max-w-full flex-wrap items-center justify-end gap-2.5 px-2.5 py-1.5 sm:px-3">
          <div v-if="hasMultipleItems" class="viewer-counter px-2.5 py-1 text-xs font-medium tabular-nums">
            {{ activeCounterLabel }}
          </div>
          <div class="flex shrink-0 flex-wrap items-center justify-end gap-1.5">
            <LabBaseButton
              icon="ic:round-swap-horiz"
              icon-only
              title="Отразить по горизонтали"
              aria-label="Отразить по горизонтали"
              variant="secondary"
              button-class="viewer-action h-8 shrink-0 px-2.5 text-sm"
              @click="toggleFlipX" />
            <LabBaseButton
              icon="ic:round-swap-vert"
              icon-only
              title="Отразить по вертикали"
              aria-label="Отразить по вертикали"
              variant="secondary"
              button-class="viewer-action h-8 shrink-0 px-2.5 text-sm"
              @click="toggleFlipY" />
            <LabBaseButton
              icon="ic:round-rotate-left"
              icon-only
              title="Повернуть против часовой стрелки на 90 градусов"
              aria-label="Повернуть против часовой стрелки на 90 градусов"
              variant="secondary"
              button-class="viewer-action h-8 shrink-0 px-2.5 text-sm"
              @click="rotateLeft" />
            <LabBaseButton
              icon="ic:round-rotate-right"
              icon-only
              title="Повернуть по часовой стрелке на 90 градусов"
              aria-label="Повернуть по часовой стрелке на 90 градусов"
              variant="secondary"
              button-class="viewer-action h-8 shrink-0 px-2.5 text-sm"
              @click="rotateRight" />
          </div>
          <span class="viewer-divider hidden h-5 w-px shrink-0 sm:block" aria-hidden="true"></span>
          <div class="viewer-meta flex shrink-0 items-center gap-1.5 text-sm font-medium">
            масштаб:
            <p class="viewer-scale-value tabular-nums">
              {{ scaleLabel }}
            </p>
          </div>
          <slot name="toolbar-extra" :active-index="activeIndex" :active-item="activeItem" :items="resolvedItems" />
          <LabBaseButton
            icon="ic:round-close"
            icon-only
            title="Закрыть окно"
            aria-label="Закрыть окно"
            variant="secondary"
            @click="close" />
        </div>
      </div>
      <div :class="viewportClass">
        <div class="pointer-events-none absolute inset-y-0 left-0 z-10 flex items-stretch">
          <LabBaseButton
            v-if="showOverlayNavButtons"
            icon="ic:round-navigate-before"
            size="xl"
            icon-only
            aria-label="Предыдущее изображение"
            variant="secondary"
            button-class="viewer-nav pointer-events-auto h-full w-16 cursor-pointer sm:w-32 rounded-none"
            @click.stop="showPreviousImage" />
        </div>
        <div class="pointer-events-none absolute inset-y-0 right-0 z-10 flex items-stretch">
          <LabBaseButton
            v-if="showOverlayNavButtons"
            icon="ic:round-navigate-next"
            size="xl"
            icon-only
            aria-label="Следующее изображение"
            variant="secondary"
            button-class="viewer-nav pointer-events-auto h-full w-16 cursor-pointer sm:w-32 rounded-none"
            @click.stop="showNextImage" />
        </div>
        <div
          ref="viewportRef"
          :class="contentViewportClass"
          @click="onViewportBackdropClick"
          @touchstart.passive="onViewportTouchStart"
          @touchmove.passive="onViewportTouchMove"
          @touchend="onViewportTouchEnd"
          @touchcancel="resetViewportTouchState">
          <div :class="contentClass" :style="contentStyle">
            <div ref="stageRef" :class="stageClass">
              <Transition :name="slideTransitionName">
                <div v-if="activeItem" :key="activeItem.src" :class="stageItemClass">
                  <div :style="imageFrameStyle" class="relative shrink-0 overflow-visible">
                    <img
                      ref="imageRef"
                      :src="activeItem.src"
                      :alt="activeItem.alt"
                      :class="[imageClass, imageLoaded ? 'opacity-100' : 'opacity-0']"
                      :style="imageStyle"
                      decoding="async"
                      @load="onImageLoad"
                      @click.stop="toggleZoom" />
                  </div>
                </div>
              </Transition>
              <div v-if="!imageLoaded" class="pointer-events-none absolute inset-0 flex items-center justify-center">
                <span class="viewer-loader inline-flex p-3">
                  <LabLoader variant="icon" />
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        v-if="hasMultipleItems || activeItemHasCredits || displayTitle"
        class="viewer-footer px-3 py-3 sm:px-4"
        @click.stop>
        <div v-if="hasMultipleItems" class="flex items-center gap-2">
          <div class="min-w-0 flex-1 overflow-x-auto">
            <div class="flex min-w-max gap-2">
              <button
                v-for="(item, index) in resolvedItems"
                :key="`${item.src}:${index}`"
                type="button"
                class="viewer-thumb relative h-16 w-16 shrink-0 overflow-hidden rounded-2xl border sm:h-18 sm:w-18"
                :class="index === activeIndex ? 'viewer-thumb-active' : 'viewer-thumb-inactive'"
                :aria-label="`Открыть изображение ${index + 1}`"
                :aria-current="index === activeIndex ? 'true' : undefined"
                @click="setActiveIndex(index)">
                <img
                  :src="item.thumbnailSrc || item.src"
                  :alt="item.alt"
                  class="h-full w-full object-cover"
                  loading="lazy"
                  decoding="async" />
                <div
                  class="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent"
                  :class="index === activeIndex ? 'opacity-100' : 'opacity-65'"></div>
                <div class="absolute inset-x-1.5 bottom-1.5 flex items-center justify-between gap-1">
                  <span class="viewer-thumb-label truncate text-[10px] font-semibold">
                    {{ index + 1 }}
                  </span>
                  <span
                    class="viewer-thumb-dot h-2.5 w-2.5 rounded-full border"
                    :class="index === activeIndex ? 'viewer-thumb-dot-active' : 'viewer-thumb-dot-inactive'"
                    aria-hidden="true"></span>
                </div>
              </button>
            </div>
          </div>
        </div>
        <LabViewerImageCredits
          v-if="activeItemHasCredits"
          :title="displayTitle"
          :author="activeItem?.author"
          :attribution="activeItem?.attribution"
          :source-url="activeItem?.sourceUrl"
          :license-url="activeItem?.licenseUrl"
          :license="activeItem?.license" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      modelValue?: boolean
      src?: string
      title?: string
      items?: ImageViewerItem[]
      initialIndex?: number
      routeQueryParam?: string
      routeQuerySync?: boolean
    }>(),
    {
      modelValue: false,
      src: '',
      title: '',
      items: () => [],
      initialIndex: 0,
      routeQueryParam: 'image',
      routeQuerySync: true
    }
  )
  const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'opened'): void
    (e: 'closed'): void
    (e: 'mode-change', mode: ViewerMode): void
    (e: 'active-index-change', index: number): void
  }>()
  const normalizeViewerText = (value: unknown) => String(value ?? '').trim() || null
  const route = useRoute()
  const router = useRouter()
  const normalizeIndex = (value: number, length: number) => {
    if (!length) return 0
    if (!Number.isFinite(value)) return 0
    const next = Math.trunc(value)
    return ((next % length) + length) % length
  }
  const resolvedItems = computed(() => {
    const items = props.items
    const normalizedItems: Array<{
      src: string
      title: string | null
      alt: string
      author: string | null
      attribution: string | null
      sourceUrl: string | null
      licenseUrl: string | null
      license: string | null
      thumbnailSrc: string | null
    }> = []
    if (Array.isArray(items) && items.length > 0) {
      for (let index = 0; index < items.length; index++) {
        const item = items[index]
        const src = item?.src.trim()
        if (!src) continue
        const title = normalizeViewerText(item?.title)
        const alt = normalizeViewerText(item?.alt)
        normalizedItems.push({
          src,
          title,
          alt: alt || title || `Изображение ${index + 1}`,
          author: normalizeViewerText(item?.author),
          attribution: normalizeViewerText(item?.attribution),
          sourceUrl: normalizeViewerText(item?.sourceUrl),
          licenseUrl: normalizeViewerText(item?.licenseUrl),
          license: normalizeViewerText(item?.license),
          thumbnailSrc: normalizeViewerText(item?.thumbnailSrc)
        })
      }
    }
    if (normalizedItems.length > 0) return normalizedItems
    const fallbackSrc = props.src.trim()
    if (!fallbackSrc) return []
    const fallbackTitle = normalizeViewerText(props.title)
    return [
      {
        src: fallbackSrc,
        title: fallbackTitle,
        alt: fallbackTitle || 'Изображение',
        author: null,
        attribution: null,
        sourceUrl: null,
        licenseUrl: null,
        license: null,
        thumbnailSrc: null
      }
    ]
  })
  const activeIndex = ref(0)
  const requestedRouteIndex = computed(() => {
    if (!props.routeQuerySync) return null
    if (!resolvedItems.value.length) return null
    const key = String(props.routeQueryParam || '').trim()
    if (!key) return null
    const raw = Array.isArray(route.query[key]) ? route.query[key][0] : route.query[key]
    const parsed = Number.parseInt(String(raw || '').trim(), 10)
    if (!Number.isFinite(parsed) || parsed <= 0) return null
    return normalizeIndex(parsed - 1, resolvedItems.value.length)
  })
  const activeItem = computed(() => resolvedItems.value[activeIndex.value] || null)
  const activeItemHasCredits = computed(() =>
    Boolean(
      activeItem.value?.title ||
      activeItem.value?.author ||
      activeItem.value?.attribution ||
      activeItem.value?.sourceUrl ||
      activeItem.value?.license
    )
  )
  const hasMultipleItems = computed(() => resolvedItems.value.length > 1)
  const showOverlayNavButtons = computed(() => hasMultipleItems.value && viewportSize.width >= 768)
  const isVisible = computed(() =>
    Boolean((props.modelValue || requestedRouteIndex.value !== null) && activeItem.value?.src)
  )
  const displayTitle = computed(() => activeItem.value?.title || normalizeViewerText(props.title))
  const activeCounterLabel = computed(() => `${activeIndex.value + 1} / ${resolvedItems.value.length}`)
  const imageNatural = reactive({
    width: 0,
    height: 0
  })
  const viewportRef = ref<HTMLElement | null>(null)
  const stageRef = ref<HTMLElement | null>(null)
  const imageRef = ref<HTMLImageElement | null>(null)
  const VIEWER_TRANSITION_MS = 300
  const viewportSize = reactive({
    width: 0,
    height: 0
  })
  const stageSize = reactive({
    width: 0,
    height: 0
  })
  const mode = ref<ViewerMode>('fit')
  const imageLoaded = ref(false)
  const navigationDirection = ref<-1 | 1>(-1)
  const flipX = ref(1)
  const flipY = ref(1)
  const rotation = ref(0)
  const effectiveMode = computed<ViewerMode>(() => {
    if (fitsViewport.value) return 'original'
    return mode.value
  })
  const scalePercent = computed(() => {
    const scale = effectiveMode.value === 'original' ? 1 : fitScale.value
    return Math.round(scale * 100)
  })
  const scaleLabel = computed(() => {
    if (!imageLoaded.value) return '...'
    return `${scalePercent.value}%`
  })
  const isQuarterTurn = computed(() => Math.abs(rotation.value / 90) % 2 === 1)
  const orientedNaturalSize = computed(() => {
    if (isQuarterTurn.value) {
      return {
        width: imageNatural.height,
        height: imageNatural.width
      }
    }
    return {
      width: imageNatural.width,
      height: imageNatural.height
    }
  })
  const VIEWPORT_PADDING_X = 80
  const VIEWPORT_PADDING_Y = 24
  const availableViewportSize = computed(() => {
    const extraX = hasMultipleItems.value ? 96 : 0
    const extraY = hasMultipleItems.value || activeItemHasCredits.value || displayTitle.value ? 12 : 0
    return {
      width: Math.max(0, viewportSize.width - VIEWPORT_PADDING_X * 2 - extraX),
      height: Math.max(0, viewportSize.height - VIEWPORT_PADDING_Y * 2 - extraY)
    }
  })
  const fitsViewport = computed(() => {
    if (
      !orientedNaturalSize.value.width ||
      !orientedNaturalSize.value.height ||
      !availableViewportSize.value.width ||
      !availableViewportSize.value.height
    )
      return false
    return (
      orientedNaturalSize.value.width <= availableViewportSize.value.width &&
      orientedNaturalSize.value.height <= availableViewportSize.value.height
    )
  })
  const fitScale = computed(() => {
    if (
      !orientedNaturalSize.value.width ||
      !orientedNaturalSize.value.height ||
      !availableViewportSize.value.width ||
      !availableViewportSize.value.height
    )
      return 1
    return Math.min(
      1,
      availableViewportSize.value.width / orientedNaturalSize.value.width,
      availableViewportSize.value.height / orientedNaturalSize.value.height
    )
  })
  const hasOriginalOverflow = computed(() => {
    if (
      !orientedNaturalSize.value.width ||
      !orientedNaturalSize.value.height ||
      !availableViewportSize.value.width ||
      !availableViewportSize.value.height
    )
      return false
    return (
      orientedNaturalSize.value.width > availableViewportSize.value.width ||
      orientedNaturalSize.value.height > availableViewportSize.value.height
    )
  })
  const canPanImage = computed(() => effectiveMode.value === 'original' && hasOriginalOverflow.value)
  const viewportClass = computed(() => {
    return ['relative min-h-0 flex-1 overflow-hidden']
  })
  const contentViewportClass = computed(() => [
    'h-full w-full',
    canPanImage.value ? 'overflow-auto' : 'overflow-hidden'
  ])
  const contentClass = computed(() => {
    const baseClass = [
      'box-border flex min-h-full min-w-full px-16 py-4 sm:px-20 sm:py-6',
      hasMultipleItems.value ? 'px-20 sm:px-24' : '',
      hasMultipleItems.value || activeItemHasCredits.value ? 'py-2 sm:py-3' : ''
    ]
      .filter(Boolean)
      .join(' ')
    if (effectiveMode.value === 'original' && hasOriginalOverflow.value) {
      return `${baseClass} h-max w-max items-center justify-center`
    }
    return `${baseClass} h-full w-full items-center justify-center`
  })
  const stageClass = computed(() => {
    if (effectiveMode.value === 'original' && hasOriginalOverflow.value) {
      return 'relative shrink-0'
    }
    return 'relative flex h-full w-full items-center justify-center overflow-hidden'
  })
  const stageItemClass = computed(() => {
    if (effectiveMode.value === 'original' && hasOriginalOverflow.value) {
      return 'relative shrink-0'
    }
    return 'flex h-full w-full items-center justify-center'
  })
  const swipeEnabled = computed(() => hasMultipleItems.value && !canPanImage.value)
  const imageFrameStyle = computed<CSSProperties>(() => {
    if (!imageNatural.width || !imageNatural.height) return {}
    const scale = effectiveMode.value === 'fit' ? fitScale.value : 1
    return {
      width: `${Math.max(1, Math.round(orientedNaturalSize.value.width * scale))}px`,
      height: `${Math.max(1, Math.round(orientedNaturalSize.value.height * scale))}px`,
      flex: '0 0 auto'
    }
  })
  const imageClass = computed(() => {
    if (effectiveMode.value === 'fit' || !imageNatural.width || !imageNatural.height) {
      return 'pointer-events-auto block object-contain select-none transition-opacity duration-300 ease-out cursor-zoom-in'
    }
    return [
      'pointer-events-auto block select-none transition-opacity duration-300 ease-out',
      fitsViewport.value ? 'max-w-full cursor-default' : 'max-w-none cursor-zoom-out'
    ].join(' ')
  })
  const imageStyle = computed<CSSProperties>(() => {
    if (!imageNatural.width || !imageNatural.height) {
      return {}
    }
    const scale = effectiveMode.value === 'fit' ? fitScale.value : 1
    return {
      position: 'absolute',
      left: '50%',
      top: '50%',
      width: `${Math.max(1, Math.round(imageNatural.width * scale))}px`,
      height: `${Math.max(1, Math.round(imageNatural.height * scale))}px`,
      maxWidth: 'none',
      maxHeight: 'none',
      transform: `translate(-50%, -50%) rotate(${rotation.value}deg) scale(${flipX.value}, ${flipY.value})`,
      transformOrigin: 'center center'
    }
  })
  const contentStyle = computed(() => {
    const transitionDuration = viewportTouchState.animating ? `${VIEWER_TRANSITION_MS}ms` : '0ms'
    const opacity = swipeEnabled.value
      ? String(Math.max(0.58, 1 - Math.abs(viewportTouchState.deltaX) / Math.max(1, viewportSize.width * 1.15)))
      : '1'
    return {
      transform: `translate3d(${viewportTouchState.deltaX}px, 0, 0)`,
      opacity,
      transition: `transform ${transitionDuration} ease, opacity ${transitionDuration} ease`
    }
  })
  const slideTransitionName = computed(() => {
    return navigationDirection.value < 0 ? 'viewer-slide-next' : 'viewer-slide-prev'
  })
  const resetMode = () => {
    const nextMode: ViewerMode = imageNatural.width && imageNatural.height && fitsViewport.value ? 'original' : 'fit'
    if (mode.value === nextMode) return
    mode.value = nextMode
    emit('mode-change', mode.value)
  }
  const syncActiveIndex = () => {
    const nextIndex =
      requestedRouteIndex.value !== null
        ? requestedRouteIndex.value
        : normalizeIndex(props.initialIndex, resolvedItems.value.length)
    if (activeIndex.value === nextIndex) return
    activeIndex.value = nextIndex
  }
  const syncRouteQuery = async (index: number | null) => {
    if (!props.routeQuerySync || !import.meta.client) return
    const key = String(props.routeQueryParam || '').trim()
    if (!key) return
    const currentValue = Array.isArray(route.query[key]) ? route.query[key][0] : route.query[key]
    const nextValue = index === null ? null : String(index + 1)
    if ((currentValue == null || String(currentValue).trim() === '') && nextValue === null) {
      return
    }
    if (nextValue !== null && String(currentValue || '').trim() === nextValue) {
      return
    }
    const nextQuery = { ...route.query }
    if (nextValue === null) {
      delete nextQuery[key]
    } else {
      nextQuery[key] = nextValue
    }
    await router.replace({
      path: route.path,
      query: nextQuery
    })
  }
  const resetActiveImageState = () => {
    cancelRevealFrame()
    imageLoaded.value = false
    imageNatural.width = 0
    imageNatural.height = 0
    flipX.value = 1
    flipY.value = 1
    rotation.value = 0
    if (props.modelValue) {
      nextTick(updateViewportSize)
      resetMode()
    }
  }
  const setActiveIndex = (value: number, direction?: -1 | 1) => {
    const nextIndex = normalizeIndex(value, resolvedItems.value.length)
    if (activeIndex.value === nextIndex) return
    navigationDirection.value = direction || (nextIndex > activeIndex.value ? -1 : 1)
    activeIndex.value = nextIndex
    emit('active-index-change', nextIndex)
    void syncRouteQuery(nextIndex)
  }
  const showPreviousImage = () => {
    if (!hasMultipleItems.value) return
    setActiveIndex(activeIndex.value - 1, 1)
  }
  const showNextImage = () => {
    if (!hasMultipleItems.value) return
    setActiveIndex(activeIndex.value + 1, -1)
  }
  const close = () => {
    void syncRouteQuery(null)
    emit('update:modelValue', false)
  }
  const viewportTouchState = reactive({
    startX: 0,
    startY: 0,
    deltaX: 0,
    deltaY: 0,
    tracking: false,
    swiped: false,
    animating: false
  })
  let swipeResetTimer: number | null = null
  const resetViewportTouchState = () => {
    if (swipeResetTimer !== null && import.meta.client) {
      window.clearTimeout(swipeResetTimer)
      swipeResetTimer = null
    }
    viewportTouchState.startX = 0
    viewportTouchState.startY = 0
    viewportTouchState.deltaX = 0
    viewportTouchState.deltaY = 0
    viewportTouchState.tracking = false
    viewportTouchState.swiped = false
    viewportTouchState.animating = false
  }
  const onViewportBackdropClick = () => {
    if (viewportTouchState.swiped) {
      viewportTouchState.swiped = false
      return
    }
    close()
  }
  const onViewportTouchStart = (event: TouchEvent) => {
    if (!swipeEnabled.value) {
      resetViewportTouchState()
      return
    }
    const touch = event.touches[0]
    if (!touch) return
    viewportTouchState.startX = touch.clientX
    viewportTouchState.startY = touch.clientY
    viewportTouchState.deltaX = 0
    viewportTouchState.deltaY = 0
    viewportTouchState.tracking = true
    viewportTouchState.swiped = false
    viewportTouchState.animating = false
  }
  const onViewportTouchMove = (event: TouchEvent) => {
    if (!viewportTouchState.tracking || !swipeEnabled.value) return
    const touch = event.touches[0]
    if (!touch) return
    viewportTouchState.deltaX = touch.clientX - viewportTouchState.startX
    viewportTouchState.deltaY = touch.clientY - viewportTouchState.startY
  }
  const onViewportTouchEnd = () => {
    if (!viewportTouchState.tracking || !swipeEnabled.value) {
      resetViewportTouchState()
      return
    }
    const absDeltaX = Math.abs(viewportTouchState.deltaX)
    const absDeltaY = Math.abs(viewportTouchState.deltaY)
    const direction = viewportTouchState.deltaX < 0 ? -1 : 1
    if (absDeltaX >= 56 && absDeltaX > absDeltaY) {
      viewportTouchState.swiped = true
      viewportTouchState.animating = true
      navigationDirection.value = direction
      viewportTouchState.deltaX = direction * Math.max(Math.round(viewportSize.width * 1.35), 420)
      const navigate = () => {
        if (direction < 0) {
          showNextImage()
        } else {
          showPreviousImage()
        }
        viewportTouchState.deltaX = 0
        viewportTouchState.deltaY = 0
        viewportTouchState.animating = false
      }
      if (import.meta.client) {
        window.setTimeout(navigate, VIEWER_TRANSITION_MS)
      } else {
        navigate()
      }
      if (import.meta.client) {
        swipeResetTimer = window.setTimeout(() => {
          viewportTouchState.swiped = false
          swipeResetTimer = null
        }, VIEWER_TRANSITION_MS + 120)
      }
    } else {
      viewportTouchState.animating = true
      viewportTouchState.deltaX = 0
      viewportTouchState.deltaY = 0
      if (import.meta.client) {
        window.setTimeout(() => {
          viewportTouchState.animating = false
        }, VIEWER_TRANSITION_MS)
      } else {
        viewportTouchState.animating = false
      }
    }
    viewportTouchState.tracking = false
  }
  const onImageLoad = (event: Event) => {
    const img = event.target as HTMLImageElement | null
    if (!img) return
    imageNatural.width = img.naturalWidth
    imageNatural.height = img.naturalHeight
    resetMode()
    nextTick(alignOriginalPosition)
    cancelRevealFrame()
    revealFrameId = window.requestAnimationFrame(() => {
      imageLoaded.value = true
      revealFrameId = null
    })
  }
  const syncLoadedImageState = () => {
    const img = imageRef.value
    if (!img?.complete) return
    if (!img.naturalWidth && !img.naturalHeight) return
    imageNatural.width = img.naturalWidth
    imageNatural.height = img.naturalHeight
    resetMode()
    alignOriginalPosition()
    cancelRevealFrame()
    imageLoaded.value = true
  }
  const toggleZoom = () => {
    if (fitsViewport.value) return
    mode.value = mode.value === 'fit' ? 'original' : 'fit'
    emit('mode-change', mode.value)
    nextTick(alignOriginalPosition)
  }
  const toggleFlipX = () => {
    flipX.value *= -1
  }
  const toggleFlipY = () => {
    flipY.value *= -1
  }
  const rotateBy = (delta: number) => {
    rotation.value = (((rotation.value + delta) % 360) + 360) % 360
    nextTick(() => {
      resetMode()
      alignOriginalPosition()
    })
  }
  const rotateLeft = () => {
    rotateBy(-90)
  }
  const rotateRight = () => {
    rotateBy(90)
  }
  const updateViewportSize = () => {
    const viewport = viewportRef.value
    if (!viewport) return
    viewportSize.width = viewport.clientWidth
    viewportSize.height = viewport.clientHeight
  }
  const alignOriginalPosition = () => {
    if (!import.meta.client) return
    const viewport = viewportRef.value
    if (!viewport) return
    if (effectiveMode.value !== 'original') return
    const maxScrollLeft = Math.max(0, viewport.scrollWidth - viewport.clientWidth)
    const maxScrollTop = Math.max(0, viewport.scrollHeight - viewport.clientHeight)
    viewport.scrollLeft = maxScrollLeft > 0 ? Math.round(maxScrollLeft / 2) : 0
    viewport.scrollTop = maxScrollTop > 0 ? Math.round(maxScrollTop / 2) : 0
  }
  let resizeObserver: ResizeObserver | null = null
  let revealFrameId: number | null = null
  let loadedStateSyncFrameId: number | null = null
  const cancelRevealFrame = () => {
    if (revealFrameId === null || !import.meta.client) return
    window.cancelAnimationFrame(revealFrameId)
    revealFrameId = null
  }
  const cancelLoadedStateSyncFrame = () => {
    if (loadedStateSyncFrameId === null || !import.meta.client) return
    window.cancelAnimationFrame(loadedStateSyncFrameId)
    loadedStateSyncFrameId = null
  }
  const queueLoadedImageStateSync = () => {
    nextTick(() => {
      if (!import.meta.client) {
        syncLoadedImageState()
        return
      }
      cancelLoadedStateSyncFrame()
      loadedStateSyncFrameId = window.requestAnimationFrame(() => {
        loadedStateSyncFrameId = null
        syncLoadedImageState()
      })
    })
  }
  const observeViewport = () => {
    const viewport = viewportRef.value
    if (!resizeObserver || !viewport) return
    resizeObserver.disconnect()
    resizeObserver.observe(viewport)
    updateViewportSize()
  }
  const onKeydown = (event: KeyboardEvent) => {
    if (!isVisible.value) return
    if (event.key === 'Escape') {
      close()
      return
    }
    if (event.key === 'ArrowLeft') {
      showPreviousImage()
      return
    }
    if (event.key === 'ArrowRight') {
      showNextImage()
    }
  }
  watch(
    () => props.modelValue,
    (next, prev) => {
      if (next && !prev) {
        syncActiveIndex()
        void syncRouteQuery(activeIndex.value)
        nextTick(updateViewportSize)
        nextTick(alignOriginalPosition)
        resetMode()
        emit('opened')
      }
      if (!next && prev) {
        emit('closed')
      }
    }
  )
  watch(
    () => requestedRouteIndex.value,
    next => {
      if (next === null) return
      activeIndex.value = next
      if (!props.modelValue) {
        emit('update:modelValue', true)
      }
    },
    { immediate: true }
  )
  watch(
    () => props.initialIndex,
    () => {
      syncActiveIndex()
    }
  )
  watch(
    resolvedItems,
    () => {
      syncActiveIndex()
    },
    { deep: true }
  )
  watch(
    () => activeItem.value?.src,
    () => {
      resetViewportTouchState()
      resetActiveImageState()
      queueLoadedImageStateSync()
    }
  )
  watch([fitsViewport, () => stageSize.width, () => stageSize.height], () => {
    if (!props.modelValue) return
    if (!imageNatural.width || !imageNatural.height) return
    if (fitsViewport.value && mode.value !== 'original') {
      mode.value = 'original'
      emit('mode-change', mode.value)
    }
    nextTick(alignOriginalPosition)
  })
  watch(isVisible, next => {
    if (!next) return
    nextTick(observeViewport)
    queueLoadedImageStateSync()
  })
  watch(
    [isVisible, () => activeItem.value?.src, () => imageRef.value],
    ([visible, src, img]) => {
      if (!visible || !src || !img) return
      queueLoadedImageStateSync()
    },
    { flush: 'post' }
  )
  onMounted(() => {
    if (!import.meta.client) return
    window.addEventListener('keydown', onKeydown)
    resizeObserver = new ResizeObserver(() => {
      updateViewportSize()
    })
    observeViewport()
  })
  onBeforeUnmount(() => {
    if (!import.meta.client) return
    resetViewportTouchState()
    cancelRevealFrame()
    cancelLoadedStateSyncFrame()
    window.removeEventListener('keydown', onKeydown)
    resizeObserver?.disconnect()
  })
</script>
<style scoped>
  .viewer-shell {
    background: var(--lab-bg-canvas);
  }
  .viewer-toolbar,
  .viewer-loader,
  .viewer-footer {
    background: var(--lab-bg-canvas);
  }
  .viewer-toolbar {
    border-bottom: 1px solid var(--lab-border);
  }
  .viewer-footer {
    border-top: 1px solid var(--lab-border);
  }
  .viewer-counter {
    background: var(--lab-bg-control);
    color: var(--lab-text-secondary);
  }
  .viewer-action,
  .viewer-nav {
    border-color: var(--lab-border) !important;
    background: var(--lab-bg-control) !important;
    color: var(--lab-text-primary);
  }
  .viewer-nav {
    border-top: 0 !important;
    border-bottom: 0 !important;
  }
  .viewer-nav:first-child {
    border-left: 0 !important;
  }
  .viewer-nav:last-child {
    border-right: 0 !important;
  }
  .viewer-action:hover,
  .viewer-nav:hover {
    border-color: var(--lab-border-strong) !important;
    background: var(--lab-bg-control-hover) !important;
    color: var(--lab-text-primary);
  }
  .viewer-divider {
    background: var(--lab-border);
  }
  .viewer-meta {
    color: var(--lab-text-secondary);
  }
  .viewer-loader {
    color: var(--lab-text-primary);
  }
  .viewer-scale-value {
    width: 4ch;
    text-align: right;
  }
  .viewer-thumb-active {
    border-color: var(--lab-text-primary);
    background: color-mix(in srgb, var(--lab-text-primary) 10%, transparent);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--lab-text-primary) 58%, transparent);
  }
  .viewer-thumb-inactive {
    border-color: color-mix(in srgb, var(--lab-border) 82%, transparent);
  }
  .viewer-thumb-inactive:hover {
    border-color: var(--lab-border-strong);
  }
  .viewer-thumb-label {
    color: rgb(255 255 255);
  }
  .viewer-thumb-dot {
    border-color: rgb(255 255 255);
  }
  .viewer-thumb-dot-active {
    background: var(--lab-accent);
    border-color: color-mix(in srgb, var(--lab-accent-hover) 64%, transparent);
  }
  .viewer-thumb-dot-inactive {
    background: transparent;
  }
  :deep(.viewer-action .iconify),
  :deep(.viewer-nav .iconify) {
    color: currentColor;
  }
  .viewer-slide-next-enter-active,
  .viewer-slide-next-leave-active,
  .viewer-slide-prev-enter-active,
  .viewer-slide-prev-leave-active {
    transition:
      transform 300ms ease,
      opacity 300ms ease;
  }
  .viewer-slide-next-enter-from,
  .viewer-slide-prev-enter-from {
    opacity: 0;
  }
  .viewer-slide-next-leave-active,
  .viewer-slide-prev-leave-active {
    position: absolute;
    inset: 0;
  }
  .viewer-slide-next-leave-to,
  .viewer-slide-prev-leave-to {
    opacity: 0;
  }
  .viewer-slide-next-enter-from,
  .viewer-slide-prev-leave-to {
    transform: translate3d(96px, 0, 0);
  }
  .viewer-slide-prev-enter-from,
  .viewer-slide-next-leave-to {
    transform: translate3d(-96px, 0, 0);
  }
</style>
