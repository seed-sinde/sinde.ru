// https://nuxt.com/docs/api/configuration/nuxt-config
import { env, isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig } from './config/nuxt-vite'
const isVitest = env.VITEST === '1' || env.VITEST === 'true' || env.NODE_ENV === 'test'
const modules = ['@nuxt/fonts', '@nuxt/icon', '@pinia/nuxt', '@vite-pwa/nuxt']
if (!isVitest) {
  modules.splice(2, 0, 'nuxt-security')
}
// Перехват стандартного потока ошибок, куда Tailwind шлет ворнинги
if (process.env.NODE_ENV === 'production') {
  const originalStderrWrite = process.stderr.write
  // @ts-ignore
  process.stderr.write = function (chunk, encoding, callback) {
    const message = chunk.toString()
    // Список фраз для блокировки
    if (
      message.includes('@tailwindcss/vite:generate:build') ||
      message.includes('Sourcemap is likely to be incorrect') ||
      message.includes('serverBundle.externalizeIconsJson')
    ) {
      return true // Просто поглощаем сообщение, не выводя его
    }
    return originalStderrWrite.call(process.stderr, chunk, encoding, callback)
  }
}
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: false,
  typescript: {
    typeCheck: true,
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
    logLevel: 1,
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
  hooks: {
    close: () => {
      setTimeout(() => process.exit(0), 100)
    }
  },
  vite:
    isVitest ?
      {
        logLevel: 'error'
      }
    : viteConfig
})
