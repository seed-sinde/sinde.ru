import tailwindcss from "@tailwindcss/vite"
import Icons from "unplugin-icons/vite"
import {isDev} from "./nuxt-env"
const viteConfig = {
  clearScreen: false,
  plugins: [
    tailwindcss(),
    Icons({
      autoInstall: true,
      scale: 1,
      defaultClass: "text-lg h-5 w-5 shrink-0"
    })
  ],
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
