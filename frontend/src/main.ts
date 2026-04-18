import './app.css'

import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'

setBrowserTimezoneCookie()
applyLocalTimeElements()
initShellSidebar()
initSettingsMemberDialog()
initCategoryModal()
