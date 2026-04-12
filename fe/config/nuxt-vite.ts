import tailwindcss from '@tailwindcss/vite'
const viteConfig = {
  sourcemap: false,
  logLevel: 'error' as const,
  clearScreen: false,
  plugins: [tailwindcss() as any],
  build: {
    sourcemap: false,
    chunkSizeWarningLimit: 3500
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

export { viteConfig }
