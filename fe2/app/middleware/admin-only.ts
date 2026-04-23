const normalizeAuthNextPath = (value?: string | null) => {
  const raw = String(value || "").trim()
  if (!raw || !raw.startsWith("/") || raw.startsWith("//")) return ""
  return raw
}

const buildLoginPath = (next?: string | null) => {
  const normalizedNext = normalizeAuthNextPath(next)
  return normalizedNext ? `/auth/login?next=${encodeURIComponent(normalizedNext)}` : "/auth/login"
}

export default defineNuxtRouteMiddleware(async to => {
  const {ensureLoaded, isAuthenticated, isAdmin} = useAuth()
  await ensureLoaded()
  if (!isAuthenticated.value) return navigateTo(buildLoginPath(to.fullPath))
  if (!isAdmin.value) return navigateTo("/auth/account/profile")
})
