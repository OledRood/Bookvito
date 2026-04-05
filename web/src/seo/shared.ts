export const FALLBACK_SITE_URL = 'https://bookvito.ru';

export const normalizeSiteUrl = (rawValue?: string | null) => {
  if (!rawValue) return null;

  try {
    const url = new URL(rawValue);
    return url.origin;
  } catch {
    return null;
  }
};

export const normalizeCanonicalPath = (pathname: string) => {
  if (!pathname || pathname === '/') return '/';

  const withoutQuery = pathname.split('?')[0]?.split('#')[0] || '/';
  const trimmed = withoutQuery.replace(/\/+$/, '');
  return trimmed || '/';
};

export const buildCanonicalUrl = (siteUrl: string, pathname: string) => {
  const normalizedSiteUrl = normalizeSiteUrl(siteUrl) || FALLBACK_SITE_URL;
  const normalizedPath = normalizeCanonicalPath(pathname);
  return normalizedPath === '/' ? `${normalizedSiteUrl}/` : `${normalizedSiteUrl}${normalizedPath}`;
};

export const buildAbsoluteUrl = (siteUrl: string, path: string) => {
  return buildCanonicalUrl(siteUrl, path);
};

export const buildRobotsTxt = (siteUrl: string) => {
  const sitemapUrl = buildCanonicalUrl(siteUrl, '/sitemap.xml');

  return [
    'User-agent: *',
    'Allow: /',
    'Disallow: /admin',
    'Disallow: /books/new',
    'Disallow: /books/my',
    'Disallow: /books/read',
    'Disallow: /books/reserved',
    'Disallow: /books/return',
    'Disallow: /books/shelf',
    'Disallow: /color-palette',
    'Disallow: /create',
    'Disallow: /forgot-password',
    'Disallow: /login',
    'Disallow: /moder',
    'Disallow: /moderation',
    'Disallow: /profile',
    'Disallow: /register',
    '',
    `Sitemap: ${sitemapUrl}`,
    '',
  ].join('\n');
};

export const buildSitemapXml = (siteUrl: string, paths: string[] = ['/']) => {
  const uniquePaths = Array.from(new Set(paths.map((path) => normalizeCanonicalPath(path))));
  const urlEntries = uniquePaths
    .map((path) => `  <url>\n    <loc>${buildCanonicalUrl(siteUrl, path)}</loc>\n  </url>`)
    .join('\n');

  return [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    urlEntries,
    '</urlset>',
    '',
  ].join('\n');
};
