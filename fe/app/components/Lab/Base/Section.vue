<template>
  <section :class="sectionClassList" v-bind="sectionAttrs">
    <div v-if="$slots.header || title || description || $slots.actions" :class="headerClassList">
      <div class="min-w-0 space-y-1">
        <slot name="eyebrow" />
        <slot name="header">
          <h2 v-if="title" class="truncate text-base font-semibold text-zinc-900 dark:text-zinc-100">
            {{ title }}
          </h2>
          <p v-if="description" class="text-sm leading-6 text-zinc-600 dark:text-zinc-400">
            {{ description }}
          </p>
        </slot>
      </div>
      <div v-if="$slots.actions" class="shrink-0">
        <slot name="actions" />
      </div>
    </div>
    <div :class="contentClassList">
      <slot />
    </div>
    <div v-if="$slots.footer" :class="footerClassList">
      <slot name="footer" />
    </div>
  </section>
</template>
<script setup lang="ts">
import { twMerge } from 'tailwind-merge'
defineOptions({ inheritAttrs: false })
const props = withDefaults(
  defineProps<{
    title?: string
    description?: string
    variant?: 'plain' | 'panel' | 'soft'
    size?: 'sm' | 'md' | 'lg'
    headerClass?: string
    contentClass?: string
    footerClass?: string
    sectionClass?: string
  }>(),
  {
    title: '',
    description: '',
    variant: 'panel',
    size: 'md',
    headerClass: '',
    contentClass: '',
    footerClass: '',
    sectionClass: ''
  }
)
const attrs = useAttrs()
const sectionAttrs = computed(() => {
  const out: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(attrs)) {
    if (key === 'class') continue
    out[key] = value
  }
  return out
})
const externalClass = computed(() => String(attrs.class || '').trim())
const variantClass = computed(() => {
  switch (props.variant) {
    case 'plain':
      return ''
    case 'soft':
      return 'bg-zinc-50/80 dark:border-white/10 dark:bg-white/4'
    case 'panel':
    default:
      return 'bg-white dark:border-white/10 dark:bg-zinc-950/70'
  }
})
const paddingClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'p-3'
    case 'lg':
      return 'p-6'
    case 'md':
    default:
      return 'p-4'
  }
})
const gapClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'gap-3'
    case 'lg':
      return 'gap-5'
    case 'md':
    default:
      return 'gap-4'
  }
})
const headerClassList = computed(() =>
  twMerge(
    'flex items-start justify-between gap-4',
    props.variant === 'plain' ? '' : 'border-b border-zinc-200 pb-4 dark:border-white/10',
    props.headerClass
  )
)
const contentClassList = computed(() => twMerge(props.variant === 'plain' ? '' : '', props.contentClass))
const footerClassList = computed(() =>
  twMerge('pt-4', props.variant === 'plain' ? '' : 'border-t border-zinc-200 dark:border-white/10', props.footerClass)
)
const sectionClassList = computed(() =>
  twMerge(
    'flex min-w-0 flex-col',
    gapClass.value,
    variantClass.value,
    props.variant === 'plain' ? '' : paddingClass.value,
    props.sectionClass,
    externalClass.value
  )
)
</script>
