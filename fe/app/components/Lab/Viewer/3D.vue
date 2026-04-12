<script setup lang="ts">
  defineOptions({
    inheritAttrs: false
  })
  const props = withDefaults(
    defineProps<{
      title?: string
      posterSrc?: string | null
      modelSrc?: string | null
      rotationPerSecond?: string | null
      compact?: boolean
    }>(),
    {
      title: '3D модель',
      posterSrc: '',
      modelSrc: '',
      rotationPerSecond: '',
      compact: false
    }
  )
  const { effectiveTheme } = useInterfacePreferences()
  const inlineFrameRef = ref<HTMLIFrameElement | null>(null)
  const fullscreenFrameRef = ref<HTMLIFrameElement | null>(null)
  const viewerBackground = computed(() => (effectiveTheme.value === 'light' ? 'transparent' : 'transparent'))
  const fullscreenOpen = ref(false)
  const webglSupported = ref(true)
  const webglChecked = ref(false)
  const normalizedModelSrc = computed(() => String(props.modelSrc || '').trim())
  const normalizedRotationPerSecond = computed(() => String(props.rotationPerSecond || '').trim())
  const iframeClass = 'block h-full w-full border-0 bg-transparent outline-none'
  const detectWebGL = (): boolean => {
    if (!import.meta.client) return true
    try {
      const canvas = document.createElement('canvas')
      const gl2 = canvas.getContext('webgl2', {
        antialias: false,
        alpha: true,
        depth: true,
        stencil: false,
        powerPreference: 'default'
      })
      if (gl2) return true
      const gl1 =
        canvas.getContext('webgl', {
          antialias: false,
          alpha: true,
          depth: true,
          stencil: false,
          powerPreference: 'default'
        }) || canvas.getContext('experimental-webgl')
      return !!gl1
    } catch {
      return false
    }
  }
  const downloadableModelSrc = computed(() => {
    return normalizedModelSrc.value
  })
  const modelFilename = computed(() => {
    const raw = normalizedModelSrc.value
    if (!raw) return 'model.glb'
    try {
      const base = import.meta.client ? window.location.origin : 'http://127.0.0.1'
      const url = new URL(raw, base)
      const name = url.pathname.split('/').filter(Boolean).pop() || 'model.glb'
      return decodeURIComponent(name)
    } catch {
      return 'model.glb'
    }
  })
  const viewerSrc = computed(() => {
    if (!normalizedModelSrc.value || !webglSupported.value) return ''
    const params = new URLSearchParams()
    params.set('src', normalizedModelSrc.value)
    if (props.title) params.set('title', props.title)
    if (normalizedRotationPerSecond.value) {
      params.set('rotationPerSecond', normalizedRotationPerSecond.value)
    }
    params.set('theme', effectiveTheme.value)
    params.set('background', viewerBackground.value)
    return `/api/3d-model-viewer?${params.toString()}`
  })
  const postViewerTheme = (frame: HTMLIFrameElement | null) => {
    if (!import.meta.client) return
    if (!frame?.contentWindow) return
    frame.contentWindow.postMessage(
      {
        type: 'lab-3d-viewer-theme',
        theme: effectiveTheme.value,
        background: viewerBackground.value
      },
      window.location.origin
    )
  }
  const syncViewerTheme = () => {
    postViewerTheme(inlineFrameRef.value)
    postViewerTheme(fullscreenFrameRef.value)
  }
  const frameHeightClass = computed(() => {
    return props.compact ? 'h-full min-h-0' : 'h-80 md:h-90'
  })
  const closeFullscreen = () => {
    fullscreenOpen.value = false
  }
  const openFullscreen = () => {
    if (!viewerSrc.value) return
    fullscreenOpen.value = true
  }
  const onKeydown = (event: KeyboardEvent) => {
    if (!fullscreenOpen.value) return
    if (event.key === 'Escape') closeFullscreen()
  }
  const onWindowMessage = (event: MessageEvent) => {
    if (!fullscreenOpen.value) return
    if (event.origin !== window.location.origin) return
    if (event.data?.type !== 'lab-3d-viewer-close') return
    closeFullscreen()
  }
  if (import.meta.client) {
    onMounted(() => {
      webglSupported.value = detectWebGL()
      webglChecked.value = true
      window.addEventListener('keydown', onKeydown)
      window.addEventListener('message', onWindowMessage)
    })
    onBeforeUnmount(() => {
      window.removeEventListener('keydown', onKeydown)
      window.removeEventListener('message', onWindowMessage)
    })
  }
  if (import.meta.client) {
    watch(
      () => effectiveTheme.value,
      () => {
        syncViewerTheme()
      }
    )

    watch(
      () => fullscreenOpen.value,
      opened => {
        if (!opened) return
        nextTick(() => {
          syncViewerTheme()
        })
      }
    )
  }
