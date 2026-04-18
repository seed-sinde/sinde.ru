import type {HTTPMethod} from "h3"
import {
  defineEventHandler,
  getHeaders,
  readRawBody,
  setResponseStatus,
  setHeader,
  appendHeader,
  createError,
  sendStream
} from "h3"
import {joinURL} from "ufo"
const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "proxy-connection",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade"
])
function shouldSkipRequestHeader(name: string) {
  const lower = name.toLowerCase()
  return lower === "host" || lower === "content-length" || HOP_BY_HOP_HEADERS.has(lower)
}
function shouldSkipResponseHeader(name: string) {
  const lower = name.toLowerCase()
  return (
    lower === "set-cookie" ||
    lower === "content-length" ||
    lower === "content-encoding" ||
    HOP_BY_HOP_HEADERS.has(lower)
  )
}
function isExpectedUnauthenticated(method: HTTPMethod, proxyPath: string, status: number) {
  if (method === "GET" && proxyPath === "/auth/me") {
    return status === 401
  }
  if (method === "GET" && (proxyPath === "/auth/summary/stream" || proxyPath === "/auth/admin/summary/stream")) {
    return status === 401 || status === 403
  }
  if (method === "POST" && (proxyPath === "/auth/refresh" || proxyPath === "/auth/logout")) {
    return status === 401 || status === 403
  }
  return false
}
function displayProxyPath(proxyPath: string) {
  const raw = String(proxyPath || "/")
  try {
    return decodeURIComponent(raw)
  } catch {
    return raw
  }
}
const processRef = globalThis as typeof globalThis & {
  process?: {
    stdout?: {isTTY?: boolean}
    stderr?: {isTTY?: boolean}
  }
}
const supportsAnsi = Boolean(processRef.process?.stdout?.isTTY && processRef.process?.stderr?.isTTY)
const color = supportsAnsi
  ? ({
      reset: "\x1b[0m",
      blue: "\x1b[34m",
      yellow: "\x1b[33m",
      red: "\x1b[31m",
      white: "\x1b[37m",
      gray: "\x1b[90m",
      italic: "\x1b[3m"
    } as const)
  : ({
      reset: "",
      blue: "",
      yellow: "",
      red: "",
      white: "",
      gray: "",
      italic: ""
    } as const)
