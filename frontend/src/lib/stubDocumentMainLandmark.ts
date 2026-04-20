/**
 * Vitest-only `document` stubs for **`resolveContentQueryRoot`** / layout behavior (not imported from `main.ts`).
 */
import { APP_MAIN_SELECTOR } from './domSelectors'

type QuerySelectorAllFn = Document['querySelectorAll']

/**
 * Minimal `document.querySelector` stub for tests: only **`APP_MAIN_SELECTOR`** resolves to `main`.
 * Matches layout pages where `resolveContentQueryRoot(document)` stops at `<main.app-main>`.
 */
export function stubDocumentMainLandmark(main: ParentNode): Pick<Document, 'querySelector'> {
  return {
    querySelector: (sel: string) => (sel === APP_MAIN_SELECTOR ? main : null),
  }
}

/**
 * No matching **`main.app-main`** — `resolveContentQueryRoot(document)` uses **`document`** as the scan root.
 * Use with **`querySelectorAll`** when init walks `document` (e.g. login layout, or tests that skip `<main>`).
 */
export function stubDocumentWithoutMainLandmark(
  extras: Partial<{
    querySelector: Document['querySelector']
    querySelectorAll: QuerySelectorAllFn
  }> = {},
): Pick<Document, 'querySelector'> & Partial<Pick<Document, 'querySelectorAll'>> {
  return {
    querySelector: extras.querySelector ?? (() => null),
    ...(extras.querySelectorAll !== undefined ? { querySelectorAll: extras.querySelectorAll } : {}),
  }
}
