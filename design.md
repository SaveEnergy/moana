# Moana UI — design system

Documentation grounded in `frontend/src/app.css`, `frontend/src/styles/*.css`, `frontend/src/lib/*.ts`, `frontend/src/boot.ts`, `internal/assets/templates/`, and `frontend/src/main.ts`.

## 1. Intent and voice

Moana’s UI is a **server-rendered, CSS-first “editorial shell”**: warm **canvas** (`#f0ebe3`) behind cool-neutral **paper** cards, with **Fraunces** for display wordmark and section titles, **Inter** for body, and **Manrope** for numeric and display emphasis. The product reads as a **personal finance dashboard** (period filters, KPIs, heatmaps, transaction history), not a generic admin UI.

## 2. Architecture

| Layer | Role |
|--------|------|
| **HTML** | Go `html/template` files under `internal/assets/templates/` (e.g. `layout.html`). |
| **CSS** | **Entry:** `frontend/src/app.css` `@import`s ordered partials under `frontend/src/styles/` (`01-tokens.css` … `07-reduced-motion.css`). **Build:** Vite concatenates/minifies to **one** `internal/assets/static/css/app.css` (single network request, unchanged URL). |
| **JS** | `frontend/src/main.ts` loads CSS then calls **`bootApp()`** from `frontend/src/boot.ts`, which iterates **`BOOT_APP_INITIALIZERS`** (**timezone cookie → local times → shell → settings dialog → category modal → history sort auto-submit → `data-confirm` submit guards**; order is regression-tested). Each module **no-ops** when its DOM is absent (e.g. no `#app-shell` on login-only pages, no category modal on dashboard-only loads). **Mobile shell width** (`1023px`) is defined once as **`MOBILE_SHELL_MAX_WIDTH_PX`** / **`MOBILE_SHELL_MEDIA_QUERY`** in `frontend/src/lib/shellBreakpoints.ts`; `onMediaQueryChange` wraps `addEventListener` with a legacy `addListener` fallback for `shellSidebar.ts`. **Mobile drawer dismiss:** **`shouldCloseMobileSidebarFromShellClick`** (`mobileShellDismiss.ts`), wired from **`shellSidebar.ts`** — **`APP_SHELL_SELECTOR`**, **`APP_SIDEBAR_TOGGLE_SELECTOR`**, **`APP_SIDEBAR_BACKDROP_SELECTOR`**, **`APP_SIDEBAR_CLOSE_SELECTOR`** (**`domSelectors.ts`**); **`queryAppShell`** on **`document`**, then **`querySidebarToggle`** / **`querySidebarBackdrop`** on **`#app-shell`** (toggle/backdrop live under the shell in **`layout.html`**); one bubbling **`click`** on **`#app-shell`** when the resolved target is the backdrop or **`closest()`** resolves the drawer close control (**`clickEventTargetElement`**). **Escape vs layered UI:** **`shouldDeferMobileShellEscape`** (`dialogKeyboard.ts`) scans **`composedPath()`** for open **`dialog`** or **`details`**, then **`isAppUserMenuDetailsOpen`** (**`APP_USER_MENU_OPEN_SELECTOR`**, **`domSelectors.ts`**) when the menu can stay **`[open]`** without the focused node in the path (the drawer can cover the bar while `[open]` stays true). **Category modal** — **`readCategoryEditRowDataset`** (**`categoryModalDataset.ts`**) mirrors **`categories.html`** list-row **Edit** `data-*` for **`openEditModal`**; default preview strip uses shared `CATEGORY_MODAL_DEFAULT_PREVIEW_BG` in `categoryColor.ts`; color/icon preset radios are indexed **once** in a `Map` by `value` for fast edit opens (**`setRadioCheckedByValue`** in **`radioMap.ts`** applies dataset-driven values with an **auto** fallback instead of chained **`get()`** + assertions); **form-level delegation** handles custom color `input` and color/icon **`change`** on `#cat-modal-form` instead of per-control listeners; **Edit** actions use **event delegation** on **`.cat-list-section`** with **`clickEventTargetElement`** then **`closest('.cat-modal-open-edit')`** (same shared target helper as dialog dismiss); **native `<dialog>` dismiss** for settings + categories uses **`attachNativeDialogDismiss`** in **`dialogDismiss.ts`** (backdrop + **`closest()`** on close/cancel controls; **`clickEventTargetElement`** handles **Text** nodes inside labels/×). **Local times:** `formatLocalTimeLabel` uses one shared `Intl.DateTimeFormat`; **`createLocalTimeLabelMemo`** keeps one pass over elements matching **`LOCAL_TIME_ELEMENTS_SELECTOR`** (**`TIME_DATETIME_ATTRIBUTE`**, **`domSelectors.ts`**), memoizing successful labels in a `Map` and recording invalid ISO strings in a `Set` so bad attributes are not re-parsed on every duplicate row. **`applyLocalTimeElements`** uses one module-level memo; the default **`document`** boot path calls **`resolveBootContentQueryRoot`** (**`contentRoot.ts`**), explicit subtrees use **`resolveContentQueryRoot`** — same **`WeakMap`** per **`document`** so **`bootApp`** shares one **`main.app-main`** resolution when present (**`APP_MAIN_SELECTOR`**, **`layout.html`**), else **`document`** (e.g. **`login.html`**). **`resolveBootContentQueryRoot`** also scopes **`initConfirmSubmitForms`** (**`form[data-confirm]`**), **`initCategoryModal`**, **`initHistoryControls`**, and **`initSettingsMemberDialog`**. |
| **Icons** | Embedded Lucide-derived SVG via Go (`internal/icons/`), classes `.moana-icon` and size modifiers. |

