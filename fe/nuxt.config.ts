// https://nuxt.com/docs/api/configuration/nuxt-config
import { isDev, runtimeConfig } from './config/nuxt-env'
import { pwaConfig } from './config/nuxt-pwa'
import { viteConfig, viteHooks } from './config/nuxt-vite'
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  sourcemap: {
    server: 'hidden',
    client: 'hidden'
  },
  experimental: {
    payloadExtraction: !isDev
  },
  imports: {
    dirs: ['data', 'constants']
  },
  modules: ['@nuxt/fonts', '@nuxt/icon', 'nuxt-security', '@pinia/nuxt', '@vite-pwa/nuxt'],
  icon: {
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
    preset: 'bun'
  },
  runtimeConfig,
  css: ['~/assets/css/main.css'],
  pwa: pwaConfig,
  hooks: viteHooks,
  vite: viteConfig
})
