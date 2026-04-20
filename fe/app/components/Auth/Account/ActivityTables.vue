<script setup lang="ts">
type SessionActivityRow = LabDataTableRow & {
  id: string
  device: string
  ip: string
  status: string
  count: number
  lastSeenAt: string
  revokableSessionIds: string[]
  hasCurrent: boolean
  action: string
  source: AuthSessionGroupView
}
type LoginAttemptRow = LabDataTableRow & {
  id: string
  createdAt: string
  outcome: string
  ip: string
  risk: string
  details: string
  source: AuthLoginAttemptView
}
type SecurityEventRow = LabDataTableRow & {
  id: string
  createdAt: string
  event: string
  ip: string
  payload: string
  source: AuthSecurityEventView
}

const props = withDefaults(
  defineProps<{
    sessions?: AuthSessionView[]
    loginAttempts?: AuthLoginAttemptView[]
    securityEvents?: AuthSecurityEventView[]
    loading?: boolean
    error?: string
    actionError?: string
    actionInfo?: string
    allowRevoke?: boolean
    user?: AuthUser | null
  }>(),
  {
    sessions: () => [],
    loginAttempts: () => [],
    securityEvents: () => [],
    loading: false,
    error: '',
    actionError: '',
    actionInfo: '',
    allowRevoke: false,
    user: null
  }
)
const emit = defineEmits<{
  revoke: [group: AuthSessionGroupView]
}>()

const { t } = useI18nSection('auth')
const { formatAbsoluteDateTime } = useLocalizedDateTime()

const formatDateTime = (value?: string | null) =>
  formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
const groupedSessions = computed<AuthSessionGroupView[]>(() => {
  const groups = new Map<
    string,
    {
      deviceLabel: string
      ip: string
      latestSession: AuthSessionView
      mfaVerified: boolean
      count: number
      revokableSessionIds: string[]
      currentSessionIds: string[]
      hasCurrent: boolean
    }
  >()
  for (const item of props.sessions) {
    if (item.revoked_at) continue
    const deviceLabel =
      String(item.device_label || '').trim() || String(item.user_agent || '').trim() || 'Неизвестное устройство'
    const ip = String(item.ip || '').trim() || '—'
    const key = `${deviceLabel}::${ip}`
    const current = groups.get(key)
    if (!current) {
      groups.set(key, {
        deviceLabel,
        ip,
        latestSession: item,
        mfaVerified: Boolean(item.mfa_verified),
        count: 1,
        revokableSessionIds: [item.session_id],
        currentSessionIds: item.is_current ? [item.session_id] : [],
        hasCurrent: Boolean(item.is_current)
      })
      continue
    }
    current.count += 1
    current.revokableSessionIds.push(item.session_id)
    if (item.is_current) {
      current.currentSessionIds.push(item.session_id)
      current.hasCurrent = true
    }
    if ((new Date(item.last_seen_at).getTime() || 0) > (new Date(current.latestSession.last_seen_at).getTime() || 0)) {
      current.latestSession = item
      current.mfaVerified = Boolean(item.mfa_verified)
    }
  }
  return Array.from(groups.entries())
    .map(([key, group]) => ({
      key,
      ip: group.ip,
      deviceLabel: group.deviceLabel,
      count: group.count,
      mfaVerified: group.mfaVerified,
      lastSeenAt: group.latestSession.last_seen_at,
      revokableSessionIds: group.revokableSessionIds,
      currentSessionIds: group.currentSessionIds,
      hasCurrent: group.hasCurrent
    }))
    .sort((a, b) => (new Date(b.lastSeenAt).getTime() || 0) - (new Date(a.lastSeenAt).getTime() || 0))
})
const sessionColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'device', label: t('account.activity.device'), cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'status', label: '2FA', nowrap: true },
  { key: 'activity', label: t('account.activity.activity'), cellClass: 'whitespace-normal wrap-break-word' },
  ...(props.allowRevoke ? [{ key: 'action', label: t('account.activity.status'), nowrap: true }] : [])
])
const sessionRows = computed<SessionActivityRow[]>(() =>
  groupedSessions.value.map(item => ({
    id: item.key,
    device: item.deviceLabel,
    ip: item.ip,
    status: item.mfaVerified ? t('account.activity.mfa_verified') : t('account.activity.mfa_unverified'),
    count: item.count,
    lastSeenAt: item.lastSeenAt,
    revokableSessionIds: item.revokableSessionIds,
    hasCurrent: item.hasCurrent,
    action:
      item.revokableSessionIds.length === 0
        ? t('account.activity.current')
        : item.hasCurrent
          ? t('account.activity.revoke_and_logout')
          : t('account.activity.revoke'),
    source: item
  }))
)
const loginAttemptColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: t('account.activity.date'), nowrap: true },
  { key: 'outcome', label: t('account.activity.result'), nowrap: true },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'risk', label: t('account.activity.risk'), nowrap: true },
  { key: 'details', label: t('account.activity.details'), cellClass: 'whitespace-normal wrap-break-word' }
])
const loginAttemptRows = computed<LoginAttemptRow[]>(() =>
  props.loginAttempts.map(item => ({
    id: item.attempt_id,
    createdAt: formatDateTime(item.created_at),
    outcome: item.outcome || '—',
    ip: item.ip || '—',
    risk: String(item.risk_score ?? '—'),
    details: item.failure_reason || item.suspicious_reason || item.user_agent || '—',
    source: item
  }))
)
const securityEventColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: t('account.activity.date'), nowrap: true },
  { key: 'event', label: t('account.activity.event'), cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'payload', label: 'Payload', cellClass: 'whitespace-normal wrap-break-word' }
])
const securityEventRows = computed<SecurityEventRow[]>(() =>
  props.securityEvents.map(item => ({
    id: item.event_id,
    createdAt: formatDateTime(item.created_at),
    event: `${item.event_type} · ${item.severity}`,
    ip: item.ip || '—',
    payload: JSON.stringify(item.payload || {}),
    source: item
  }))
)
</script>

