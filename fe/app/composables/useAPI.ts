import {joinURL} from "ufo"
import {hasAuthSessionHint, syncAuthSessionHint} from "~/utils/authSessionHint"
import {interfaceLocaleToTag, normalizeInterfaceLocale} from "~/data/themePreferences"
type HttpMethod = NonNullable<ApiJsonOptions["method"]>
type AuthOptions = {
  requiresSession?: boolean
  allowAutoRefresh?: boolean
}
type ApiJsonOptions = NonNullable<Parameters<typeof $fetch>[1]> & {
  auth?: AuthOptions
}
type FetchOptions = NonNullable<Parameters<typeof $fetch>[1]>
type BasicFetch = <T>(request: string, options?: FetchOptions) => Promise<T>
type HeaderMap = Record<string, string>
type FetchErrorLike = {
  status?: number
  statusCode?: number
  message?: string
  data?: {message?: string}
  response?: {
    status?: number
    _data?: {message?: string}
  }
}
const SAFE_METHODS = new Set(["GET", "HEAD", "OPTIONS"])
const AUTH_PUBLIC_PATHS = new Set([
  "/auth/register",
  "/auth/verify-email/request",
  "/auth/verify-email/confirm",
  "/auth/login",
  "/auth/login/2fa",
  "/auth/refresh",
  "/auth/logout",
  "/auth/password/forgot",
  "/auth/password/reset"
])
const AUTH_STALE_SKIP_PATHS = new Set(["/auth/summary", "/auth/admin/summary"])
const CSRF_COOKIE_NAMES = ["csrf_token", "__Host-csrf_token"] as const
const SESSION_STATE_MESSAGES = new Set([
  "authentication required",
  "invalid or expired session",
  "unauthorized",
  "требуется авторизация.",
  "требуется авторизация",
  "недействительная или просроченная сессия.",
  "недействительная или просроченная сессия"
])
const KITCHEN_PROTECTED_PREFIXES = [
  "/kitchen/ingredients/custom",
  "/kitchen/ingredients/favorites",
  "/kitchen/recipes/manage/",
  "/kitchen/recipes/mine",
  "/kitchen/admin/"
] as const
const KITCHEN_RECIPE_FAVORITE_RE = /^\/kitchen\/recipes\/[^/]+\/favorite(?:\?.*)?$/
const KITCHEN_RECIPE_ITEM_RE = /^\/kitchen\/recipes\/[^/?]+(?:\?.*)?$/

const toPath = (p: string) => (p.startsWith("/") ? p : `/${p}`)
const getMethod = (m?: ApiJsonOptions["method"]): HttpMethod => String(m ?? "GET").toUpperCase() as HttpMethod
const decodeCookieValue = (v: string) => {
  try {
    return decodeURIComponent(v)
  } catch {
    return v
  }
}
const getStatusCodeFromError = (e: FetchErrorLike) => Number(e?.status || e?.statusCode || e?.response?.status || 0)
const getErrorMessage = (e: FetchErrorLike) =>
  String(e?.data?.message || e?.response?._data?.message || e?.message || "").trim()
const shouldNotifyAuthStateStale = (path: string, method: string) =>
  Boolean(path) && !(method === "GET" && AUTH_STALE_SKIP_PATHS.has(path))
const isSessionStateError = (msg: string) => SESSION_STATE_MESSAGES.has(msg.trim().toLowerCase())
const notifyAuthStateStale = () => {
  if (!import.meta.client) return
  window.dispatchEvent(new Event("auth-state-stale"))
}
const readCookieValue = (raw: string, names: readonly string[]) => {
  if (!raw) return ""
  let fallback = ""
  for (const chunk of raw.split(";")) {
    const [namePart, ...valueParts] = chunk.split("=")
    const name = String(namePart || "").trim()
    if (!name || !names.includes(name)) continue
    const value = valueParts.join("=").trim()
    if (!value) return ""
    const decoded = decodeCookieValue(value)
    if (name === "csrf_token") return decoded
    fallback = decoded
  }
  return fallback
}
const readCSRFFromCookieHeader = (cookieHeader?: string) =>
  readCookieValue(String(cookieHeader || ""), CSRF_COOKIE_NAMES)
