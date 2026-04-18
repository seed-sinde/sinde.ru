export const asRecord = (value: unknown): Record<string, any> => {
  if (!value || typeof value !== "object" || Array.isArray(value)) return {}
  return {...(value as Record<string, any>)}
}
