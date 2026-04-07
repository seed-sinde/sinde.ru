<script setup lang="ts">
  const route = useRoute()
  const router = useRouter()
  const runtimeConfig = useRuntimeConfig()
  const requestUrl = useRequestURL()
  const DESKTOP_MENU_STATE_KEY = 'layout-desktop-menu-collapsed'
  const desktopMenuStateCookie = useCookie<string | null>(DESKTOP_MENU_STATE_KEY, {
    path: '/',
    sameSite: 'lax'
  })
  const isMenuOpen = ref(false)
  const normalizeDesktopMenuState = (raw: unknown): boolean | null => {
    if (raw === true || raw === 1) return true
    if (raw === false || raw === 0) return false
    if (typeof raw !== 'string') return null
    const normalized = raw.trim().replace(/^"+|"+$/g, '')
    if (normalized === '1' || normalized.toLowerCase() === 'true') return true
    if (normalized === '0' || normalized.toLowerCase() === 'false') return false
    return null
  }
  const isDesktopMenuCollapsed = ref(normalizeDesktopMenuState(desktopMenuStateCookie.value) === true)
  const isDesktopMenuAnimationEnabled = ref(false)
  const mobileBreadcrumbItems = useState<BreadcrumbItem[]>('mobile-header-breadcrumb-items', () => [])
  const pageScroll = ref<HTMLElement | null>(null)
  const scrollPositions = useState<Record<string, number>>('page-scroll-positions', () => ({}))
  const windowScrollPositions = useState<Record<string, number>>('window-scroll-positions', () => ({}))
  const isPopNavigation = useState<boolean>('is-pop-navigation', () => false)
  let removeBeforeEach: (() => void) | null = null
  let removeAfterEach: (() => void) | null = null
  const normalizePath = (rawPath: string) => {
    const path = String(rawPath || '/').replace(/\/{2,}/g, '/')
    if (path.length <= 1) return '/'
    return path.replace(/\/+$/, '')
  }
  const siteUrl = computed(() => {
    const configured = String(runtimeConfig.public.baseURL || '').trim()
    if (configured) return configured.replace(/\/+$/, '')
    return `${requestUrl.protocol}//${requestUrl.host}`.replace(/\/+$/, '')
  })
  const normalizedPath = computed(() => normalizePath(route.path))
  const robotsContent = computed(() => {
    const metaRobots = route.meta.robots
    if (typeof metaRobots === 'string' && metaRobots.trim()) return metaRobots.trim()
    return import.meta.dev ? 'noindex, nofollow' : 'index, follow'
  })
  const canonicalUrl = computed(() => `${siteUrl.value}${normalizedPath.value}`)
  useHead({
    link: [{ key: 'canonical', rel: 'canonical', href: canonicalUrl }],
    meta: [
      { key: 'robots', name: 'robots', content: robotsContent },
      { key: 'og:type', property: 'og:type', content: 'website' },
      { key: 'og:url', property: 'og:url', content: canonicalUrl }
    ]
  })
  const sidebar_items = useSidebarItems()
  const persistDesktopMenuState = (collapsed: boolean) => {
    const next = collapsed ? '1' : '0'
    desktopMenuStateCookie.value = next
    if (!import.meta.client) return
    try {
      window.localStorage.setItem(DESKTOP_MENU_STATE_KEY, next)
    } catch {
      // ignore storage errors
    }
  }
  const closeMenu = () => {
    isMenuOpen.value = false
  }
  const toggleDesktopMenu = () => {
    const next = !isDesktopMenuCollapsed.value
    isDesktopMenuCollapsed.value = next
    persistDesktopMenuState(next)
  }
  watch(
    () => route.fullPath,
    () => {
      closeMenu()
    }
  )
  const onPopState = () => {
    isPopNavigation.value = true
  }
  const shouldPreserveScrollOnNavigation = (
    to: { path: string; hash?: string | null },
    from: { path: string; hash?: string | null }
  ) => {
    return normalizePath(to.path) === normalizePath(from.path) && String(to.hash || '') === String(from.hash || '')
  }
  const saveScrollPosition = (path: string) => {
    if (!import.meta.client || !path) return
    const el = pageScroll.value
    if (el) {
      scrollPositions.value[path] = el.scrollTop
    }
    windowScrollPositions.value[path] = window.scrollY
  }
  const restoreScrollPosition = (path: string) => {
    if (!import.meta.client) return
    const el = pageScroll.value
    const savedContainerTop = scrollPositions.value[path]
    const savedWindowTop = windowScrollPositions.value[path]
    const shouldRestoreContainer = Number.isFinite(savedContainerTop)
    const shouldRestoreWindow = Number.isFinite(savedWindowTop)
    if (el) {
      el.scrollTop = shouldRestoreContainer ? Number(savedContainerTop) : 0
    }
    window.scrollTo({
      top: shouldRestoreWindow ? Number(savedWindowTop) : 0,
      left: 0,
      behavior: 'auto'
    })
  }
  const resetScrollPosition = () => {
    if (!import.meta.client) return
    const el = pageScroll.value
    if (el) {
      el.scrollTop = 0
    }
    window.scrollTo({
      top: 0,
      left: 0,
      behavior: 'auto'
    })
  }
  const scheduleScrollUpdate = (cb: () => void) => {
    requestAnimationFrame(() => {
      requestAnimationFrame(cb)
    })
  }
  onMounted(() => {
    if ('scrollRestoration' in window.history) {
      window.history.scrollRestoration = 'manual'
    }
    window.addEventListener('popstate', onPopState)
    removeBeforeEach = router.beforeEach((_, from) => {
      if (from.fullPath) {
        saveScrollPosition(from.fullPath)
      }
      return true
    })
    removeAfterEach = router.afterEach(async (to, from) => {
      await nextTick()
      scheduleScrollUpdate(() => {
        if (isPopNavigation.value) {
          restoreScrollPosition(to.fullPath)
        } else if (!shouldPreserveScrollOnNavigation(to, from)) {
          resetScrollPosition()
        }
        isPopNavigation.value = false
      })
    })
    scheduleScrollUpdate(() => {
      if (isPopNavigation.value) {
        restoreScrollPosition(route.fullPath)
        isPopNavigation.value = false
      }
      isDesktopMenuAnimationEnabled.value = true
    })
  })
  onBeforeUnmount(() => {
    if (import.meta.client) {
      saveScrollPosition(route.fullPath)
      if ('scrollRestoration' in window.history) {
        window.history.scrollRestoration = 'auto'
      }
    }
    window.removeEventListener('popstate', onPopState)
    removeBeforeEach?.()
    removeAfterEach?.()
  })
