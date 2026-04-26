type ApiResponse = {
  ok: boolean
  message?: string
  details?: string
}

type ApiResponseWithData<T> = ApiResponse & {
  data: T
}

type AuthUser = {
  user_id: string
  email: string
  status: string
  email_verified_at?: string | null
  display_name: string
  locale: string
  timezone: string
  roles: string[]
  is_two_factor_enabled: boolean
  last_login_at?: string | null
  blocked_reason?: string
  blocked_at?: string | null
  created_at: string
  profile: Record<string, unknown>
  settings: Record<string, unknown>
}

type AuthUserPatch = Partial<
  Pick<AuthUser, "display_name" | "locale" | "timezone" | "profile" | "settings">
>

type ApiError = Error & {
  status?: number
  statusCode?: number
  data?: {
    message?: string
  }
}

const AUTH_SERVER_COOKIE_NAMES = [
  "__Host-access_token",
  "access_token",
  "__Host-refresh_token",
  "refresh_token"
] as const
const AUTH_CLIENT_COOKIE_NAMES = ["__Host-auth_session_hint", "auth_session_hint"] as const
const AUTH_SESSION_RUNTIME_KEY = "__sindeAuthSessionHint__"

const hasNamedCookie = (cookieHeader: string | null | undefined, names: readonly string[]) => {
  const raw = String(cookieHeader || "").trim()
  if (!raw) return false

  for (const chunk of raw.split(";")) {
    const [namePart, ...valueParts] = chunk.split("=")
    const name = String(namePart || "").trim()
    if (name && names.includes(name) && valueParts.join("=").trim()) return true
  }

  return false
}

export const hasAuthSessionHint = (cookieHeader?: string | null) => {
  if (import.meta.server) return hasNamedCookie(cookieHeader, AUTH_SERVER_COOKIE_NAMES)

  const scope = globalThis as typeof globalThis & {
    [AUTH_SESSION_RUNTIME_KEY]?: boolean
  }
  return (
    Boolean(scope[AUTH_SESSION_RUNTIME_KEY]) ||
    hasNamedCookie(document.cookie, AUTH_CLIENT_COOKIE_NAMES)
  )
}

const syncAuthSessionHint = (enabled: boolean) => {
  if (!import.meta.client) return

  const scope = globalThis as typeof globalThis & {
    [AUTH_SESSION_RUNTIME_KEY]?: boolean
  }
  scope[AUTH_SESSION_RUNTIME_KEY] = enabled

  const secure = window.location.protocol === "https:" ? "; Secure" : ""
  const preferredName =
    window.location.protocol === "https:" ? "__Host-auth_session_hint" : "auth_session_hint"

  if (enabled) {
    document.cookie = `${preferredName}=1; Path=/; SameSite=Strict${secure}`
    return
  }

  for (const name of AUTH_CLIENT_COOKIE_NAMES) {
    document.cookie = `${name}=; Max-Age=0; Path=/; SameSite=Strict${secure}`
  }
}

const errorStatus = (err: unknown) =>
  Number((err as ApiError)?.status || (err as ApiError)?.statusCode || 0)

const applyAuthUserPatch = (user: AuthUser, patch: AuthUserPatch): AuthUser => ({
  ...user,
  ...patch,
  profile: patch.profile ?? user.profile,
  settings: patch.settings ?? user.settings
})

let loadMeInFlight: Promise<AuthUser | null> | null = null

