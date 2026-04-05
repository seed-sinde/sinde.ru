import tailwindcss from '@tailwindcss/vite'
import { createLogger } from 'vite'
import type { InlineConfig, Logger } from 'vite'
import type { LoggingFunction, RollupLog } from 'rollup'
function suppressSourcemapWarnings(logger: Logger): Logger {
  return {
    ...logger,
    warn(msg, options) {
      // Temporary workaround for Nuxt 4.4.2 sourcemap noise.
      // Remove this filter once these warnings stop appearing in normal builds.
      if (
        msg.includes('Sourcemap is likely to be incorrect') ||
        msg.includes('[plugin @tailwindcss/vite:generate:build]')
      ) {
        return
      }
      logger.warn(msg, options)
    }
  }
}
const viteLogger = createLogger()
const viteCustomLogger = suppressSourcemapWarnings({
  ...viteLogger,
  warn(msg, options) {
    // Temporary workaround for Nuxt 4.4.2 + @tailwindcss/vite sourcemap noise.
    // Remove this filter once this warning stops appearing in normal builds.
    if (msg.includes('[plugin @tailwindcss/vite:generate:build] Sourcemap is likely to be incorrect')) {
      return
    }
    viteLogger.warn(msg, options)
  }
})
const viteHooks = {
  'vite:extendConfig'(config: InlineConfig) {
    const mutableConfig = config as InlineConfig & { customLogger?: Logger }
    if (mutableConfig.customLogger) {
      mutableConfig.customLogger = suppressSourcemapWarnings(mutableConfig.customLogger)
    }
  }
}
const viteConfig = {
  customLogger: viteCustomLogger,
  logLevel: 'error' as const,
  plugins: [tailwindcss() as any],
  build: {
    chunkSizeWarningLimit: 3500,
    rollupOptions: {
      onwarn(warning: RollupLog, warn: LoggingFunction) {
        // Temporary workaround for Nuxt 4.4.2 sourcemap noise.
        // Remove this filter once this warning stops appearing in normal builds.
        if (warning.code === 'SOURCEMAP_BROKEN') {
          return
        }
        warn(warning)
      }
    }
  },
  optimizeDeps: {
    include: [
      'tailwind-merge',
      'cropperjs',
      'highlight.js/lib/common',
      'highlightjs-vue',
      'highlight.js/lib/languages/stylus'
    ]
  }
}
export { viteConfig, viteHooks }
