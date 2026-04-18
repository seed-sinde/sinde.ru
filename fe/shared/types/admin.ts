import type {AuthLoginAttemptView, AuthSecurityEventView, AuthSessionView, AuthUser} from "./auth"

export type AdminTab = "users" | "moderation" | "keys" | "analysis"
export type AdminModerationStatus = "draft" | "pending" | "approved" | "rejected"
export type AdminUserView = {
  user_id: string
  email: string
  status: string
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
export type PublicUserProfileView = {
  user_id: string
  display_name: string
  primary_trait_uuid?: string | null
  profile: Record<string, unknown>
}
export type AdminUserDetailView = {
  user: AuthUser
  sessions: AuthSessionView[]
  login_attempts: AuthLoginAttemptView[]
  security_events: AuthSecurityEventView[]
}
export type AdminTraitKeySearchItem = {
  key_id: number
  syn: string
  meta: Record<string, unknown>
  trait_count: number
}
export type AdminTraitsSetsAnalysisTopKey = {
  key_id: number
  syn: string
  trait_count: number
}
export type AdminTraitsSetsAnalysis = {
  total_traits: number
  unique_trait_pairs: number
  unique_trait_keys: number
  traits_referenced_in_sets: number
  orphan_traits: number
  total_sets: number
  unique_sets_by_children: number
  derived_sets: number
  derived_set_rate: number
  set_uniqueness_rate: number
  trait_coverage_in_sets_rate: number
  top_keys: AdminTraitsSetsAnalysisTopKey[]
}
