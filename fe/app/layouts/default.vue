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

const isMobileMenuOpen = ref(false)
const isDesktopMenuAnimationEnabled = ref(false)
const pageScrollRef = ref<HTMLElement | null>(null)

const mobileBreadcrumbItems = computed<BreadcrumbItem[]>(() =>
  Array.isArray(route.meta.mobileBreadcrumbItems) ? (route.meta.mobileBreadcrumbItems as BreadcrumbItem[]) : []
)
const mobileHeaderActions = computed<MobileHeaderAction[]>(() =>
  Array.isArray(route.meta.mobileHeaderActions) ? (route.meta.mobileHeaderActions as MobileHeaderAction[]) : []
)

const pageScrollPositions = useState<Record<string, number>>('page-scroll-positions', () => ({}))
const windowScrollPositions = useState<Record<string, number>>('window-scroll-positions', () => ({}))
const isHistoryPopNavigation = useState<boolean>('is-pop-navigation', () => false)

let removeBeforeEachGuard: (() => void) | null = null
let removeAfterEachHook: (() => void) | null = null

const normalizePath = (rawPath: string) => {
  const normalizedPath = String(rawPath || '/').replace(/\/{2,}/g, '/')
  return normalizedPath.length <= 1 ? '/' : normalizedPath.replace(/\/+$/, '')
}

const parseDesktopMenuState = (value: unknown): boolean | null => {
  if (value === true || value === 1) return true
  if (value === false || value === 0) return false
  if (typeof value !== 'string') return null
  const normalizedValue = value
    .trim()
    .replace(/^"+|"+$/g, '')
    .toLowerCase()
  return normalizedValue === '1' || normalizedValue === 'true'
    ? true
    : normalizedValue === '0' || normalizedValue === 'false'
      ? false
      : null
}

const runAfterTwoFrames = (callback: () => void) => requestAnimationFrame(() => requestAnimationFrame(callback))

const scrollWindowTo = (top = 0) =>
  window.scrollTo({
    top,
    left: 0,
    behavior: 'auto'
  })

const isDesktopMenuCollapsed = ref(parseDesktopMenuState(desktopMenuStateCookie.value) === true)

const siteUrl = computed(() => {
  const configuredBaseUrl = String(runtimeConfig.public.baseURL || '').trim()
  return (configuredBaseUrl ? configuredBaseUrl : `${requestUrl.protocol}//${requestUrl.host}`).replace(/\/+$/, '')
})

const normalizedRoutePath = computed(() => normalizePath(route.path))
const canonicalUrl = computed(() => `${siteUrl.value}${normalizedRoutePath.value}`)
const robotsContent = computed(() => {
  const routeRobots = route.meta.robots
  return typeof routeRobots === 'string' && routeRobots.trim()
    ? routeRobots.trim()
    : import.meta.dev
      ? 'noindex, nofollow'
      : 'index, follow'
})

useHead({
  link: [{ key: 'canonical', rel: 'canonical', href: canonicalUrl }],
  meta: [
    { key: 'robots', name: 'robots', content: robotsContent },
    { key: 'og:type', property: 'og:type', content: 'website' },
    { key: 'og:url', property: 'og:url', content: canonicalUrl }
  ]
})

const sidebarItems = useSidebarItems()

