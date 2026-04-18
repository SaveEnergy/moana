import { describe, expect, it } from 'vitest'

import { setRadioCheckedByValue } from './radioMap'

function radio(checked = false): HTMLInputElement {
  return { checked } as HTMLInputElement
}

describe('setRadioCheckedByValue', () => {
  it('checks preferred when the key exists', () => {
    const a = radio()
    const b = radio()
    const map = new Map([
      ['', b],
      ['x', a],
    ])
    expect(setRadioCheckedByValue(map, 'x', '')).toBe(true)
    expect(a.checked).toBe(true)
    expect(b.checked).toBe(false)
  })

  it('uses fallback when preferred is missing', () => {
    const auto = radio()
    const map = new Map([['', auto]])
    expect(setRadioCheckedByValue(map, 'unknown', '')).toBe(true)
    expect(auto.checked).toBe(true)
  })

  it('returns false when neither preferred nor fallback exists', () => {
    const map = new Map<string, HTMLInputElement>()
    expect(setRadioCheckedByValue(map, 'x', '')).toBe(false)
  })

  it('prefers exact empty string over fallback when both could apply', () => {
    const empty = radio()
    const map = new Map([['', empty]])
    expect(setRadioCheckedByValue(map, '', '')).toBe(true)
    expect(empty.checked).toBe(true)
  })
})
