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
  const fullscreenOpen = ref(false)
  const webglSupported = ref(true)
  const webglChecked = ref(false)
  const normalizedModelSrc = computed(() => String(props.modelSrc || '').trim())
  const normalizedRotationPerSecond = computed(() => String(props.rotationPerSecond || '').trim())
  const iframeClass =
    'block h-full w-full border-0 bg-zinc-950 outline-none ring-0 focus:outline-none focus:ring-0 focus-visible:outline-none focus-visible:ring-0'
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
      const base = import.meta.client ? window.location.origin : 'http://localhost'
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
    if (normalizedRotationPerSecond.value) params.set('rotationPerSecond', normalizedRotationPerSecond.value)
    return `/api/3d-model-viewer?${params.toString()}`
  })
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
</script>
<template>
  <div
    v-bind="$attrs"
    class="overflow-hidden"
    :class="props.compact ? 'h-full bg-black/20' : 'border border-zinc-800 bg-zinc-900/80'">
    <div v-if="!props.compact" class="flex items-center justify-between gap-3 border-b border-zinc-800 px-3 py-2">
        <div class="text-[10px] uppercase tracking-widest text-zinc-500">3D Viewer</div>
        <div class="flex items-center gap-2">
          <LabViewerModelDownloadLink
            v-if="downloadableModelSrc"
            :href="downloadableModelSrc"
            :filename="modelFilename"
            class-name="border border-cyan-500/35 bg-cyan-500/10 text-[11px] font-medium text-cyan-100 transition hover:bg-cyan-500/16" />
          <LabBaseButton
            v-if="viewerSrc"
            icon-only
            icon="ic:round-open-in-full"
            button-class="inline-flex h-8 w-8 items-center justify-center rounded-full border border-zinc-700 bg-zinc-800 text-zinc-100 transition hover:bg-zinc-700"
            aria-label="Открыть просмотр модели на весь экран"
            title="Открыть на весь экран"
            @click="openFullscreen" />
        </div>
      </div>
    <div class="relative bg-zinc-950" :class="frameHeightClass">
      <div v-if="props.compact" class="absolute right-3 top-3 z-10 flex items-center gap-2">
        <LabViewerModelDownloadLink
          v-if="downloadableModelSrc && webglChecked && webglSupported"
          :href="downloadableModelSrc"
          :filename="modelFilename"
          icon-only
          class-name="rounded-full border border-white/12 bg-black/45 text-white backdrop-blur-sm transition hover:bg-black/60" />
        <LabBaseButton
          v-if="viewerSrc"
          icon-only
          icon="ic:round-open-in-full"
          button-class="h-8 w-8 rounded-full border border-white/12 bg-black/45 text-white backdrop-blur-sm transition hover:bg-black/60"
          aria-label="Открыть просмотр модели на весь экран"
          title="Открыть на весь экран"
          @click="openFullscreen" />
      </div>
      <div v-if="!webglChecked" class="flex h-full w-full p-2 items-center justify-center text-sm text-zinc-400">
        Проверка поддержки WebGL…
      </div>
      <div
        v-else-if="!webglSupported"
        class="flex h-full w-full flex-col items-center justify-center gap-3 px-6 text-center">
        <Icon name="ic:round-warning-amber" class="h-8 w-8 text-amber-400" />
        <div class="text-sm font-medium text-zinc-100">3D-просмотр недоступен</div>
        <p class="max-w-md text-sm leading-6 text-zinc-400">
          В браузере отключён WebGL или аппаратное ускорение графики. Включите аппаратное ускорение и перезапустите
          браузер.
        </p>
      </div>
      <iframe
        v-else-if="viewerSrc"
        :src="viewerSrc"
        :title="`${title} 3D viewer`"
        :class="iframeClass"
        loading="lazy"
        referrerpolicy="same-origin"
        allowfullscreen></iframe>
    </div>
  </div>
  <Teleport to="body">
    <div v-if="fullscreenOpen && viewerSrc" class="fixed inset-0 z-70 bg-black/92 backdrop-blur-sm">
      <div class="flex h-full flex-col">
        <div class="z-10">
          <div
            class="flex max-w-full flex-wrap items-center justify-end gap-2.5 bg-black/45 px-2.5 py-1.5 text-zinc-100 backdrop-blur-sm sm:px-3">
            <p class="min-w-0 flex-1 truncate text-sm font-semibold text-zinc-100 sm:text-base">
              {{ title }}
            </p>
            <LabViewerModelDownloadLink
              v-if="downloadableModelSrc"
              :href="downloadableModelSrc"
              :filename="modelFilename"
              class-name="rounded-lg bg-white/8 text-zinc-100 transition hover:bg-white/14" />
            <LabBaseButton
              icon="ic:round-close"
              label="Закрыть"
              button-class="h-8 shrink-0 rounded-lg bg-white/8 px-2.5 text-sm text-zinc-100 hover:bg-white/14"
              @click="closeFullscreen" />
          </div>
        </div>
        <div class="relative min-h-0 flex-1 bg-zinc-950">
          <iframe
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
