import { joinURL } from 'ufo'
import { hasAuthSessionHint, syncAuthSessionHint } from '~/utils/authSessionHint'
import { interfaceLocaleToTag, normalizeInterfaceLocale } from '~/data/interfacePreferences'
type ApiJsonOptions = NonNullable<Parameters<typeof $fetch>[1]> & {
  auth?: {
    requiresSession?: boolean
    allowAutoRefresh?: boolean
  }
}
type BasicFetch = <Response>(request: string, options?: NonNullable<Parameters<typeof $fetch>[1]>) => Promise<Response>
const shouldNotifyAuthStateStale = (path: string, method: string) => {
  const normalizedPath = String(path || '').trim()
  const normalizedMethod = String(method || 'GET').toUpperCase()
  if (!normalizedPath) return false
  if (normalizedMethod === 'GET' && (normalizedPath === '/auth/summary' || normalizedPath === '/auth/admin/summary')) {
    return false
  }
  return true
}
function readCSRFFromCookieHeader(cookieHeader: string | undefined) {
  const raw = String(cookieHeader || '')
  if (!raw) return ''
  let hostCsrfToken = ''
  for (const chunk of raw.split(';')) {
    const [namePart, ...valueParts] = chunk.split('=')
    const name = String(namePart || '').trim()
    if (name !== 'csrf_token' && name !== '__Host-csrf_token') continue
    const value = valueParts.join('=').trim()
    if (!value) return ''
    const decoded = (() => {
      try {
        return decodeURIComponent(value)
      } catch {
        return value
      }
    })()
    if (name === 'csrf_token') {
      return decoded
    }
    hostCsrfToken = decoded
  }
  return hostCsrfToken
}
function readCookieFromDocument(names: readonly string[]) {
  if (!import.meta.client) return ''
  const raw = String(document.cookie || '')
  if (!raw) return ''
  let fallback = ''
  for (const chunk of raw.split(';')) {
    const [namePart, ...valueParts] = chunk.split('=')
    const name = String(namePart || '').trim()
    if (!name || !names.includes(name)) continue
    const value = valueParts.join('=').trim()
    if (!value) return ''
    const decoded = (() => {
      try {
        return decodeURIComponent(value)
      } catch {
        return value
      }
    })()
    if (name === 'csrf_token') {
      return decoded
    }
    fallback = decoded
  }
  return fallback
}
function copyHeaderIfMissing(
  headers: Record<string, string>,
  requestHeaders: Record<string, string | undefined>,
  sourceName: string,
  targetName = sourceName
) {
  if (headers[targetName] || headers[targetName.toLowerCase()]) return
  const value = requestHeaders[sourceName]
  if (!value) return
  headers[targetName] = value
}
const getStatusCodeFromError = (err: any) => Number(err?.status || err?.statusCode || err?.response?.status || 0)
const getErrorMessage = (err: any) =>
  String(err?.data?.message || err?.response?._data?.message || err?.message || '').trim()