</script>
<template>
  <div class="bg-(--lab-bg-canvas) text-(--lab-text-secondary) min-h-screen lg:h-screen lg:overflow-hidden">
    <main class="flex min-h-screen w-full flex-col lg:h-full lg:flex-row lg:overflow-hidden">
      <div class="bg-(--lab-bg-overlay) sticky top-0 z-40 flex items-center gap-2 border-b lg:hidden">
        <LabBaseButton
          icon="ic:round-menu"
          icon-only
          size="sm"
          variant="ghost"
          button-class="m-1"
          @click="isMenuOpen = true" />
        <div class="min-w-0 flex-1">
          <ClientOnly>
            <LabNavBreadcrumb
              v-if="mobileBreadcrumbItems.length"
              :items="mobileBreadcrumbItems"
              container-class="min-w-0"
              list-class="text-(--lab-text-secondary) flex min-w-0 flex-nowrap items-center gap-2 overflow-hidden text-sm"
              separator-class="text-(--lab-text-soft) h-4 w-4 shrink-0"
              class="min-w-0" />
          </ClientOnly>
        </div>
        <div id="mobile-header-actions" class="flex shrink-0 items-center gap-1"></div>
        <LabAvatar
          :show-label="false"
          link-class="hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover) my-1 mr-1 inline-flex h-8 w-8 shrink-0 items-center justify-center bg-transparent p-1 text-(--lab-text-primary) transition-colors" />
      </div>
      <LabNavDrawer v-model="isMenuOpen" :items="sidebar_items" />
      <aside class="hidden h-full min-h-0 shrink-0 border-r lg:block">
        <LabNavSidebar
          :items="sidebar_items"
          :collapsed="isDesktopMenuCollapsed"
          :show-toggle="true"
          :animate="isDesktopMenuAnimationEnabled"
          @toggle-collapse="toggleDesktopMenu"
          @request-close="closeMenu" />
      </aside>
      <div ref="pageScroll" class="min-h-0 min-w-0 lg:flex-1 lg:overflow-y-auto">
        <slot></slot>
        <div class="m-4 flex flex-wrap items-center gap-x-3 gap-y-2">
          <div class="flex shrink-0 items-center gap-2">
            <LabLocaleSwitcher />
            <LabThemeSwitcher />
          </div>
          <LabNavFooter :show-controls="false" :show-links="true" />
        </div>
        <div class="h-16 lg:h-12" aria-hidden="true"></div>
      </div>
    </main>
  </div>
</template>
