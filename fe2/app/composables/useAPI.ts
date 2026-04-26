import {joinURL} from "ufo"

type FetchOptions = NonNullable<Parameters<typeof $fetch>[1]>
type BasicFetch = <T>(request: string, options?: FetchOptions) => Promise<T>
type HeaderMap = Record<string, string>

const SSR_HEADER_NAMES = [
  "cookie",
  "user-agent",
  "accept-language",
  "x-forwarded-for",
  "x-real-ip"
] as const

const toPath = (p: string) => (p.startsWith("/") ? p : `/${p}`)

const mergeHeaders = (src?: HeadersInit): HeaderMap => {
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
  for (const key of keys) {
    const v = reqHeaders[key]
    if (v) setHeaderIfMissing(headers, key, v)
  }
}

const createHttpError = async (res: Response) => {
  const text = await res.text().catch(() => "")
  const err = new Error(`HTTP ${res.status}`) as Error & {
    status?: number
    data?: string
  }
  err.status = res.status
  err.data = text
  return err
}

const parseNDJSON = async (res: Response, onLine: (line: unknown) => void) => {
  const reader = res.body?.getReader()
  if (!reader) throw new Error("Stream reader unavailable")

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

      if (line) onLine(JSON.parse(line))
      i = buf.indexOf("\n")
    }
  }

  const tail = buf.trim()
  if (tail) onLine(JSON.parse(tail))
}

/**
 * Transport layer only: proxy + SSR headers + credentials
 */
export const useAPI = () => {
  const basicFetch = $fetch as BasicFetch

  const nuxtApp = tryUseNuxtApp()
  const config = nuxtApp ? useRuntimeConfig() : null

  const ssrHeaders = import.meta.server && nuxtApp ? useRequestHeaders([...SSR_HEADER_NAMES]) : {}

  const ssrOrigin = import.meta.server && nuxtApp ? useRequestURL().origin : ""

  const json = async <T>(path: string, options: FetchOptions = {}): Promise<T> => {
    const p = toPath(path)
    const headers = mergeHeaders(options.headers)

    if (import.meta.server) {
      copyRequestHeadersIfMissing(headers, ssrHeaders, SSR_HEADER_NAMES)
    }

    return basicFetch<T>(`/api/proxy${p}`, {
      ...options,
      method: options.method ?? "GET",
      credentials: "include",
      headers
    })
  }

  const stream = async (
    path: string,
    onLine: (line: unknown) => void,
    options: RequestInit = {}
  ): Promise<void> => {
    const p = toPath(path)

    const base = import.meta.client ? "" : ssrOrigin || config?.public.baseURL || ""

    const url = joinURL(base, "/api/proxy", p)

    const headers = mergeHeaders(options.headers)

    setHeaderIfMissing(headers, "accept", "application/x-ndjson, application/json")

    if (import.meta.server) {
      copyRequestHeadersIfMissing(headers, ssrHeaders, SSR_HEADER_NAMES)
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
      const data = await res.json()
      if (Array.isArray(data)) data.forEach(onLine)
      else onLine(data)
      return
    }

    if (!ct.includes("application/x-ndjson")) {
      throw new Error(`Unexpected content-type: ${ct || "empty"}`)
    }

    await parseNDJSON(res, onLine)
  }

  return {json, stream}
}
