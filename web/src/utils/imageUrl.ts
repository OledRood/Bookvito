// Resolve image URLs coming from the backend.
// If the backend returns an absolute URL (http...), return as-is.
// If it returns a root-relative path like "/images/xxx.png" or "images/xxx.png",
// prefix it with the API server root derived from VITE_API_BASE_URL (fallback to http://localhost:8080).
export default function resolveImageUrl(imageUrl?: string | null): string {
  if (!imageUrl) return '';
  // Already absolute
  if (/^https?:\/\//i.test(imageUrl)) return imageUrl;

  const meta: any = (import.meta as any) || {};
  const rawBase = (meta && meta.env && meta.env.VITE_API_BASE_URL)
    ? meta.env.VITE_API_BASE_URL
    : 'http://localhost:8080/api/v1/';

  // Derive server root (remove trailing /api/v1/ if present)
  let serverRoot = rawBase.replace(/\/api\/v1\/?$/, '').replace(/\/$/, '');

  if (imageUrl.startsWith('/')) return serverRoot + imageUrl;
  return serverRoot + '/' + imageUrl;
}
