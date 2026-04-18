export function initCategoryModal(): void {
  const dialog = document.getElementById('cat-modal') as HTMLDialogElement | null
  const form = document.getElementById('cat-modal-form') as HTMLFormElement | null
  const idInput = document.getElementById('cat-modal-id') as HTMLInputElement | null
  const titleEl = document.getElementById('cat-modal-title')
  const submitBtn = document.getElementById('cat-modal-submit')
  const preview = document.getElementById('cat-modal-preview')
  const iconWrap = document.getElementById('cat-modal-preview-icon')
  const nameInput = document.getElementById('cat-modal-name') as HTMLInputElement | null
  const closeBtn = document.getElementById('cat-modal-close')
  const addCategoryBtn = document.getElementById('cat-modal-open-create')

  if (!dialog || !form || !idInput || !titleEl || !submitBtn || !preview || !iconWrap || !nameInput) {
    return
  }

  function syncCatModalPreview() {
    const cr = form.querySelector<HTMLInputElement>('input[name="color"]:checked')
    let bg = 'color-mix(in srgb, var(--primary) 12%, #fff8f0)'
    if (cr?.value === 'custom') {
      const nat = form.querySelector<HTMLInputElement>('#cat-modal-color-native')
      bg = nat?.value?.trim() || '#818cf8'
    } else if (cr?.value) {
      bg = cr.value
    }
    preview.style.background = bg

    const ir = form.querySelector<HTMLInputElement>('input[name="icon"]:checked')
    iconWrap.innerHTML = ''
    if (!ir?.value) {
      iconWrap.classList.add('cat-modal-preview-icon--auto')
      iconWrap.textContent = 'A'
      return
    }
    iconWrap.classList.remove('cat-modal-preview-icon--auto')
    const label = ir.closest('label')
    const svg = label?.querySelector('svg.moana-icon')
    if (svg) {
      const clone = svg.cloneNode(true) as SVGElement
      clone.classList.add('moana-icon--cat-preview')
      iconWrap.appendChild(clone)
    }
  }

  function wireColorNative() {
    form.querySelectorAll('.cat-color-native').forEach((pc) => {
      pc.addEventListener('input', () => {
        const wrap = (pc as HTMLElement).closest('.cat-color-swatch--custom')
        const r = wrap?.querySelector<HTMLInputElement>('input[type="radio"][value="custom"]')
        if (r) {
          r.checked = true
          syncCatModalPreview()
        }
      })
    })
  }

  form.querySelectorAll('input[name="color"], input[name="icon"]').forEach((el) => {
    el.addEventListener('change', () => syncCatModalPreview())
  })
  wireColorNative()

  function openCreateModal() {
    form.action = '/categories'
    idInput.value = ''
    titleEl.textContent = 'New category'
    submitBtn.textContent = 'Create category'
    form.reset()
    form.querySelector<HTMLInputElement>('input[name="color"][value=""]')!.checked = true
    form.querySelector<HTMLInputElement>('input[name="icon"][value=""]')!.checked = true
    const nat = form.querySelector<HTMLInputElement>('#cat-modal-color-native')
    if (nat) nat.value = '#818cf8'
    syncCatModalPreview()
    nameInput.focus()
    dialog.showModal()
  }

  function openEditModal(btn: HTMLElement) {
    const id = btn.dataset.id
    if (!id) return
    form.action = '/categories/update'
    idInput.value = id
    titleEl.textContent = 'Edit category'
    submitBtn.textContent = 'Save changes'

    nameInput.value = btn.dataset.name ?? ''

    const rawColor = (btn.dataset.color ?? '').trim()
    const isCustom = btn.dataset.custom === '1'
    const customHex = (btn.dataset.customHex ?? '#818cf8').trim()

    if (!rawColor) {
      form.querySelector<HTMLInputElement>('input[name="color"][value=""]')!.checked = true
    } else if (isCustom) {
      form.querySelector<HTMLInputElement>('input[name="color"][value="custom"]')!.checked = true
      const nat = form.querySelector<HTMLInputElement>('#cat-modal-color-native')
      if (nat) nat.value = /^#[0-9a-fA-F]{6}$/.test(customHex) ? customHex : '#818cf8'
    } else {
      const preset = Array.from(form.querySelectorAll<HTMLInputElement>('input[name="color"]')).find(
        (r) => r.value === rawColor,
      )
      if (preset) {
        preset.checked = true
      } else {
        form.querySelector<HTMLInputElement>('input[name="color"][value=""]')!.checked = true
      }
    }

    const iconVal = (btn.dataset.icon ?? '').trim()
    if (!iconVal) {
      form.querySelector<HTMLInputElement>('input[name="icon"][value=""]')!.checked = true
    } else {
      const iconRadio = Array.from(form.querySelectorAll<HTMLInputElement>('input[name="icon"]')).find(
        (r) => r.value === iconVal,
      )
      if (iconRadio) {
        iconRadio.checked = true
      } else {
        form.querySelector<HTMLInputElement>('input[name="icon"][value=""]')!.checked = true
      }
    }

    syncCatModalPreview()
    nameInput.focus()
    dialog.showModal()
  }

  addCategoryBtn?.addEventListener('click', () => openCreateModal())
  document.querySelectorAll('.cat-modal-open-edit').forEach((b) => {
    b.addEventListener('click', () => openEditModal(b as HTMLElement))
  })

  closeBtn?.addEventListener('click', () => dialog.close())
  dialog.addEventListener('click', (e) => {
    if (e.target === dialog) {
      dialog.close()
    }
  })
}
