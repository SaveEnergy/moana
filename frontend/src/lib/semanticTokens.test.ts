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
})
