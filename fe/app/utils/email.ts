const EMAIL_MAX_LENGTH = 254
const EMAIL_LOCAL_MAX_LENGTH = 64
const EMAIL_PATTERN =
  /^(?<local>[a-z0-9.!#$%&'*+/=?^_`{|}~-]+)@(?<domain>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+)$/i
export const normalizeEmail = (value: string) =>
  String(value ?? "")
    .trim()
    .toLowerCase()
export const isValidEmail = (value: string) => {
  const email = normalizeEmail(value)
  if (!email || email.length > EMAIL_MAX_LENGTH || email.includes("..")) return false
  const match = EMAIL_PATTERN.exec(email)
  if (!match?.groups) return false
  const localPart = String(match.groups.local || "")
  return Boolean(localPart) && localPart.length <= EMAIL_LOCAL_MAX_LENGTH
}
export const isInvalidEmailInput = (value: string) => {
  const raw = String(value ?? "").trim()
  return Boolean(raw) && !isValidEmail(raw)
}
