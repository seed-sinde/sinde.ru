<template>
  <section
    :class="[
      'min-w-0 w-full space-y-3',
      isOpen ? 'bg-(--lab-bg-surface)' : '',
      containerClass
    ]">
    <button
      type="button"
      :aria-expanded="isOpen ? 'true' : 'false'"
      :aria-label="isOpen ? hideLabel : showLabel"
      :class="[
        'flex min-w-0 w-full items-center justify-between gap-3 text-left text-(--lab-text-primary) transition focus:outline-none focus-visible:ring-2 focus-visible:ring-(--lab-accent)',
        headerClass
      ]"
      @click="toggle">
      <span class="min-w-0 flex-1">
        <span :class="['inline-flex max-w-full border-b border-(--lab-text-primary) pb-0.5', titleClass]">
          {{ title }}
        </span>
      </span>
      <span class="shrink-0 text-(--lab-text-secondary)" aria-hidden="true">
        <Icon :name="isOpen ? collapseIcon : expandIcon" class="text-xl" />
      </span>
    </button>

    <div v-if="isOpen" :class="['min-w-0', contentClass]">
      <slot :expanded="isOpen"></slot>
    </div>
  </section>
</template>

<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      title: string
      modelValue?: boolean
      defaultExpanded?: boolean
      showLabel?: string
      hideLabel?: string
      expandIcon?: string
      collapseIcon?: string
      containerClass?: string
      headerClass?: string
      titleClass?: string
      contentClass?: string
    }>(),
    {
      title: 'Спойлер',
      defaultExpanded: false,
      showLabel: 'Показать',
      hideLabel: 'Скрыть',
      expandIcon: 'ic:round-keyboard-double-arrow-down',
      collapseIcon: 'ic:round-keyboard-double-arrow-up',
      containerClass: '',
      headerClass: '',
      titleClass: '',
      contentClass: ''
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
