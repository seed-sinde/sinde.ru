<script setup lang="ts">
  import Cropper from 'cropperjs'
  import 'cropperjs/dist/cropper.css'
  const props = withDefaults(
    defineProps<{
      file: File | null
      title?: string
      loading?: boolean
      aspectPreset?: CropperAspectPreset
      aspectLocked?: boolean
    }>(),
    {
      title: '',
      loading: false,
      aspectPreset: 'free',
      aspectLocked: false
    }
  )
  const emit = defineEmits<{
    (e: 'cancel'): void
    (e: 'confirm', file: File): void
  }>()
  const imageRef = ref<HTMLImageElement | null>(null)
  const cropperViewportRef = ref<HTMLElement | null>(null)
  const headerScrollerRef = ref<HTMLElement | null>(null)
  const controlsScrollerRef = ref<HTMLElement | null>(null)
  const objectUrl = ref('')
  const cropError = ref<string | null>(null)
  const aspectPreset = ref<CropperAspectPreset>('free')
  const customAspectWidth = ref('1')
  const customAspectHeight = ref('1')
  const scaleX = ref(1)
  const scaleY = ref(1)
  const rotation = ref(0)
  const zoomRatio = ref(1)
  const originalImageAspect = ref<number | null>(null)
  let cropper: Cropper | null = null
  let cropperResizeObserver: ResizeObserver | null = null
  const { edges: headerScrollEdges, sync: syncHeaderScrollEdges } = useScrollableEdges(headerScrollerRef, { axis: 'x' })
  const { edges: controlsScrollEdges, sync: syncControlsScrollEdges } = useScrollableEdges(controlsScrollerRef, {
    axis: 'x'
  })
  const parsedCustomAspect = computed(() => {
    const w = Number.parseFloat(String(customAspectWidth.value || '').replace(',', '.'))
    const h = Number.parseFloat(String(customAspectHeight.value || '').replace(',', '.'))
    if (!Number.isFinite(w) || !Number.isFinite(h) || w <= 0 || h <= 0) return null
    return w / h
  })
  const currentAspectRatio = computed(() => {
    if (aspectPreset.value === '1:1') return 1
    if (aspectPreset.value === '4:3') return 4 / 3
    if (aspectPreset.value === '16:9') return 16 / 9
    if (aspectPreset.value === 'custom') return parsedCustomAspect.value
    return null
  })
  const revokeObjectUrl = () => {
    if (!objectUrl.value) return
    URL.revokeObjectURL(objectUrl.value)
    objectUrl.value = ''
  }
  const destroyCropper = () => {
    if (!cropper) return
    cropper.destroy()
    cropper = null
  }
  const readCurrentZoomRatio = () => {
    if (!cropper) return 1
    const image = cropper.getImageData()
    const naturalWidth = Number(image?.naturalWidth || 0)
    const renderedWidth = Number(image?.width || 0)
    if (!(naturalWidth > 0 && renderedWidth > 0)) return 1
    const ratio = renderedWidth / naturalWidth
    return ratio > 0 && Number.isFinite(ratio) ? ratio : 1
  }
  const applyAspectRatio = () => {
    if (!cropper) return
    const ratio = currentAspectRatio.value
    if (ratio && Number.isFinite(ratio) && ratio > 0) {
      cropper.setAspectRatio(ratio)
      cropError.value = null
      return
    }
    if (aspectPreset.value === 'custom') {
      cropError.value = 'Для ручного соотношения укажите оба значения больше нуля.'
      return
    }
    cropError.value = null
    cropper.setAspectRatio(NaN)
  }
  const initCropper = async () => {
    destroyCropper()
    if (!props.file || !imageRef.value) return
    cropper = new Cropper(imageRef.value, {
      viewMode: 1,
      dragMode: 'move',
      autoCrop: true,
      autoCropArea: 0.92,
      background: false,
      center: true,
      guides: true,
      highlight: false,
      movable: true,
      zoomable: true,
      scalable: true,
      rotatable: true,
      responsive: true,
      zoom: (event: Cropper.ZoomEvent<HTMLImageElement>) => {
        const ratio = Number(event?.detail?.ratio || 0)
        if (ratio > 0 && Number.isFinite(ratio)) {
          zoomRatio.value = ratio
        }
      },
      ready: () => {
        applyAspectRatio()
        cropper?.scaleX(1)
        cropper?.scaleY(1)
        const imageData = cropper?.getImageData()
        const naturalWidth = Number(imageData?.naturalWidth || 0)
        const naturalHeight = Number(imageData?.naturalHeight || 0)
        if (naturalWidth > 0 && naturalHeight > 0) {
          originalImageAspect.value = naturalWidth / naturalHeight
          // Keep "Своб." selected, but initialize crop box to source image proportions.
          if (aspectPreset.value === 'free') {
            cropper?.setAspectRatio(originalImageAspect.value)
            setTimeout(() => {
              cropper?.setAspectRatio(NaN)
            }, 0)
          }
        } else {
          originalImageAspect.value = null
        }
        zoomRatio.value = readCurrentZoomRatio()
      }
    })
  }
  const resetCropper = () => {
    if (!cropper) return
    cropper.reset()
    scaleX.value = 1
    scaleY.value = 1
    rotation.value = 0
    zoomRatio.value = readCurrentZoomRatio()
    applyAspectRatio()
  }
  const normalizeAngle = (value: number) => {
    const normalized = value % 360
    return normalized < 0 ? normalized + 360 : normalized
  }
  const fitImageToView = () => {
    if (!cropper) return
    const container = cropper.getContainerData()
    const image = cropper.getImageData()
    const containerWidth = Number(container?.width || 0)
    const containerHeight = Number(container?.height || 0)
    const naturalWidth = Number(image?.naturalWidth || 0)
    const naturalHeight = Number(image?.naturalHeight || 0)
    if (!(containerWidth > 0 && containerHeight > 0 && naturalWidth > 0 && naturalHeight > 0)) return
    const angle = (normalizeAngle(rotation.value) * Math.PI) / 180
    const rotatedWidth = Math.abs(naturalWidth * Math.cos(angle)) + Math.abs(naturalHeight * Math.sin(angle))
    const rotatedHeight = Math.abs(naturalWidth * Math.sin(angle)) + Math.abs(naturalHeight * Math.cos(angle))
    if (!(rotatedWidth > 0 && rotatedHeight > 0)) return
    const fitRatio = Math.min((containerWidth * 0.96) / rotatedWidth, (containerHeight * 0.96) / rotatedHeight)
    if (!(fitRatio > 0)) return
    cropper.zoomTo(fitRatio)
    zoomRatio.value = fitRatio
    const next = cropper.getImageData()
    cropper.moveTo((containerWidth - next.width) / 2, (containerHeight - next.height) / 2)
  }
  const fitToViewport = () => {
    fitImageToView()
  }
  const rotateBy = (delta: number) => {
    if (!cropper) return
    rotation.value = normalizeAngle(rotation.value + delta)
    cropper.rotateTo(rotation.value)
    fitImageToView()
  }
  const rotateLeft = () => rotateBy(-90)
  const rotateRight = () => rotateBy(90)
  const flipHorizontal = () => {
    if (!cropper) return
    scaleX.value = scaleX.value === 1 ? -1 : 1
    cropper.scaleX(scaleX.value)
  }
  const flipVertical = () => {
    if (!cropper) return
    scaleY.value = scaleY.value === 1 ? -1 : 1
    cropper.scaleY(scaleY.value)
  }
  const zoomIn = () => cropper?.zoom(0.1)
  const zoomOut = () => cropper?.zoom(-0.1)
  const zoomPercentLabel = computed(() => `${Math.round(Math.max(0, zoomRatio.value) * 100)}%`)
  const aspectPresetOptions = [
    { value: '1:1', label: '1:1' },
    { value: '4:3', label: '4:3' },
    { value: '16:9', label: '16:9' },
    { value: 'custom', label: 'Ручн.' },
    { value: 'free', label: 'Своб.' }
  ] satisfies SelectOptionInput[]
  const controlButtonClass = 'h-9 w-9 shrink-0 rounded-full p-0 text-sm'
  const actionButtonClass = 'h-9 shrink-0 px-3 text-sm'
  const onZoomSliderInput = () => {
    if (!cropper) return
    const ratio = Number(zoomRatio.value || 0)
    if (!(ratio > 0 && Number.isFinite(ratio))) return
    cropper.zoomTo(ratio)
  }
  const invertAspectRatio = () => {
    if (aspectPreset.value === 'free') return
    if (aspectPreset.value === '1:1') return
    if (aspectPreset.value === '4:3') {
      customAspectWidth.value = '3'
      customAspectHeight.value = '4'
      aspectPreset.value = 'custom'
      return
    }
    if (aspectPreset.value === '16:9') {
      customAspectWidth.value = '9'
      customAspectHeight.value = '16'
      aspectPreset.value = 'custom'
      return
    }
    const nextWidth = customAspectHeight.value || '1'
    const nextHeight = customAspectWidth.value || '1'
    customAspectWidth.value = nextWidth
    customAspectHeight.value = nextHeight
  }
  const chooseOutputMime = (raw: string) => {
    const safe = String(raw || '')
      .trim()
      .toLowerCase()
    if (safe.startsWith('image/')) return 'image/webp'
    return 'image/webp'
  }
  const fileNameBase = (name: string) => {
    const raw = String(name || '').trim()
    const idx = raw.lastIndexOf('.')
    if (idx <= 0) return raw || 'image'
    return raw.slice(0, idx)
  }
  const fileExtFromMime = (mime: string) => {
    if (mime === 'image/webp') return 'webp'
    return 'jpg'
  }
  const applyCrop = async () => {
    if (!cropper || !props.file || props.loading) return
    const ratio = currentAspectRatio.value
    if (aspectPreset.value === 'custom' && (!ratio || !Number.isFinite(ratio) || ratio <= 0)) {
      cropError.value = 'Для ручного соотношения укажите корректные значения.'
      return
    }
    cropError.value = null
    const canvas = cropper.getCroppedCanvas({
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high'
    })
    if (!canvas) {
      cropError.value = 'Не удалось получить область кадрирования.'
      return
    }
    const mime = chooseOutputMime(props.file.type)
    const quality = mime === 'image/jpeg' || mime === 'image/webp' ? 0.92 : undefined
    const blob = await new Promise<Blob | null>(resolve => {
      canvas.toBlob(resolve, mime, quality)
    })
    if (!blob) {
      cropError.value = 'Не удалось подготовить изображение к загрузке.'
      return
    }
    const ext = fileExtFromMime(mime)
    const nextFile = new File([blob], `${fileNameBase(props.file.name)}-cropped.${ext}`, {
      type: mime,
      lastModified: Date.now()
    })
    emit('confirm', nextFile)
  }
  const scheduleCropperLayoutSync = () => {
    if (!import.meta.client) {
      fitImageToView()
      return
    }
    requestAnimationFrame(() => {
      fitImageToView()
    })
  }
  const onEscKeydown = (event: KeyboardEvent) => {
    if (event.key !== 'Escape') return
    if (props.loading) return
    emit('cancel')
  }
  watch(
    () => props.file,
    async file => {
      destroyCropper()
      revokeObjectUrl()
      cropError.value = null
      scaleX.value = 1
      scaleY.value = 1
      rotation.value = 0
      zoomRatio.value = 1
      originalImageAspect.value = null
      aspectPreset.value = props.aspectPreset
      customAspectWidth.value = '1'
      customAspectHeight.value = '1'
      if (!file) return
      objectUrl.value = URL.createObjectURL(file)
      await nextTick()
      await initCropper()
    },
    { immediate: true }
  )
  watch([aspectPreset, customAspectWidth, customAspectHeight], () => {
    applyAspectRatio()
  })
  watch(
    () => props.aspectPreset,
    nextPreset => {
      if (!props.file) return
      aspectPreset.value = nextPreset
      applyAspectRatio()
    }
  )
  onMounted(() => {
    if (!import.meta.client) return
    window.addEventListener('keydown', onEscKeydown)
    cropperResizeObserver = new ResizeObserver(() => {
      syncHeaderScrollEdges()
      syncControlsScrollEdges()
      scheduleCropperLayoutSync()
    })
    if (cropperViewportRef.value) {
      cropperResizeObserver.observe(cropperViewportRef.value)
    }
  })
  onBeforeUnmount(() => {
    destroyCropper()
    revokeObjectUrl()
    cropperResizeObserver?.disconnect()
    cropperResizeObserver = null
    if (import.meta.client) {
      window.removeEventListener('keydown', onEscKeydown)
    }
  })