const isSessionStateError = (message: string) => {
  const normalized = message.trim().toLowerCase()
  if (!normalized) return false
  return (
    normalized === 'authentication required' ||
    normalized === 'invalid or expired session' ||
    normalized === 'unauthorized' ||
    normalized === 'требуется авторизация.' ||
    normalized === 'требуется авторизация' ||
    normalized === 'недействительная или просроченная сессия.' ||
    normalized === 'недействительная или просроченная сессия'
  )
}
const isProtectedKitchenPath = (path: string, method: string) => {
  if (path === '/kitchen/ingredients/account') return true
  if (path.startsWith('/kitchen/ingredients/custom')) return true
  if (path.startsWith('/kitchen/ingredients/favorites')) return true
  if (path.startsWith('/kitchen/recipes/manage/')) return true
  if (path.startsWith('/kitchen/recipes/mine')) return true
  if (path.startsWith('/kitchen/admin/')) return true
  if (/^\/kitchen\/recipes\/[^/]+\/favorite(?:\?.*)?$/.test(path)) return true
  if (path === '/kitchen/recipes' && method === 'POST') return true
  if (/^\/kitchen\/recipes\/[^/?]+(?:\?.*)?$/.test(path) && ['PATCH', 'DELETE'].includes(method)) return true
  return false
}
const isProtectedPath = (path: string, method: string) => {
  const normalized = String(path || '').trim()
  if (!normalized) return false
  if (normalized.startsWith('/auth/')) {
    return ![
      '/auth/register',
      '/auth/verify-email/request',
      '/auth/verify-email/confirm',
      '/auth/login',
      '/auth/login/2fa',
      '/auth/refresh',
      '/auth/logout',
      '/auth/password/forgot',
      '/auth/password/reset'
    ].includes(normalized)
  }
  if (normalized === '/media/upload') return true
  return isProtectedKitchenPath(normalized, method)
}
const canAutoRefreshByPath = (path: string, method: string, auth: ApiJsonOptions['auth'], err?: any) => {
  const normalized = String(path || '').trim()
  if (!normalized) return false
  if (normalized === '/auth/refresh') return false
  if (auth?.allowAutoRefresh === false) return false
  if (!isProtectedPath(normalized, method) && auth?.requiresSession !== true) return false
  if (err && !isSessionStateError(getErrorMessage(err))) return false
  return true
}
const notifyAuthStateStale = () => {
  if (!import.meta.client) return
  window.dispatchEvent(new Event('auth-state-stale'))
}
/**
 * Provides shared API helpers for JSON and NDJSON requests.
 */
