const isPatchRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null && !Array.isArray(value)
const mergePatchRecord = (
  current: Record<string, unknown>,
  patch: Record<string, unknown>
): Record<string, unknown> => {
  return Object.entries(patch).reduce<Record<string, unknown>>(
    (next, [key, value]) => {
      if (value === null) {
        const { [key]: _removed, ...rest } = next
        return rest
      }
      if (isPatchRecord(value) && isPatchRecord(next[key])) {
        return {
          ...next,
          [key]: mergePatchRecord(next[key] as Record<string, unknown>, value)
        }
      }
      return {
        ...next,
        [key]: value
      }
    },
    { ...current }
  )
}
export const applyAuthUserPatch = (user: AuthUser, patch?: AuthUserPatch | null): AuthUser => {
  if (!patch) return user
  return {
    ...user,
    display_name: patch.display_name ?? user.display_name,
    locale: patch.locale ?? user.locale,
    timezone: patch.timezone ?? user.timezone,
    profile: isPatchRecord(patch.profile) ? mergePatchRecord(user.profile, patch.profile) : user.profile,
    settings: isPatchRecord(patch.settings) ? mergePatchRecord(user.settings, patch.settings) : user.settings
  }
}
