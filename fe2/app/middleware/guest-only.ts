export default defineNuxtRouteMiddleware(async () => {
  const {ensureLoaded, isAuthenticated} = useAuth()
  await ensureLoaded()
  if (isAuthenticated.value) return navigateTo("/auth/account/profile")
})
