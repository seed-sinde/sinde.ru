import { createError, getQuery, getRequestURL, setHeader, type H3Event } from 'h3'
const escapeHtml = (value: string) => value
  .replaceAll('&', '&amp;')
  .replaceAll('<', '&lt;')
  .replaceAll('>', '&gt;')
  .replaceAll('"', '&quot;')
  .replaceAll("'", '&#39;')
const normalizeAssetPath = (value: unknown, origin: string) => {
  const raw = String(value || '').trim()
  if (!raw) return ''
  let url: URL
  try {
    url = new URL(raw, origin)
  } catch {
    return ''
  }
  if (url.origin !== origin) return ''
  if (!/^\/[a-z0-9/_-]+\.(glb|gltf)$/i.test(url.pathname)) return ''
  return `${url.pathname}${url.search}`
}
const normalizeRotationPerSecond = (value: unknown) => {
  const raw = String(value || '').trim()
  if (!raw) return '18deg'
  if (!/^[0-9]+(?:\.[0-9]+)?deg$/i.test(raw)) return '18deg'
  return raw.toLowerCase()
}
const getQueryValue = (value: unknown, fallback: string) => {
  const raw = String(Array.isArray(value) ? value[0] : value || '').trim()
  return raw || fallback
}
export const renderModelViewerPage = (event: H3Event, fallbackTitle = '3D model') => {
  const query = getQuery(event)
  const origin = getRequestURL(event).origin
  const src = normalizeAssetPath(query.src, origin)
  const title = escapeHtml(getQueryValue(query.title, fallbackTitle))
  const rotationPerSecond = escapeHtml(normalizeRotationPerSecond(query.rotationPerSecond))
  if (!src) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Bad Request',
      message: 'Missing 3D model source'
    })
  }
  const escapedSrc = escapeHtml(src)
  setHeader(event, 'content-type', 'text/html; charset=utf-8')
  setHeader(event, 'cache-control', 'public, max-age=3600')
  setHeader(event, 'x-frame-options', 'SAMEORIGIN')
  setHeader(event, 'cross-origin-embedder-policy', 'require-corp')
  setHeader(event, 'cross-origin-resource-policy', 'same-origin')
  setHeader(event, 'content-security-policy', [
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
  ].join('; '))
  return `<!doctype html>
<html lang="ru">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>${title}</title>
    <script type="module" src="/vendor/model-viewer.min.js"></script>
    <style>
      :root {
        color-scheme: dark;
        font-family: ui-sans-serif, system-ui, sans-serif;
        background:
          radial-gradient(circle at top, rgba(56, 189, 248, 0.18), transparent 42%),
          linear-gradient(180deg, #111827 0%, #020617 100%);
      }
      * {
        box-sizing: border-box;
      }
      html, body {
        margin: 0;
        width: 100%;
        height: 100%;
        overflow: hidden;
        background: transparent;
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
        background: rgba(9, 9, 11, 0.92);
        outline: none;
      }
      .viewer-label {
        position: absolute;
        top: 12px;
        left: 12px;
        z-index: 2;
        padding: 6px 10px;
        border: 1px solid rgba(63, 63, 70, 0.95);
        border-radius: 999px;
        background: rgba(9, 9, 11, 0.72);
        color: #d4d4d8;
        font-size: 11px;
        letter-spacing: 0.08em;
        text-transform: uppercase;
      }
      .viewer-loader {
        position: absolute;
        inset: 0;
        z-index: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 12px;
        background:
          radial-gradient(circle at 30% 20%, rgba(34, 197, 94, 0.10), transparent 24%),
          radial-gradient(circle at 70% 78%, rgba(56, 189, 248, 0.12), transparent 26%),
          linear-gradient(180deg, rgba(24, 24, 27, 0.98) 0%, rgba(9, 9, 11, 1) 100%);
        color: #d4d4d8;
        transition: opacity 220ms ease, visibility 220ms ease;
      }
      .viewer-loader.is-hidden {
        opacity: 0;
        visibility: hidden;
        pointer-events: none;
      }
      .viewer-spinner {
        width: 40px;
        height: 40px;
        border-radius: 999px;
        border: 2px solid rgba(63, 63, 70, 0.9);
        border-top-color: #22d3ee;
        animation: spin 0.9s linear infinite;
      }
      .viewer-loader-text {
        font-size: 11px;
        letter-spacing: 0.12em;
        text-transform: uppercase;
        color: #71717a;
      }
      model-viewer {
        width: 100%;
        height: 100%;
        outline: none;
        background:
          radial-gradient(circle at 30% 20%, rgba(34, 197, 94, 0.15), transparent 24%),
          radial-gradient(circle at 70% 78%, rgba(56, 189, 248, 0.16), transparent 26%),
          linear-gradient(180deg, rgba(24, 24, 27, 0.96) 0%, rgba(9, 9, 11, 0.98) 100%);
      }
      model-viewer:focus,
      model-viewer:focus-visible {
        outline: none;
      }
      @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
      }
    </style>
  </head>
  <body>
    <div class="viewer-shell">
      <div class="viewer-label">${title}</div>
      <div class="viewer-loader" id="viewer-loader">
        <div class="viewer-spinner"></div>
        <div class="viewer-loader-text">Загрузка модели</div>
      </div>
      <model-viewer
        id="model-viewer"
        src="${escapedSrc}"
        camera-controls
        camera-target="0m 0m 0m"
        touch-action="pan-y"
        autoplay
        auto-rotate
        rotation-per-second="${rotationPerSecond}"
        interaction-prompt="none"
        exposure="1.1"
        disable-tap
      ></model-viewer>
      <script>
        const modelViewer = document.getElementById('model-viewer');
        const loader = document.getElementById('viewer-loader');
        const hideLoader = () => {
          if (loader) loader.classList.add('is-hidden');
        };
        if (modelViewer) {
          modelViewer.addEventListener('load', hideLoader, { once: true });
          modelViewer.addEventListener('error', hideLoader, { once: true });
          window.addEventListener('keydown', (event) => {
            if (event.key !== 'Escape') return;
            window.parent.postMessage({ type: 'lab-3d-viewer-close' }, window.location.origin);
          });
        } else {
          hideLoader();
        }
      </script>
    </div>
  </body>
</html>`
}
