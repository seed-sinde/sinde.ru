// https://nuxt.com/docs/api/configuration/nuxt-config
import { env, isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig } from './config/nuxt-vite'
const isVitest = env.VITEST === '1' || env.VITEST === 'true' || env.NODE_ENV === 'test'
const modules = ['@nuxt/fonts', '@nuxt/icon', '@pinia/nuxt', '@vite-pwa/nuxt']
if (!isVitest) {
  modules.splice(2, 0, 'nuxt-security')
}
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: false,
  typescript: {
    typeCheck: false,
    strict: true,
    tsConfig: {
      exclude: ['../node_modules', '../.git', '../.output', '../dist', '../coverage', '../.cache'],
      compilerOptions: {
        skipLibCheck: true,
        noUnusedLocals: true,
        noUnusedParameters: true,
        noUncheckedIndexedAccess: true,
        exactOptionalPropertyTypes: true,
        useUnknownInCatchVariables: true,
        noFallthroughCasesInSwitch: true,
        noImplicitOverride: true
      }
    },
    sharedTsConfig: {
      exclude: ['../node_modules', '../.git', '../.output', '../dist', '../coverage', '../.cache']
    },
    nodeTsConfig: {
      exclude: ['../node_modules', '../.git', '../.output', '../dist', '../coverage', '../.cache']
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
      externalizeIconsJson: true
    },
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
        exclude: ['../node_modules', '../.git', '../.output', '../dist', '../coverage', '../.cache'],
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
  vite: isVitest
    ? {
        logLevel: 'error'
      }
    : viteConfig
})
