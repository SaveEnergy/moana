import { describe, expect, it, vi } from 'vitest'

import { buildRadioMapByValue, getFormRadioGroupValue, setRadioCheckedByValue } from './radioMap'

function radio(checked = false): HTMLInputElement {
  return { checked } as HTMLInputElement
}

describe('buildRadioMapByValue', () => {
  it('forwards radiosSelector to querySelectorAll on scope', () => {
    let seen = ''
    const scope = {
      querySelectorAll: (sel: string) => {
        seen = sel
        return [] as unknown as NodeListOf<HTMLInputElement>
      },
    } as unknown as ParentNode
    buildRadioMapByValue(scope, 'input[type="radio"][name="color"]')
    expect(seen).toBe('input[type="radio"][name="color"]')
  })

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

  it('indexes by trimmed value so spaced duplicates collapse (last wins)', () => {
    const first = { value: '  x  ' } as HTMLInputElement
    const second = { value: 'x' } as HTMLInputElement
    const scope = {
      querySelectorAll: () => [first, second] as unknown as NodeListOf<HTMLInputElement>,
    } as unknown as ParentNode
    const m = buildRadioMapByValue(scope, 'input[type="radio"]')
    expect(m.get('x')).toBe(second)
    expect(m.size).toBe(1)
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

  it('does not hit the map twice when preferred and fallback trim to the same key', () => {
    const target = radio()
    const map = new Map<string, HTMLInputElement>([['x', target]])
    const getSpy = vi.spyOn(map, 'get')
    expect(setRadioCheckedByValue(map, 'x', 'x')).toBe(true)
    expect(target.checked).toBe(true)
    expect(getSpy).toHaveBeenCalledTimes(1)
    getSpy.mockRestore()
  })

  it('trims preferred and fallback keys when resolving map entries', () => {
    const target = radio()
    const map = new Map([['preset', target]])
    expect(setRadioCheckedByValue(map, '  preset  ', '  other  ')).toBe(true)
    expect(target.checked).toBe(true)
  })

  it('does not assign checked when the resolved input is already selected', () => {
    let checkedSets = 0
    let c = true
    const input = {
      get checked() {
        return c
      },
      set checked(v: boolean) {
        checkedSets++
        c = v
      },
    } as unknown as HTMLInputElement
    const map = new Map([['x', input]])
    expect(setRadioCheckedByValue(map, 'x', '')).toBe(true)
    expect(checkedSets).toBe(0)
    expect(input.checked).toBe(true)
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

  it('returns value when namedItem is a single control with a string value (not only RadioNodeList)', () => {
    const form = {
      elements: {
        namedItem: () => ({ value: 'one-off' }) as HTMLInputElement,
      },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('one-off')
  })

  it('trims string values (radio group / single control)', () => {
    const form = {
      elements: {
        namedItem: () => ({ value: '  preset  ' } as RadioNodeList),
      },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('preset')
  })

  it('returns empty string when value is present but not a string', () => {
    const form = {
      elements: {
        namedItem: () => ({ value: 42 }) as unknown as Element,
      },
    } as unknown as HTMLFormElement
    expect(getFormRadioGroupValue(form, 'color')).toBe('')
  })

  it('uses trusted change target when name matches (skips namedItem)', () => {
    const namedItem = vi.fn()
    const form = {
      elements: { namedItem },
    } as unknown as HTMLFormElement
    const target = { name: 'color', value: 'from-target' } as HTMLInputElement
    expect(getFormRadioGroupValue(form, 'color', target)).toBe('from-target')
    expect(namedItem).not.toHaveBeenCalled()
  })

  it('falls back to namedItem when trusted target name differs', () => {
    const namedItem = vi.fn(() => ({ value: 'from-form' } as RadioNodeList))
    const form = {
      elements: { namedItem },
    } as unknown as HTMLFormElement
    const target = { name: 'icon', value: 'wrong' } as HTMLInputElement
    expect(getFormRadioGroupValue(form, 'color', target)).toBe('from-form')
    expect(namedItem).toHaveBeenCalledWith('color')
  })

  it('trims trusted change target value', () => {
    const namedItem = vi.fn()
    const form = {
      elements: { namedItem },
    } as unknown as HTMLFormElement
    const target = { name: 'color', value: '  x  ' } as HTMLInputElement
    expect(getFormRadioGroupValue(form, 'color', target)).toBe('x')
    expect(namedItem).not.toHaveBeenCalled()
  })
})
