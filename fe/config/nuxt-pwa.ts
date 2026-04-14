import { env } from './nuxt-env'
const manifestDisplay = (env.NUXT_PWA_MANIFEST_DISPLAY || 'fullscreen') as
  | 'fullscreen'
  | 'standalone'
  | 'minimal-ui'
  | 'browser'
const manifestDisplayOverride = (
  env.NUXT_PWA_MANIFEST_DISPLAY_OVERRIDE || 'fullscreen,standalone,window-controls-overlay,browser'
)
  .split(',')
  .map((item) => item.trim())
  .filter(Boolean) as Array<'window-controls-overlay' | 'fullscreen' | 'standalone' | 'minimal-ui' | 'browser'>
const manifestOrientation = (env.NUXT_PWA_MANIFEST_ORIENTATION || 'portrait') as
  | 'any'
  | 'natural'
  | 'landscape'
  | 'landscape-primary'
  | 'landscape-secondary'
  | 'portrait'
  | 'portrait-primary'
  | 'portrait-secondary'
const manifestCategories = (env.NUXT_PWA_MANIFEST_CATEGORIES || 'education,productivity,utilities')
  .split(',')
  .map((item) => item.trim())
  .filter(Boolean)
const pwaConfig = {
  strategies: 'injectManifest' as const,
  srcDir: '.',
  filename: 'sw.ts',
  registerType: 'autoUpdate' as const,
  injectRegister: 'auto' as const,
  manifest: {
    id: env.NUXT_PWA_MANIFEST_ID || '/',
    name: env.NUXT_PWA_MANIFEST_NAME || 'sinde',
    short_name: env.NUXT_PWA_MANIFEST_SHORT_NAME || 'sinde',
    description: env.NUXT_PWA_MANIFEST_DESCRIPTION || 'sinde: знания, практики и прикладные инструменты.',
    lang: env.NUXT_PWA_MANIFEST_LANG || 'ru',
    scope: env.NUXT_PWA_MANIFEST_SCOPE || '/',
    start_url: env.NUXT_PWA_MANIFEST_START_URL || '/',
    display: manifestDisplay,
    display_override: manifestDisplayOverride,
    background_color: env.NUXT_PWA_MANIFEST_BACKGROUND_COLOR || '#0c0c0f',
    theme_color: env.NUXT_PWA_MANIFEST_THEME_COLOR || '#0c0c0f',
    orientation: manifestOrientation,
    categories: manifestCategories,
    icons: [
      {
        src: env.NUXT_PWA_ICON_192_SRC || '/pwa-192.png',
        sizes: '192x192',
        type: 'image/png'
      },
      {
        src: env.NUXT_PWA_ICON_512_SRC || '/pwa-512.png',
        sizes: '512x512',
        type: 'image/png'
      },
      {
        src: env.NUXT_PWA_ICON_MASKABLE_512_SRC || '/pwa-maskable-512.png',
        sizes: '512x512',
        type: 'image/png',
        purpose: 'maskable' as const
      }
    ]
  },
  injectManifest: {
    maximumFileSizeToCacheInBytes: 2 * 1024 * 1024,
    globPatterns: ['**/*.{js,css,ico,svg,webmanifest,txt,woff2}'],
    globIgnores: ['**/*.glb', '**/avatars/**/*', '**/recipes/**/*', '**/uploads/**/*', '**/user-content/**/*']
  },
  devOptions: {
    enabled: false,
    suppressWarnings: true,
    type: 'module' as const
  }
}
export { pwaConfig }
