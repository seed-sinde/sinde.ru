const AUTH_SERVER_COOKIE_NAMES = [
  '__Host-access_token',
  'access_token',
  '__Host-refresh_token',
  'refresh_token'
] as const
const AUTH_CLIENT_COOKIE_NAMES = [
  '__Host-auth_session_hint',
  'auth_session_hint'
] as const
const AUTH_SESSION_RUNTIME_KEY = '__traitsAuthSessionHint__'
const hasNamedCookie = (cookieHeader: string | null | undefined, names: readonly string[]) => {
  const raw = String(cookieHeader || '').trim()
  if (!raw) return false
  for (const chunk of raw.split(';')) {
    const [namePart, ...valueParts] = chunk.split('=')
    const name = String(namePart || '').trim()
    if (!name || !names.includes(name)) continue
    const value = valueParts.join('=').trim()
    if (value) return true
  }
  return false
}
export const hasAuthSessionHint = (cookieHeader?: string | null) => {
  if (import.meta.server) {
    return hasNamedCookie(cookieHeader, AUTH_SERVER_COOKIE_NAMES)
  }
  const scope = globalThis as typeof globalThis & {
    [AUTH_SESSION_RUNTIME_KEY]?: boolean
  }
  return Boolean(scope[AUTH_SESSION_RUNTIME_KEY]) || hasNamedCookie(document.cookie, AUTH_CLIENT_COOKIE_NAMES)
}
export const syncAuthSessionHint = (enabled: boolean) => {
  if (!import.meta.client) return
  const scope = globalThis as typeof globalThis & {
    [AUTH_SESSION_RUNTIME_KEY]?: boolean
  }
  scope[AUTH_SESSION_RUNTIME_KEY] = enabled
  const secure = window.location.protocol === 'https:' ? '; Secure' : ''
  const preferredName = window.location.protocol === 'https:' ? '__Host-auth_session_hint' : 'auth_session_hint'
  if (enabled) {
    document.cookie = `${preferredName}=1; Path=/; SameSite=Strict${secure}`
    return
  }
  for (const name of AUTH_CLIENT_COOKIE_NAMES) {
    document.cookie = `${name}=; Max-Age=0; Path=/; SameSite=Strict${secure}`
  }
}