const readCookieFromDocument = (names: readonly string[]) =>
  !import.meta.client ? "" : readCookieValue(String(document.cookie || ""), names)
const mergeHeaders = (src?: HeadersInit) => {
  const out: HeaderMap = {}
  if (!src) return out
  if (src instanceof Headers) {
    src.forEach((v, k) => (out[k.toLowerCase()] = v))
    return out
  }
  if (Array.isArray(src)) {
    for (const [k, v] of src) {
      if (typeof v === "string") out[String(k).toLowerCase()] = v
    }
    return out
  }
  for (const [k, v] of Object.entries(src)) {
    if (typeof v === "string") out[k.toLowerCase()] = v
  }
  return out
}
const setHeaderIfMissing = (headers: HeaderMap, key: string, value?: string) => {
  const k = key.toLowerCase()
  if (!value || headers[k]) return
  headers[k] = value
}
const copyRequestHeadersIfMissing = (
  headers: HeaderMap,
  reqHeaders: Record<string, string | undefined>,
  keys: readonly string[]
) => {
  for (const key of keys) setHeaderIfMissing(headers, key, reqHeaders[key])
}
const isProtectedKitchenPath = (path: string, method: string) =>
  path === "/kitchen/ingredients/account" ||
  KITCHEN_PROTECTED_PREFIXES.some(p => path.startsWith(p)) ||
  KITCHEN_RECIPE_FAVORITE_RE.test(path) ||
  (path === "/kitchen/recipes" && method === "POST") ||
  (KITCHEN_RECIPE_ITEM_RE.test(path) && ["PATCH", "DELETE"].includes(method))
const isProtectedPath = (path: string, method: string) =>
  !path
    ? false
    : path.startsWith("/auth/")
      ? !AUTH_PUBLIC_PATHS.has(path)
      : path === "/media/upload"
        ? true
        : isProtectedKitchenPath(path, method)
const canAutoRefreshByPath = (path: string, method: string, auth?: AuthOptions, err?: FetchErrorLike) =>
  Boolean(path) &&
  path !== "/auth/refresh" &&
  auth?.allowAutoRefresh !== false &&
  (isProtectedPath(path, method) || auth?.requiresSession === true) &&
  (!err || isSessionStateError(getErrorMessage(err)))
