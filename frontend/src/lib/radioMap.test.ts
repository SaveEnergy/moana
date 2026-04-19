import { describe, expect, it } from 'vitest'

import { buildRadioMapByValue, getFormRadioGroupValue, setRadioCheckedByValue } from './radioMap'

function radio(checked = false): HTMLInputElement {
  return { checked } as HTMLInputElement
}

describe('buildRadioMapByValue', () => {
  it('indexes radios by value for the given selector', () => {
    const a = { value: 'x' } as HTMLInputElement
    const b = { value: '' } as HTMLInputElement
    const scope = {
      querySelectorAll: () => [a, b] as unknown as NodeListOf<HTMLInputElement>,
    } as unknown as ParentNode
    const m = buildRadioMapByValue(scope, 'input[type="radio"]')
    expect(m.get('x')).toBe(a)
    expect(m.get('')).toBe(b)
    expect(m.size).toBe(2)
  })

  it('last duplicate value wins', () => {
    const first = { value: 'dup' } as HTMLInputElement
    const second = { value: 'dup' } as HTMLInputElement
    const scope = {
      querySelectorAll: () => [first, second] as unknown as NodeListOf<HTMLInputElement>,
    } as unknown as ParentNode
    const m = buildRadioMapByValue(scope, 'input[type="radio"]')
    expect(m.get('dup')).toBe(second)
  })
})

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

  it('uses a non-default fallback key when preferred is missing', () => {
    const y = radio()
    const map = new Map<string, HTMLInputElement>([
      ['z', radio()],
      ['y', y],
    ])
    expect(setRadioCheckedByValue(map, 'missing', 'y')).toBe(true)
    expect(y.checked).toBe(true)
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

describe('getFormRadioGroupValue', () => {
  it('returns RadioNodeList.value when namedItem resolves a radio group', () => {
    const form = {
      elements: {
        namedItem: () => ({ value: 'preset' } as RadioNodeList),
      },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('preset')
  })

  it('returns empty string when no control matches the name', () => {
    const form = {
      elements: { namedItem: () => null },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('')
  })

  it('returns empty string when namedItem does not expose a string value', () => {
    const form = {
      elements: {
        namedItem: () => ({}) as unknown as Element,
      },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('')
  })
})