Templates load `/static/css/app.css` and Google Fonts in `<head>`; **`rel="modulepreload"`** on `/static/js/app.js` starts fetch early while HTML parses. **`app.js` executes at the end of `<body>`** (module deferred, runs after parse). The login template also includes `app.js` so the **timezone cookie** is set before the first authenticated request; subsequent loads **skip rewriting** `document.cookie` when `moana_tz` already matches the browser zone (`parseMoanaTimezoneCookie` in `timezoneCookie.ts`). **`parseMoanaTimezoneCookie`** walks semicolon-separated segments **without** `String.prototype.split`, so large `document.cookie` values avoid allocating a full segment array.

There is **no React component library** in-repo; “components” are **CSS class contracts** in HTML. **Structure:** ordered style partials preserve cascade; behavior is split into small TypeScript modules for testing and navigation.

## 3. Design tokens (`:root`)

Tokens follow **Material-style naming** (surface / on-surface / primary / container) plus app-specific **dashboard** aliases.

### Core palette

- **Primary (brand / actions):** `--primary: #306369`, `--primary-container: #b2e6ec`, `--on-primary: #ffffff`
- **Secondary (muted chrome):** `--secondary`, `--secondary-container`, `--on-secondary-container`
- **Tertiary (accent — “expense” affordances):** `--tertiary: #a03a0f`, `--tertiary-container: #ff946e`
- **Neutrals:** `--on-surface`, `--on-surface-variant`, `--outline-variant`, surface stops (`--surface` … `--surface-lowest`)
- **Muted copy:** `--text-muted` — alias of `--on-surface-variant` (alerts, secondary lines)
- **Error:** `--error: #b31b25`

### App frame

- `--dashboard-canvas` / `--dashboard-paper` / `--dashboard-edge` — page background, card fill, hairline borders
- `--app-canvas` / `--app-paper` — aliases; shell rails use canvas; cards and primary content use paper

### Form surfaces (on warm canvas)

- `--input-surface`, `--input-surface-strong`, `--input-surface-muted`, `--input-surface-readonly` — `color-mix()` blends canvas and paper so controls avoid flat cool greys.

### Layout and chrome

- `--app-sidebar-width: 13.75rem`, sidebar text / hover / active tokens, `--shell-topbar`, `--shell-search-icon`

### Radii and elevation

- `--radius-lg: 0.5rem` — default control radius  
- Many cards use literal `border-radius: 1rem` (not tokenized)
- `--shadow-float`, `--shadow-card`

### Typography tokens

