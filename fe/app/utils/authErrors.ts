type ApiErrorMessageSource = {
  data?: {message?: unknown} | string
  message?: unknown
}
export const extractApiErrorMessage = (err: unknown, fallback = "") => {
  const source = err as ApiErrorMessageSource
  const candidate =
    source?.data && typeof source.data === "object"
      ? source.data.message || source.message || fallback
      : typeof source?.data === "string"
        ? source.data || source.message || fallback
        : source?.message || fallback
  return String(candidate || fallback).trim()
}
export const MFA_TICKET_EXPIRED_MESSAGE = "Время подтверждения истекло. Войдите заново."
const normalizeErrorMessage = (message: string) => message.trim().toLowerCase()
const isCsrfErrorMessage = (message: string) => {
  const normalized = normalizeErrorMessage(message)
  return normalized.includes("csrf") || normalized.includes("проверка csrf") || normalized.includes("сессия защиты")
}
const isInvalidTwoFactorCodeMessage = (message: string) => {
  const normalized = normalizeErrorMessage(message)
  return (
    normalized.includes("invalid 2fa code") ||
    normalized.includes("неверный код 2fa") ||
    normalized.includes("неверный код из приложения") ||
    normalized.includes("неверный код приложения") ||
    normalized.includes("неверный код сброса")
  )
}
export const isMfaTicketExpiredMessage = (message: string) => {
  const normalized = normalizeErrorMessage(message)
  return (
    normalized === "mfa ticket expired" ||
    normalized === "mfa session expired" ||
    normalized === "время подтверждения истекло" ||
    normalized === "время подтверждения истекло. войдите заново."
  )
}
export const isMfaTicketExpiredError = (err: unknown) => {
  return isMfaTicketExpiredMessage(extractApiErrorMessage(err, ""))
}
export const getTwoFactorErrorMessage = (err: unknown, kind: "totp" | "backup", fallback: string) => {
  const message = extractApiErrorMessage(err, fallback)
  if (isMfaTicketExpiredMessage(message)) {
    return MFA_TICKET_EXPIRED_MESSAGE
  }
  if (isCsrfErrorMessage(message)) {
    return "Сессия защиты устарела. Обновите страницу и попробуйте снова."
  }
  if (isInvalidTwoFactorCodeMessage(message)) {
    return kind === "backup"
      ? "Неверный код сброса. Попробуйте ещё раз."
      : "Неверный код из приложения. Попробуйте ещё раз."
  }
  return message || fallback
}
