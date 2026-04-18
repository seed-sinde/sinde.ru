export const normalizeMediaKey = (value?: string) => {
  const key = String(value || "")
    .trim()
    .replace(/^\/+/, "")
  if (!key) return ""
  return key
}
export const buildMediaFileUrl = (storageKey?: string) => {
  const key = normalizeMediaKey(storageKey)
  if (!key) return ""
  const encodedPath = key
    .split("/")
    .map(segment => {
      try {
        return encodeURIComponent(decodeURIComponent(segment))
      } catch {
        return encodeURIComponent(segment)
      }
    })
    .join("/")
  return `/api/proxy/media/files/${encodedPath}`
}
