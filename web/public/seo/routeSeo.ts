import { isBookDetailPath } from '../../src/routing/paths';
import { normalizeCanonicalPath } from '../../src/seo/shared';

export type RouteSeo = {
  routeKey: string;
  title: string;
  description: string;
  robots: 'index,follow' | 'noindex,follow' | 'noindex,nofollow';
  canonicalPath: string;
  ogType?: 'website' | 'article';
};

const HOME_SEO: RouteSeo = {
  routeKey: 'home',
  title: 'Bookvito | Каталог книг',
  description: 'Bookvito помогает находить, просматривать и бронировать книги в едином каталоге.',
  robots: 'index,follow',
  canonicalPath: '/',
  ogType: 'website',
};

export const getRouteSeo = (pathname: string): RouteSeo => {
  const normalizedPath = normalizeCanonicalPath(pathname);

  if (normalizedPath === '/') {
    return HOME_SEO;
  }

  if (normalizedPath === '/about') {
    return {
      routeKey: 'about',
      title: 'О проекте | Bookvito',
      description: 'Информация о проекте Bookvito.',
      robots: 'noindex,follow',
      canonicalPath: normalizedPath,
      ogType: 'website',
    };
  }

  if (normalizedPath === '/search') {
    return {
      routeKey: 'search',
      title: 'Поиск книг | Bookvito',
      description: 'Расширенный поиск книг в каталоге Bookvito.',
      robots: 'noindex,follow',
      canonicalPath: normalizedPath,
      ogType: 'website',
    };
  }

  if (isBookDetailPath(normalizedPath)) {
    return {
      routeKey: 'book-detail',
      title: 'Книга | Bookvito',
      description: 'Карточка книги в каталоге Bookvito.',
      robots: 'index,follow',
      canonicalPath: normalizedPath,
      ogType: 'article',
    };
  }

  if (normalizedPath === '/login' || normalizedPath === '/register' || normalizedPath === '/forgot-password') {
    return {
      routeKey: 'auth',
      title: 'Авторизация | Bookvito',
      description: 'Страница входа и регистрации в Bookvito.',
      robots: 'noindex,nofollow',
      canonicalPath: normalizedPath,
      ogType: 'website',
    };
  }

  if (
    normalizedPath === '/profile' ||
    normalizedPath === '/admin' ||
    normalizedPath === '/moderation' ||
    normalizedPath === '/books/new' ||
    normalizedPath === '/color-palette'
  ) {
    return {
      routeKey: 'private-single',
      title: 'Служебная страница | Bookvito',
      description: 'Служебный раздел Bookvito.',
      robots: 'noindex,nofollow',
      canonicalPath: normalizedPath,
      ogType: 'website',
    };
  }

  if (normalizedPath === '/books' || normalizedPath.startsWith('/books/')) {
    return {
      routeKey: 'books-private',
      title: 'Раздел книг | Bookvito',
      description: 'Приватный раздел управления книгами в Bookvito.',
      robots: 'noindex,nofollow',
      canonicalPath: normalizedPath,
      ogType: 'website',
    };
  }

  return {
    routeKey: 'not-found',
    title: 'Страница не найдена | Bookvito',
    description: 'Запрошенная страница Bookvito не найдена.',
    robots: 'noindex,nofollow',
    canonicalPath: normalizedPath,
    ogType: 'website',
  };
};