const buildJsonHeaders = (path: string, method: string, options: ApiJsonOptions) => {
  const headers = mergeHeaders(options.headers)
  if (import.meta.server) {
    const reqHeaders = useRequestHeaders([
      "cookie",
      "x-csrf-token",
      "user-agent",
      "accept-language",
      "x-forwarded-for",
      "x-real-ip"
    ])
    copyRequestHeadersIfMissing(headers, reqHeaders, [
      "cookie",
      "x-csrf-token",
      "user-agent",
      "accept-language",
      "x-forwarded-for",
      "x-real-ip"
    ])
    if (!headers["x-csrf-token"]) {
      setHeaderIfMissing(headers, "x-csrf-token", readCSRFFromCookieHeader(reqHeaders.cookie))
    }
  }
  if (import.meta.client && !headers["accept-language"]) {
    const ui = useUiPreferencesStore()
    headers["accept-language"] = interfaceLocaleToTag(normalizeInterfaceLocale(ui.interfaceLocale))
  }
  if (import.meta.client && !SAFE_METHODS.has(method) && !headers["x-csrf-token"]) {
    const csrf = readCookieFromDocument(CSRF_COOKIE_NAMES).trim()
    if (csrf) {
      headers["x-csrf-token"] = csrf
    }
  }
  return headers
}
const buildJsonRequestConfig = (options: ApiJsonOptions, headers: HeaderMap) => {
  const {auth: _auth, ...rest} = options
  return {
    ...rest,
    credentials: "include" as const,
    headers
  }
}
const emitAuthUserRefreshed = (user: AuthUser | null) => {
  if (!import.meta.client) return
  window.dispatchEvent(new CustomEvent("auth-user-refreshed", {detail: {user}}))
}
const refreshSession = async (basicFetch: BasicFetch, headers: HeaderMap) => {
  const res = await basicFetch<ApiResponseWithData<{user?: AuthUser}>>("/api/proxy/auth/refresh", {
    method: "POST",
    credentials: "include",
    headers
  })
  const user = res?.data?.user || null
  if (import.meta.client) {
    syncAuthSessionHint(Boolean(user))
    emitAuthUserRefreshed(user)
    const csrf = readCookieFromDocument(CSRF_COOKIE_NAMES).trim()
    if (csrf) headers["x-csrf-token"] = csrf
  }
  return user
}
const handleJsonError = (path: string, method: string, err: FetchErrorLike) => {
  if (
    import.meta.client &&
    hasAuthSessionHint() &&
    getStatusCodeFromError(err) === 401 &&
    isSessionStateError(getErrorMessage(err)) &&
    shouldNotifyAuthStateStale(path, method)
  )
    notifyAuthStateStale()
}
const parseJsonLikeStreamPayload = async (res: Response, onLine: (line: unknown) => void) => {
  const text = await res.text()
  const payload = JSON.parse(text)
  if (!Array.isArray(payload)) {
    onLine(payload)
    return
  }
  const isTuple = payload.length === 2 && typeof payload[0] === "string" && typeof payload[1] === "string"
  if (isTuple) onLine(payload)
  else payload.forEach(item => onLine(item))
}
const createHttpError = async (res: Response) => {
  const text = await res.text().catch(() => "")
  const err = new Error(`Stream request failed: HTTP ${res.status}`) as Error & {
    status?: number
    statusCode?: number
    data?: string
  }
  err.status = res.status
  err.statusCode = res.status
  err.data = text
  return err
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
    const p = toPath(path)
    const method = getMethod(options.method)
    const headers = buildJsonHeaders(p, method, options)
    const cfg = buildJsonRequestConfig({...options, method}, headers)
    try {
      return await basicFetch<T>(`/api/proxy${p}`, cfg)
    } catch (e) {
      const err = e as FetchErrorLike
      const canRetry =
        import.meta.client &&
        hasAuthSessionHint() &&
        getStatusCodeFromError(err) === 401 &&
        canAutoRefreshByPath(p, method, options.auth, err)
      if (!canRetry) {
        handleJsonError(p, method, err)
        throw err
      }
      try {
        await refreshSession(basicFetch, headers)
        return await basicFetch<T>(`/api/proxy${p}`, cfg)
      } catch {
        if (import.meta.client) {
          syncAuthSessionHint(false)
        }
        notifyAuthStateStale()
        throw err
      }
    }
  }
  /**
   * Executes a proxied NDJSON stream request and emits parsed lines.
   */
  const stream = async (path: string, onLine: (line: unknown) => void, options: RequestInit = {}): Promise<void> => {
    const p = toPath(path)
    const config = useRuntimeConfig()
    const base = import.meta.client ? "" : useRequestURL().origin || config.public.baseURL
    const url = joinURL(base, "/api/proxy", p)
    const headers = mergeHeaders(options.headers)
    setHeaderIfMissing(headers, "accept", "application/x-ndjson, application/json")
    if (import.meta.server) {
      const req = useRequestHeaders(["cookie"])
      setHeaderIfMissing(headers, "cookie", req.cookie)
    }
    const res = await fetch(url, {
      ...options,
      method: options.method ?? "GET",
      credentials: "include",
      headers
    })
    if (!res.ok) throw await createHttpError(res)
    const ct = String(res.headers.get("content-type") || "").toLowerCase()
    if (ct.includes("application/json") && !ct.includes("application/x-ndjson")) {
      await parseJsonLikeStreamPayload(res, onLine)
      return
    }
    if (!ct.includes("application/x-ndjson")) throw new Error(`Unexpected stream content-type: ${ct || "empty"}`)
    const reader = res.body?.getReader()
    if (!reader) throw new Error("Response body reader is unavailable")
    const decoder = new TextDecoder()
    let buf = ""
    while (true) {
      const {value, done} = await reader.read()
      if (done) break
      buf += decoder.decode(value, {stream: true})
      let i = buf.indexOf("\n")
      while (i !== -1) {
        const line = buf.slice(0, i).trim()
        buf = buf.slice(i + 1)
        if (line) {
          onLine(JSON.parse(line))
        }
        i = buf.indexOf("\n")
      }
    }
    const tail = buf.trim()
    if (!tail) return
    onLine(JSON.parse(tail))
  }
  return {json, stream}
}
