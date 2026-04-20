<template>
  <section :class="['min-w-0', isOpen ? 'bg-(--lab-bg-surface)' : '', containerClass]">
    <button type="button" :class="headerClass" :aria-expanded="isOpen ? 'true' : 'false'" @click="toggle">
      <span class="min-w-0">
        <span :class="labelClass">
          {{ label }}
        </span>
        <span v-if="description || $slots.description" :class="descriptionClass">
          <slot name="description">{{ description }}</slot>
        </span>
      </span>
      <span v-if="$slots.meta" :class="metaClass">
        <slot name="meta" :expanded="isOpen" />
      </span>
      <Icon :name="isOpen ? 'ic:round-expand-less' : 'ic:round-expand-more'" :class="iconClass" aria-hidden="true" />
    </button>
    <div v-if="isOpen" :class="['min-w-0', contentClass]">
      <slot :expanded="isOpen" />
    </div>
  </section>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    description?: string
    modelValue?: boolean
    defaultExpanded?: boolean
    containerClass?: string
    headerClass?: string
    labelClass?: string
    descriptionClass?: string
    metaClass?: string
    iconClass?: string
    contentClass?: string
  }>(),
  {
    label: 'Спойлер',
    description: '',
    defaultExpanded: false,
    containerClass: '',
    headerClass:
      'lab-focus flex w-full min-w-0 items-center justify-between gap-3 border border-(--lab-border) bg-(--lab-bg-surface) px-3 py-2 text-left text-(--lab-text-primary) transition hover:ring-1 hover:ring-(--lab-border)',
    labelClass: 'block text-sm font-semibold',
    descriptionClass: 'mt-1 block text-xs leading-5 text-(--lab-text-secondary)',
    metaClass: 'ml-auto hidden shrink-0 text-xs text-(--lab-text-muted) sm:block',
    iconClass: 'shrink-0 text-xl text-(--lab-text-secondary)',
    contentClass: 'pt-3'
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  toggle: [value: boolean]
}>()

const internalExpanded = ref(Boolean(props.defaultExpanded))
const controlled = computed(() => typeof props.modelValue === 'boolean')
const isOpen = computed(() => (controlled.value ? Boolean(props.modelValue) : internalExpanded.value))

const setOpen = (value: boolean) => {
  if (!controlled.value) {
    internalExpanded.value = value
  }
  emit('update:modelValue', value)
  emit('toggle', value)
}

const toggle = () => {
  setOpen(!isOpen.value)
}
</script>
