/**
 * Set a DOM attribute only when **`getAttribute`** differs — fewer redundant mutations
 * (ARIA sync, repeat close paths, duplicate boot).
 */
export function setAttributeIfChanged(el: Element, name: string, value: string): void {
  if (el.getAttribute(name) !== value) {
    el.setAttribute(name, value)
  }
}
