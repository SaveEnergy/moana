import { APP_MAIN_SELECTOR } from './domSelectors'

/**
 * Minimal `document.querySelector` stub for tests: only **`APP_MAIN_SELECTOR`** resolves to `main`.
 * Matches layout pages where `resolveContentQueryRoot(document)` stops at `<main.app-main>`.
 */
export function stubDocumentMainLandmark(main: ParentNode): Pick<Document, 'querySelector'> {
  return {
    querySelector: (sel: string) => (sel === APP_MAIN_SELECTOR ? main : null),
  }
}
