import tailwindcss from "@tailwindcss/vite"
import {isDev} from "./nuxt-env"
const viteConfig = {
  clearScreen: false,
  plugins: [tailwindcss()],
  build: {
    chunkSizeWarningLimit: 3500
  },
  ...(isDev
    ? {
        optimizeDeps: {
          include: ["cropperjs"]
        }
      }
    : {})
}
export {viteConfig}