export const useAPI = () => {
  const basicFetch = $fetch as BasicFetch
  /**
   * Executes a proxied JSON request with auth-aware retry logic.
   */
  const json = async <T>(path: string, options: ApiJsonOptions = {}): Promise<T> => {
    const p = path.startsWith('/') ? path : `/${path}`
    const method = String((options as Record<string, any>).method ?? 'GET')
    const normalizedMethod = method.toUpperCase()
    const auth = options.auth
    const headers: Record<string, string> = {}
    if (options.headers) {
      for (const [key, value] of Object.entries(options.headers as Record<string, string>)) {
        if (typeof value === 'string') headers[key] = value
      }
    }
    if (import.meta.server) {
      const reqHeaders = useRequestHeaders([
        'cookie',
        'x-csrf-token',
        'user-agent',
        'accept-language',
        'x-forwarded-for',
        'x-real-ip'
      ])
      copyHeaderIfMissing(headers, reqHeaders, 'cookie')
      copyHeaderIfMissing(headers, reqHeaders, 'x-csrf-token')
      copyHeaderIfMissing(headers, reqHeaders, 'user-agent')
      copyHeaderIfMissing(headers, reqHeaders, 'accept-language')
      copyHeaderIfMissing(headers, reqHeaders, 'x-forwarded-for')
      copyHeaderIfMissing(headers, reqHeaders, 'x-real-ip')
      if (!headers['X-CSRF-Token'] && !headers['x-csrf-token']) {
        const csrfToken = readCSRFFromCookieHeader(reqHeaders.cookie)
        if (csrfToken) headers['X-CSRF-Token'] = csrfToken
      }
    }
    if (import.meta.client && !headers['Accept-Language'] && !headers['accept-language']) {
      const uiPreferences = useUiPreferencesStore()
      headers['Accept-Language'] = interfaceLocaleToTag(normalizeInterfaceLocale(uiPreferences.interfaceLocale))
    }
    if (import.meta.client && !['GET', 'HEAD', 'OPTIONS'].includes(normalizedMethod)) {
      const csrfToken = readCookieFromDocument(['csrf_token', '__Host-csrf_token']).trim()
      if (csrfToken && !headers['X-CSRF-Token'] && !headers['x-csrf-token']) {
        headers['X-CSRF-Token'] = csrfToken
      }
    }
    const requestConfig = {
      ...(options as Record<string, any>),
      method: method as any,
      credentials: 'include' as const,
      headers
    }
    delete (requestConfig as Record<string, any>).auth
    let retriedAfterRefresh = false
    while (true) {
      try {
        const response = await basicFetch<T>(`/api/proxy${p}`, requestConfig)
        return response as T
      } catch (err: any) {
        const status = getStatusCodeFromError(err)
        const sessionHint = import.meta.client ? hasAuthSessionHint() : false
        const canRetry =
          import.meta.client &&
          status === 401 &&
          !retriedAfterRefresh &&
          sessionHint &&
          canAutoRefreshByPath(p, normalizedMethod, auth, err)
        if (!canRetry) {
          if (
            import.meta.client &&
            sessionHint &&
            status === 401 &&
            isSessionStateError(getErrorMessage(err)) &&
            shouldNotifyAuthStateStale(p, normalizedMethod)
          ) {
            notifyAuthStateStale()
          }
          throw err
        }
        retriedAfterRefresh = true
        try {
          const refreshRes = await basicFetch<ApiResponseWithData<{ user?: AuthUser }>>('/api/proxy/auth/refresh', {
            method: 'POST',
            credentials: 'include',
            headers
          })
          const refreshedUser = refreshRes?.data?.user || null
          if (import.meta.client) {
            syncAuthSessionHint(Boolean(refreshedUser))
            window.dispatchEvent(
              new CustomEvent('auth-user-refreshed', {
                detail: {
                  user: refreshedUser
                }
              })
            )
            const csrfToken = readCookieFromDocument(['csrf_token', '__Host-csrf_token']).trim()
            if (csrfToken) {
              headers['X-CSRF-Token'] = csrfToken
            }
          }
        } catch {
          syncAuthSessionHint(false)
          notifyAuthStateStale()
          throw err
        }
      }
    }
  }
  /**
   * Executes a proxied NDJSON stream request and emits parsed lines.
   */
  const stream = async (path: string, onLine: (line: Trait) => void, options: RequestInit = {}): Promise<void> => {
    const p = path.startsWith('/') ? path : `/${path}`
    const config = useRuntimeConfig()
    const base = import.meta.client ? '' : useRequestURL().origin || config.public.baseURL
    const url = joinURL(base, '/api/proxy', p)
    const headers: Record<string, string> = {
      Accept: 'application/x-ndjson, application/json',
      ...((options.headers as Record<string, string>) || {})
    }
    if (!import.meta.client) {
      const req = useRequestHeaders(['cookie'])
      if (req.cookie) headers.cookie = req.cookie
    }
    const res = await fetch(url, {
      ...options,
      method: options.method ?? 'GET',
      credentials: 'include',
      headers
    })
    if (!res.ok) {
      return
    }
    const ct = res.headers.get('content-type') || ''
    if (!ct.includes('application/x-ndjson')) return
    const reader = res.body?.getReader()
    if (!reader) return
    const decoder = new TextDecoder()
    let buffer = ''
    while (true) {
      const { value, done } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      let idx: number
      while ((idx = buffer.indexOf('\n')) !== -1) {
        const line = buffer.slice(0, idx).trim()
        buffer = buffer.slice(idx + 1)
        if (!line) continue
        try {
          const parsed = JSON.parse(line) as Trait
          if (
            parsed &&
            typeof parsed.t_uuid === 'string' &&
            typeof parsed.t_key === 'string' &&
            typeof parsed.t_value === 'string'
          ) {
            onLine(parsed)
          }
        } catch {}
      }
    }
    const tail = buffer.trim()
    if (tail) {
      try {
        const parsed = JSON.parse(tail) as Trait
        if (
          parsed &&
          typeof parsed.t_uuid === 'string' &&
          typeof parsed.t_key === 'string' &&
          typeof parsed.t_value === 'string'
        ) {
          onLine(parsed)
        }
      } catch {
        // ignore
      }
    }
  }
  return { json, stream }
}
