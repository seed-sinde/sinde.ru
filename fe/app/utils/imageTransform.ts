const fileNameBase = (name: string) => {
  const raw = String(name || '').trim()
  const idx = raw.lastIndexOf('.')
  if (idx <= 0) return raw || 'image'
  return raw.slice(0, idx)
}
const fileExtFromMime = (mime: string) => {
  if (mime === 'image/webp') return 'webp'
  return 'jpg'
}
const chooseOutputMime = (raw: string) => {
  const safe = String(raw || '')
    .trim()
    .toLowerCase()
  if (safe.startsWith('image/')) return 'image/webp'
  return 'image/webp'
}
const loadImage = async (file: File) => {
  const objectUrl = URL.createObjectURL(file)
  try {
    const image = await new Promise<HTMLImageElement>((resolve, reject) => {
      const el = new Image()
      el.onload = () => resolve(el)
      el.onerror = () => reject(new Error('Не удалось прочитать изображение.'))
      el.src = objectUrl
    })
    return image
  } finally {
    URL.revokeObjectURL(objectUrl)
  }
}
export const resizeImageToSquareFile = async (
  file: File,
  size: number,
  options?: {
    outputMime?: string
    quality?: number
    suffix?: string
  }
) => {
  const safeSize = Math.max(1, Math.round(size))
  const image = await loadImage(file)
  const sourceSize = Math.min(image.naturalWidth || image.width, image.naturalHeight || image.height)
  if (!(sourceSize > 0)) {
    throw new Error('Не удалось определить размеры изображения.')
  }
  const sx = Math.round(((image.naturalWidth || image.width) - sourceSize) / 2)
  const sy = Math.round(((image.naturalHeight || image.height) - sourceSize) / 2)
  const canvas = document.createElement('canvas')
  canvas.width = safeSize
  canvas.height = safeSize
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('Canvas недоступен.')
  }
  ctx.imageSmoothingEnabled = true
  ctx.imageSmoothingQuality = 'high'
  ctx.drawImage(image, sx, sy, sourceSize, sourceSize, 0, 0, safeSize, safeSize)
  const outputMime = chooseOutputMime(options?.outputMime || file.type)
  const quality = typeof options?.quality === 'number' ? options.quality : 0.92
  const blob = await new Promise<Blob | null>((resolve) => {
    canvas.toBlob(resolve, outputMime, quality)
  })
  if (!blob) {
    throw new Error('Не удалось сформировать файл изображения.')
  }
  const ext = fileExtFromMime(outputMime)
  const suffix = String(options?.suffix || `${safeSize}`).trim()
  const base = fileNameBase(file.name)
  const filename = `${base}-${suffix}.${ext}`
  return new File([blob], filename, {
    type: outputMime,
    lastModified: Date.now()
  })
}
