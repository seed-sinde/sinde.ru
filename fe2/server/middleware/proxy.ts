import type {H3Event, HTTPMethod} from "h3"
import {
  defineEventHandler,
  getHeaders,
  readRawBody,
  setResponseStatus,
  setHeader,
  appendHeader,
  sendStream,
  createError
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

const shouldSkipRequestHeader = (name: string) => {
  const lower = name.toLowerCase()
  return lower === "host" || lower === "content-length" || HOP_BY_HOP_HEADERS.has(lower)
}

const shouldSkipResponseHeader = (name: string) => {
  const lower = name.toLowerCase()
  return lower === "content-length" || HOP_BY_HOP_HEADERS.has(lower)
}

const appendSetCookies = (event: H3Event, response: Response) => {
  const raw = response.headers.get("set-cookie")
  if (!raw) return

  raw.split(/,(?=\s*\w+=)/g).forEach(v => {
    if (v) appendHeader(event, "set-cookie", v.trim())
  })
}

export default defineEventHandler(async (event: H3Event) => {
  const url = event.node.req.url
  if (!url?.startsWith("/api/proxy")) return

  const config = useRuntimeConfig()

  const apiBase = String(config.apiInternalUrl || "").replace(/\/+$/, "")
  const apiVersion = String(config.apiVersion || "v1").replace(/^\/+/, "")

  if (!apiBase) {
    throw createError({statusCode: 500, message: "API base not configured"})
  }

  const proxyPath = url.replace(/^\/api\/proxy/, "")
  const targetUrl = joinURL(apiBase, "api", apiVersion, proxyPath)

  const method = (event.method || "GET").toUpperCase() as HTTPMethod
  const headers = getHeaders(event)

  const requestHeaders = Object.fromEntries(
    Object.entries(headers).filter(([k, v]) => !shouldSkipRequestHeader(k) && typeof v === "string")
  ) as Record<string, string>

  const hasBody = ["POST", "PUT", "PATCH"].includes(method)
  const body = hasBody ? await readRawBody(event, false) : undefined

  let res: Response

  try {
    res = await fetch(targetUrl, {
      method,
      headers: requestHeaders,
      body: body ? new Blob([new Uint8Array(body)]) : undefined,
      credentials: "include"
    })
    console.log(method, targetUrl, res.status, res.statusText)
  } catch (e) {
    throw createError({
      statusCode: 502,
      message: "Backend unavailable"
    })
  }

  setResponseStatus(event, res.status)

  appendSetCookies(event, res)

  res.headers.forEach((value, key) => {
    if (!value || shouldSkipResponseHeader(key)) return
    setHeader(event, key, value)
  })

  const contentType = String(res.headers.get("content-type") || "")

  if (contentType.includes("text/event-stream")) {
    if (!res.body) {
      throw createError({statusCode: 502, message: "Empty stream"})
    }
    return sendStream(event, res.body)
  }

  return new Uint8Array(await res.arrayBuffer())
})
