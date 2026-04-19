import { describe, expect, it } from 'vitest'

import viteConfig from '../../vite.config'

describe('vite.config production esbuild', () => {
  function resolveConfig(mode: 'production' | 'development') {
    if (typeof viteConfig !== 'function') {
      throw new Error('expected defineConfig callback export')
    }
    return viteConfig({ command: 'build', mode, isSsrBuild: false })
  }

  it('drops console and debugger only in production', () => {
    const prod = resolveConfig('production')
    expect(prod.esbuild).toMatchObject({
      legalComments: 'none',
      drop: ['console', 'debugger'],
    })
    const dev = resolveConfig('development')
    expect(dev.esbuild).toEqual({ legalComments: 'none' })
  })
})
