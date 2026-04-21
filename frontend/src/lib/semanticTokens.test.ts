import { readFileSync } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const tokensPath = path.join(path.dirname(fileURLToPath(import.meta.url)), '../styles/01-tokens.css')

describe('01-tokens.css semantic + indicator tokens', () => {
  it('defines --semantic-positive, --semantic-negative, --semantic-negative-muted with stable hex', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(/--semantic-positive:\s*#1b5e20\s*;/)
    expect(src).toMatch(/--semantic-negative:\s*#b91c1c\s*;/)
    expect(src).toMatch(/--semantic-negative-muted:\s*#b85454\s*;/)
  })

  it('defines --indicator-live for dashboard pill-live dot', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(/--indicator-live:\s*#2196f3\s*;/)
  })

  it('defines --surface-tint-warm for category modal preview mix', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(/--surface-tint-warm:\s*#fff8f0\s*;/)
  })

  it('defines pill rail elevation tokens (history vs kind toggle keep distinct alpha)', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(/--shadow-raise-xs:\s*0 1px 3px rgba\(44,\s*47,\s*48,\s*0\.06\)\s*;/)
    expect(src).toMatch(
      /--shadow-raise-xs-strong:\s*0 1px 3px rgba\(44,\s*47,\s*48,\s*0\.08\)\s*;/,
    )
  })

  it('defines primary-tinted control shadows (soft / action / fab)', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(
      /--shadow-primary-soft:\s*0 1px 2px color-mix\(in srgb, var\(--primary\) 28%, transparent\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-primary-action:\s*0 2px 8px color-mix\(in srgb, var\(--primary\) 35%, transparent\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-primary-fab:\s*0 2px 12px color-mix\(in srgb, var\(--primary\) 35%, transparent\)\s*;/,
    )
  })

  it('defines mobile shell drawer and close-control shadows (on-surface color-mix)', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(
      /--shadow-drawer-edge:\s*8px 0 32px color-mix\(in srgb, var\(--on-surface\) 18%, transparent\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-shell-control:\s*0 1px 4px color-mix\(in srgb, var\(--on-surface\) 10%, transparent\)\s*;/,
    )
  })

  it('defines modal, picker stack, and category tile elevation tokens', () => {
    const src = readFileSync(tokensPath, 'utf8')
    expect(src).toMatch(
      /--shadow-modal-panel:\s*0 24px 56px rgba\(44,\s*47,\s*48,\s*0\.14\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-picker-stack:\s*0 10px 40px rgba\(44,\s*47,\s*48,\s*0\.06\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-category-tile:\s*0 1px 2px rgba\(44,\s*47,\s*48,\s*0\.04\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-category-tile-hover:\s*0 4px 14px rgba\(44,\s*47,\s*48,\s*0\.07\)\s*;/,
    )
    expect(src).toMatch(
      /--shadow-category-tile-selected:\s*0 6px 18px rgba\(44,\s*47,\s*48,\s*0\.09\)\s*;/,
    )
  })
})
