import { describe, expect, it } from 'vitest'

import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT } from './categoryColor'
import { readCategoryEditRowDataset } from './categoryModalDataset'

describe('readCategoryEditRowDataset', () => {
  it('returns null when id is missing or whitespace-only', () => {
    expect(readCategoryEditRowDataset({} as DOMStringMap)).toBeNull()
    expect(readCategoryEditRowDataset({ id: '  \t  ' } as unknown as DOMStringMap)).toBeNull()
    expect(readCategoryEditRowDataset({ id: '\u00a0' } as unknown as DOMStringMap)).toBeNull()
  })

  it('parses edit row fields and trims strings', () => {
    const ds = {
      id: '  42  ',
      name: '  Books  ',
      color: '  #112233  ',
      custom: '1',
      customHex: '  #aabbcc  ',
      icon: '  star  ',
    } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)).toEqual({
      id: '42',
      name: '  Books  ',
      rawColor: '#112233',
      isCustom: true,
      customHex: '#aabbcc',
      iconVal: 'star',
    })
  })

  it('defaults customHex when absent', () => {
    const ds = { id: '1', custom: '0' } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)?.customHex).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
  })

  it('accepts id "0" (only empty/whitespace id is rejected)', () => {
    const ds = { id: '0', name: 'edge' } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)?.id).toBe('0')
  })

  it('sets isCustom when trimmed data-custom is exactly "1"', () => {
    const ds = { id: '1', custom: 'yes' } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)?.isCustom).toBe(false)
    expect(readCategoryEditRowDataset({ id: '1', custom: ' 1 ' } as unknown as DOMStringMap)?.isCustom).toBe(
      true,
    )
    expect(readCategoryEditRowDataset({ id: '1', custom: '1\t' } as unknown as DOMStringMap)?.isCustom).toBe(
      true,
    )
    expect(readCategoryEditRowDataset({ id: '1', custom: '01' } as unknown as DOMStringMap)?.isCustom).toBe(
      false,
    )
    expect(readCategoryEditRowDataset({ id: '1', custom: ' 0 ' } as unknown as DOMStringMap)?.isCustom).toBe(
      false,
    )
    expect(readCategoryEditRowDataset({ id: '1', custom: '' } as unknown as DOMStringMap)?.isCustom).toBe(false)
  })

  it('keeps empty string customHex when attribute is present but blank (dataset key exists)', () => {
    const ds = { id: '1', customHex: '' } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)?.customHex).toBe('')
  })

  it('trims whitespace-only customHex to empty (openEditModal sanitizes to fallback hex)', () => {
    const ds = { id: '1', customHex: '  \t  ' } as unknown as DOMStringMap
    expect(readCategoryEditRowDataset(ds)?.customHex).toBe('')
  })
})
