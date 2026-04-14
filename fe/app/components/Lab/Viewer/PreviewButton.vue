<template>
  <LabBaseButton
    size="none"
    icon-size="none"
    variant="plain"
    :button-class="resolvedButtonClass"
    @click="$emit('preview')"
  >
    <span :class="resolvedFrameClass">
      <img :src="src" :alt="alt" :class="resolvedImageClass" >
      <span class="image-preview-overlay">
        <span class="image-preview-overlay-bar">
          <Icon name="ic:round-open-in-full" class="h-4 w-4 shrink-0" />
          <span class="truncate">{{ label }}</span>
        </span>
      </span>
    </span>
  </LabBaseButton>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    src: string
    alt?: string
    label?: string
    compact?: boolean
    stepFormat?: boolean
    freeHeight?: boolean
    buttonClass?: string
    frameClass?: string
    imageClass?: string
  }>(),
  {
    alt: '',
    label: 'Посмотреть',
    compact: false,
    stepFormat: false,
    freeHeight: false,
    buttonClass: '',
    frameClass: '',
    imageClass: ''
  }
)
defineEmits<{
  (e: 'preview'): void
}>()
const resolvedButtonClass = computed(() =>
  [
    'group relative min-w-0 self-start bg-transparent p-0',
    props.stepFormat && props.compact
      ? 'w-52 max-w-full'
      : props.stepFormat
        ? 'w-full max-w-sm'
        : props.compact
          ? 'w-32 max-w-full'
          : 'w-auto max-w-full',
    props.buttonClass
  ]
    .join(' ')
    .trim()
)
const resolvedFrameClass = computed(() =>
  [
    'relative block w-full overflow-hidden rounded-md',
    props.freeHeight ? '' : props.stepFormat ? 'aspect-4/3' : props.compact ? 'aspect-square' : 'aspect-4/3',
    props.frameClass
  ]
    .join(' ')
    .trim()
)
const resolvedImageClass = computed(() =>
  ['block h-full w-full object-cover align-middle', props.imageClass].join(' ').trim()
)
</script>
<style scoped>
.image-preview-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  padding: 0.5rem;
  opacity: 0;
  transition: opacity 0.18s ease;
  background: linear-gradient(to top, rgb(9 9 11 / 0.56), rgb(9 9 11 / 0.12) 34%, transparent 62%);
  pointer-events: none;
}
.group:hover .image-preview-overlay,
.group:focus-visible .image-preview-overlay {
  opacity: 1;
}
.image-preview-overlay-bar {
  display: inline-flex;
  width: 100%;
  min-width: 0;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  border: 1px solid rgb(63 63 70 / 0.8);
  background: rgb(9 9 11 / 0.78);
  padding: 0.5rem 0.625rem;
  font-size: 0.75rem;
  line-height: 1rem;
  color: rgb(228 228 231);
  backdrop-filter: blur(4px);
}
</style>
