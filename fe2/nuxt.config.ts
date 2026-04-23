// https://nuxt.com/docs/api/configuration/nuxt-config
import {defineNuxtConfig} from "nuxt/config"
import {isDev, runtimeConfig} from "./config/nuxt-env"
import {pwaConfig} from "./config/nuxt-pwa"
import {viteConfig} from "./config/nuxt-vite"
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  modules: ["unplugin-icons/nuxt", "nuxt-security", "@pinia/nuxt", "@vite-pwa/nuxt", "@nuxt/eslint"],
  ssr: true,
  devtools: {enabled: false},
  experimental: {
    serverAppConfig: false
  },
  css: ["~/assets/css/main.css"],
  eslint: {
    checker: !isDev,
    config: {
      stylistic: false
    }
  },
  security: {
    hidePoweredBy: true,
    sri: !isDev,
    removeLoggers: true,
    headers: {
      contentSecurityPolicy: {
        "upgrade-insecure-requests": !isDev,
        "img-src": ["'self'", "data:", "blob:"]
      }
    }
  },
  runtimeConfig,
  pwa: pwaConfig,
  vite: viteConfig
})
