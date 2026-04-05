import React from 'react';
import { Helmet } from 'react-helmet-async';
import { useLocation } from 'react-router-dom';
import { buildAbsoluteUrl, buildCanonicalUrl } from '../src/seo/shared';
import { getConfiguredSiteUrl } from './seo/runtime';
import { getRouteSeo } from './seo/routeSeo';

const SeoManager: React.FC = () => {
  const location = useLocation();
  const seo = getRouteSeo(location.pathname);
  const siteUrl = getConfiguredSiteUrl();
  const canonicalUrl = buildCanonicalUrl(siteUrl, seo.canonicalPath);
  const fallbackImageUrl = buildAbsoluteUrl(siteUrl, '/images/default-book.png');

  return (
    <Helmet prioritizeSeoTags>
      <html lang="ru" />
      <title>{seo.title}</title>
      <meta name="description" content={seo.description} />
      <meta name="robots" content={seo.robots} />
      <link rel="canonical" href={canonicalUrl} />

      <meta property="og:site_name" content="Bookvito" />
      <meta property="og:locale" content="ru_RU" />
      <meta property="og:type" content={seo.ogType || 'website'} />
      <meta property="og:title" content={seo.title} />
      <meta property="og:description" content={seo.description} />
      <meta property="og:url" content={canonicalUrl} />
      <meta property="og:image" content={fallbackImageUrl} />
      <meta property="og:image:alt" content="Bookvito" />

      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={seo.title} />
      <meta name="twitter:description" content={seo.description} />
      <meta name="twitter:image" content={fallbackImageUrl} />
      <meta name="twitter:image:alt" content="Bookvito" />
    </Helmet>
  );
};

export default SeoManager;