</script>
<template>
  <div
    class="bg-(--lab-bg-canvas) fixed inset-0 z-60"
    @click.self="emit('cancel')">
    <div class="bg-(--lab-bg-canvas) flex h-full w-full flex-col">
      <div class="z-10" @click.stop>
        <div class="relative">
          <div
            ref="headerScrollerRef"
            class="bg-(--lab-bg-canvas) text-(--lab-text-secondary) lab-scroll-hidden w-full overflow-x-auto overflow-y-hidden border-b px-2.5 py-1.5 sm:px-3">
            <div class="flex min-w-max items-center justify-end gap-2.5">
              <div class="flex shrink-0 items-center gap-1.5 text-sm font-medium">
                масштаб:
                <p class="tabular-nums">{{ zoomPercentLabel }}</p>
              </div>
              <LabBaseButton
                icon="ic:round-close"
                icon-only
                aria-label="Закрыть окно"
                :disabled="loading"
                @click="emit('cancel')" />
            </div>
          </div>
          <div
            class="lab-scroll-fade lab-scroll-fade-x-left"
            :class="{ 'lab-scroll-fade-visible': headerScrollEdges.left }"
            aria-hidden="true"></div>
          <div
            class="lab-scroll-fade lab-scroll-fade-x-right"
            :class="{ 'lab-scroll-fade-visible': headerScrollEdges.right }"
            aria-hidden="true"></div>
        </div>
      </div>
      <div
        ref="cropperViewportRef"
        class="bg-(--lab-bg-canvas) relative min-h-0 flex-1 overflow-hidden px-2 py-2 sm:px-3 sm:py-3">
        <div
          class="bg-(--lab-bg-surface-muted) h-full w-full overflow-hidden border [&_.cropper-container]:h-full [&_.cropper-container]:w-full [&_.cropper-dashed]:border-white/20 [&_.cropper-face]:bg-transparent [&_.cropper-line]:bg-white/85 [&_.cropper-point]:bg-white/85 [&_.cropper-view-box]:outline [&_.cropper-view-box]:outline-white/40"
          @click.stop>
          <img
            ref="imageRef"
            :src="objectUrl"
            alt="Изображение для кадрирования"
            class="block h-full max-h-full w-full max-w-full object-contain" />
        </div>
      </div>
      <div
        class="bg-(--lab-bg-canvas) border-t px-3 py-3 backdrop-blur-sm sm:px-4"
        @click.stop>
        <div class="relative">
          <div ref="controlsScrollerRef" class="lab-scroll-hidden overflow-x-auto overflow-y-hidden">
            <div class="flex min-w-max items-center gap-3">
            <div class="flex shrink-0 items-center gap-2">
              <LabBaseButton
                icon="ic:round-rotate-left"
                icon-only
                title="Повернуть влево"
                aria-label="Повернуть влево"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="rotateLeft" />
              <LabBaseButton
                icon="ic:round-rotate-right"
                icon-only
                title="Повернуть вправо"
                aria-label="Повернуть вправо"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="rotateRight" />
              <LabBaseButton
                icon="ic:round-swap-horiz"
                icon-only
                title="Отразить по горизонтали"
                aria-label="Отразить по горизонтали"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="flipHorizontal" />
              <LabBaseButton
                icon="ic:round-swap-vert"
                icon-only
                title="Отразить по вертикали"
                aria-label="Отразить по вертикали"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="flipVertical" />
              <LabBaseButton
                icon="ic:round-minus"
                icon-only
                title="Уменьшить"
                aria-label="Уменьшить"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="zoomOut" />
              <LabBaseButton
                icon="ic:round-plus"
                icon-only
                title="Увеличить"
                aria-label="Увеличить"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="zoomIn" />
              <div class="flex h-9 shrink-0 items-center gap-2 px-1">
                <input
                  id="image-crop-zoom-ratio"
                  v-model.number="zoomRatio"
                  name="image_crop_zoom_ratio"
                  type="range"
                  min="0.1"
                  max="3"
                  step="0.01"
                  class="accent-(--lab-info) w-28 border-0 bg-transparent"
                  @input="onZoomSliderInput" />
              </div>
              <LabBaseButton
                icon="ic:round-refresh"
                icon-only
                title="Сбросить"
                aria-label="Сбросить"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="resetCropper" />
              <LabBaseButton
                icon="ic:round-fit-screen"
                icon-only
                title="Вместить"
                aria-label="Вместить"
                :button-class="controlButtonClass"
                :disabled="loading"
                @click="fitToViewport" />
              <template v-if="!props.aspectLocked">
                <LabBaseSelect
                  id="image-crop-aspect-preset"
                  v-model="aspectPreset"
                  name="image_crop_aspect_preset"
                  aria-label="Формат кадрирования"
                  title="Формат кадрирования"
                  :options="aspectPresetOptions"
                  select-class="h-9 w-20 shrink-0 px-2.5 text-sm"
                  :disabled="loading" />
                <LabBaseButton
                  icon="ic:round-swap-calls"
                  icon-only
                  title="Инвертировать формат"
                  aria-label="Инвертировать формат"
                  :button-class="controlButtonClass"
                  :disabled="loading"
                  @click="invertAspectRatio" />
              </template>
              <div v-else class="inline-flex h-9 shrink-0 items-center border px-3 text-sm">
                {{ aspectPreset }}
              </div>
              <div
                v-if="!props.aspectLocked && aspectPreset === 'custom'"
                class="inline-flex shrink-0 items-center gap-2">
                <LabBaseInput
                  id="image-crop-custom-aspect-width"
                  v-model="customAspectWidth"
                  name="image_crop_custom_aspect_width"
                  type="number"
                  min="0.01"
                  step="0.01"
                  aria-label="Ширина пользовательского формата"
                  input-class="h-9 w-16 px-2 text-sm" />
                <span class="text-(--lab-text-muted) text-sm">:</span>
                <LabBaseInput
                  id="image-crop-custom-aspect-height"
                  v-model="customAspectHeight"
                  name="image_crop_custom_aspect_height"
                  type="number"
                  min="0.01"
                  step="0.01"
                  aria-label="Высота пользовательского формата"
                  input-class="h-9 w-16 px-2 text-sm" />
              </div>
            </div>
            <div class="ml-auto flex shrink-0 items-center gap-2">
              <LabBaseButton
                label="Отмена"
                variant="secondary"
                :button-class="actionButtonClass"
                :disabled="loading"
                @click="emit('cancel')" />
              <LabBaseButton
                label="Загрузить"
                loading-label="Загрузка..."
                variant="success"
                :button-class="actionButtonClass"
                :loading="loading"
                @click="applyCrop" />
            </div>
          </div>
          </div>
          <div
            class="lab-scroll-fade lab-scroll-fade-x-left"
            :class="{ 'lab-scroll-fade-visible': controlsScrollEdges.left }"
            aria-hidden="true"></div>
          <div
            class="lab-scroll-fade lab-scroll-fade-x-right"
            :class="{ 'lab-scroll-fade-visible': controlsScrollEdges.right }"
            aria-hidden="true"></div>
        </div>
        <LabNotify :text="cropError" tone="error" size="xs" class="mt-3" />
      </div>
    </div>
  </div>
</template>
<style scoped>
  :deep(.cropper-container) {
    background: var(--lab-bg-surface-muted);
  }
  :deep(.cropper-bg) {
    background: var(--lab-bg-surface-subtle);
  }
  :deep(.cropper-modal) {
    background: var(--lab-bg-canvas);
  }
</style>
