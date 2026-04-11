import tailwindcss from '@tailwindcss/vite'

const viteConfig = {
  logLevel: 'error' as const,
  plugins: [tailwindcss() as any],
  build: {
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
