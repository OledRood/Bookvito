import { FALLBACK_SITE_URL, normalizeSiteUrl } from '../../src/seo/shared';

export const getConfiguredSiteUrl = () => {
  const configured = normalizeSiteUrl(import.meta.env.VITE_SITE_URL);
  if (configured) return configured;

  return FALLBACK_SITE_URL;
};