- `--font-body`: Inter  
- `--font-serif`: Fraunces  
- `--font-display`: Manrope  

### Fluid color

Heavy use of **`color-mix(in srgb, …)`** for hovers, selection, and tinted backgrounds.

## 4. Typography

| Token | Usage |
|--------|--------|
| `--font-body` | Default UI text, inputs, buttons, tables |
| `--font-serif` | Brand wordmark (`.brand-wordmark`), page and section titles (dashboard hero, card titles, entry titles, history title, settings cards, modal titles) |
| `--font-display` | Numbers and KPIs: balances, amounts, emphasis in charts-related UI |

**Base scale:** `html.app-html` sets `font-size: 16px`; `body` uses **`font-size: 0.875rem` (14px)** and `line-height: 1.5` — dense, product-style rhythm.

**Utility classes:** `.headline-sm`, `.display-lg`, `.display-sm`, `.label`, `.muted`, `.small`, `.text-error`.

## 5. Color semantics (product meaning)

| Meaning | Implementation |
|----------|----------------|
| **Primary / brand** | `--primary`, navigation active states, income-flavored accents |
| **Expense accent** | `--tertiary` — expense icons, expense category chrome, expense side of kind toggles |
| **Positive money** | **Hardcoded green** `#1b5e20` in several rules (KPI trends, income amounts, pills) — not a CSS variable |
| **Negative / danger** | `--error` and additional reds (e.g. `#b91c1c`, `#b85454`) in trend text |
| **Category accents** | `--cat-accent` set per row/card; fallback to `--primary` |

**Maintainability note:** Consider tokenizing success/negative money colors (e.g. `--semantic-positive`, `--semantic-negative`) to replace scattered hex values.

## 6. Layout — app shell

**Structure:** `.app-shell` → `.app-sidebar` + `.app-content`.

- **Sidebar:** Sticky, full height, `z-index: 50`, wordmark, `.sidebar-nav` links (current route: **`.sidebar-link-active`** from `layout.html` **`Active`**), `.sidebar-fab` primary CTA (“Add transaction”).
- **Top bar:** `.app-topbar` — mobile menu control, global search (`.app-search`, GET `/history`), notifications link (`.app-topbar-icon-btn`; **`aria-current="page"`** when **Active** is notifications), `<details class="app-user-menu">` for account (panel **Settings** link **`aria-current="page"`** when **Active** is settings).  
- **Mobile drawer (≤1023px):** Open via `#app-sidebar-toggle`; close via toggle, **Escape**, or **delegated `click` on `#app-shell`** (backdrop `#app-sidebar-backdrop` or `#app-sidebar-close`; predicate **`mobileShellDismiss.ts`**).
- **Main:** `.app-main`; optional **`.dashboard-page`** narrows `.app-content-container` to `max-width: 72rem` versus default `80rem`.

**Content column:** `.app-content-container` uses responsive horizontal padding (notable breakpoints at **640px** and **1024px**), max width **80rem**, vertical padding for page rhythm.

**Footer:** `.site-footer` with wordmark and legal links; `.app-content-container--footer` adjusts spacing.

## 7. Responsive behavior

Breakpoints are **declared per component** (not one shared map). Representative values:

- **520px** — form columns, hero row stacking  
- **540px** — entry field rows  
- **640px** — two-column grids, history controls, login padding  
- **720px** — category icon grid density  
- **768px** — footer switches to row layout  
- **800px** — dashboard outflow split  
- **880px** — transaction entry two-column layout  
- **900px** — main vertical padding  
- **1023px** — **mobile shell:** sidebar off-canvas, backdrop, toggle visible; Escape closes drawer (`shellSidebar.ts`)

Pattern: **desktop-first shell** with a **1023px** cutoff for drawer navigation.

## 8. Component patterns

### Buttons (`.btn`)

- `.btn-primary` / `.btn-indigo` — both use primary fill (`btn-indigo` is legacy naming; color is teal, not indigo)
- `.btn-secondary` — neutral filled  
- `.btn-tertiary` — text style, primary color  
- `.btn-ghost` — muted text + hover wash  
- `.btn-small`, `.btn-block`

