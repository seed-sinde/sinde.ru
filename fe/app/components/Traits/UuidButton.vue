<script setup lang="ts">
const { locale, key, load, t } = useI18nSection('traits')
await useAsyncData(key.value, load, { watch: [locale] })
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
const actionLabel = computed(
  () => props.title || props.label || (isCopy.value ? t('uuid_button.copy_uuid') : t('uuid_button.paste_uuid'))
)
const copyIdleLabel = computed(() => props.label || formatShortUuid(props.uuid, 5) || t('detail.copy_uuid_idle'))
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
    :hover="t('uuid_button.copy_hover')"
    :done="t('uuid_button.copied')"
    :has-icon="true"
    @click="onCopyClick"
  >
    <template #default="{ display }">
      <LabBaseButton
        :aria-label="compact ? actionLabel : display"
        :class="sharedButtonClass"
        :disabled="resolvedDisabled"
        :icon-only="compact"
        :icon="iconName"
        :label="compact ? '' : display"
        :size="size"
        :variant="variant"
      />
    </template>
  </LabCopyHover>
  <LabBaseButton
    v-else
    :aria-label="actionLabel"
    :class="sharedButtonClass"
    :disabled="resolvedDisabled"
    :icon-only="compact"
    :icon="iconName"
    :label="compact ? '' : actionLabel"
    :size="size"
    :variant="variant"
    @click="onPasteClick"
  />
</template>
