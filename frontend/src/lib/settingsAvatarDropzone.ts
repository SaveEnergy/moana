/**
 * Wires the profile photo modal: no visible file “Upload” / browse field — a round preview area
 * opens the file picker (click) or accepts drops, then the user uses **Save photo** to submit.
 */

import { isLikelyImageFile, setInputFilesAndNotify } from './settingsAvatarFileInput'

const settingsAvatarDropzoneWired = new WeakSet<HTMLFormElement>()

export function initSettingsAvatarDropzone(dialog: HTMLDialogElement): void {
  const form = dialog.querySelector<HTMLFormElement>('#settings-avatar-upload-form')
  const fileInput = dialog.querySelector<HTMLInputElement>('#settings-avatar-file')
  const dropzone = dialog.querySelector<HTMLElement>('#settings-avatar-dropzone')
  const saveBtn = dialog.querySelector<HTMLButtonElement>('#settings-avatar-save')
  if (!form || !fileInput || !dropzone) {
    return
  }
  if (settingsAvatarDropzoneWired.has(form)) {
    return
  }
  settingsAvatarDropzoneWired.add(form)

  const inner = dropzone.querySelector<HTMLElement>('.settings-avatar-dropzone-inner')
  if (!inner) {
    return
  }
  const initialPreviewHTML = inner.innerHTML
  let localPreviewURL: string | null = null

  const revokeLocalPreviewURL = () => {
    if (localPreviewURL) {
      URL.revokeObjectURL(localPreviewURL)
      localPreviewURL = null
    }
  }

  const restoreServerPreview = () => {
    revokeLocalPreviewURL()
    inner.innerHTML = initialPreviewHTML
  }

  const showFileInPreview = (file: File) => {
    if (!isLikelyImageFile(file)) {
      return
    }
    revokeLocalPreviewURL()
    localPreviewURL = URL.createObjectURL(file)
    inner.innerHTML = ''
    const img = document.createElement('img')
    img.className = 'settings-avatar-dialog-preview-img'
    img.src = localPreviewURL
    img.alt = ''
    img.width = 200
    img.height = 200
    img.decoding = 'async'
    img.setAttribute('draggable', 'false')
    img.loading = 'eager'
    inner.appendChild(img)
  }

  const setSaveState = () => {
    if (saveBtn) {
      saveBtn.disabled = !fileInput.files || fileInput.files.length === 0
    }
  }
  setSaveState()

  const onPickOrDrop = (file: File) => {
    if (!isLikelyImageFile(file)) {
      return
    }
    setInputFilesAndNotify(fileInput, file)
  }

  dropzone.addEventListener('click', () => {
    fileInput.click()
  })
  dropzone.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault()
      fileInput.click()
    }
  })

  fileInput.addEventListener('change', () => {
    setSaveState()
    const f = fileInput.files?.[0]
    if (f) {
      showFileInPreview(f)
    } else {
      restoreServerPreview()
    }
  })

  let dragDepth = 0
  const dragActive = 'settings-avatar-dropzone--dragover'
  dropzone.addEventListener('dragenter', (e) => {
    e.preventDefault()
    dragDepth += 1
    dropzone.classList.add(dragActive)
  })
  dropzone.addEventListener('dragleave', (e) => {
    e.preventDefault()
    dragDepth -= 1
    if (dragDepth <= 0) {
      dragDepth = 0
      dropzone.classList.remove(dragActive)
    }
  })
  dropzone.addEventListener('dragover', (e) => {
    e.preventDefault()
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'copy'
    }
  })
  dropzone.addEventListener('drop', (e) => {
    e.preventDefault()
    dragDepth = 0
    dropzone.classList.remove(dragActive)
    const f = e.dataTransfer?.files?.[0]
    if (f) {
      onPickOrDrop(f)
    }
  })

  const resetPickedFile = () => {
    fileInput.value = ''
    setSaveState()
    restoreServerPreview()
  }

  dialog.addEventListener('close', () => {
    resetPickedFile()
  })
}
