<template>
  <section :class="['space-y-3', isOpen ? 'bg-(--lab-bg-surface)' : '', containerClass]">
    <LabBaseButton
      :class="headerClass"
      :aria-expanded="isOpen ? 'true' : 'false'"
      :label="label"
      :label-class="labelClass"
      @click="toggle" />
    <div v-if="isOpen" :class="['min-w-0', contentClass]">
      <slot :expanded="isOpen"></slot>
    </div>
  </section>
</template>

<script setup lang="ts">
  const props = withDefaults(
    defineProps<{
      label: string
      modelValue?: boolean
      defaultExpanded?: boolean
      containerClass?: string
      headerClass?: string
      labelClass?: string
      contentClass?: string
    }>(),
    {
      label: 'Спойлер',
      defaultExpanded: false,
      containerClass: '',
      headerClass: '',
      labelClass: '',
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
