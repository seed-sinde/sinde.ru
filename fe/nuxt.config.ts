// https://nuxt.com/docs/api/configuration/nuxt-config
import { env, isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig } from './config/nuxt-vite'
const isVitest = env.VITEST === '1' || env.VITEST === 'true' || env.NODE_ENV === 'test'
const modules = ['@nuxt/fonts', '@nuxt/icon', ...(!isVitest ? ['nuxt-security'] : []), '@pinia/nuxt', '@vite-pwa/nuxt']
const ignoreWarnings = [
  '@tailwindcss/vite:generate:build',
  'Sourcemap is likely to be incorrect',
  'serverBundle.externalizeIconsJson'
]
// @ts-ignore
process.stderr.write = (w => (c, e, cb) => (ignoreWarnings.some(s => c.toString().includes(s)) ? true : w(c, e, cb)))(
  process.stderr.write.bind(process.stderr)
)
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: false,
  typescript: {
    typeCheck: false,
    strict: true
  },
  experimental: {
    payloadExtraction: !isDev
  },
  imports: {
    dirs: ['data', 'constants']
  },
  modules,
  icon: {
    serverBundle: {
      collections: ['ic'],
      externalizeIconsJson: true
    },
    localApiEndpoint: '/_nuxt_icon',
    clientBundle: {
      scan: true
    }
  },
  security:
    isVitest ? undefined : (
      {
        hidePoweredBy: true,
        sri: true,
        removeLoggers: true,
        headers: {
          contentSecurityPolicy: {
            'img-src': ["'self'", 'data:', 'blob:']
          }
        }
      }
    ),
  nitro: {
    sourceMap: false,
    preset: 'bun',
    logLevel: 1
  },
  runtimeConfig,
  css: ['~/assets/css/main.css'],
  pwa: pwaConfig,
  hooks:
    !isDev ?
      {
        close: () => {
          setTimeout(() => process.exit(0), 100)
        }
      }
    : undefined,
  vite:
    isVitest ?
      {
        logLevel: 'error'
      }
    : viteConfig
})
