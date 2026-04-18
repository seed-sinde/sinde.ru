import {isMfaTicketExpiredError, MFA_TICKET_EXPIRED_MESSAGE} from "~/utils/authErrors"
import {hasAuthSessionHint, syncAuthSessionHint} from "~/utils/authSessionHint"
import {emitAuthSyncEvent, subscribeAuthSyncEvents} from "~/utils/authSyncBus"
type AuthApiResult<T> = ApiResponseWithData<T>
let authSummaryRuntimeInitialized = false
let authSummaryUserStream: EventSource | null = null
let authSummaryAdminStream: EventSource | null = null
let authSummarySyncCleanup: (() => void) | null = null
let authLoadMeInFlight: Promise<AuthUser | null> | null = null
let authUserSummaryInFlight: Promise<UserSummary | null> | null = null
let authAdminSummaryInFlight: Promise<AdminSummary | null> | null = null
function closeSummaryStream(stream: EventSource | null) {
  stream?.close()
}
function toSharedAdminSummary(summary: UserSummary | null, current: AdminSummary | null): AdminSummary | null {
  const admin = summary?.admin
  if (!admin) return null
  return {
    users_total: Number(admin.users_total || current?.users_total || 0),
    pending_recipes_total: Number(admin.pending_recipes_total || current?.pending_recipes_total || 0),
    recipe_status_totals: admin.recipe_status_totals || current?.recipe_status_totals || {},
    new_users_since_last_login: Number(admin.new_users_since_last_login || 0),
    new_pending_recipes_since_last_login: Number(admin.new_pending_recipes_since_last_login || 0),
    last_login_before_current_session_at:
      summary?.last_login_before_current_session_at ?? current?.last_login_before_current_session_at ?? null,
    notifications_read_at: summary?.notifications_read_at ?? current?.notifications_read_at ?? null,
    notifications_since_at: summary?.notifications_since_at ?? current?.notifications_since_at ?? null,
    has_unread:
      Number(admin.new_users_since_last_login || 0) > 0 || Number(admin.new_pending_recipes_since_last_login || 0) > 0,
    checked_at: summary?.checked_at || current?.checked_at || ""
  }
}
export const useAuth = () => {
  const {json: useApiJson} = useAPI()
  const nuxtApp = useNuxtApp()
  const user = useState<AuthUser | null>("auth-user", () => null)
  const loaded = useState<boolean>("auth-loaded", () => false)
  const sessionRestoreHint = useState<boolean>("auth-session-restore-hint", () => {
    if (import.meta.server) {
      const cookieHeader = useRequestHeaders(["cookie"])?.cookie
      return hasAuthSessionHint(cookieHeader)
    }
    return hasAuthSessionHint()
  })
  const mfaTicket = useState<string | null>("auth-mfa-ticket", () => null)
  const mfaExpiresAt = useState<string | null>("auth-mfa-expires-at", () => null)
  const sharedUserSummary = useState<UserSummary | null>("auth-user-summary", () => null)
  const sharedUserSummaryLoading = useState<boolean>("auth-user-summary-loading", () => false)
  const sharedAdminSummary = useState<AdminSummary | null>("auth-admin-summary", () => null)
  const sharedAdminSummaryLoading = useState<boolean>("auth-admin-summary-loading", () => false)
  const isAuthenticated = computed(() => Boolean(user.value))
  const isAdmin = computed(() => Boolean(user.value?.roles?.includes("admin")))
  const canAttemptSessionRestore = computed(() => {
    if (import.meta.server) {
      const cookieHeader = useRequestHeaders(["cookie"])?.cookie
      return hasAuthSessionHint(cookieHeader)
    }
    return sessionRestoreHint.value || hasAuthSessionHint()
  })
  const applyUserSummary = (summary: UserSummary | null) => {
    sharedUserSummary.value = summary
    sharedAdminSummary.value = isAdmin.value ? toSharedAdminSummary(summary, sharedAdminSummary.value) : null
  }
  const applyAdminSummary = (summary: AdminSummary | null) => {
    sharedAdminSummary.value = isAdmin.value ? summary : null
  }
  const refreshSharedSummaries = async () => {
    if (!isAuthenticated.value) {
      applyUserSummary(null)
      applyAdminSummary(null)
      return null
    }
    const summary = await loadSharedUserSummary()
    if (isAdmin.value && !sharedAdminSummary.value) {
      await loadSharedAdminSummary()
    }
    return summary
  }
  const clearMfaState = () => {
    mfaTicket.value = null
    mfaExpiresAt.value = null
  }
  const syncSessionRestoreHint = (enabled: boolean) => {
    sessionRestoreHint.value = enabled
    syncAuthSessionHint(enabled)
  }
  const isMfaTicketExpired = () => {
    const expiresAt = Date.parse(String(mfaExpiresAt.value || ""))
    if (!Number.isFinite(expiresAt)) return true
    return Date.now() >= expiresAt
  }
  const setAnonymousState = () => {
    user.value = null
    loaded.value = true
    syncSessionRestoreHint(false)
    clearMfaState()
    applyUserSummary(null)
    applyAdminSummary(null)
    closeSummaryStream(authSummaryUserStream)
    closeSummaryStream(authSummaryAdminStream)
    authSummaryUserStream = null
    authSummaryAdminStream = null
    return null
  }
  const loadMe = async () => {
    if (authLoadMeInFlight) return await authLoadMeInFlight
    authLoadMeInFlight = (async () => {
      const restoreAllowed = canAttemptSessionRestore.value
      if (!restoreAllowed) {
        return setAnonymousState()
      }
      try {
        const res = await useApiJson<AuthApiResult<{user: AuthUser}>>("/auth/me", {
          auth: {
            requiresSession: true,
            allowAutoRefresh: false
          }
        })
        user.value = res.data.user
        loaded.value = true
        syncSessionRestoreHint(true)
        return user.value
      } catch (err: any) {
        const status = Number(err?.status || err?.statusCode || 0)
        if (status === 401 && canAttemptSessionRestore.value) {
          try {
            const res = await useApiJson<
              AuthApiResult<{
                user: AuthUser
                csrf_token?: string
              }>
            >("/auth/refresh", {
              method: "POST",
              auth: {
                requiresSession: true,
                allowAutoRefresh: false
              }
            })
            user.value = res.data.user
            loaded.value = true
            syncSessionRestoreHint(true)
            return user.value
          } catch {
            return setAnonymousState()
          }
        }
        return setAnonymousState()
      }
    })()
    try {
      return await authLoadMeInFlight
    } finally {
      authLoadMeInFlight = null
    }
  }
  const ensureLoaded = async () => {
    if (loaded.value) return user.value
    if (!canAttemptSessionRestore.value) {
      return setAnonymousState()
    }
    return await loadMe()
  }
  const register = async (input: {
    email: string
    password: string
    display_name?: string
    locale?: string
    timezone?: string
    captcha_token?: string
    website?: string
  }) => {
    return await useApiJson<
      AuthApiResult<{
        user: AuthUser
        verification_ttl: string
      }>
    >("/auth/register", {
      method: "POST",
      body: input
    })
  }
  const requestEmailVerification = async (email: string) => {
    return await useApiJson<AuthApiResult<{queued: boolean}>>("/auth/verify-email/request", {
      method: "POST",
      body: {email}
    })
  }
  const requestEmailChange = async (email: string) => {
    return await useApiJson<AuthApiResult<{queued: boolean}>>("/auth/email/change/request", {
      method: "POST",
      body: {email}
    })
  }
  const verifyEmail = async (token: string) => {
    return await useApiJson<
      AuthApiResult<{verified: boolean; action?: string; email?: string; session_hints?: boolean}>
    >("/auth/verify-email/confirm", {
      method: "POST",
      body: {token}
    })
  }
  const login = async (email: string, password: string) => {
    clearMfaState()
    const res = await useApiJson<
      AuthApiResult<{
        user?: AuthUser
        csrf_token?: string
        mfa_required: boolean
        mfa_ticket?: string
        mfa_expires_at?: string
        mfa_methods?: string[]
      }>
    >("/auth/login", {
      method: "POST",
      body: {email, password}
    })
    loaded.value = true
    if (res.data.mfa_required) {
      mfaTicket.value = String(res.data.mfa_ticket || "").trim() || null
      mfaExpiresAt.value = String(res.data.mfa_expires_at || "").trim() || null
    } else {
      syncSessionRestoreHint(true)
      clearMfaState()
    }
    if (res.data.user) user.value = res.data.user
    emitAuthSyncEvent("summary-refresh")
    return res
  }
  const completeMfa = async (code: string) => {
    const ticket = String(mfaTicket.value || "").trim()
    if (!ticket) throw new Error("Сессия 2FA не найдена.")
    if (isMfaTicketExpired()) {
      clearMfaState()
      throw new Error(MFA_TICKET_EXPIRED_MESSAGE)
    }
    try {
      const res = await useApiJson<
        AuthApiResult<{
          user: AuthUser
          csrf_token?: string
        }>
      >("/auth/login/2fa", {
        method: "POST",
        body: {ticket, code}
      })
      user.value = res.data.user
      loaded.value = true
      syncSessionRestoreHint(true)
      clearMfaState()
      emitAuthSyncEvent("summary-refresh")
      return res
    } catch (err) {
      if (isMfaTicketExpiredError(err)) {
        clearMfaState()
      }
      throw err
    }
  }
  const resetMfaTicket = () => {
    clearMfaState()
  }
  const refresh = async () => {
    const res = await useApiJson<
      AuthApiResult<{
        user: AuthUser
        csrf_token?: string
      }>
    >("/auth/refresh", {
      method: "POST",
      auth: {
        requiresSession: true,
        allowAutoRefresh: false
      }
    })
    user.value = res.data.user
    loaded.value = true
    syncSessionRestoreHint(true)
    emitAuthSyncEvent("summary-refresh")
    return res
  }
  const logout = async () => {
    try {
      await useApiJson<AuthApiResult<{logged_out: boolean}>>("/auth/logout", {
        method: "POST"
      })
    } finally {
      setAnonymousState()
      emitAuthSyncEvent("summary-refresh")
    }
  }
  const logoutAll = async () => {
    try {
      await useApiJson<AuthApiResult<{logged_out: boolean; all_devices: boolean}>>("/auth/logout-all", {
        method: "POST"
      })
    } finally {
      setAnonymousState()
      emitAuthSyncEvent("summary-refresh")
    }
  }
  const listSessions = async () => {
    return await useApiJson<AuthApiResult<{items: AuthSessionView[]}>>("/auth/sessions", {
      method: "GET"
    })
  }
  const listLoginAttempts = async () => {
    return await useApiJson<AuthApiResult<{items: AuthLoginAttemptView[]}>>("/auth/login-attempts", {
      method: "GET"
    })
  }
  const listSecurityEvents = async () => {
    return await useApiJson<AuthApiResult<{items: AuthSecurityEventView[]}>>("/auth/security-events", {
      method: "GET"
    })
  }
  const revokeSession = async (sessionId: string) => {
    return await useApiJson<AuthApiResult<{revoked: boolean}>>(`/auth/sessions/${sessionId}`, {
      method: "DELETE"
    })
  }
  const forgotPassword = async (email: string) => {
    return await useApiJson<AuthApiResult<{queued: boolean}>>("/auth/password/forgot", {
      method: "POST",
      body: {email}
    })
  }
  const resetPassword = async (token: string, newPassword: string) => {
    return await useApiJson<AuthApiResult<{password_reset: boolean}>>("/auth/password/reset", {
      method: "POST",
      body: {
        token,
        new_password: newPassword
      }
    })
  }
  const changePassword = async (currentPassword: string, newPassword: string) => {
    return await useApiJson<AuthApiResult<{password_changed: boolean}>>("/auth/password/change", {
      method: "POST",
      body: {
        current_password: currentPassword,
        new_password: newPassword
      }
    })
  }
  const updateProfile = async (input: {
    display_name?: string
    locale?: string
    timezone?: string
    profile?: Record<string, any>
    settings?: Record<string, any>
  }) => {
    const res = await useApiJson<AuthApiResult<{changes: AuthUserPatch}>>("/auth/me", {
      method: "PATCH",
      body: input
    })
    if (user.value) {
      user.value = applyAuthUserPatch(user.value, res.data.changes)
    }
    return res
  }
  const listTraitSets = async () => {
    const res = await useApiJson<
      AuthApiResult<{
        primary_trait_uuid?: string | null
        items: SavedTraitSetView[]
      }>
    >("/auth/traits/sets", {
      method: "GET"
    })
    if (user.value) {
      user.value = {
        ...user.value,
        primary_trait_uuid: res.data.primary_trait_uuid ?? null
      }
    }
    return res
  }
  const saveTraitSet = async (input: {set_uuid: string; name: string; description: string}) => {
    return await useApiJson<AuthApiResult<SavedTraitSetView>>("/auth/traits/sets", {
      method: "POST",
      body: input
    })
  }
  const updateTraitSet = async (
    savedSetId: string,
    input: {
      name: string
      description: string
    }
  ) => {
    return await useApiJson<AuthApiResult<SavedTraitSetView>>(`/auth/traits/sets/${savedSetId}`, {
      method: "PATCH",
      body: input
    })
  }
  const deleteTraitSet = async (savedSetId: string) => {
    return await useApiJson<AuthApiResult<{deleted: boolean}>>(`/auth/traits/sets/${savedSetId}`, {
      method: "DELETE"
    })
  }
  const setPrimaryTraitUuid = async (setUuid: string | null) => {
    const res = await useApiJson<AuthApiResult<{changes: AuthUserPatch}>>("/auth/traits/primary", {
      method: "POST",
      body: {set_uuid: setUuid || ""}
    })
    if (user.value) {
      user.value = applyAuthUserPatch(user.value, res.data.changes)
    }
    return res
  }
  const setupTwoFactor = async () => {
    return await useApiJson<
      AuthApiResult<{
        secret: string
        otpauth_url: string
        qr_data_url: string
        backup_codes?: string[]
      }>
    >("/auth/2fa/setup", {
      method: "POST"
    })
  }
  const enableTwoFactor = async (code: string) => {
    const res = await useApiJson<
      AuthApiResult<{
        enabled: boolean
        backup_codes: string[]
      }>
    >("/auth/2fa/enable", {
      method: "POST",
      body: {code}
    })
    await loadMe()
    return res
  }
  const disableTwoFactor = async (password: string, code: string) => {
    const res = await useApiJson<AuthApiResult<{disabled: boolean}>>("/auth/2fa/disable", {
      method: "POST",
      body: {password, code}
    })
    await loadMe()
    return res
  }
  const adminListUsers = async (params?: {
    q?: string
    status?: string
    role?: string
    limit?: number
    offset?: number
  }) => {
    const query = new URLSearchParams()
    if (params?.q) query.set("q", String(params.q))
    if (params?.status) query.set("status", String(params.status))
    if (params?.role) query.set("role", String(params.role))
    if (typeof params?.limit === "number") query.set("limit", String(params.limit))
    if (typeof params?.offset === "number") query.set("offset", String(params.offset))
    const suffix = query.toString()
    return await useApiJson<AuthApiResult<{items: AdminUserView[]; total: number; limit: number; offset: number}>>(
      `/auth/admin/users${suffix ? `?${suffix}` : ""}`,
      {method: "GET"}
    )
  }
  const publicUserProfile = async (userId: string) => {
    return await useApiJson<AuthApiResult<PublicUserProfileView>>(`/users/${userId}`, {
      method: "GET",
      auth: {
        allowAutoRefresh: true,
        requiresSession: false
      }
    })
  }
  const adminUserDetail = async (userId: string) => {
    return await useApiJson<AuthApiResult<AdminUserDetailView>>(`/auth/admin/users/${userId}`, {
      method: "GET"
    })
  }
  const userSummary = async () => {
    return await useApiJson<AuthApiResult<UserSummary>>("/auth/summary", {
      method: "GET"
    })
  }
  const markUserSummaryRead = async () => {
    const res = await useApiJson<AuthApiResult<{read: boolean}>>("/auth/summary/read", {
      method: "POST"
    })
    await loadSharedUserSummary()
    emitAuthSyncEvent("summary-refresh")
    return res
  }
  const adminSummary = async () => {
    return await useApiJson<AuthApiResult<AdminSummary>>("/auth/admin/summary", {
      method: "GET"
    })
  }
  const adminMarkSummaryRead = async () => {
    const res = await useApiJson<AuthApiResult<{read: boolean}>>("/auth/admin/summary/read", {
      method: "POST"
    })
    await loadSharedAdminSummary()
    emitAuthSyncEvent("summary-refresh")
    return res
  }
  const adminSetUserRole = async (userId: string, role: "admin" | "user") => {
    return await useApiJson<AuthApiResult<{user: AuthUser}>>(`/auth/admin/users/${userId}/role`, {
      method: "PATCH",
      body: {role}
    })
  }
  const adminBlockUser = async (userId: string, reason = "") => {
    return await useApiJson<AuthApiResult<{blocked: boolean}>>(`/auth/admin/users/${userId}/block`, {
      method: "POST",
      body: {reason}
    })
  }
  const adminUnblockUser = async (userId: string) => {
    return await useApiJson<AuthApiResult<{unblocked: boolean}>>(`/auth/admin/users/${userId}/unblock`, {
      method: "POST"
    })
  }
  const adminForceLogoutUser = async (userId: string) => {
    return await useApiJson<AuthApiResult<{forced_logout: boolean}>>(`/auth/admin/users/${userId}/force-logout`, {
      method: "POST"
    })
  }
  const adminDeleteUser = async (userId: string) => {
    return await useApiJson<AuthApiResult<{deleted: boolean}>>(`/auth/admin/users/${userId}`, {
      method: "DELETE"
    })
  }
  const adminSearchKeys = async (q: string, limit = 30) => {
    const query = new URLSearchParams()
    if (q) query.set("q", q)
    query.set("limit", String(limit))
    return await useApiJson<AuthApiResult<{items: AdminTraitKeySearchItem[]; query: string; limit: number}>>(
      `/auth/admin/keys/search?${query.toString()}`,
      {method: "GET"}
    )
  }
  const adminTraitsSetsAnalysis = async () => {
    return await useApiJson<AuthApiResult<AdminTraitsSetsAnalysis>>("/auth/admin/analysis/traits-sets", {
      method: "GET"
    })
  }
  const loadSharedUserSummary = async () => {
    if (!isAuthenticated.value) {
      applyUserSummary(null)
      return null
    }
    if (sharedUserSummaryLoading.value) return sharedUserSummary.value
    if (authUserSummaryInFlight) return await authUserSummaryInFlight
    sharedUserSummaryLoading.value = true
    authUserSummaryInFlight = (async () => {
      try {
        const res = await userSummary()
        applyUserSummary(res.data || null)
        return sharedUserSummary.value
      } catch {
        applyUserSummary(null)
        return null
      } finally {
        sharedUserSummaryLoading.value = false
        authUserSummaryInFlight = null
      }
    })()
    return await authUserSummaryInFlight
  }
  const loadSharedAdminSummary = async () => {
    if (!isAdmin.value) {
      applyAdminSummary(null)
      return null
    }
    if (sharedAdminSummaryLoading.value) return sharedAdminSummary.value
    if (authAdminSummaryInFlight) return await authAdminSummaryInFlight
    sharedAdminSummaryLoading.value = true
    authAdminSummaryInFlight = (async () => {
      try {
        const res = await adminSummary()
        applyAdminSummary(res.data || null)
        return sharedAdminSummary.value
      } catch {
        applyAdminSummary(null)
        return null
      } finally {
        sharedAdminSummaryLoading.value = false
        authAdminSummaryInFlight = null
      }
    })()
    return await authAdminSummaryInFlight
  }
  const ensureSummaryRealtime = () => {
    if (!import.meta.client || authSummaryRuntimeInitialized) return
    authSummaryRuntimeInitialized = true
    const openSummaryStream = <T>(path: string, apply: (payload: T | null) => void) => {
      const stream = new EventSource(`/api/proxy${path}`, {withCredentials: true})
      stream.addEventListener("summary", event => {
        const message = event as MessageEvent<string>
        try {
          const payload = JSON.parse(String(message.data || "null")) as T | null
          apply(payload)
        } catch {
          // Ignore malformed SSE payloads and keep the last good snapshot.
        }
      })
      stream.onerror = () => {
        // Native EventSource reconnect is enough here; avoid noisy client logs.
      }
      return stream
    }
    const ensureUserSummaryStream = () => {
      if (authSummaryUserStream) return
      authSummaryUserStream = openSummaryStream<UserSummary>("/auth/summary/stream", payload => {
        applyUserSummary(payload)
      })
    }
    const ensureAdminSummaryStream = () => {
      if (!isAdmin.value || authSummaryAdminStream) return
      authSummaryAdminStream = openSummaryStream<AdminSummary>("/auth/admin/summary/stream", payload => {
        applyAdminSummary(payload)
      })
    }
    authSummarySyncCleanup?.()
    authSummarySyncCleanup = subscribeAuthSyncEvents(event => {
      if (event.type !== "summary-refresh") return
      void nuxtApp.runWithContext(async () => {
        await refreshSharedSummaries()
      })
    })
    watch(
      [loaded, isAuthenticated, isAdmin],
      async ([isLoaded, authed, admin]) => {
        if (!isLoaded) return
        if (!authed) {
          applyUserSummary(null)
          applyAdminSummary(null)
          closeSummaryStream(authSummaryUserStream)
          closeSummaryStream(authSummaryAdminStream)
          authSummaryUserStream = null
          authSummaryAdminStream = null
          return
        }
        if (!sharedUserSummary.value) {
          await loadSharedUserSummary()
        }
        ensureUserSummaryStream()
        if (!admin) {
          closeSummaryStream(authSummaryAdminStream)
          authSummaryAdminStream = null
          applyAdminSummary(null)
          return
        }
        if (!sharedAdminSummary.value) {
          await loadSharedAdminSummary()
        }
        ensureAdminSummaryStream()
      },
      {immediate: true}
    )
  }
  ensureSummaryRealtime()
  return {
    user,
    loaded,
    isAuthenticated,
    isAdmin,
    hasAuthSessionHint,
    canAttemptSessionRestore,
    mfaTicket,
    mfaExpiresAt,
    resetMfaTicket,
    ensureLoaded,
    loadMe,
    register,
    requestEmailVerification,
    requestEmailChange,
    verifyEmail,
    login,
    completeMfa,
    refresh,
    logout,
    logoutAll,
    listSessions,
    listLoginAttempts,
    listSecurityEvents,
    revokeSession,
    forgotPassword,
    resetPassword,
    changePassword,
    updateProfile,
    listTraitSets,
    saveTraitSet,
    updateTraitSet,
    deleteTraitSet,
    setAnonymousState,
    setPrimaryTraitUuid,
    setupTwoFactor,
    enableTwoFactor,
    disableTwoFactor,
    adminListUsers,
    publicUserProfile,
    adminUserDetail,
    userSummary,
    markUserSummaryRead,
    adminSummary,
    adminMarkSummaryRead,
    adminSetUserRole,
    adminBlockUser,
    adminUnblockUser,
    adminForceLogoutUser,
    adminDeleteUser,
    adminSearchKeys,
    adminTraitsSetsAnalysis,
    adminListKitchenModerationRecipes,
    adminModerateKitchenRecipe,
    adminChangeKitchenRecipeOwner,
    sharedUserSummary,
    sharedAdminSummary,
    loadSharedUserSummary,
    loadSharedAdminSummary,
    refreshSharedSummaries,
    ensureSummaryRealtime
  }
}