function logProxyResult(method: HTTPMethod, proxyPath: string, status: number, message: string, details?: unknown) {
  if (status === 304 || status < 400) {
    return
  }
  const shortUrl = displayProxyPath(proxyPath)
  const statusColor = status >= 500 ? color.red : color.yellow
  const line =
    `[proxy] ${color.blue}${method}${color.reset} ${statusColor}-> ${status}${color.reset} ${color.white}${shortUrl}${color.reset} - ${color.gray}${message}${color.reset}` +
    (details ? ` ${color.italic}${color.gray}(${String(details)})${color.reset}` : "")
  if (status >= 500) {
    console.error(line)
    return
  }
  console.warn(line)
}
function splitSetCookieHeader(raw: string) {
  return raw
    .split(/,(?=\s*[!#$%&'*+\-.^_`|~0-9A-Za-z]+=)/g)
    .map(value => value.trim())
    .filter(Boolean)
}
function appendSetCookies(
  event: Parameters<typeof defineEventHandler>[0] extends (event: infer E) => any ? E : never,
  response: Response
) {
  const headers = response.headers as Headers & {
    getSetCookie?: () => string[]
    getAll?: (name: string) => string[]
    raw?: () => Record<string, string[]>
  }
  const values =
    (typeof headers.getSetCookie === "function" && headers.getSetCookie()) ||
    (typeof headers.getAll === "function" && headers.getAll("set-cookie")) ||
    headers.raw?.()["set-cookie"] ||
    splitSetCookieHeader(headers.get("set-cookie") || "")
  for (const value of values) {
    if (value) appendHeader(event, "set-cookie", value)
  }
}
function normalizeApiOrigin(value: string) {
  const trimmed = String(value || "")
    .trim()
    .replace(/\/+$/, "")
  if (!trimmed) return ""
  return trimmed.replace(/\/api(?:\/v\d+)?$/i, "")
}
function shouldForwardSetCookie(method: HTTPMethod, proxyPath: string, response: Response) {
  if (isExpectedUnauthenticated(method, proxyPath, response.status)) {
    return false
  }
  return true
}
function normalizeMediaProxyPath(proxyPath: string) {
  const rawProxyPath = String(proxyPath || "")
  const [rawPathname = "", search = ""] = rawProxyPath.split("?")
  const pathname = rawPathname
  if (!pathname.startsWith("/media/files/")) return proxyPath
  const prefix = "/media/files/"
  const rawKey = pathname.slice(prefix.length)
  const normalizedKey = rawKey
    .split("/")
    .map(segment => {
      try {
        return encodeURIComponent(decodeURIComponent(segment))
      } catch {
        return encodeURIComponent(segment)
      }
    })
    .join("/")
  return `${prefix}${normalizedKey}${search ? `?${search}` : ""}`
}
export default defineEventHandler(async event => {
  const {url} = event.node.req
  if (!url?.startsWith("/api/proxy")) return
  const config = useRuntimeConfig()
  const headers = getHeaders(event)
  const method = (event.method?.toUpperCase() || "GET") as HTTPMethod
  const proxyPath = url.replace(/^\/api\/proxy/, "") // e.g., /sets/some-uuid
  const normalizedProxyPath = normalizeMediaProxyPath(proxyPath)
  const apiBase = normalizeApiOrigin(String(config.apiInternalUrl || ""))
  const apiVersion =
    String(config.apiVersion || "v1")
      .trim()
      .replace(/^\/+/, "")
      .replace(/^api\//, "") || "v1"
  if (!apiBase) {
    throw createError({
      statusCode: 500,
      message: "API_INTERNAL_URL is not configured"
    })
  }
  const targetUrl = joinURL(apiBase, "api", apiVersion, normalizedProxyPath)
  const hasBody = ["POST", "PUT", "PATCH"].includes(method)
  const requestContentType = String(headers["content-type"] || "").toLowerCase()
  const isMultipartRequest = hasBody && requestContentType.includes("multipart/form-data")
  const isMediaFilePath = normalizedProxyPath.startsWith("/media/files/")
  const rawBody = hasBody ? await readRawBody(event, false) : undefined
  const requestHeaders = Object.fromEntries(
    Object.entries(headers).filter(([key]) => !shouldSkipRequestHeader(key))
  ) as Record<string, string>
  if (isMultipartRequest || isMediaFilePath) {
    const doRawRequest = async (url: string) =>
      await fetch(url, {
        method,
        body: hasBody ? (rawBody as any) : undefined,
        headers: requestHeaders,
        credentials: "include"
      })
    const response = await doRawRequest(targetUrl)
    setResponseStatus(event, response.status)
    if (shouldForwardSetCookie(method, proxyPath, response)) {
      appendSetCookies(event, response)
    }
    response.headers.forEach((value, key) => {
      if (!value) return
      const lower = key.toLowerCase()
      if (shouldSkipResponseHeader(lower)) {
        return
      }
      setHeader(event, key, value)
    })
    const bytes = new Uint8Array(await response.arrayBuffer())
    return bytes
  }
  try {
    const doRequest = (url: string) =>
      fetch(url, {
        method,
        body: hasBody ? (rawBody as any) : undefined,
        credentials: "include",
        headers: requestHeaders
      })
    const response = await doRequest(targetUrl)
    setResponseStatus(event, response.status)
    if (shouldForwardSetCookie(method, proxyPath, response)) {
      appendSetCookies(event, response)
    }
    setResponseStatus(event, response.status)
    const contentType = String(response.headers.get("content-type") || "").toLowerCase()
    const isEventStream = contentType.includes("text/event-stream")
    if (!isEventStream || response.ok) {
      if (shouldForwardSetCookie(method, proxyPath, response)) {
        appendSetCookies(event, response)
      }
    }
    response.headers.forEach((value, key) => {
      if (!value) return
      const lower = key.toLowerCase()
      if (shouldSkipResponseHeader(lower)) {
        return
      }
      setHeader(event, key, value)
    })
    if (isEventStream) {
      if (!response.ok || !response.body) {
        throw createError({
          statusCode: response.status,
          statusMessage: response.statusText || "Bad Gateway",
          message: response.statusText || "Bad Gateway"
        })
      }
      return sendStream(event, response.body)
    }
    if (contentType.includes("application/json")) {
      const raw = await response.text()
      let payload: any = null
      if (raw) {
        try {
          payload = JSON.parse(raw)
        } catch {
          throw createError({
            statusCode: 502,
            statusMessage: "Bad Gateway",
            message: `Upstream returned invalid JSON for ${proxyPath || "/"}`
          })
        }
      }
      if (!response.ok) {
        const message =
          payload && typeof payload === "object"
            ? (payload as any)?.message || response.statusText || "Unknown error"
            : response.statusText || "Unknown error"
        const details = payload && typeof payload === "object" ? (payload as any)?.details : undefined
        throw createError({
          statusCode: response.status,
          statusMessage: response.statusText || message,
          message,
          data: details
        })
      }
      return payload
    }
    const text = await response.text()
    if (!response.ok) {
      throw createError({
        statusCode: response.status,
        message: response.statusText || "Unknown error"
      })
    }
    return text
  } catch (error: any) {
    const status = error?.response?.status || error?.statusCode || error?.status || 500
    const message = error?.response?._data?.message || error?.message || "Unknown error"
    const details = error?.response?._data?.details
    // special case: keys/meta not found → return "empty" instead of error
    if (status === 404 && proxyPath.startsWith("/keys/meta")) {
      setResponseStatus(event, 200)
      return {notFound: true}
    }
    if (!isExpectedUnauthenticated(method, proxyPath, status)) {
      logProxyResult(method, proxyPath, status, message, details)
    }
    throw createError({
      statusCode: status,
      statusMessage: error?.statusMessage || error?.message || message,
      message,
      data: details
    })
  }
})
