import { describe, expect, it } from 'vitest'

import { stripUnicodeFormatChars } from './unicodeCf'

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