const persistDesktopMenuState = (isCollapsed: boolean) => {
  const nextValue = isCollapsed ? '1' : '0'
  desktopMenuStateCookie.value = nextValue
  if (!import.meta.client) return
  try {
    window.localStorage.setItem(DESKTOP_MENU_STATE_KEY, nextValue)
  } catch {
    return
  }
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

const toggleDesktopMenu = () => {
  isDesktopMenuCollapsed.value = !isDesktopMenuCollapsed.value
  persistDesktopMenuState(isDesktopMenuCollapsed.value)
}

const markHistoryPopNavigation = () => {
  isHistoryPopNavigation.value = true
}

const shouldPreserveScrollPosition = (
  to: { path: string; hash?: string | null },
  from: { path: string; hash?: string | null }
) => normalizePath(to.path) === normalizePath(from.path) && String(to.hash || '') === String(from.hash || '')

const saveScrollPosition = (fullPath: string) => {
  if (!import.meta.client || !fullPath) return
  const scrollContainer = pageScrollRef.value
  if (scrollContainer) pageScrollPositions.value[fullPath] = scrollContainer.scrollTop
  windowScrollPositions.value[fullPath] = window.scrollY
}

const restoreScrollPosition = (fullPath: string) => {
  if (!import.meta.client) return
  const scrollContainer = pageScrollRef.value
  const savedContainerScrollTop = pageScrollPositions.value[fullPath]
  const savedWindowScrollTop = windowScrollPositions.value[fullPath]
  if (scrollContainer) {
    scrollContainer.scrollTop = Number.isFinite(savedContainerScrollTop) ? Number(savedContainerScrollTop) : 0
  }
  scrollWindowTo(Number.isFinite(savedWindowScrollTop) ? Number(savedWindowScrollTop) : 0)
}

const resetScrollPosition = () => {
  if (!import.meta.client) return
  if (pageScrollRef.value) pageScrollRef.value.scrollTop = 0
  scrollWindowTo()
}

watch(() => route.fullPath, closeMobileMenu)

onMounted(() => {
  if ('scrollRestoration' in window.history) window.history.scrollRestoration = 'manual'
  window.addEventListener('popstate', markHistoryPopNavigation)

  removeBeforeEachGuard = router.beforeEach((_, from) => (from.fullPath && saveScrollPosition(from.fullPath), true))

  removeAfterEachHook = router.afterEach(async (to, from) => {
    await nextTick()
    runAfterTwoFrames(() => {
      if (isHistoryPopNavigation.value) restoreScrollPosition(to.fullPath)
      else if (!shouldPreserveScrollPosition(to, from)) resetScrollPosition()
      isHistoryPopNavigation.value = false
    })
  })

  runAfterTwoFrames(() => {
    if (isHistoryPopNavigation.value) restoreScrollPosition(route.fullPath)
    isHistoryPopNavigation.value = false
    isDesktopMenuAnimationEnabled.value = true
  })
})

onBeforeUnmount(() => {
  if (import.meta.client) {
    saveScrollPosition(route.fullPath)
    if ('scrollRestoration' in window.history) window.history.scrollRestoration = 'auto'
    window.removeEventListener('popstate', markHistoryPopNavigation)
  }
  removeBeforeEachGuard?.()
  removeAfterEachHook?.()
})
</script>

<template>
  <div class="min-h-screen bg-(--lab-bg-canvas) text-(--lab-text-secondary) lg:h-screen lg:overflow-hidden">
    <main class="flex min-h-screen w-full flex-col lg:h-full lg:flex-row lg:overflow-hidden">
      <div
        class="sticky top-0 z-40 flex items-center gap-2 border-b border-(--lab-border) bg-(--lab-bg-overlay) lg:hidden"
      >
        <LabBaseButton
          icon="ic:round-menu"
          icon-only
          size="sm"
          variant="ghost"
          class="m-1 rounded-full"
          @click="isMobileMenuOpen = true"
        />
        <div class="min-w-0 flex-1">
          <LabNavBreadcrumb
            v-if="mobileBreadcrumbItems.length"
            :items="mobileBreadcrumbItems"
            container-class="min-w-0"
            list-class="flex min-w-max flex-nowrap items-center gap-2 text-sm text-(--lab-text-secondary)"
            separator-class="h-4 w-4 shrink-0 text-(--lab-text-soft)"
            class="min-w-0"
          >
            <template #append>
              <li class="inline-flex shrink-0 items-center gap-1">
                <LabNavMobileActions v-if="mobileHeaderActions.length" :items="mobileHeaderActions" />
                <div id="mobile-header-actions" class="flex shrink-0 items-center gap-1" />
              </li>
            </template>
          </LabNavBreadcrumb>
        </div>
        <LabAvatar
          :show-label="false"
          link-class="my-1 mr-1 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-transparent p-1 text-(--lab-text-primary) transition-colors hover:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
        />
      </div>
      <LabNavDrawer v-model="isMobileMenuOpen" :items="sidebarItems" />
      <aside class="hidden h-full min-h-0 shrink-0 border-r border-(--lab-border) lg:block">
        <LabNavSidebar
          :items="sidebarItems"
          :collapsed="isDesktopMenuCollapsed"
          :show-toggle="true"
          :animate="isDesktopMenuAnimationEnabled"
          @toggle-collapse="toggleDesktopMenu"
          @request-close="closeMobileMenu"
        />
      </aside>
      <div ref="pageScrollRef" class="min-h-0 min-w-0 lg:flex-1 lg:overflow-y-auto">
        <slot />
        <LabNavFooter :show-controls="false" :show-links="true" />
      </div>
    </main>
  </div>
</template>
