import fs from 'node:fs'
import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { buildBookPath } from './src/routing/paths'
import { buildRobotsTxt, buildSitemapXml, FALLBACK_SITE_URL, normalizeSiteUrl } from './src/seo/shared'

type SitemapBookItem = {
  id: string
  title?: string
}

type BookListResponse = {
  items?: SitemapBookItem[]
  has_more?: boolean
}

const normalizeApiBaseUrl = (rawValue?: string) => {
  if (!rawValue) return null
  return rawValue.endsWith('/') ? rawValue : `${rawValue}/`
}

const fetchSitemapBookPaths = async (apiBaseUrl?: string) => {
  const normalizedApiBaseUrl = normalizeApiBaseUrl(apiBaseUrl)
  if (!normalizedApiBaseUrl) return []

  const collectedPaths: string[] = []
  const seenIds = new Set<string>()
  let offset = 0
  const limit = 100
  let hasMore = true

  while (hasMore) {
    const endpoint = `${normalizedApiBaseUrl}books/list?limit=${limit}&offset=${offset}`
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 3000)
    const response = await fetch(endpoint, {
      headers: { Accept: 'application/json' },
      signal: controller.signal,
    })
    clearTimeout(timeout)

    if (!response.ok) {
      throw new Error(`Failed to fetch sitemap books: ${response.status} ${response.statusText}`)
    }

    const payload = (await response.json()) as BookListResponse
    const items = Array.isArray(payload.items) ? payload.items : []

    for (const item of items) {
      if (!item?.id || seenIds.has(item.id)) continue
      seenIds.add(item.id)
      collectedPaths.push(buildBookPath(item.id, item.title))
    }

    hasMore = Boolean(payload.has_more)
    offset += limit
  }

  return collectedPaths
}

// Project uses a `public/` directory that contains `index.html` and app entry files.
// Set Vite root to that folder and output the build to ../dist so Docker can copy /app/dist.
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, '')
  const siteUrl = normalizeSiteUrl(env.VITE_SITE_URL) || FALLBACK_SITE_URL
  const apiBaseUrl = env.VITE_API_BASE_URL || `${siteUrl}/api/v1/`

  return {
    root: 'public',
    plugins: [
      react(),
      {
        name: 'generate-seo-files',
        async closeBundle() {
          const outDir = path.resolve(__dirname, 'dist')
          let sitemapPaths = ['/']

          try {
            const bookPaths = await fetchSitemapBookPaths(apiBaseUrl)
            sitemapPaths = ['/', ...bookPaths]
          } catch (error) {
            console.warn('[generate-seo-files] Failed to fetch book URLs for sitemap, homepage only will be emitted.', error)
          }

          fs.mkdirSync(outDir, { recursive: true })
          fs.writeFileSync(path.join(outDir, 'robots.txt'), buildRobotsTxt(siteUrl), 'utf8')
          fs.writeFileSync(path.join(outDir, 'sitemap.xml'), buildSitemapXml(siteUrl, sitemapPaths), 'utf8')
        },
      },
    ],
    build: {
      outDir: '../dist'
    }
  }
})
