import { defineConfig } from 'vite'
import { dirname, resolve } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  root: __dirname,
  build: {
    outDir: resolve(__dirname, '../internal/assets/static'),
    emptyOutDir: false,
    rollupOptions: {
      input: resolve(__dirname, 'src/main.ts'),
      output: {
        entryFileNames: 'js/app.js',
        assetFileNames: (assetInfo) => {
          const n = assetInfo.names?.[0] ?? assetInfo.name
          if (typeof n === 'string' && n.endsWith('.css')) {
            return 'css/app.css'
          }
          return 'assets/[name][extname]'
        },
      },
    },
  },
})