**Interaction:** `:active { transform: translateY(1px) }`; hover via **opacity** or **`filter: brightness(1.06)`** on solid fills.

### Forms

- **Underline style:** `.input` — focus draws a bottom edge in `--primary`
- **Float pattern:** `.float-field`, `.float-label`, `.float-input` inside bordered surfaces; focus uses border + outer shadow
- **Amount entry:** `.amount-input-wrap`, `.input-amount` — large Manrope numerics; wrap highlights with `--primary-container` on focus  
- **Confirm-before POST:** `form[data-confirm="…"]` (**`FORM_DATA_CONFIRM_SELECTOR`**, **`DATA_CONFIRM_ATTRIBUTE`**, **`domSelectors.ts`**) — `window.confirm` on **`submit`** (`confirmSubmitForms.ts`; **`initConfirmSubmitForms`** and **`findDataConfirmForms`** share one `form[data-confirm]` walk under **`resolveBootContentQueryRoot`** (**`contentRoot.ts`**, **`main.app-main`** when present); **`readDataConfirmMessage`** trims the attribute and skips blank / whitespace-only values, then **`attachConfirmBeforeSubmit`** wires **`submit`** — idempotent per form via **`WeakSet`**); avoid inline handlers  
- **Bootstrap idempotency:** **`WeakSet`** guards in **`shellSidebar`**, **`confirmSubmitForms`**, **`historyControls`**, **`dialogDismiss`** (`attachNativeDialogDismiss`), **`settingsMemberDialog`**, and **`categoryModal`** record the element **after** listener registration (where applicable), so a failed **`addEventListener`** does not permanently skip wiring on a later **`bootApp()`**; duplicate **`bootApp()`** still does not stack listeners (tests or future hot reload)

### Segmented controls

- **Dashboard period:** underline tabs — `.dashboard-period-seg`, `.dashboard-period-opt.is-active`
- **History:** pill control — `.history-segmented`, `.history-seg`, `.history-seg-active`
- **Transaction kind:** pill radios — `.kind-toggle`, `.kind-option`

### Cards

- **Dashboard:** `.dashboard-card`, `.dashboard-kpi-row`, heatmap-specific wrappers  
- **Generic:** `.entry-card`, `.layer-card`, `.settings-card`  
- **Elevation:** `var(--shadow-card)` or float shadow

### Alerts

`.alert` with `.alert-error`, `.alert-success`, `.alert-info`

### Dialogs

Native `<dialog>`: `.admin-add-dialog` (settings add member), `.cat-modal` (categories) with `.cat-modal-panel`, header, footer, full-width primary actions.

### Chips / status

`.pill`, `.pill-trend`, `.pill-trend-neg`, `.pill-live`

### Lists

`.dashboard-outflow-*`, `.history-card-list`, `.tx-list`, `.dashboard-recent-*` — recurring pattern: **icon tile + body + right-aligned amount**

## 9. Icons

- **Base:** `.moana-icon`; modifiers: `--sm`, `--grid`, `--nav`, `--cat-preview`
- **Rendering:** Inline SVG from template helpers; see `internal/icons/data_gen.go`
- **Color:** `currentColor` — parent sets semantic color

## 10. Motion and depth

- Transitions roughly **0.12–0.28s** on hovers; sidebar uses **`cubic-bezier(0.4, 0, 0.2, 1)`**
- Heatmap budget bar width transitions  
- Category modal uses **backdrop blur**
- **`prefers-reduced-motion: reduce`:** `frontend/src/styles/07-reduced-motion.css` disables the heaviest UI transitions (sidebar slide, chevron spin).

## 11. Accessibility

- **`.sr-only`** for visually hidden labels  
- **Landmarks:** `<main>`, `<nav>` with `aria-label`, `aria-current` where used  
- **`<details>`/`<summary>`** for user menu  
- **`:focus-visible`** on some controls (e.g. sidebar brand)  
- **Mobile shell Escape:** `shellSidebar.ts` uses a **capture-phase** document listener and defers when **`shouldDeferMobileShellEscape`** (`dialogKeyboard.ts`) is true. It scans **`composedPath()`** once for an open **`dialog`** or **`details`**, then **`isAppUserMenuDetailsOpen`** (**`querySelector`**) when the account menu can stay **`[open]`** while focus is outside the disclosure (e.g. mobile drawer over the top bar). **`dialogKeyboard`** uppercases **`tagName`** when scanning paths (consistent with HTML’s uppercase `tagName`, tolerant in tests).

