window.addEventListener('DOMContentLoaded', async () => {
  const modelViewer = document.getElementById('model-viewer')
  const loader = document.getElementById('viewer-loader')

  const hideLoader = () => {
    if (loader) loader.classList.add('is-hidden')
  }

  if (!modelViewer) {
    hideLoader()
    return
  }

  try {
    await customElements.whenDefined('model-viewer')
  } catch {}

  modelViewer.addEventListener('load', hideLoader, { once: true })
  modelViewer.addEventListener('error', hideLoader, { once: true })

  window.addEventListener('keydown', event => {
    if (event.key !== 'Escape') return
    window.parent.postMessage({ type: 'lab-3d-viewer-close' }, window.location.origin)
  })

  setTimeout(hideLoader, 2000)
})
