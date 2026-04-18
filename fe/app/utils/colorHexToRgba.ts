export const colorHexToRgba = (hexColor: string, alpha: number) => {
  const normalized = String(hexColor || "").replace("#", "")
  if (!/^[0-9a-f]{6}$/i.test(normalized)) {
    return `rgba(82, 82, 91, ${alpha})`
  }
  const [red, green, blue] = [0, 2, 4].map(offset => Number.parseInt(normalized.slice(offset, offset + 2), 16))
  return `rgba(${red}, ${green}, ${blue}, ${alpha})`
}
