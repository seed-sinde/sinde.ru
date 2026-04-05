const AUTH_NEXT_EXCLUDE_PREFIXES = [
  '/auth/login',
  '/auth/register',
  '/auth/verify-email',
  '/auth/reset-password',
] as const
export const normalizeAuthNextPath = (value?: string | null) => {
  const raw = String(value || '').trim()
  if (!raw || !raw.startsWith('/') || raw.startsWith('//')) return ''
  const pathOnly = raw.split('?')[0]?.split('#')[0] || '/'
  if (AUTH_NEXT_EXCLUDE_PREFIXES.some(prefix => pathOnly === prefix || pathOnly.startsWith(`${prefix}/`))) {
    return ''
  }
  return raw
}
export const buildLoginPath = (next?: string | null) => {
  const normalizedNext = normalizeAuthNextPath(next)
  if (!normalizedNext) return '/auth/login'
  return `/auth/login?next=${encodeURIComponent(normalizedNext)}`
}
