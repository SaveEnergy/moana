import { describe, expect, it } from 'vitest'

import { isLikelyImageFile } from './settingsAvatarFileInput'

describe('isLikelyImageFile', () => {
  it('accepts common rasters and rejects SVG', () => {
    expect(isLikelyImageFile(new File(['x'], 'a.png', { type: 'image/png' }))).toBe(true)
    expect(isLikelyImageFile(new File(['x'], 'a.jpg', { type: 'image/jpeg' }))).toBe(true)
    expect(isLikelyImageFile(new File(['<svg'], 'a.svg', { type: 'image/svg+xml' }))).toBe(false)
  })

  it('uses extension when MIME is missing (some drag-drop / OS paths)', () => {
    expect(isLikelyImageFile(new File(['x'], 'b.jpeg', { type: '' }))).toBe(true)
    expect(isLikelyImageFile(new File(['x'], 'b.webp', { type: 'application/octet-stream' }))).toBe(true)
  })
})