</script>
<template>
  <div
    v-bind="$attrs"
    class="overflow-hidden"
    :class="props.compact ? 'h-full bg-(--lab-bg-canvas)' : 'border bg-(--lab-bg-surface)'">
    <div v-if="!props.compact" class="flex items-center justify-between gap-3 border-b px-3 py-2">
      <div class="text-[10px] uppercase tracking-widest opacity-60">3D Viewer</div>
      <div class="flex items-center gap-2">
        <LabViewerModelDownloadLink
          v-if="downloadableModelSrc"
          :href="downloadableModelSrc"
          :filename="modelFilename"
          class-name="border text-xs font-medium transition" />
        <LabBaseButton
          v-if="viewerSrc"
          icon-only
          icon="ic:round-open-in-full"
          button-class="inline-flex h-8 w-8 items-center justify-center rounded-full border transition"
          aria-label="Открыть просмотр модели на весь экран"
          title="Открыть на весь экран"
          @click="openFullscreen" />
      </div>
    </div>
    <div class="relative bg-(--lab-bg-canvas)" :class="frameHeightClass">
      <div v-if="props.compact" class="absolute right-3 top-3 z-10 flex items-center gap-2">
        <LabViewerModelDownloadLink
          v-if="downloadableModelSrc && webglChecked && webglSupported"
          :href="downloadableModelSrc"
          :filename="modelFilename"
          icon-only
          class-name="rounded-full border bg-[color-mix(in_srgb,var(--lab-bg-overlay)_80%,transparent)] text-(--lab-text-primary) transition" />
        <LabBaseButton
          v-if="viewerSrc"
          icon-only
          icon="ic:round-open-in-full"
          button-class="h-8 w-8 rounded-full border bg-[color-mix(in_srgb,var(--lab-bg-overlay)_80%,transparent)] transition"
          aria-label="Открыть просмотр модели на весь экран"
          title="Открыть на весь экран"
          @click="openFullscreen" />
      </div>
      <div v-if="!webglChecked" class="flex h-full w-full items-center justify-center p-2 text-sm opacity-70">
        Проверка поддержки WebGL…
      </div>
      <div
        v-else-if="!webglSupported"
        class="flex h-full w-full flex-col items-center justify-center gap-3 px-6 text-center">
        <Icon name="ic:round-warning-amber" class="h-8 w-8 text-amber-400" />
        <div class="text-sm font-medium">3D-просмотр недоступен</div>
        <p class="max-w-md text-sm leading-6 opacity-70">
          В браузере отключён WebGL или аппаратное ускорение графики. Включите аппаратное ускорение и перезапустите
          браузер.
        </p>
      </div>
      <iframe
        v-else-if="viewerSrc"
        ref="inlineFrameRef"
        :src="viewerSrc"
        :title="`${title} 3D previewer`"
        :class="iframeClass"
        loading="lazy"
        referrerpolicy="same-origin"
        allowfullscreen></iframe>
    </div>
  </div>
  <Teleport to="body">
    <div v-if="fullscreenOpen && viewerSrc" class="fixed inset-0 z-70 bg-(--lab-bg-canvas)">
      <div class="flex h-full flex-col">
        <div class="z-10">
          <div
            class="flex max-w-full flex-wrap items-center justify-end gap-2.5 border-b bg-(--lab-bg-surface) px-2.5 py-1.5 sm:px-3">
            <p class="min-w-0 flex-1 truncate text-sm font-semibold sm:text-base">
              {{ title }}
            </p>
            <LabViewerModelDownloadLink
              v-if="downloadableModelSrc"
              :href="downloadableModelSrc"
              :filename="modelFilename"
              class-name="border bg-(--lab-bg-surface) text-(--lab-text-primary) transition" />
            <LabBaseButton
              icon="ic:round-close"
              label="Закрыть"
              button-class="h-8 shrink-0 border bg-(--lab-bg-surface) px-2.5 text-sm transition"
              @click="closeFullscreen" />
          </div>
        </div>
        <div class="relative min-h-0 flex-1 bg-(--lab-bg-canvas)">
          <iframe
            ref="fullscreenFrameRef"
            :src="viewerSrc"
            :title="`${title} fullscreen viewer`"
            :class="`absolute inset-0 ${iframeClass}`"
            loading="eager"
            referrerpolicy="same-origin"
            allowfullscreen></iframe>
        </div>
      </div>
    </div>
  </Teleport>
</template>