export const useAuth = () => {
  const {json} = useAPI()
  const user = useState<AuthUser | null>("auth-user", () => null)
  const loaded = useState("auth-loaded", () => false)
  const sessionRestoreHint = useState("auth-session-restore-hint", () =>
    import.meta.server
      ? hasAuthSessionHint(useRequestHeaders(["cookie"])?.cookie)
      : hasAuthSessionHint()
  )
  const mfaTicket = useState<string | null>("auth-mfa-ticket", () => null)
  const mfaExpiresAt = useState<string | null>("auth-mfa-expires-at", () => null)

  const isAuthenticated = computed(() => Boolean(user.value))
  const isAdmin = computed(() => Boolean(user.value?.roles?.includes("admin")))
  const canAttemptSessionRestore = computed(() =>
    import.meta.server
      ? hasAuthSessionHint(useRequestHeaders(["cookie"])?.cookie)
      : sessionRestoreHint.value || hasAuthSessionHint()
  )

  const clearMfa = () => {
    mfaTicket.value = null
    mfaExpiresAt.value = null
  }

  const setSessionHint = (enabled: boolean) => {
    sessionRestoreHint.value = enabled
    syncAuthSessionHint(enabled)
  }

  const setAnonymousState = () => {
    user.value = null
    loaded.value = true
    setSessionHint(false)
    clearMfa()
    return null
  }

  const refresh = async () => {
    const res = await json<ApiResponseWithData<{user: AuthUser; csrf_token?: string}>>(
      "/auth/refresh",
      {
        method: "POST"
      }
    )
    user.value = res.data.user
    loaded.value = true
    setSessionHint(true)
    return res
  }

  const loadMe = async () => {
    if (loadMeInFlight) return await loadMeInFlight

    loadMeInFlight = (async () => {
      if (!canAttemptSessionRestore.value) return setAnonymousState()

      try {
        const res = await json<ApiResponseWithData<{user: AuthUser}>>("/auth/me")
        user.value = res.data.user
        loaded.value = true
        setSessionHint(true)
        return user.value
      } catch (err: unknown) {
        if (errorStatus(err) !== 401) return setAnonymousState()

        try {
          const res = await refresh()
          return res.data.user
        } catch {
          return setAnonymousState()
        }
      }
    })()

    try {
      return await loadMeInFlight
    } finally {
      loadMeInFlight = null
    }
  }

  const ensureLoaded = async () => (loaded.value ? user.value : await loadMe())

  const register = async (input: {
    email: string
    password: string
    display_name?: string
    locale?: string
    timezone?: string
  }) =>
    await json<ApiResponseWithData<{user: AuthUser; verification_ttl: string}>>("/auth/register", {
      method: "POST",
      body: input
    })

  const requestEmailVerification = async (email: string) =>
    await json<ApiResponseWithData<{queued: boolean}>>("/auth/verify-email/request", {
      method: "POST",
      body: {email}
    })

  const login = async (email: string, password: string) => {
    clearMfa()
    const res = await json<
      ApiResponseWithData<{
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
    mfaTicket.value = res.data.mfa_required
      ? String(res.data.mfa_ticket || "").trim() || null
      : null
    mfaExpiresAt.value = res.data.mfa_required
      ? String(res.data.mfa_expires_at || "").trim() || null
      : null
    if (!res.data.mfa_required) setSessionHint(true)
    if (res.data.user) user.value = res.data.user
    return res
  }

  const completeMfa = async (code: string) => {
    const ticket = String(mfaTicket.value || "").trim()
    if (!ticket) throw new Error("Сессия 2FA не найдена.")

    const res = await json<ApiResponseWithData<{user: AuthUser; csrf_token?: string}>>(
      "/auth/login/2fa",
      {
        method: "POST",
        body: {ticket, code}
      }
    )

    user.value = res.data.user
    loaded.value = true
    setSessionHint(true)
    clearMfa()
    return res
  }

  const logout = async () => {
    try {
      await json<ApiResponseWithData<{logged_out: boolean}>>("/auth/logout", {method: "POST"})
    } finally {
      setAnonymousState()
    }
  }

  const forgotPassword = async (email: string) =>
    await json<ApiResponseWithData<{queued: boolean}>>("/auth/password/forgot", {
      method: "POST",
      body: {email}
    })

  const resetPassword = async (token: string, newPassword: string) =>
    await json<ApiResponseWithData<{password_reset: boolean}>>("/auth/password/reset", {
      method: "POST",
      body: {
        token,
        new_password: newPassword
      }
    })

  const verifyEmail = async (token: string) =>
    await json<
      ApiResponseWithData<{
        verified: boolean
        action?: string
        email?: string
        session_hints?: boolean
      }>
    >("/auth/verify-email/confirm", {
      method: "POST",
      body: {token}
    })

  const requestEmailChange = async (email: string) =>
    await json<ApiResponseWithData<{queued: boolean}>>("/auth/email/change/request", {
      method: "POST",
      body: {email}
    })

  const updateProfile = async (input: {
    display_name?: string
    locale?: string
    timezone?: string
    profile?: Record<string, unknown>
    settings?: Record<string, unknown>
  }) => {
    const res = await json<ApiResponseWithData<{changes: AuthUserPatch}>>("/auth/me", {
      method: "PATCH",
      body: input
    })
    if (user.value) user.value = applyAuthUserPatch(user.value, res.data.changes)
    return res
  }

  const changePassword = async (currentPassword: string, newPassword: string) =>
    await json<ApiResponseWithData<{password_changed: boolean}>>("/auth/password/change", {
      method: "POST",
      body: {
        current_password: currentPassword,
        new_password: newPassword
      }
    })

  return {
    user,
    loaded,
    isAuthenticated,
    isAdmin,
    canAttemptSessionRestore,
    mfaTicket,
    mfaExpiresAt,
    setAnonymousState,
    ensureLoaded,
    loadMe,
    refresh,
    register,
    requestEmailVerification,
    login,
    completeMfa,
    logout,
    forgotPassword,
    resetPassword,
    verifyEmail,
    requestEmailChange,
    updateProfile,
    changePassword
  }
}
