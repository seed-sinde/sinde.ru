const palette = [
  '#22c55e',
  '#f59e0b',
  '#ec4899',
  '#38bdf8',
  '#ef4444',
  '#84cc16',
  '#14b8a6',
  '#c084fc',
  '#f97316',
  '#94a3b8'
] as const
const hashString = (value: string) => {
  let hash = 0
  for (let i = 0; i < value.length; i += 1) {
    hash = (hash << 5) - hash + value.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash)
}
const hexToRgb = (hex: string) => {
  const clean = hex.replace('#', '')
  if (clean.length !== 6) return { r: 148, g: 163, b: 184 }
  const num = Number.parseInt(clean, 16)
  return {
    r: (num >> 16) & 255,
    g: (num >> 8) & 255,
    b: num & 255
  }
}
export const kitchenCategoryColor = (category: string) => {
  const key = String(category || '')
    .trim()
    .toLowerCase()
  if (!key) return '#94a3b8'
  return palette[hashString(key) % palette.length] || '#94a3b8'
}
export const kitchenColorWithAlpha = (hex: string, alpha: number) => {
  const { r, g, b } = hexToRgb(hex)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}
