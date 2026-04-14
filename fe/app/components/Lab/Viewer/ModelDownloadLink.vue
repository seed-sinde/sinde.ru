<template>
  <a :href="href" :download="filename" :class="resolvedClass" :aria-label="label" :title="label">
    <Icon name="ic:round-download" />
    <span v-if="!iconOnly">{{ label }}</span>
  </a>
</template>
<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    href: string
    filename?: string
    label?: string
    iconOnly?: boolean
    className?: string
  }>(),
  {
    filename: '',
    label: 'Скачать glb',
    iconOnly: false,
    className: ''
  }
)
const resolvedClass = computed(() =>
  [
    props.iconOnly
      ? 'inline-flex h-8 w-8 items-center justify-center'
      : 'inline-flex h-8 items-center gap-2 px-2.5 text-sm',
    props.className
  ]
    .join(' ')
    .trim()
)
</script>
