import { CATEGORY_MODAL_DEFAULT_PREVIEW_BG, sanitizeCategoryCustomHex } from './categoryColor'

export function initCategoryModal(): void {
  const dialog = document.getElementById('cat-modal') as HTMLDialogElement | null
  if (!dialog) {
    return
  }

  const form = document.getElementById('cat-modal-form') as HTMLFormElement | null
  const idInput = document.getElementById('cat-modal-id') as HTMLInputElement | null
  const titleEl = document.getElementById('cat-modal-title')
  const submitBtn = document.getElementById('cat-modal-submit')
  const preview = document.getElementById('cat-modal-preview')
  const iconWrap = document.getElementById('cat-modal-preview-icon')
  const nameInput = document.getElementById('cat-modal-name') as HTMLInputElement | null
  const addCategoryBtn = document.getElementById('cat-modal-open-create')

  if (!form || !idInput || !titleEl || !submitBtn || !preview || !iconWrap || !nameInput) {
    return
  }

  /* Narrow once for nested functions (TS does not always narrow captured consts in closures). */
  const modal = dialog
  const catForm = form
  const catId = idInput
  const catTitle = titleEl
  const catSubmit = submitBtn
  const catPreview = preview
  const catIconWrap = iconWrap
  const catName = nameInput

  function radiosByValue(name: 'color' | 'icon'): Map<string, HTMLInputElement> {
    const m = new Map<string, HTMLInputElement>()
    for (const r of catForm.querySelectorAll<HTMLInputElement>(`input[name="${name}"]`)) {
      m.set(r.value, r)
    }
    return m
  }

  const colorRadioByValue = radiosByValue('color')
  const iconRadioByValue = radiosByValue('icon')

  function syncCatModalPreview() {
    const cr = catForm.querySelector<HTMLInputElement>('input[name="color"]:checked')
    let bg: string = CATEGORY_MODAL_DEFAULT_PREVIEW_BG
    if (cr?.value === 'custom') {
      const nat = catForm.querySelector<HTMLInputElement>('#cat-modal-color-native')
      bg = nat?.value?.trim() || '#818cf8'
    } else if (cr?.value) {
      bg = cr.value
    }
    catPreview.style.background = bg

    const ir = catForm.querySelector<HTMLInputElement>('input[name="icon"]:checked')
    catIconWrap.innerHTML = ''
    if (!ir?.value) {
      catIconWrap.classList.add('cat-modal-preview-icon--auto')
      catIconWrap.textContent = 'A'
      return
    }
    catIconWrap.classList.remove('cat-modal-preview-icon--auto')
    const label = ir.closest('label')
    const svg = label?.querySelector('svg.moana-icon')
    if (svg) {
      const clone = svg.cloneNode(true) as SVGElement
      clone.classList.add('moana-icon--cat-preview')
      catIconWrap.appendChild(clone)
    }
  }

  function wireCategoryFormPreview() {
    catForm.addEventListener('input', (e) => {
      const t = e.target
      if (!(t instanceof Element) || !t.classList.contains('cat-color-native')) {
        return
      }
      const wrap = t.closest('.cat-color-swatch--custom')
      const r = wrap?.querySelector<HTMLInputElement>('input[type="radio"][value="custom"]')
      if (r) {
        r.checked = true
        syncCatModalPreview()
      }
    })
    catForm.addEventListener('change', (e) => {
      const t = e.target
      if (!(t instanceof HTMLInputElement)) {
        return
      }
      if (t.name === 'color' || t.name === 'icon') {
        syncCatModalPreview()
      }
    })
  }

  wireCategoryFormPreview()

  function openCreateModal() {
    catForm.action = '/categories'
    catId.value = ''
    catTitle.textContent = 'New category'
    catSubmit.textContent = 'Create category'
    catForm.reset()
    colorRadioByValue.get('')!.checked = true
    iconRadioByValue.get('')!.checked = true
    const nat = catForm.querySelector<HTMLInputElement>('#cat-modal-color-native')
    if (nat) nat.value = '#818cf8'
    syncCatModalPreview()
    catName.focus()
    modal.showModal()
  }

  function openEditModal(btn: HTMLElement) {
    const id = btn.dataset.id
    if (!id) return
    catForm.action = '/categories/update'
    catId.value = id
    catTitle.textContent = 'Edit category'
    catSubmit.textContent = 'Save changes'

    catName.value = btn.dataset.name ?? ''

    const rawColor = (btn.dataset.color ?? '').trim()
    const isCustom = btn.dataset.custom === '1'
    const customHex = (btn.dataset.customHex ?? '#818cf8').trim()

    if (!rawColor) {
      colorRadioByValue.get('')!.checked = true
    } else if (isCustom) {
      colorRadioByValue.get('custom')!.checked = true
      const nat = catForm.querySelector<HTMLInputElement>('#cat-modal-color-native')
      if (nat) nat.value = sanitizeCategoryCustomHex(customHex)
    } else if (colorRadioByValue.has(rawColor)) {
      colorRadioByValue.get(rawColor)!.checked = true
    } else {
      colorRadioByValue.get('')!.checked = true
    }

    const iconVal = (btn.dataset.icon ?? '').trim()
    if (!iconVal) {
      iconRadioByValue.get('')!.checked = true
    } else if (iconRadioByValue.has(iconVal)) {
      iconRadioByValue.get(iconVal)!.checked = true
    } else {
      iconRadioByValue.get('')!.checked = true
    }

    syncCatModalPreview()
    catName.focus()
    modal.showModal()
  }

  addCategoryBtn?.addEventListener('click', () => openCreateModal())

  /** Scoped to the list card so topbar/sidebar clicks do not run this handler. */
  const editDelegationRoot = document.querySelector('.cat-list-section')
  if (editDelegationRoot) {
    editDelegationRoot.addEventListener('click', (e) => {
      const el = e.target
      if (!(el instanceof Element)) {
        return
      }
      const btn = el.closest('.cat-modal-open-edit')
      if (!btn) {
        return
      }
      openEditModal(btn as HTMLElement)
    })
  }

  /** Backdrop and × — one listener; closest covers inner nodes on the close control. */
  modal.addEventListener('click', (e) => {
    const el = e.target
    if (!(el instanceof Element)) {
      return
    }
    if (el === modal) {
      modal.close()
      return
    }
    if (el.closest('#cat-modal-close') !== null) {
      modal.close()
    }
  })
}
