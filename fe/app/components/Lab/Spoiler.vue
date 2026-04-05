<template>
  <section :class="['lab-spoiler', containerClass]">
    <div
      :class="[
        'flex cursor-pointer flex-wrap items-center gap-3',
        inlineToggle ? 'justify-start' : 'justify-between',
        headerClass
      ]"
      role="button"
      tabindex="0"
      :aria-expanded="isOpen ? 'true' : 'false'"
      @click="toggle"
      @keydown.enter.prevent="toggle"
      @keydown.space.prevent="toggle">
      <h2 :class="titleClass">{{ title }}</h2>
      <LabBaseButton
        :aria-expanded="isOpen ? 'true' : 'false'"
        :label="isOpen ? hideLabel : showLabel"
        :icon="isOpen ? collapseIcon : expandIcon"
        :button-class="toggleButtonClass"
        icon-class="text-sm"
        @click.stop="toggle" />
    </div>
    <div v-if="isOpen" :class="contentClass">
      <slot :expanded="isOpen" />
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
      toggleButtonClass?: string
      contentClass?: string
      inlineToggle?: boolean
    }>(),
    {
      defaultExpanded: false,
      showLabel: 'Показать',
      hideLabel: 'Свернуть',
      expandIcon: 'ic:round-keyboard-double-arrow-down',
      collapseIcon: 'ic:round-keyboard-double-arrow-up',
      containerClass: 'space-y-3',
      headerClass: '',
      titleClass: 'lab-text-primary text-base font-semibold',
      toggleButtonClass: 'h-8 rounded-lg px-2.5 text-xs',
      contentClass: 'space-y-3',
      inlineToggle: false
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
<style scoped>
  .lab-spoiler {
    color: var(--lab-text-primary);
  }
</style>
