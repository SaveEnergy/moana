/**
 * File input helpers for the settings avatar form (kept in a small module for unit tests).
 */

/** Rejects non-images and SVG; matches server-supported raster types. */
export function isLikelyImageFile(file: File): boolean {
  if (file.type) {
    if (file.type === 'image/svg+xml') {
      return false
    }
    if (file.type.startsWith('image/')) {
      return true
    }
  }
  return /\.(png|jpe?g|gif|webp)$/i.test(file.name)
}

export function setInputFilesAndNotify(fileInput: HTMLInputElement, file: File): void {
  const dt = new DataTransfer()
  dt.items.add(file)
  fileInput.files = dt.files
  fileInput.dispatchEvent(new Event('change', { bubbles: true }))
}
