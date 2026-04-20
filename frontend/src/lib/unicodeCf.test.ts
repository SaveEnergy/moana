import { describe, expect, it } from 'vitest'

import { stripCfTrimEdges, stripUnicodeFormatChars } from './unicodeCf'

describe('stripUnicodeFormatChars', () => {
  it('removes Cf code points', () => {
    expect(stripUnicodeFormatChars('\u200bhi\u200b')).toBe('hi')
    expect(stripUnicodeFormatChars('\ufeffx')).toBe('x')
  })

  it('returns empty when the string is only Cf', () => {
    expect(stripUnicodeFormatChars('\u200b\u200c')).toBe('')
  })

  it('leaves strings without Cf unchanged', () => {
    expect(stripUnicodeFormatChars('  ok  ')).toBe('  ok  ')
  })

  it('returns the same instance when there is nothing to strip', () => {
    const s = 'plain'
    expect(stripUnicodeFormatChars(s)).toBe(s)
  })
})

describe('stripCfTrimEdges', () => {
  it('strips Cf then removes edge whitespace', () => {
    expect(stripCfTrimEdges('  \u200bhello\u200b  ')).toBe('hello')
  })

  it('returns the same instance when there is no Cf and edges are tight', () => {
    const s = 'Sure?'
    expect(stripCfTrimEdges(s)).toBe(s)
  })

  it('returns empty when Cf-only and/or whitespace-only after normalization', () => {
    expect(stripCfTrimEdges('\u200b\u200c')).toBe('')
    expect(stripCfTrimEdges('  \t  ')).toBe('')
  })
})
