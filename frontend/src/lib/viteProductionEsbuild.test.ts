import { describe, expect, it } from 'vitest'

import viteConfig from '../../vite.config'

function resolveViteConfig(mode: 'production' | 'development') {
  if (typeof viteConfig !== 'function') {
    throw new Error('expected defineConfig callback export')
  }
  return viteConfig({ command: 'build', mode, isSsrBuild: false })
}

describe('vite.config production esbuild', () => {
  it('drops console and debugger only in production', () => {
    const prod = resolveViteConfig('production')
    expect(prod.esbuild).toMatchObject({
      legalComments: 'none',
      drop: ['console', 'debugger'],
    })
    const dev = resolveViteConfig('development')
    expect(dev.esbuild).toEqual({ legalComments: 'none' })
  })

  it('keeps es2022 target, css minify, and non-destructive static outDir (design.md build section)', () => {
    const prod = resolveViteConfig('production')
    expect(prod.build?.target).toBe('es2022')
    expect(prod.build?.cssMinify).toBe(true)
    expect(prod.build?.emptyOutDir).toBe(false)
    expect(prod.build?.reportCompressedSize).toBe(false)
    expect(String(prod.build?.outDir).replace(/\\/g, '/')).toMatch(/internal\/assets\/static$/)
  })

  it('resolves main.ts input and fixed js/app.js + css/app.css output names', () => {
    const prod = resolveViteConfig('production')
    expect(String(prod.build?.rollupOptions?.input)).toMatch(/main\.ts$/)
    const output = prod.build?.rollupOptions?.output
    const out = Array.isArray(output) ? output[0] : output
    expect(out).toMatchObject({
      entryFileNames: 'js/app.js',
    })
    const assetFileNames = out?.assetFileNames
    expect(typeof assetFileNames).toBe('function')
    const fn = assetFileNames as (info: { names?: string[]; name?: string }) => string
    expect(fn({ names: ['app.css'] })).toBe('css/app.css')
    expect(fn({ name: 'app.css' })).toBe('css/app.css')
    expect(fn({ name: 'other.png' })).toBe('assets/[name][extname]')
  })
})

describe('vite.config vitest block', () => {
  it('runs unit tests from frontend root in node without allowing empty suites', () => {
    const cfg = resolveViteConfig('development')
    expect(cfg.test).toMatchObject({
      environment: 'node',
      include: ['src/**/*.test.ts'],
      passWithNoTests: false,
    })
    expect(String(cfg.test?.root).replace(/\\/g, '/')).toMatch(/\/frontend$/)
  })
})
