<script setup lang="ts">
  const formatShortUuid = shortUuid
  const props = withDefaults(
    defineProps<{
      action: 'copy' | 'paste'
      uuid?: string
      compact?: boolean
      disabled?: boolean
      variant?: LabButtonVariant
      size?: LabButtonSize
      buttonClass?: string
      label?: string
      title?: string
    }>(),
    {
      uuid: '',
      compact: false,
      disabled: false,
      variant: 'ghost',
      size: 'sm',
      buttonClass: '',
      label: '',
      title: ''
    }
  )
  const emit = defineEmits<{
    (e: 'click'): void
  }>()
  const { copyFrom } = useClipboard()
  const isCopy = computed(() => props.action === 'copy')
  const resolvedDisabled = computed(() => props.disabled || (isCopy.value && !props.uuid))
  const iconName = computed(() => (isCopy.value ? 'ic:round-content-copy' : 'ic:round-content-paste'))
  const actionLabel = computed(() => props.title || props.label || (isCopy.value ? 'Копировать UUID' : 'Вставить UUID'))
  const copyIdleLabel = computed(() => props.label || formatShortUuid(props.uuid, 5) || '—')
  const sharedButtonClass = computed(() =>
    ['text-zinc-200 hover:text-amber-300', isCopy.value ? 'truncate' : '', props.buttonClass].filter(Boolean).join(' ')
  )
  const onCopyClick = () => {
    if (resolvedDisabled.value || !props.uuid) return
    copyFrom(props.uuid)
    emit('click')
  }
  const onPasteClick = () => {
    if (resolvedDisabled.value) return
    emit('click')
  }
</script>
<template>
  <LabCopyHover
    v-if="isCopy"
    :idle="copyIdleLabel"
    hover="Копировать"
    done="Скопировано!"
    :has-icon="true"
    @click="onCopyClick">
    <template #default="{ display, labelStyle }">
      <LabBaseButton
        :aria-label="compact ? actionLabel : display"
        :button-class="sharedButtonClass"
        :disabled="resolvedDisabled"
        :icon-only="compact"
        :icon="iconName"
        :label="compact ? '' : display"
        label-class="truncate font-mono text-xs"
        :label-style="compact ? undefined : labelStyle"
        :size="size"
        :variant="variant" />
    </template>
  </LabCopyHover>
  <LabBaseButton
    v-else
    :aria-label="actionLabel"
    :button-class="sharedButtonClass"
    :disabled="resolvedDisabled"
    :icon-only="compact"
    :icon="iconName"
    :label="compact ? '' : actionLabel"
    :size="size"
    :variant="variant"
    @click="onPasteClick" />
</template>
