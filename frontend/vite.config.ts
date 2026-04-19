/// <reference types="vitest/config" />

import { defineConfig } from 'vite'
import { dirname, resolve } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => ({
  root: __dirname,
  esbuild: {
    legalComments: 'none',
    /* Production bundle: smaller parse + no stray debugger pauses in shipped `app.js`. */
    ...(mode === 'production' ? { drop: ['console', 'debugger'] as const } : {}),
  },
  build: {
    target: 'es2022',
    cssMinify: true,
    reportCompressedSize: false,
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
  test: {
    root: __dirname,
    environment: 'node',
    include: ['src/**/*.test.ts'],
    passWithNoTests: false,
  },
}))
