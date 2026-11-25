import { defineConfig, devices } from '@playwright/test';
import { launchOptions } from './browser.js';

// AELIG Mail root address (see ../config.toml).
const BASE_URL = process.env.AELIG_MAIL_URL || 'http://localhost:9000';
const MAILHOG_URL = process.env.MAILHOG_URL || 'http://localhost:8025';

export default defineConfig({
  testDir: './specs',
  workers: 1,
  reporter: process.env.CI ? 'list' : [['list'], ['html', { open: 'never' }]],

  // Log in once and save the session before any tests run.
  globalSetup: './global-setup.js',

  // Kill the AELIG Mail instance resetDB().
  globalTeardown: './global-teardown.js',

  webServer: [
    // Run mailhog local SMTP test server.
    ...(process.env.MAILHOG_EXTERNAL === '1' ? [] : [{
      command: 'pkill -9 -- "${MAILHOG:-mailhog}"; "${MAILHOG:-mailhog}"',
      url: MAILHOG_URL,
      reuseExistingServer: false,
    }]),
    // Create a fresh AELIG Mail database instance for every test run.
    {
      command: 'pkill -9 aelig-mail; ./aelig-mail --config "${AELIG_MAIL_CONFIG:-${DEV_CONFIG:-config.toml}}" --install --yes && ./aelig-mail --config "${AELIG_MAIL_CONFIG:-${DEV_CONFIG:-config.toml}}" --static-dir static',
      cwd: '..',
      stdout: 'pipe',
      stderr: 'pipe',
      url: `${BASE_URL}/health`,
      reuseExistingServer: false,
      env: { AELIG_MAIL_ADMIN_USER: 'admin', AELIG_MAIL_ADMIN_PASSWORD: 'aelig-mail' },
    },
  ],

  use: {
    launchOptions,
    baseURL: BASE_URL,
    storageState: './.auth/admin.json',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },

  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'], viewport: { width: 1400, height: 950 } } },
  ],
});
