import { describe, expect, it } from 'vitest'

import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT } from './categoryColor'
import { readCategoryEditRowDataset } from './categoryModalDataset'

describe('readCategoryEditRowDataset', () => {
  it('returns null when id is missing or whitespace-only', () => {
    expect(readCategoryEditRowDataset({} as DOMStringMap)).toBeNull()
    expect(readCategoryEditRowDataset({ id: '  \t  ' } as unknown as DOMStringMap)).toBeNull()
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
})
