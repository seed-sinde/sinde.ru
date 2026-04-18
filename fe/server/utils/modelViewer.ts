import {createError, getQuery, getRequestURL, setHeader, type H3Event} from "h3"
const escapeHtml = (value: string) =>
  value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;")
const normalizeAssetPath = (value: unknown, origin: string) => {
  const raw = String(value || "").trim()
  if (!raw) return ""
  let url: URL
  try {
    url = new URL(raw, origin)
  } catch {
    return ""
  }
  if (url.origin !== origin) return ""
  if (!/^\/[a-z0-9/_-]+\.(glb|gltf)$/i.test(url.pathname)) return ""
  return `${url.pathname}${url.search}`
}
const normalizeBackground = (value: unknown, fallback: string) => {
  const raw = String(value || "").trim()
  if (!raw) return fallback
  if (/^#[0-9a-f]{3,8}$/i.test(raw)) return raw
  return fallback
}
const normalizeRotationPerSecond = (value: unknown) => {
  const raw = String(value || "").trim()
  if (!raw) return "18deg"
  if (!/^[0-9]+(?:\.[0-9]+)?deg$/i.test(raw)) return "18deg"
  return raw.toLowerCase()
}
const getQueryValue = (value: unknown, fallback: string) => {
  const raw = String(Array.isArray(value) ? value[0] : value || "").trim()
  return raw || fallback
}
export const renderModelViewerPage = (event: H3Event, fallbackTitle = "3D model") => {
  const query = getQuery(event)
  const origin = getRequestURL(event).origin
  const src = normalizeAssetPath(query.src, origin)
  const title = escapeHtml(getQueryValue(query.title, fallbackTitle))
  const background = normalizeBackground(query.background, "#18181b")
  const rotationPerSecond = escapeHtml(normalizeRotationPerSecond(query.rotationPerSecond))
  if (!src) {
    throw createError({
      statusCode: 400,
      statusMessage: "Bad Request",
      message: "Missing 3D model source"
    })
  }
  const escapedSrc = escapeHtml(src)
  setHeader(event, "content-type", "text/html; charset=utf-8")
  // setHeader(event, 'cache-control', 'public, max-age=3600')
  setHeader(event, "cache-control", "no-store")
  setHeader(event, "x-frame-options", "SAMEORIGIN")
  setHeader(event, "cross-origin-embedder-policy", "require-corp")
  setHeader(event, "cross-origin-resource-policy", "same-origin")
  setHeader(
    event,
    "content-security-policy",
    [
      "default-src 'self' data: blob:",
      "script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'",
      "style-src 'self' 'unsafe-inline'",
      "img-src 'self' data: blob:",
      "font-src 'self' data:",
      "connect-src 'self' data: blob:",
      "media-src 'self' data: blob:",
      "worker-src 'self' blob:",
      "frame-ancestors 'self'",
      "base-uri 'none'"
    ].join("; ")
  )
  return `<!doctype html>
<html lang="ru">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>${title}</title>
    <script type="module" src="/vendor/model-viewer.min.js"></script>
    <script defer src="/vendor/model-viewer-init.js"></script>
    <style>
      html, body {
        margin: 0;
        width: 100%;
        height: 100%;
        overflow: hidden;
        background: ${background};
      }
      body {
        padding: 0;
        outline: none;
      }
      .viewer-shell {
        position: relative;
        width: 100%;
        height: 100%;
        overflow: hidden;
        border: 0;
        border-radius: 0;
        background: ${background};
        outline: none;
      }
      model-viewer {
        width: 100%;
        height: 100%;
        outline: none;
        background: ${background};
      }
      .viewer-loader { position: absolute; inset: 0; z-index: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; color: #d4d4d8; transition: opacity 220ms ease, visibility 220ms ease; } .viewer-loader.is-hidden { opacity: 0; visibility: hidden; pointer-events: none; } .viewer-spinner { width: 40px; height: 40px; border-radius: 999px; border: 2px solid rgba(63, 63, 70, 0.9); border-top-color: #22d3ee; animation: spin 0.9s linear infinite; } .viewer-loader-text { font-size: 11px; letter-spacing: 0.12em; text-transform: uppercase; color: #71717a; }
    </style>
  </head>
  <body>
    <div class="viewer-shell">
      <div class="viewer-loader" id="viewer-loader">
        <div class="viewer-spinner"></div>
      </div>
      <model-viewer
        id="model-viewer"
        src="${escapedSrc}"
        camera-controls
        camera-target="auto auto auto"
        field-of-view="30deg"
        min-camera-orbit="auto auto auto"
        max-camera-orbit="auto auto auto"
        touch-action="pan-y"
        autoplay
        auto-rotate
        rotation-per-second="${rotationPerSecond}"
        interaction-prompt="none"
        exposure="1.1"
        shadow-intensity="0"
        disable-tap
      ></model-viewer>
    </div>
  </body>
</html>`
}
