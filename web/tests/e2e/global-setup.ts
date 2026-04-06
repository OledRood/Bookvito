import fs from 'node:fs/promises'
import path from 'node:path'
import { execFile } from 'node:child_process'
import { promisify } from 'node:util'

const execFileAsync = promisify(execFile)

type SeededState = {
  access_token: string
  refresh_token: string
  user: {
    id: string
    email: string
    name: string
    role: string
    avatar?: string
  }
}

export default async function globalSetup() {
  const authDir = path.resolve(__dirname, '.auth')
  const seedOutputPath = path.resolve(authDir, 'seed-output.json')
  const backendDir = path.resolve(__dirname, '../../../back')
  const password = process.env.E2E_SEED_PASSWORD || 'password123'
  const baseURL = process.env.PLAYWRIGHT_TEST_BASE_URL || 'http://127.0.0.1:4173'

  await fs.mkdir(authDir, { recursive: true })

  await execFileAsync('go', ['run', './cmd/testseed'], {
    cwd: backendDir,
    env: {
      ...process.env,
      APP_ENV: 'test',
      DB_HOST: process.env.TEST_DB_HOST || 'localhost',
      DB_PORT: process.env.TEST_DB_PORT || '15432',
      DB_USER: process.env.TEST_DB_USER || 'postgres',
      DB_PASSWORD: process.env.TEST_DB_PASSWORD || 'postgres',
      DB_NAME: process.env.TEST_DB_NAME || 'bookvito_test',
      DB_SSLMODE: process.env.TEST_DB_SSLMODE || 'disable',
      JWT_SECRET: process.env.TEST_JWT_SECRET || 'test-secret',
      E2E_SEED_PASSWORD: password,
      E2E_SEED_OUTPUT: seedOutputPath,
      ALLOW_TEST_SEED: 'true',
    },
  })

  const raw = await fs.readFile(seedOutputPath, 'utf8')
  const states = JSON.parse(raw) as Record<string, SeededState>

  for (const [role, state] of Object.entries(states)) {
    const storageState = {
      cookies: [],
      origins: [
        {
          origin: baseURL,
          localStorage: [
            { name: 'accessToken', value: state.access_token },
            { name: 'token', value: state.access_token },
            { name: 'refreshToken', value: state.refresh_token },
            { name: 'user', value: JSON.stringify(state.user) },
          ],
        },
      ],
    }
    await fs.writeFile(path.join(authDir, `${role}.json`), JSON.stringify(storageState, null, 2))
  }
}
