export type AuthUser = {
  user_id: string
  email: string
  status: string
  email_verified_at?: string | null
  display_name: string
  locale: string
  timezone: string
  roles: string[]
  is_two_factor_enabled: boolean
  primary_trait_uuid?: string | null
  last_login_at?: string | null
  blocked_reason?: string
  blocked_at?: string | null
  created_at: string
  profile: Record<string, unknown>
  settings: Record<string, unknown>
}
export type AuthUserPatch = {
  display_name?: string
  locale?: string
  timezone?: string
  primary_trait_uuid?: string | null
  profile?: Record<string, unknown>
  settings?: Record<string, unknown>
}
export type SavedTraitSetView = {
  saved_set_id: string
  set_uuid: string
  name: string
  description: string
  created_at: string
  updated_at: string
}
export type AuthSessionView = {
  session_id: string
  device_label: string
  ip: string
  user_agent: string
  mfa_verified: boolean
  last_seen_at: string
  expires_at: string
  revoked_at?: string | null
  revoke_reason?: string
  created_at: string
  is_current: boolean
}
export type AuthLoginAttemptView = {
  attempt_id: string
  ip: string
  user_agent: string
  outcome: string
  failure_reason?: string
  risk_score: number
  suspicious_reason?: string
  created_at: string
}
export type AuthSecurityEventView = {
  event_id: string
  category: string
  event_type: string
  severity: string
  ip: string
  user_agent: string
  session_id?: string | null
  payload: Record<string, unknown>
  created_at: string
}
export type UserSummaryAdmin = {
  users_total: number
  pending_recipes_total: number
  recipe_status_totals: Record<string, number>
  new_users_since_last_login: number
  new_pending_recipes_since_last_login: number
}
export type UserSummary = {
  new_approved_recipes_since_last_login: number
  new_rejected_recipes_since_last_login: number
  last_login_before_current_session_at?: string | null
  notifications_read_at?: string | null
  notifications_since_at?: string | null
  admin?: UserSummaryAdmin | null
  has_unread: boolean
  checked_at: string
}
export type AdminSummary = UserSummaryAdmin & {
  last_login_before_current_session_at?: string | null
  notifications_read_at?: string | null
  notifications_since_at?: string | null
  has_unread: boolean
  checked_at: string
}
export type AuthSessionGroupView = {
  key: string
  ip: string
  deviceLabel: string
  count: number
  mfaVerified: boolean
  lastSeenAt: string
  revokableSessionIds: string[]
  currentSessionIds: string[]
  hasCurrent: boolean
}