Custom widgets (segmented groups, FAB-as-link) may need extra ARIA depending on future audit scope.

## 12. Domain-specific surfaces

- **Dashboard:** Hero balance, period filter (**`.dashboard-period-opt.is-active`** on the link matching **`period=`** — `dashboard.html`), KPI grid (1 / 2 / 4 columns by breakpoint), donut outflow, heatmap levels `--lv0`–`lv4`, recent transactions  
- **Transactions / entry:** New transaction **POST** can re-render the form with **`role="alert"`** for validation errors (e.g. zero amount — see `internal/txform/parse.go`); category picker cards with **`--cat-accent`** left border and icon well  
- **History:** Day groupings, `.history-card` rows, search/sort/date range (**from** / **to** + **Apply dates** submits **`GET /history`** with `from=` & `to=`); kind filter tabs — **`.history-seg-active`** on the tab matching **`kind=`** (`all` / `expense` / `income`); **`#history-sort`** (**`HISTORY_SORT_SELECTOR`** / **`HISTORY_SORT_ELEMENT_ID`** in **`domSelectors.ts`**; **`initHistoryControls`** uses **`queryHistorySortSelect(resolveBootContentQueryRoot())`** (**`contentRoot.ts`**); skips **`wireHistorySortAutoSubmit`** when that **`<select>`** is already wired; **`wireHistorySortAutoSubmit`** in **`historyControls.ts`**) **`change`** → **`form.requestSubmit()`** (no inline handler)  
- **Categories:** Modal create/edit (**`categoryModal.ts`** — **`resolveBootContentQueryRoot().querySelector`** (**`CATEGORY_MODAL_SELECTOR`**), duplicate **`init`** returns before further queries, then **`dialog.querySelector`** for in-modal **`CATEGORY_MODAL_*`** hooks (**Add category** via **`CATEGORY_PAGE_INTRO_SECTION_SELECTOR`** + **`CATEGORY_MODAL_OPEN_CREATE_SELECTOR`**, else **`contentRoot`**; list **`Edit`** — **`intro.parentElement?.querySelector(CATEGORY_LIST_SECTION_SELECTOR)`** else **`contentRoot.querySelector(CATEGORY_LIST_SECTION_SELECTOR)`** so delegation stays on the categories page subtree); **`readCategoryEditRowDataset`** in **`categoryModalDataset.ts`** for row **Edit** `data-*`; **`CATEGORY_MODAL_*`**, **`CATEGORY_LIST_SECTION_SELECTOR`**, **`MOANA_ICON_SVG_SELECTOR`**, **`CATEGORY_COLOR_*`**, **`CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT`**, **`resolveCategoryModalPreviewBackground`** (trims checked radio `value`), **`sanitizeCategoryCustomHex`** in **`categoryColor.ts`**, **`CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME`** / **`CATEGORY_MODAL_ICON_RADIO_GROUP_NAME`**, **`CATEGORY_MODAL_COLOR_RADIOS_SELECTOR`** / **`CATEGORY_MODAL_ICON_RADIOS_SELECTOR`**, **`buildRadioMapByValue`** + **`getFormRadioGroupValue`** for live preview (icon clone falls back to **`CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR`**), **`CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR`** / **`CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR`** in **`domSelectors.ts`** for tests, **`CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM`** / **`CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR`** for the custom swatch radio, one resolved **`#cat-modal-color-native`** reference for preview + open flows, preview classes, **`domSelectors.ts`**); color presets from `internal/category/colors.go`; icon grid; preview; **dismiss** — **`attachNativeDialogDismiss`** on **`#cat-modal`** (backdrop + **`#cat-modal-close`**)  
- **Settings:** Float-field forms, readonly **Email** field `#settings-email` bound to the signed-in user, add-member (**`domSelectors.ts`**: **`SETTINGS_ADD_MEMBER_*`**, **`SETTINGS_ADD_MEMBER_DIALOG_SELECTOR`**, **`SETTINGS_ADD_MEMBER_OPEN_SELECTOR`**; **`querySettingsAddMemberDialog(resolveBootContentQueryRoot())`**, **`querySettingsAddMemberOpenButton`** on **`dialog.parentElement`** when attached (else **`document`**) in **`settingsMemberDialog.ts`**; **`attachNativeDialogDismiss`**: backdrop, close, Cancel), member list with role badges  
- **Notifications:** `notifications.html` — intro line + empty inbox **`role="status"`**; in **`layout.html`** the topbar bell link adds **`aria-current="page"`** when that view is active  
- **Login:** `.login-split` — stacked on small screens; split layout with hero image from **1024px**; **Remember me** checkbox; **OAuth** row uses **`role="group"`** with **`aria-label="OAuth sign-in (not available in this build)"`** and **disabled** stub buttons (Google, GitHub) matching `login.html`

