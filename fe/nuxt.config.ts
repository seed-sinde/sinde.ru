// https://nuxt.com/docs/api/configuration/nuxt-config
import { isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig } from './config/nuxt-vite'
const modules = ['@nuxt/fonts', '@nuxt/icon', 'nuxt-security', '@pinia/nuxt', '@vite-pwa/nuxt', '@nuxt/eslint']
const ignoreWarnings = [
  '@tailwindcss/vite:generate:build',
  'Sourcemap is likely to be incorrect',
  'serverBundle.externalizeIconsJson'
]
// @ts-expect-error stderr.write is monkey-patched to drop known noisy build warnings.
process.stderr.write = ((w) => (c, e, cb) =>
  ignoreWarnings.some((s) => c.toString().includes(s)) ? true : w(c, e, cb))(process.stderr.write.bind(process.stderr))
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: false,
  typescript: {
    typeCheck: true,
    strict: true,
    shim: false,
    tsConfig: {
      compilerOptions: {
        resolveJsonModule: true,
        skipLibCheck: true,
        noUncheckedIndexedAccess: true,
        exactOptionalPropertyTypes: true
      }
    }
  },
  eslint: {
    checker: !isDev,
    config: {
      stylistic: false
    }
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
      externalizeIconsJson: !isDev
    },
    localApiEndpoint: '/_nuxt_icon',
    clientBundle: {
      scan: true
    }
  },
  security: {
    hidePoweredBy: true,
    sri: true,
    removeLoggers: true,
    headers: {
      contentSecurityPolicy: {
        'img-src': ["'self'", 'data:', 'blob:']
      }
    }
  },
  nitro: {
    sourceMap: false,
    preset: 'bun',
    logLevel: 1
  },
  runtimeConfig,
  css: ['~/assets/css/main.css'],
  pwa: pwaConfig,
  hooks: !isDev
    ? {
        close: () => {
          setTimeout(() => process.exit(0), 100)
        }
      }
    : undefined,
  vite: viteConfig
})
