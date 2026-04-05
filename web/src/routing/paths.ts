const CYRILLIC_TO_LATIN: Record<string, string> = {
  а: 'a',
  б: 'b',
  в: 'v',
  г: 'g',
  д: 'd',
  е: 'e',
  ё: 'e',
  ж: 'zh',
  з: 'z',
  и: 'i',
  й: 'y',
  к: 'k',
  л: 'l',
  м: 'm',
  н: 'n',
  о: 'o',
  п: 'p',
  р: 'r',
  с: 's',
  т: 't',
  у: 'u',
  ф: 'f',
  х: 'h',
  ц: 'ts',
  ч: 'ch',
  ш: 'sh',
  щ: 'sch',
  ъ: '',
  ы: 'y',
  ь: '',
  э: 'e',
  ю: 'yu',
  я: 'ya',
};

const UUID_PREFIX_RE = /^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})(?:-(.+))?$/i;

export const BOOKS_NEW_PATH = '/books/new';
export const MODERATION_PATH = '/moderation';

export const normalizePathname = (pathname: string) => {
  if (!pathname) return '/';
  const withoutQuery = pathname.split('?')[0]?.split('#')[0] || '/';
  const trimmed = withoutQuery.replace(/\/+$/, '');
  return trimmed || '/';
};

export const slugifySegment = (value: string) => {
  return value
    .trim()
    .toLowerCase()
    .split('')
    .map((char) => CYRILLIC_TO_LATIN[char] ?? char)
    .join('')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .replace(/-{2,}/g, '-');
};

export const buildBookPath = (bookId: string | number, title?: string) => {
  const id = String(bookId).trim();
  const slug = slugifySegment(title || '');
  return slug ? `/books/${id}-${slug}` : `/books/${id}`;
};

export const extractBookId = (bookIdSlug?: string | null) => {
  if (!bookIdSlug) return '';

  const cleaned = bookIdSlug.trim();
  const match = cleaned.match(UUID_PREFIX_RE);
  if (match?.[1]) {
    return match[1].toLowerCase();
  }

  return cleaned;
};

export const isBookDetailPath = (pathname: string) => {
  const normalized = normalizePathname(pathname);
  const segment = normalized.startsWith('/books/') ? normalized.slice('/books/'.length) : '';
  return UUID_PREFIX_RE.test(segment);
};
