import fs from 'node:fs'
import path from 'node:path'

export type E2ERole = 'user' | 'moder' | 'admin'

type LocalStorageEntry = { name: string; value: string }
type StorageState = {
  origins?: Array<{
    localStorage?: LocalStorageEntry[]
  }>
}

const authDir = path.resolve(__dirname, '.auth')

export const backendBaseURL = () => `http://127.0.0.1:${process.env.E2E_BACKEND_PORT || '18080'}/api/v1`

export const readRoleState = (role: E2ERole) => {
  const raw = fs.readFileSync(path.join(authDir, `${role}.json`), 'utf8')
  const parsed = JSON.parse(raw) as StorageState
  const entries = parsed.origins?.[0]?.localStorage || []
  const map = new Map(entries.map((entry) => [entry.name, entry.value]))
  const accessToken = map.get('accessToken') || map.get('token') || ''
  const refreshToken = map.get('refreshToken') || ''
  const userRaw = map.get('user') || '{}'

  return {
    accessToken,
    refreshToken,
    user: JSON.parse(userRaw) as { id?: string; email?: string; role?: string },
  }
}

export const authHeader = (role: E2ERole) => {
  const state = readRoleState(role)
  return { Authorization: `Bearer ${state.accessToken}` }
}
