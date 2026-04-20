<template>
  <nav
    class="relative flex h-full min-h-0 flex-col gap-1 overflow-x-hidden overflow-y-auto"
    :class="[
      animate ? 'transition-[width] duration-200' : '',
      collapsed ? 'w-10' : 'w-56',
      showToggle && collapsed ? 'cursor-ew-resize' : ''
    ]"
    @click="onNavBackgroundClick"
  >
    <div v-if="showToggle" class="relative flex h-10 w-full items-center overflow-hidden">
      <LabNavHomeBtn />
      <LabBaseButton
        size="sm"
        icon-only
        variant="ghost"
        :class="[
          'absolute right-1 m-1 cursor-ew-resize rounded-full text-(--lab-text-primary)',
          animate ? 'transition-[opacity,background-color] duration-200' : '',
          collapsed ? 'pointer-events-none opacity-0' : 'opacity-100'
        ]"
        :aria-label="t('collapse_menu')"
        @click="emit('toggle-collapse')"
      >
        <span aria-hidden="true" :class="collapseGlyphClass">
          <span class="absolute inset-y-0 left-1 w-px bg-current" />
        </span>
      </LabBaseButton>
    </div>
    <div class="space-y-1 px-1">
      <NuxtLink
        v-for="item in props.items"
        :key="item.to"
        :to="item.to"
        :title="item.label"
        :class="navItemClass(item.to)"
      >
        <span class="relative inline-flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden">
          <Icon
            v-if="item.icon"
            :name="item.icon"
            :class="['h-4.5 w-4.5 shrink-0', item.iconColor || 'text-current']"
          />
        </span>
        <span
          v-if="!collapsed"
          class="ml-1 truncate text-sm whitespace-nowrap"
          :class="animate ? 'transition-opacity duration-150' : ''"
        >
          {{ item.label }}
        </span>
      </NuxtLink>
    </div>
    <div class="mx-auto mt-6 pb-6">
      <LabAvatar
        :show-label="!collapsed"
        :link-class="
          collapsed
            ? 'mx-auto h-8 w-8 justify-center hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)'
            : 'min-h-8 px-2 py-1 hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)'
        "
        @request-close="emit('request-close')"
      />
    </div>
  </nav>
</template>
<script setup lang="ts">
const emit = defineEmits<{
  (e: 'toggle-collapse' | 'request-close'): void
}>()
const route = useRoute()
const { locale, key, load, t } = useI18nSection('nav')
await useAsyncData(key.value, load, { watch: [locale] })
const props = withDefaults(
  defineProps<{
    items: MenuItem[]
    collapsed?: boolean
    showToggle?: boolean
    animate?: boolean
  }>(),
  {
    collapsed: false,
    showToggle: false,
    animate: false
  }
)
const { collapsed, showToggle, animate } = toRefs(props)
const collapseGlyphClass = 'relative inline-block h-3.5 w-3.5 border border-current opacity-90'
const onNavBackgroundClick = (event: MouseEvent) => {
  if (!showToggle.value || !collapsed.value) return
  const target = event.target as HTMLElement | null
  if (!target) return
  if (target.closest('a, button, [role="button"], input, select, textarea, label')) return
  emit('toggle-collapse')
}
const normalizePath = (rawPath: string) => {
  const path = String(rawPath || '/').replace(/\/{2,}/g, '/')
  if (path.length <= 1) return '/'
  return path.replace(/\/+$/, '')
}
const isItemActive = (itemTo: string) => {
  const current = normalizePath(route.path)
  const target = normalizePath(itemTo)
  if (target === '/') return current === '/'
  return current === target || current.startsWith(`${target}/`)
}
const navItemClass = (itemTo: string) => [
  'text-(--lab-text-secondary) hover:bg-(--lab-bg-surface-hover) hover:text-(--lab-text-primary) focus-visible:bg-(--lab-bg-surface-hover) focus-visible:text-(--lab-text-primary) rounded-full flex h-8 w-full items-center justify-start overflow-hidden px-0 select-none transition-colors',
  isItemActive(itemTo) ? 'bg-[color-mix(in_srgb,var(--lab-accent)_16%,transparent)] text-(--lab-accent)' : ''
]
</script>
