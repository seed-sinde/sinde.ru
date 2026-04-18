import {buildLoginPath} from "~/utils/authNavigation"
export default defineNuxtRouteMiddleware(async to => {
  const {ensureLoaded, isAuthenticated, isAdmin} = useAuth()
  await ensureLoaded()
  if (!isAuthenticated.value) {
    return navigateTo(buildLoginPath(to.fullPath))
  }
  if (!isAdmin.value) {
    return navigateTo("/auth/account/profile")
  }
})