## 13. Build and source of truth

- **Edit CSS in:** `frontend/src/styles/*.css` (preserve **import order** in `frontend/src/app.css`)  
- **Build:** `bun run build` — Vite outputs to `internal/assets/static/css/app.css` and `js/app.js` (see `frontend/vite.config.ts`; `es2022` target, CSS minify, **esbuild** strips **`legalComments`** from JS, and in **production** drops **`console`** / **`debugger`** for a smaller bundle, compressed-size reporting off for faster builds).  
- **Do not hand-edit** generated files under `internal/assets/static/`.

## 14. Tests

- **Typecheck:** `bun run typecheck` — `tsc --project frontend/tsconfig.json`.  
- **Unit:** `bun run test:unit` — Vitest, config **merged into `frontend/vite.config.ts`** (`test` block). **Production `esbuild`:** `console` / `debugger` **drop** only in prod (`viteProductionEsbuild.test.ts`). Covers `frontend/src/lib/*`, **`bootApp()`** / **`BOOT_APP_INITIALIZERS`** wiring and initializer order (`boot.test.ts` (**`BOOT_INITIALIZER_COUNT`**), `bootInitializers.test.ts`, **`bootInitializerNames.ts`** / **`bootInitializerNames.test.ts`** — **`BootInitializer`**, **`BOOT_INITIALIZER_NAMES`** / **`BOOT_APP_INITIALIZERS`** `function.name` order, duplicate-name guard), timezone cookie segment + **`parseMoanaTimezoneCookie`**, **`setBrowserTimezoneCookie`**, **`shouldDeferMobileShellEscape`** and related helpers (`dialogKeyboard.test.ts`; **`APP_USER_MENU_OPEN_SELECTOR`**, **`dialogKeyboard.ts`** composed-path **`DIALOG`/`DETAILS`** checks), `formatLocalTimeLabel`, **`createLocalTimeLabelMemo`**, `applyLocalTimeElements` (`localTime.test.ts`; **`LOCAL_TIME_ELEMENTS_SELECTOR`**, **`APP_MAIN_SELECTOR`**, default-root + explicit-root when **`document`** is absent (Vitest Node), **`TIME_DATETIME_ATTRIBUTE`**, invalid **`datetime`** values, **`domSelectors.ts`**), **`resolveContentQueryRoot`**, **`resolveBootContentQueryRoot`** (`contentRoot.test.ts`; **`APP_MAIN_SELECTOR`**, per-`document` memo), **`stubDocumentMainLandmark`** / **`stubDocumentWithoutMainLandmark`** (`stubDocumentMainLandmark.test.ts`; Vitest **`document`** stubs for layout **`<main>`** vs login-style fallback), category modal preview background + custom hex (`categoryColor.test.ts`); **`readCategoryEditRowDataset`** (`categoryModalDataset.test.ts`; blank **`customHex`**); category modal radios, **data-confirm**, local time + history sort + settings add-member + mobile shell + category modal `#` selectors (`domSelectors.test.ts`; **`APP_SHELL_SELECTOR`**, **`APP_SIDEBAR_TOGGLE_SELECTOR`**, **`APP_SIDEBAR_BACKDROP_SELECTOR`**, **`APP_MAIN_SELECTOR`**, **`HISTORY_SORT_SELECTOR`**, **`SETTINGS_ADD_MEMBER_DIALOG_SELECTOR`**, **`SETTINGS_ADD_MEMBER_OPEN_SELECTOR`**, **`CATEGORY_MODAL_SELECTOR`**, **`CATEGORY_MODAL_FORM_SELECTOR`**, **`CATEGORY_MODAL_ID_INPUT_SELECTOR`**, **`CATEGORY_MODAL_TITLE_SELECTOR`**, **`CATEGORY_MODAL_SUBMIT_SELECTOR`**, **`CATEGORY_MODAL_PREVIEW_SELECTOR`**, **`CATEGORY_MODAL_PREVIEW_ICON_SELECTOR`**, **`CATEGORY_MODAL_NAME_SELECTOR`**, **`CATEGORY_MODAL_OPEN_CREATE_SELECTOR`**); **`CATEGORY_MODAL_*`** (**`categoryModal.ts`**, **`domSelectors.ts`**), **`buildRadioMapByValue`**, **`setRadioCheckedByValue`**, **`getFormRadioGroupValue`** (`radioMap.test.ts`; **`RadioNodeList`** + single **`namedItem`** control), **`viteProductionEsbuild.test.ts`** ( **`vite.config.ts`** production **`esbuild.drop`**), **`MOBILE_SHELL_MAX_WIDTH_PX`** / **`MOBILE_SHELL_MEDIA_QUERY`**, **`onMediaQueryChange`** (`shellBreakpoints.test.ts`), **`queryAppShell`**, **`querySidebarToggle`**, **`querySidebarBackdrop`**, **`initShellSidebar`** (`shellSidebar.test.ts`; **`APP_SHELL_SELECTOR`**, **`APP_SIDEBAR_TOGGLE_SELECTOR`**, **`APP_SIDEBAR_BACKDROP_SELECTOR`**, **`MOBILE_SHELL_MEDIA_QUERY`**, **`domSelectors.ts`), **`clickEventTargetElement`** (`clickTarget.test.ts`), **`shouldCloseNativeDialogFromClick`**, **`attachNativeDialogDismiss`** (`dialogDismiss.test.ts`; backdrop-only / empty dismiss selectors), **`shouldCloseMobileSidebarFromShellClick`** (`mobileShellDismiss.test.ts`; **`APP_SIDEBAR_CLOSE_SELECTOR`** + shell ids, **`domSelectors.ts`**), **`queryHistorySortSelect`**, **`wireHistorySortAutoSubmit`**, **`initHistoryControls`** (`historyControls.test.ts`; **`HISTORY_SORT_SELECTOR`**, **`HISTORY_SORT_ELEMENT_ID`**, **`APP_MAIN_SELECTOR`**, **`domSelectors.ts`**), **`querySettingsAddMemberDialog`**, **`querySettingsAddMemberOpenButton`**, **`initSettingsMemberDialog`** (`settingsMemberDialog.test.ts`; **`SETTINGS_ADD_MEMBER_DIALOG_SELECTOR`**, **`SETTINGS_ADD_MEMBER_OPEN_SELECTOR`**, **`APP_MAIN_SELECTOR`**, **`SETTINGS_ADD_MEMBER_*`**), **`findDataConfirmForms`**, **`readDataConfirmMessage`**, **`attachConfirmBeforeSubmit`**, **`initConfirmSubmitForms`** (`confirmSubmitForms.test.ts`; **`FORM_DATA_CONFIRM_SELECTOR`**, **`DATA_CONFIRM_ATTRIBUTE`**, **`APP_MAIN_SELECTOR`**, **`domSelectors.ts`**); **`initCategoryModal`** (`categoryModal.test.ts`; **`APP_MAIN_SELECTOR`** / main landmark, content-root wiring, duplicate-boot guard).  
- **E2E:** `bun run test:e2e` — Playwright (`test/e2e/tests/`). Shared login: `test/e2e/helpers/auth.ts` (`signInAsTestUser`) matches the seeded user from `test/e2e/scripts/start-server.sh`. `test/e2e/helpers/dates.ts` (`todayLocalISODate`) supplies local calendar **date** values for transaction form tests; **`test/e2e/helpers/modulePreload.ts`** exports **`STATIC_APP_CSS_PATH`**, **`STATIC_APP_JS_PATH`**, **`APP_CSS_STYLESHEET`**, and **`APP_JS_MODULE_PRELOAD`** (Playwright selectors for the bundled **`app.css`** / **`app.js`** head links; paths match **`layout.html`** / **`login.html`**). Covers **`/health`**, **`/login`** form (**Email address**, **Password**, **Sign in**, **Remember me**, **OAuth** stub group with disabled Google / GitHub), sign-in (`signInAsTestUser`), dashboard shell (**`--primary`** token, **Statistics and outflow period** group, **Total Income** KPI, **Money movement** section, **period** links **30 days / 12 months** → **`?period=30d|12m`** (**`.is-active`** on the matching **`.dashboard-period-opt`** for **`?period=`**), **sidebar `#app-sidebar-nav`** (**`.sidebar-link-active`** on Transactions / History / Categories / Dashboard routes), **`.sidebar-fab`** (“Add transaction” → `/transactions`), **global search** `role="search"`, **topbar search** → `/history` with `q`, mobile **Escape** does not collapse drawer while **`details.app-user-menu`** is open), **user menu → Settings** link (**`aria-current="page"`** on `/settings`), notifications page (empty state copy, topbar **`aria-current="page"`** on the bell link), mobile sidebar (toggle, **Escape**, **backdrop** dismiss, **#app-sidebar-close**), **footer** `navigation` “Legal and source” (MIT + GitHub), history (**GET** search `q=`; **sort** `<select>` → **`sort=`**; **Apply dates** → **`from=`** / **`to=`**; **Clear filters** → **`/history`** with empty query; **kind** tablist **All / Income / Expenses** — **`.history-seg-active`** matches **`kind=`**; preserves other query params and sets **`kind=`**), route smoke (categories, new transaction + **kind** Income/Expense via **`.kind-option`** labels — inputs are visually hidden; **new transaction** zero amount shows **`role="alert"`** (`Amount must be greater than zero.`); **edit entry** from **history** → **`/transactions/{id}/edit`** (**Edit entry** + **Save changes** → **`/history`**); settings — **`#settings-email`** matches seed user), **categories modal** (open/close, **Escape**, **backdrop** click outside **`.cat-modal-panel`**, mobile **Escape** with drawer + `dialog` open — sidebar stays open; **showModal** in test harness when backdrop blocks UI; **delete** runs **`window.confirm`** via **`data-confirm`** — dismiss keeps row, accept removes), **settings add-member dialog** (open/cancel, **Escape**, **backdrop** click outside card, mobile **Escape** with drawer + `dialog` open — sidebar stays open), **`time.js-local-time` hydration** after saving a transaction (→ `/history`), **logout** (user menu → `/login`), **auth gate** (anonymous `/` → `/login`), **static assets** (`/static/css/app.css`, `/static/js/app.js` non-trivial bodies), **`modulepreload`** for `app.js` on **login** and on authenticated **layout.html** routes (dashboard `/`, **categories**, **history**, **new transaction** `/transactions`, **transaction edit** `/transactions/{id}/edit`, **notifications**, **settings** `/settings`).  
- **Gate:** `bun run test:frontend` — typecheck, unit tests, then production Vite build.  
- **Full repo check:** `bun run verify` — runs `test:frontend`, Playwright E2E (including **`GET /health`**), then `go test -race ./...`.

## 15. Risk snapshot

| Item | Severity | Note |
|------|----------|------|
| Magic hex for income/loss greens/reds | Low–medium | Drift risk; consider semantic tokens |
| Many CSS partials | Low | Import order matters; edit via `app.css` import list |
| `btn-indigo` naming | Low | Misleading for contributors |

---

*Last reviewed against repository UI sources.*
