import path from 'node:path'
import { defineConfig, devices } from '@playwright/test'

const BACKEND_PORT = process.env.E2E_BACKEND_PORT || '18080'
const FRONTEND_PORT = process.env.E2E_FRONTEND_PORT || '4173'
const GOOGLE_STUB_PORT = process.env.E2E_GOOGLE_STUB_PORT || '18081'
const API_BASE_URL = `http://127.0.0.1:${BACKEND_PORT}/api/v1/`

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  timeout: 60_000,
  retries: process.env.CI ? 2 : 0,
  globalSetup: './tests/e2e/global-setup.ts',
  use: {
    baseURL: `http://127.0.0.1:${FRONTEND_PORT}`,
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command: `node ./tests/e2e/google-books-stub.cjs --port=${GOOGLE_STUB_PORT}`,
      cwd: __dirname,
      env: {
        ...process.env,
        GOOGLE_STUB_PORT,
      },
      port: Number(GOOGLE_STUB_PORT),
      reuseExistingServer: !process.env.CI,
      timeout: 30_000,
    },
    {
      command: 'go run ./cmd/api',
      cwd: path.resolve(__dirname, '../back'),
      env: {
        ...process.env,
        APP_ENV: 'test',
        DB_HOST: process.env.TEST_DB_HOST || 'localhost',
        DB_PORT: process.env.TEST_DB_PORT || '15432',
        DB_USER: process.env.TEST_DB_USER || 'postgres',
        DB_PASSWORD: process.env.TEST_DB_PASSWORD || 'postgres',
        DB_NAME: process.env.TEST_DB_NAME || 'bookvito_test',
        DB_SSLMODE: process.env.TEST_DB_SSLMODE || 'disable',
        SERVER_PORT: BACKEND_PORT,
        JWT_SECRET: process.env.TEST_JWT_SECRET || 'test-secret',
        STORAGE_DRIVER: process.env.TEST_STORAGE_DRIVER || 's3',
        STORAGE_BUCKET: process.env.TEST_STORAGE_BUCKET || 'book-images-test',
        STORAGE_ENDPOINT: process.env.TEST_STORAGE_ENDPOINT || 'localhost:19000',
        STORAGE_PUBLIC_BASE_URL: process.env.TEST_STORAGE_PUBLIC_BASE_URL || 'http://localhost:19000/book-images-test/',
        STORAGE_ACCESS_KEY: process.env.TEST_STORAGE_ACCESS_KEY || 'minioadmin',
        STORAGE_SECRET_KEY: process.env.TEST_STORAGE_SECRET_KEY || 'minioadmin',
        GOOGLE_BOOKS_API_KEY: process.env.TEST_GOOGLE_BOOKS_API_KEY || 'test-key',
        GOOGLE_BOOKS_BASE_URL: process.env.TEST_GOOGLE_BOOKS_BASE_URL || `http://127.0.0.1:${GOOGLE_STUB_PORT}`,
      },
      port: Number(BACKEND_PORT),
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
    {
      command: `npm run dev -- --host 127.0.0.1 --port ${FRONTEND_PORT} --mode test`,
      cwd: __dirname,
      env: {
        ...process.env,
        VITE_API_BASE_URL: API_BASE_URL,
        VITE_SITE_URL: `http://127.0.0.1:${FRONTEND_PORT}`,
      },
      port: Number(FRONTEND_PORT),
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})
