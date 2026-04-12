import tailwindcss from '@tailwindcss/vite'
import { isDev } from './nuxt-env'
const viteConfig = {
  logLevel: 'error' as const,
  clearScreen: false,
  plugins: [tailwindcss() as any],
  build: {
    chunkSizeWarningLimit: 3500
  },
  ...(isDev ?
    {
      optimizeDeps: {
        include: ['tailwind-merge', 'cropperjs']
      }
    }
  : {})
}
export { viteConfig }
