/** Mobile drawer shell — matches CSS (`02-shell.css`) and design.md. */
export const MOBILE_SHELL_MAX_WIDTH_PX = 1023 as const

/** Media query for when the sidebar is off-canvas with backdrop. */
export const MOBILE_SHELL_MEDIA_QUERY = `(max-width: ${MOBILE_SHELL_MAX_WIDTH_PX}px)` as const