<template>
  <section class="space-y-4">
    <div v-if="user" class="flex flex-wrap items-baseline gap-x-6 gap-y-2">
      <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
        <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">
          {{ t('account.activity.last_login') }}
        </span>
        <span class="text-sm text-(--lab-text-primary)">{{ formatDateTime(user.last_login_at) }}</span>
      </div>
      <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
        <span class="shrink-0 text-xs tracking-wide text-(--lab-text-muted) uppercase">
          {{ t('account.activity.created_at') }}
        </span>
        <span class="text-sm text-(--lab-text-primary)">{{ formatDateTime(user.created_at) }}</span>
      </div>
    </div>
    <LabNotify :text="actionError" tone="error" size="xs" />
    <LabNotify :text="actionInfo" tone="success" size="xs" />
    <p v-if="loading" class="text-xs text-(--lab-text-muted)">Загрузка сессий и активности…</p>
    <LabNotify v-else-if="error" :text="error" tone="error" size="xs" />
    <template v-else>
      <LabDataTable
        :title="t('account.activity.sessions')"
        :columns="sessionColumns"
        :rows="sessionRows"
        empty-text="Сессии не найдены."
        nowrap
      >
        <template #cell-device="{ row }">
          <div class="space-y-1">
            <p class="text-sm">{{ (row as SessionActivityRow).device }}</p>
            <p class="text-xs text-(--lab-text-muted)">{{ (row as SessionActivityRow).ip }}</p>
          </div>
        </template>
        <template #cell-activity="{ row }">
          <div class="space-y-1 text-xs">
            <p>{{ t('account.activity.sessions_count', { count: (row as SessionActivityRow).count }) }}</p>
            <p class="text-(--lab-text-muted)">
              {{ t('account.activity.last_activity') }}
              <LabRelativeTime :datetime="(row as SessionActivityRow).lastSeenAt" compact />
            </p>
          </div>
        </template>
        <template v-if="allowRevoke" #cell-action="{ row }">
          <LabBaseButton
            variant="secondary"
            size="lg"
            class="text-xs"
            :disabled="(row as SessionActivityRow).revokableSessionIds.length === 0"
            @click="emit('revoke', (row as SessionActivityRow).source)"
          >
            {{ (row as SessionActivityRow).action }}
          </LabBaseButton>
        </template>
      </LabDataTable>
      <LabDataTable
        :title="t('account.activity.attempts')"
        :columns="loginAttemptColumns"
        :rows="loginAttemptRows"
        empty-text="Попыток входа пока нет."
        nowrap
      >
        <template #cell-createdAt="{ row }">
          <div class="space-y-1 text-xs">
            <p>{{ (row as LoginAttemptRow).createdAt }}</p>
            <p class="text-(--lab-text-muted)">
              <LabRelativeTime :datetime="(row as LoginAttemptRow).source.created_at" compact />
            </p>
          </div>
        </template>
        <template #cell-details="{ row }">
          <p class="text-xs wrap-break-word">{{ (row as LoginAttemptRow).details }}</p>
        </template>
      </LabDataTable>
      <LabDataTable
        :title="t('account.activity.events')"
        :columns="securityEventColumns"
        :rows="securityEventRows"
        empty-text="Событий безопасности пока нет."
        nowrap
      >
        <template #cell-createdAt="{ row }">
          <div class="space-y-1 text-xs">
            <p>{{ (row as SecurityEventRow).createdAt }}</p>
            <p class="text-(--lab-text-muted)">
              <LabRelativeTime :datetime="(row as SecurityEventRow).source.created_at" compact />
            </p>
          </div>
        </template>
        <template #cell-payload="{ row }">
          <p class="text-xs wrap-break-word">{{ (row as SecurityEventRow).payload }}</p>
        </template>
      </LabDataTable>
    </template>
  </section>
</template>
