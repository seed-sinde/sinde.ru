<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      title: string
      viewerTitle?: string | undefined
      modelSrc?: string | null | undefined
      modelSizeLabel?: string | null | undefined
      rotationPerSecond?: string | null | undefined
      hint?: string
      compact?: boolean
    }>(),
    {
      viewerTitle: '',
      modelSrc: '',
      modelSizeLabel: '',
      rotationPerSecond: '',
      hint: '',
      compact: false
    }
  )
  const spoilerOpen = ref(false)
  const normalizedModelSrc = computed(() => String(props.modelSrc || '').trim())
  const resolvedViewerTitle = computed(() => props.viewerTitle || props.title)
  const showLabel = computed(() => {
    if (!props.modelSizeLabel) return 'Открыть 3D'
    return `Открыть 3D · ${props.modelSizeLabel}`
  })
  watch([normalizedModelSrc, resolvedViewerTitle], () => {
    spoilerOpen.value = false
  })
</script>
<template>
  <LabSpoiler
    v-if="normalizedModelSrc"
    v-model="spoilerOpen"
    :title="title"
    :show-label="showLabel"
    hide-label="Скрыть 3D"
    container-class="space-y-3 border px-4 py-4 lab-surface"
    title-class="text-sm font-semibold lab-text-primary"
    toggle-button-class="lab-button lab-button-secondary inline-flex h-8 items-center gap-1.5 px-2.5 text-xs"
    content-class="space-y-3">
    <p v-if="hint" class="text-sm leading-6 lab-text-muted wrap-break-word">
      {{ hint }}
    </p>
    <LabViewer3D
      :title="resolvedViewerTitle"
      :model-src="normalizedModelSrc"
      :rotation-per-second="rotationPerSecond"
      :compact="compact" />
  </LabSpoiler>
</template>
