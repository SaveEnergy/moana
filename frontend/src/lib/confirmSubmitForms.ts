/**
 * Destructive or irreversible POST forms: confirm in the capture/bubble submit path.
 * Templates use `data-confirm="message"` instead of inline `onsubmit="return confirm(...)"`.
 */
export function attachConfirmBeforeSubmit(form: HTMLFormElement, message: string): void {
  form.addEventListener('submit', (e) => {
    /* global `confirm` — Vitest stubs via `globalThis` / `stubGlobal`; avoid `window` (Node tests). */
    if (!confirm(message)) {
      e.preventDefault()
    }
  })
}

/** Wire all `form[data-confirm]` in `document`. Idempotent only on first boot (one listener per form). */
export function initConfirmSubmitForms(): void {
  for (const form of document.querySelectorAll<HTMLFormElement>('form[data-confirm]')) {
    const msg = form.getAttribute('data-confirm')?.trim()
    if (!msg) {
      continue
    }
    attachConfirmBeforeSubmit(form, msg)
  }
}
