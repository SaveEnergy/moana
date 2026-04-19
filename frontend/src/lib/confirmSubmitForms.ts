import { DATA_CONFIRM_ATTRIBUTE, FORM_DATA_CONFIRM_SELECTOR } from './domSelectors'

/**
 * Returns trimmed `data-confirm` or `null` if missing or whitespace-only.
 * Skips wiring so templates cannot accidentally show an empty `confirm()` dialog.
 */
export function readDataConfirmMessage(form: Element): string | null {
  const raw = form.getAttribute(DATA_CONFIRM_ATTRIBUTE)
  if (raw == null) {
    return null
  }
  const msg = raw.trim()
  return msg === '' ? null : msg
}

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

/**
 * Forms that opt in with `form[data-confirm]` and a non-blank message after trim.
 * Single pass; used by `initConfirmSubmitForms` and unit tests (regression: skip empty messages).
 */
export function findDataConfirmForms(root: ParentNode): Array<{ form: HTMLFormElement; message: string }> {
  const out: Array<{ form: HTMLFormElement; message: string }> = []
  for (const form of root.querySelectorAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR)) {
    const msg = readDataConfirmMessage(form)
    if (msg !== null) {
      out.push({ form, message: msg })
    }
  }
  return out
}

/** Forms that already have a capture/bubble `submit` guard (safe if `bootApp` runs twice). */
const confirmSubmitWiredForms = new WeakSet<HTMLFormElement>()

/** Wire all matching `form[data-confirm]` in `document`. Skips forms already wired in this session. */
export function initConfirmSubmitForms(): void {
  for (const { form, message } of findDataConfirmForms(document)) {
    if (confirmSubmitWiredForms.has(form)) {
      continue
    }
    confirmSubmitWiredForms.add(form)
    attachConfirmBeforeSubmit(form, message)
  }
}
