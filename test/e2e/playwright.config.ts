import { defineConfig } from '@playwright/test'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const e2eDir = path.dirname(fileURLToPath(import.meta.url))
const root = path.resolve(e2eDir, '../..')

export default defineConfig({
  testDir: path.join(e2eDir, 'tests'),
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? 'github' : 'list',
  use: {
    baseURL: 'http://127.0.0.1:18080',
    trace: 'on-first-retry',
  },
  webServer: {
    command: `bash ${path.join(e2eDir, 'scripts', 'start-server.sh')}`,
    cwd: root,
    url: 'http://127.0.0.1:18080/health',
    reuseExistingServer: !process.env.CI,
    timeout: 180_000,
  },
})
