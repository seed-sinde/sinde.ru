export const PASSWORD_MAX_LENGTH = 128
const COMMON_PASSWORDS = new Set(['password', '12345678', '123456789', 'qwerty123', 'letmein', 'adminadmin'])
export type PasswordPolicyRuleKey = 'length' | 'email' | 'common'
export type PasswordPolicyRuleState = 'idle' | 'passed' | 'failed'
export type PasswordPolicyRule = {
  key: PasswordPolicyRuleKey
  passed: boolean
  state: PasswordPolicyRuleState
}
export const getPasswordPolicyRules = (password: unknown, email: unknown, minLength: number): PasswordPolicyRule[] => {
  const normalizedPassword = String(password ?? '').trim()
  const normalizedPasswordLower = normalizedPassword.toLowerCase()
  const normalizedEmail = normalizeEmail(email)
  const hasPassword = normalizedPassword.length > 0
  const hasEmail = normalizedEmail.length > 0
  return [
    {
      key: 'length',
      passed: normalizedPassword.length >= minLength && normalizedPassword.length <= PASSWORD_MAX_LENGTH,
      state: !hasPassword
        ? 'idle'
        : normalizedPassword.length >= minLength && normalizedPassword.length <= PASSWORD_MAX_LENGTH
          ? 'passed'
          : 'failed'
    },
    {
      key: 'email',
      passed: hasPassword && hasEmail && !normalizedPasswordLower.includes(normalizedEmail),
      state: !hasPassword || !hasEmail
        ? 'idle'
        : !normalizedPasswordLower.includes(normalizedEmail)
          ? 'passed'
          : 'failed'
    },
    {
      key: 'common',
      passed: hasPassword && !COMMON_PASSWORDS.has(normalizedPasswordLower),
      state: !hasPassword
        ? 'idle'
        : !COMMON_PASSWORDS.has(normalizedPasswordLower)
          ? 'passed'
          : 'failed'
    }
  ]
}
