/// <reference lib="webworker" />
import { cacheNames, clientsClaim, setCacheNameDetails } from 'workbox-core'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'
import { ExpirationPlugin } from 'workbox-expiration'
import type { WorkboxPlugin } from 'workbox-core/types.js'
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { StaleWhileRevalidate } from 'workbox-strategies'
declare let self: ServiceWorkerGlobalScope & {
  __WB_MANIFEST: Array<string | { url: string; revision?: string | null }>
}
const CACHE_VERSION = '20260325'
const ICON_CACHE_NAME = `sinde-ui-icons-${CACHE_VERSION}`
const CURRENT_CACHE_NAMES = new Set<string>([ICON_CACHE_NAME])
const OBSOLETE_CACHE_PATTERNS = [
  /^static-assets$/i,
  /^html-pages$/i,
  /^workbox-precache/i,
  /^sinde-(?:precache|runtime|ui-icons)-/i
]
const ICON_CACHE_PLUGINS: WorkboxPlugin[] = [
  new CacheableResponsePlugin({
    statuses: [0, 200]
  }) as WorkboxPlugin,
  new ExpirationPlugin({
    maxEntries: 64,
    maxAgeSeconds: 60 * 60 * 24 * 30,
    purgeOnQuotaError: true
  }) as WorkboxPlugin
]
setCacheNameDetails({
  prefix: 'sinde',
  precache: 'precache',
  runtime: 'runtime',
  suffix: CACHE_VERSION
})
CURRENT_CACHE_NAMES.add(cacheNames.precache)
CURRENT_CACHE_NAMES.add(cacheNames.runtime)
self.skipWaiting()
clientsClaim()
precacheAndRoute(self.__WB_MANIFEST, {
  ignoreURLParametersMatching: [/^utm_/, /^fbclid$/]
})
cleanupOutdatedCaches()
registerRoute(
  ({ url }) => url.origin === self.location.origin && url.pathname.startsWith('/_nuxt_icon'),
  new StaleWhileRevalidate({
    cacheName: ICON_CACHE_NAME,
    plugins: ICON_CACHE_PLUGINS
  })
)
self.addEventListener('activate', event => {
  event.waitUntil(
    (async () => {
      const keys = await caches.keys()
      await Promise.all(
        keys.map(cacheKey => {
          const isCurrent = CURRENT_CACHE_NAMES.has(cacheKey)
          const isObsolete = OBSOLETE_CACHE_PATTERNS.some(pattern => pattern.test(cacheKey))
          if (isCurrent || !isObsolete) return Promise.resolve(false)
          return caches.delete(cacheKey)
        })
      )
    })()
  )
})
export {}
