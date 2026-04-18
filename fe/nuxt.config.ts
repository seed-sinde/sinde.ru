// https://nuxt.com/docs/api/configuration/nuxt-config
import {defineNuxtConfig} from "nuxt/config"
import {isDev, runtimeConfig} from "./config/nuxt-env"
import {pwaConfig} from "./config/nuxt-pwa"
import {viteConfig} from "./config/nuxt-vite"
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  modules: ["@nuxt/fonts", "@nuxt/icon", "nuxt-security", "@pinia/nuxt", "@vite-pwa/nuxt", "@nuxt/eslint"],
  debug: false,
  devtools: {enabled: false},
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
    payloadExtraction: !isDev,
    serverAppConfig: false
  },
  imports: {
    dirs: ["data", "constants"]
  },
  icon: {
    serverBundle: {
      collections: ["ic"],
      externalizeIconsJson: !isDev
    },
    localApiEndpoint: "/_nuxt_icon",
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
        "img-src": ["'self'", "data:", "blob:"]
      }
    }
  },
  nitro: {
    sourceMap: false
  },
  runtimeConfig,
  css: ["~/assets/css/main.css"],
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
