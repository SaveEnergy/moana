import { resolveBootContentQueryRoot } from './contentRoot'
import { DATA_CONFIRM_ATTRIBUTE, FORM_DATA_CONFIRM_SELECTOR } from './domSelectors'
import { stripCfTrimEdges } from './unicodeCf'

/**
 * Returns `data-confirm` text or `null` if missing or visually blank after normalizing.
 * Strips Unicode **Cf** characters, then **`trimEdgesIfNeeded`** (**`stripCfTrimEdges`**); skips wiring so templates cannot show an empty `confirm()` dialog.
 */
export function readDataConfirmMessage(form: Element): string | null {
  const raw = form.getAttribute(DATA_CONFIRM_ATTRIBUTE)
  if (raw == null) {
    return null
  }
  /* Exact empty attribute — skip Cf/trim pipeline (common on malformed templates). */
  if (raw === '') {
    return null
  }
  const msg = stripCfTrimEdges(raw)
  return msg === '' ? null : msg
}

/** Forms that already have a `submit` guard (`attachConfirmBeforeSubmit` is safe to call more than once). */
const confirmSubmitWiredForms = new WeakSet<HTMLFormElement>()

/**
 * Destructive or irreversible POST forms: confirm on **`submit`** in **capture** phase
 * (before target/bubble listeners) so **`preventDefault`** always wins if the user declines.
 * Templates use `data-confirm="message"` instead of inline `onsubmit="return confirm(...)"`.
 * No-ops if `form` is already wired (same `WeakSet` as `initConfirmSubmitForms`).
 */
export function attachConfirmBeforeSubmit(form: HTMLFormElement, message: string): void {
  if (confirmSubmitWiredForms.has(form)) {
    return
  }
  form.addEventListener(
    'submit',
    (e) => {
      /* global `confirm` — Vitest stubs via `globalThis` / `stubGlobal`; avoid `window` (Node tests). */
      if (!confirm(message)) {
        e.preventDefault()
      }
    },
    { capture: true },
  )
  confirmSubmitWiredForms.add(form)
}

/** One `querySelectorAll` + Cf strip / trim filter — boot wires without building an intermediate list. */
function forEachConfirmableForm(
  root: ParentNode,
  fn: (form: HTMLFormElement, message: string) => void,
): void {
  const forms = root.querySelectorAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR)
  for (let i = 0, n = forms.length; i < n; i++) {
    const form = forms[i]
    const msg = readDataConfirmMessage(form)
    if (msg !== null) {
      fn(form, msg)
    }
  }
}

/**
 * Forms that opt in with `form[data-confirm]` and a non-blank message after Cf strip + trim.
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
