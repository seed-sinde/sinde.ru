<template>
  <nav
    class="relative flex h-full min-h-0 flex-col gap-1 overflow-x-hidden overflow-y-auto py-1"
    :class="[
      animate ? 'transition-[width] duration-200' : '',
      collapsed ? 'w-10' : 'w-56',
      showToggle && collapsed ? 'cursor-ew-resize' : ''
    ]"
    @click="onNavBackgroundClick">
    <div v-if="showToggle" class="relative flex items-center px-1">
      <LabBaseButton
        size="sm"
        icon-only
        variant="ghost"
        :button-class="[
          'text-(--lab-text-primary) rounded-full',
          isCollapsedHoverToggle ? 'cursor-ew-resize' : 'cursor-pointer'
        ]"
        :aria-label="primaryButtonAriaLabel"
        @mouseenter="isPrimaryControlHovered = true"
        @mouseleave="isPrimaryControlHovered = false"
        @focus="isPrimaryControlHovered = true"
        @blur="isPrimaryControlHovered = false"
        @click="onPrimaryButtonClick">
        <template v-if="!isCollapsedHoverToggle">
          <picture>
            <source :srcset="faviconDarkSrc" media="(prefers-color-scheme: dark)" />
            <img :src="faviconLightSrc" alt="" class="h-4.5 w-4.5 object-contain" />
          </picture>
        </template>
        <span v-else aria-hidden="true" :class="collapseGlyphClass">
          <span class="absolute inset-y-0 left-1 w-px bg-current"></span>
        </span>
      </LabBaseButton>
      <LabBaseButton
        size="sm"
        icon-only
        variant="ghost"
        :button-class="[
          'text-(--lab-text-primary) rounded-full absolute right-1 cursor-ew-resize',
          animate ? 'transition-[opacity,background-color] duration-200' : '',
          collapsed ? 'pointer-events-none opacity-0' : 'opacity-100'
        ]"
        :aria-label="t('nav.collapse_menu')"
        @click="emit('toggle-collapse')">
        <span aria-hidden="true" :class="collapseGlyphClass">
          <span class="absolute inset-y-0 left-1 w-px bg-current"></span>
        </span>
      </LabBaseButton>
    </div>
    <template v-for="section in sections" :key="section.key">
      <div v-if="section.label" class="px-1">
        <div
          v-if="collapsed && section.marker"
          class="relative mx-1 flex h-5 items-center"
          :class="animate ? 'transition-opacity duration-200' : ''">
          <span class="bg-(--lab-border) h-px flex-1"></span>
          <span
            class="border-(--lab-bg-canvas) absolute left-1/2 top-1/2 inline-flex h-2.5 w-2.5 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border"
            :class="devMarkerClass"></span>
        </div>
        <div v-else-if="collapsed" class="bg-(--lab-border) mx-1 h-px"></div>
        <div
          v-else
          class="text-(--lab-text-soft) flex w-full items-center gap-2 overflow-hidden text-[10px] uppercase tracking-[0.08em]"
          :class="animate ? 'transition-opacity duration-200' : ''">
          <span class="bg-(--lab-border) h-px flex-1"></span>
          <span class="truncate whitespace-nowrap">{{ section.label }}</span>
          <span class="bg-(--lab-border) h-px flex-1"></span>
        </div>
      </div>
      <div class="space-y-1 px-1" :class="section.muted ? 'opacity-55' : ''">
        <NuxtLink
          v-for="item in section.items"
          :key="item.to"
          :to="item.to"
          :title="item.label"
          :class="navItemClass(item.to)">
          <span class="relative inline-flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden">
            <Icon
              v-if="item.icon"
              :name="item.icon"
              :class="['h-4.5 w-4.5 shrink-0', item.iconColor || 'text-current']" />
          </span>
          <span
            v-if="!collapsed"
            class="ml-1 truncate whitespace-nowrap text-sm"
            :class="animate ? 'transition-opacity duration-150' : ''">
            {{ item.label }}
          </span>
        </NuxtLink>
      </div>
    </template>
    <div class="mt-auto space-y-1">
      <div v-if="!collapsed" class="px-1">
        <LabNavFooter :collapsed="collapsed" :show-controls="true" :show-links="false" />
      </div>
      <div :class="collapsed ? 'px-1' : 'px-2'">
        <LabAvatar
          :show-label="!collapsed"
          :link-class="
            collapsed ?
              'mx-auto h-8 w-8 justify-center hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)'
            : 'min-h-8 px-2 py-1 hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)'
          "
          @request-close="emit('request-close')" />
      </div>
    </div>
  </nav>
</template>
<script setup lang="ts">
  const emit = defineEmits<{
    (e: 'toggle-collapse'): void
    (e: 'request-close'): void
  }>()
  const router = useRouter()
  const route = useRoute()
  const { t, faviconLightSrc, faviconDarkSrc } = useInterfacePreferences()
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
  const { items, collapsed, showToggle, animate } = toRefs(props)
  const isPrimaryControlHovered = ref(false)
  const isCollapsedHoverToggle = computed(() => showToggle.value && collapsed.value && isPrimaryControlHovered.value)
  const primaryButtonAriaLabel = computed(() => (isCollapsedHoverToggle.value ? t('nav.expand_menu') : t('nav.home')))
  const collapseGlyphClass = 'relative inline-block h-3.5 w-3.5 border border-current opacity-90'
  const onPrimaryButtonClick = async () => {
    if (isCollapsedHoverToggle.value) {
      emit('toggle-collapse')
      return
    }
    await router.push('/')
  }
  watch(collapsed, next => {
    if (!next) {
      isPrimaryControlHovered.value = false
    }
  })
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
  const primaryItems = computed(() => items.value.filter(item => !item.isDev && !item.group))
  const toolItems = computed(() => items.value.filter(item => !item.isDev && item.group === 'tools'))
  const knowledgeItems = computed(() => items.value.filter(item => !item.isDev && item.group === 'knowledge'))
  const devItems = computed(() => items.value.filter(item => item.isDev))
  const sections = computed(() =>
    [
      {
        key: 'primary',
        label: '',
        items: primaryItems.value,
        muted: false,
        marker: false
      },
      {
        key: 'tools',
        label: t('nav.tools'),
        items: toolItems.value,
        muted: false,
        marker: false
      },
      {
        key: 'knowledge',
        label: t('nav.knowledge'),
        items: knowledgeItems.value,
        muted: false,
        marker: false
      },
      {
        key: 'dev',
        label: t('nav.not_ready'),
        items: devItems.value,
        muted: true,
        marker: true
      }
    ].filter(section => section.items.length)
  )
  const devMarkerClass = 'bg-(--lab-warning) animate-pulse'
  const navItemClass = (itemTo: string) => [
    'text-(--lab-text-secondary) hover:bg-(--lab-bg-surface-hover) hover:text-(--lab-text-primary) focus-visible:bg-(--lab-bg-surface-hover) focus-visible:text-(--lab-text-primary) flex h-8 w-full items-center justify-start overflow-hidden px-0 select-none transition-colors',
    isItemActive(itemTo) ? 'bg-[color-mix(in_srgb,var(--lab-accent)_16%,transparent)] text-(--lab-accent)' : ''
  ]
</script>
