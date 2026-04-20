import { resolveBootContentQueryRoot } from './contentRoot'
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

/** Forms that already have a `submit` guard (`attachConfirmBeforeSubmit` is safe to call more than once). */
const confirmSubmitWiredForms = new WeakSet<HTMLFormElement>()

/**
 * Destructive or irreversible POST forms: confirm in the capture/bubble submit path.
 * Templates use `data-confirm="message"` instead of inline `onsubmit="return confirm(...)"`.
 * No-ops if `form` is already wired (same `WeakSet` as `initConfirmSubmitForms`).
 */
export function attachConfirmBeforeSubmit(form: HTMLFormElement, message: string): void {
  if (confirmSubmitWiredForms.has(form)) {
    return
  }
  form.addEventListener('submit', (e) => {
    /* global `confirm` — Vitest stubs via `globalThis` / `stubGlobal`; avoid `window` (Node tests). */
    if (!confirm(message)) {
      e.preventDefault()
    }
  })
  confirmSubmitWiredForms.add(form)
}

/** One `querySelectorAll` + trim filter — boot wires without building an intermediate list. */
function forEachConfirmableForm(
  root: ParentNode,
  fn: (form: HTMLFormElement, message: string) => void,
): void {
  for (const form of root.querySelectorAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR)) {
    const msg = readDataConfirmMessage(form)
    if (msg !== null) {
      fn(form, msg)
    }
  }
}

/**
 * Forms that opt in with `form[data-confirm]` and a non-blank message after trim.
 * Uses the same discovery path as `initConfirmSubmitForms` (regression: skip empty messages).
 */
export function findDataConfirmForms(root: ParentNode): Array<{ form: HTMLFormElement; message: string }> {
  const out: Array<{ form: HTMLFormElement; message: string }> = []
  forEachConfirmableForm(root, (form, message) => {
    out.push({ form, message })
  })
  return out
}

/** Wire all matching `form[data-confirm]` under the layout main landmark when present. Idempotent via `attachConfirmBeforeSubmit`. */
export function initConfirmSubmitForms(): void {
  forEachConfirmableForm(resolveBootContentQueryRoot(), (form, message) => {
    attachConfirmBeforeSubmit(form, message)
  })
}
