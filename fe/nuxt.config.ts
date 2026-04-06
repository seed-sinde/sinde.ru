// https://nuxt.com/docs/api/configuration/nuxt-config
import { env, isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig, viteHooks } from './config/nuxt-vite'

const isVitest = env.VITEST === '1' || env.VITEST === 'true' || env.NODE_ENV === 'test'
const modules = ['@nuxt/fonts', '@nuxt/icon', '@pinia/nuxt', '@vite-pwa/nuxt']

if (!isVitest) {
  modules.splice(2, 0, 'nuxt-security')
}

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: {
    server: 'hidden',
    client: 'hidden'
  },
  typescript: {
    typeCheck: true,
    strict: true,
    tsConfig: {
      compilerOptions: {
        noUnusedLocals: true,
        noUnusedParameters: true,
        noUncheckedIndexedAccess: true,
        exactOptionalPropertyTypes: true,
        useUnknownInCatchVariables: true,
        noFallthroughCasesInSwitch: true,
        noImplicitOverride: true
      }
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
    localApiEndpoint: '/_nuxt_icon',
    clientBundle: {
      scan: true
    }
  },
  security: isVitest
    ? undefined
    : {
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
    preset: 'bun',
    typescript: {
      strict: true,
      tsConfig: {
        compilerOptions: {
          noUnusedLocals: true,
          noUnusedParameters: true,
          noUncheckedIndexedAccess: true,
          exactOptionalPropertyTypes: true,
          useUnknownInCatchVariables: true,
          noFallthroughCasesInSwitch: true
        }
      }
    }
  },
  runtimeConfig,
  css: ['~/assets/css/main.css'],
  pwa: pwaConfig,
  hooks: isVitest ? {} : viteHooks,
  vite: isVitest
    ? {
        logLevel: 'error'
      }
    : viteConfig
})
