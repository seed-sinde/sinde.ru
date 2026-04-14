<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const targetUserId = computed(() => String(route.params.user || '').trim())
const {
  ensureLoaded,
  isAdmin,
  publicUserProfile,
  adminUserDetail,
  adminBlockUser,
  adminUnblockUser,
  adminForceLogoutUser,
  adminDeleteUser
} = useAuth()
const { adminUserAccess, adminUserOrders } = usePayments()
const { formatAbsoluteDateTime } = useLocalizedDateTime()

await ensureLoaded()

const adminTab = ref<'overview' | 'payments' | 'activity'>('overview')
const activityTab = ref<'sessions' | 'attempts' | 'events'>('sessions')
const clientAdminView = ref(false)
const noticesReady = ref(false)

const loading = ref(true)
const pageError = ref('')
const pageInfo = ref('')
const publicProfile = ref<PublicUserProfileView | null>(null)
const detail = ref<AdminUserDetailView | null>(null)
const access = ref<PaymentAccessSummary | null>(null)
const orders = ref<PaymentOrderView[]>([])
const ordersLoading = ref(false)

const title = computed(() => {
  const displayName =
    String(detail.value?.user.display_name || '').trim() || String(publicProfile.value?.display_name || '').trim()
  return displayName || 'Пользователь'
})
usePageSeo({
  title,
  description: 'Профиль пользователя'
})

const adminTabItems: LabTabItem[] = [
  { value: 'overview', label: 'Профиль' },
  { value: 'payments', label: 'Транзакции' },
  { value: 'activity', label: 'Активность' }
]
const activityTabItems: LabTabItem[] = [
  { value: 'sessions', label: 'Сессии' },
  { value: 'attempts', label: 'Попытки входа' },
  { value: 'events', label: 'События' }
]

const displayUser = computed<AuthUser | null>(() => {
  if (detail.value?.user) return detail.value.user
  if (!publicProfile.value) return null
  return {
    user_id: publicProfile.value.user_id,
    email: '',
    status: 'active',
    display_name: publicProfile.value.display_name,
    locale: 'ru-RU',
    timezone: 'Europe/Moscow',
    roles: [],
    is_two_factor_enabled: false,
    primary_trait_uuid: publicProfile.value.primary_trait_uuid ?? null,
    created_at: new Date(0).toISOString(),
    profile: publicProfile.value.profile || {},
    settings: {}
  }
})

const paymentAccessUntilText = computed(() => {
  if (!access.value?.access_until) return ''
  return formatAbsoluteDateTime(access.value.access_until)
})
const subscriptionTooltipText = computed(() => {
  if (!access.value?.has_active_access || !paymentAccessUntilText.value) return ''
  return `Подписка активна до ${paymentAccessUntilText.value}`
})
const adminDetailUser = computed(() => detail.value?.user || null)
const targetUserIsAdmin = computed(() => adminDetailUser.value?.roles.includes('admin') === true)
const showAdminView = computed(() => clientAdminView.value && isAdmin.value && Boolean(detail.value))
const paymentHistoryColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'plan', label: 'План' },
  { key: 'amount', label: 'Сумма', nowrap: true },
  { key: 'status', label: 'Статус', nowrap: true },
  { key: 'access', label: 'Доступ' }
])
const paymentHistoryRows = computed(() =>
  orders.value.map((item) => ({
    id: item.order_id,
    createdAt: formatDateTime(item.created_at),
    plan: paymentPlanLabel(item.plan_code),
    amount: formatPaymentAmount(item.amount),
    status: paymentStatusLabel(item.status),
    access: item.access_until
      ? `До ${formatDateTime(item.access_until)}`
      : item.access_from
        ? `С ${formatDateTime(item.access_from)}`
        : '—'
  }))
)
const sessionColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'device', label: 'Устройство', cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'status', label: '2FA', nowrap: true },
  { key: 'activity', label: 'Активность', cellClass: 'whitespace-normal wrap-break-word' }
])
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
  for (const item of detail.value?.sessions || []) {
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
    const currentTs = new Date(current.latestSession.last_seen_at).getTime() || 0
    const nextTs = new Date(item.last_seen_at).getTime() || 0
    if (nextTs > currentTs) {
      current.latestSession = item
      current.mfaVerified = Boolean(item.mfa_verified)
    }
  }
  return Array.from(groups.entries()).map(([key, group]) => ({
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
})
const sessionRows = computed(() =>
  groupedSessions.value.map((item) => ({
    id: item.key,
    device: item.deviceLabel,
    ip: item.ip,
    status: item.mfaVerified ? '2FA подтверждена' : '2FA не подтверждена',
    count: item.count,
    lastSeenAt: item.lastSeenAt
  }))
)
const loginAttemptColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'outcome', label: 'Результат', nowrap: true },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'risk', label: 'Риск', nowrap: true },
  { key: 'details', label: 'Детали', cellClass: 'whitespace-normal wrap-break-word' }
])
const loginAttemptRows = computed(() =>
  (detail.value?.login_attempts || []).map((item) => ({
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
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'event', label: 'Событие', cellClass: 'whitespace-normal wrap-break-word' },
  { key: 'ip', label: 'IP', nowrap: true },
  { key: 'payload', label: 'Payload', cellClass: 'whitespace-normal wrap-break-word' }
])
const securityEventRows = computed(() =>
  (detail.value?.security_events || []).map((item) => ({
    id: item.event_id,
    createdAt: formatDateTime(item.created_at),
    event: `${item.event_type} · ${item.severity}`,
    ip: item.ip || '—',
    payload: JSON.stringify(item.payload || {}),
    source: item
  }))
)

const formatDateTime = (value?: string | null) =>
  formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
const formatPaymentAmount = (value?: number | null) => {
  const amount = Number(value || 0)
  return new Intl.NumberFormat('ru-RU').format(Math.floor(amount / 100)) + ' ₽'
}
const paymentPlanLabel = (planCode?: string | null) => (String(planCode || '').trim() === 'donation' ? 'Донат' : 'Pro')
const paymentStatusLabel = (status?: string | null) => {
  switch (String(status || '').trim()) {
    case 'success':
      return 'Успешно'
    case 'pending':
      return 'В обработке'
    case 'failed':
      return 'Ошибка'
    case 'canceled':
      return 'Отменён'
    case 'refunded':
      return 'Возврат'
    default:
      return '—'
  }
}
const loadOrders = async () => {
  if (!isAdmin.value) return
  ordersLoading.value = true
  try {
    const res = await adminUserOrders(targetUserId.value)
    orders.value = res.data.items || []
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось загрузить историю транзакций.'
  } finally {
    ordersLoading.value = false
  }
}
const loadPage = async () => {
  loading.value = true
  pageError.value = ''
  try {
    const publicRes = await publicUserProfile(targetUserId.value)
    publicProfile.value = publicRes.data
    if (isAdmin.value) {
      const [detailRes, accessRes] = await Promise.all([
        adminUserDetail(targetUserId.value),
        adminUserAccess(targetUserId.value)
      ])
      detail.value = detailRes.data
      access.value = accessRes.data || null
    }
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось загрузить профиль пользователя.'
  } finally {
    loading.value = false
  }
}
const refreshAll = async () => {
  await loadPage()
  if (isAdmin.value && adminTab.value === 'payments') {
    await loadOrders()
  }
}
const blockUser = async () => {
  if (!detail.value) return
  if (targetUserIsAdmin.value) {
    pageError.value = 'Нельзя блокировать пользователя с ролью admin.'
    return
  }
  pageError.value = ''
  pageInfo.value = ''
  try {
    await adminBlockUser(detail.value.user.user_id, '')
    pageInfo.value = `Пользователь ${detail.value.user.email} заблокирован.`
    await refreshAll()
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось заблокировать пользователя.'
  }
}
const unblockUser = async () => {
  if (!detail.value) return
  pageError.value = ''
  pageInfo.value = ''
  try {
    await adminUnblockUser(detail.value.user.user_id)
    pageInfo.value = `Пользователь ${detail.value.user.email} разблокирован.`
    await refreshAll()
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось разблокировать пользователя.'
  }
}
const forceLogout = async () => {
  if (!detail.value) return
  if (targetUserIsAdmin.value) {
    pageError.value = 'Нельзя сбрасывать сессии пользователя с ролью admin.'
    return
  }
  pageError.value = ''
  pageInfo.value = ''
  try {
    await adminForceLogoutUser(detail.value.user.user_id)
    pageInfo.value = `Все сессии пользователя ${detail.value.user.email} завершены.`
    await refreshAll()
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось завершить сессии пользователя.'
  }
}
const deleteUser = async () => {
  if (!detail.value) return
  if (targetUserIsAdmin.value) {
    pageError.value = 'Нельзя удалять пользователя с ролью admin.'
    return
  }
  pageError.value = ''
  pageInfo.value = ''
  try {
    await adminDeleteUser(detail.value.user.user_id)
    await router.push('/auth/admin?tab=users')
  } catch (err: any) {
    pageError.value = err?.data?.message || err?.message || 'Не удалось удалить пользователя.'
  }
}

await loadPage()

onMounted(() => {
  const hasMissingPageData = !publicProfile.value || (isAdmin.value && !detail.value)
  if (pageError.value) {
    pageError.value = ''
    if (hasMissingPageData) {
      void loadPage()
    }
  }
  clientAdminView.value = isAdmin.value && Boolean(detail.value)
  noticesReady.value = true
})

watch(adminTab, async (next) => {
  if (next === 'payments' && isAdmin.value && orders.value.length === 0 && !ordersLoading.value) {
    await loadOrders()
  }
})

watch(
  () => [isAdmin.value, detail.value] as const,
  ([admin, nextDetail]) => {
    clientAdminView.value = admin && Boolean(nextDetail)
  }
)
</script>

<template>
  <div>
    <LabNavHeader :title="title" />
    <LabLoader v-if="loading" variant="inline" label="Загрузка пользователя…" />
    <template v-else-if="displayUser">
      <LabNotify v-if="noticesReady && pageError" :text="pageError" tone="error" size="xs" />
      <LabNotify v-if="noticesReady && pageInfo" :text="pageInfo" tone="success" size="xs" />
      <section v-if="showAdminView && adminDetailUser">
        <LabNavTabs v-model="adminTab" :items="adminTabItems">
          <template #panel-overview>
            <section class="p-4">
              <div class="flex flex-wrap items-start gap-4">
                <LabAvatar version="profile" :user="adminDetailUser" />
                <div class="flex min-w-0 flex-1 flex-col gap-3">
                  <div class="flex flex-wrap items-start gap-2">
                    <h2 class="text-xl font-semibold">{{ adminDetailUser.display_name || 'Пользователь' }}</h2>
                    <LabBaseTooltip
                      v-if="subscriptionTooltipText"
                      :text="subscriptionTooltipText"
                      side="right"
                      align="left"
                      :offset="10"
                      :cross-axis-offset="0"
                    >
                      <template #trigger>
                        <LabBaseButton
                          icon="ic:round-auto-awesome"
                          icon-only
                          variant="ghost"
                          size="sm"
                          button-class="h-8 w-8 rounded-full border-transparent text-orange-300 hover:bg-(--lab-bg-surface-hover) focus:bg-(--lab-bg-surface-hover) focus-visible:bg-(--lab-bg-surface-hover)"
                          aria-label="Статус подписки"
                        />
                      </template>
                    </LabBaseTooltip>
                  </div>
                  <div class="flex flex-wrap items-baseline gap-x-6 gap-y-2 text-sm">
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="lab-text-muted shrink-0 text-xs tracking-wide uppercase">Email</span>
                      <NuxtLink :to="`mailto:${adminDetailUser.email}`" external>{{ adminDetailUser.email }}</NuxtLink>
                    </div>
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="lab-text-muted shrink-0 text-xs tracking-wide uppercase">Последний вход</span>
                      <span>{{ formatDateTime(adminDetailUser.last_login_at) }}</span>
                    </div>
                    <div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                      <span class="lab-text-muted shrink-0 text-xs tracking-wide uppercase">Регистрация</span>
                      <span>{{ formatDateTime(adminDetailUser.created_at) }}</span>
                    </div>
                  </div>
                  <div class="flex flex-nowrap items-center gap-1 overflow-x-auto whitespace-nowrap">
                    <span v-if="adminDetailUser.status !== 'blocked' && !targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Блок"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить блокировку аккаунта?"
                        icon="ic:round-block"
                        @confirm="blockUser"
                      />
                    </span>
                    <span v-if="adminDetailUser.status !== 'blocked' && !targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-block"
                        aria-label="Заблокировать"
                        confirm-aria-label="Подтвердить блокировку"
                        confirm-label="Ок"
                        tooltip="Подтвердить блокировку аккаунта?"
                        @confirm="blockUser"
                      />
                    </span>
                    <LabBaseButton
                      v-if="adminDetailUser.status === 'blocked'"
                      class="max-sm:hidden"
                      variant="secondary"
                      size="xs"
                      icon="ic:round-lock-open"
                      label="Разблок"
                      @click="unblockUser"
                    />
                    <LabBaseButton
                      v-if="adminDetailUser.status === 'blocked'"
                      class="sm:hidden"
                      variant="secondary"
                      size="xs"
                      icon="ic:round-lock-open"
                      icon-only
                      aria-label="Разблокировать"
                      @click="unblockUser"
                    />
                    <span v-if="!targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Сброс сессий"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить принудительный сброс всех сессий?"
                        icon="ic:round-logout"
                        @confirm="forceLogout"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-logout"
                        aria-label="Сбросить сессии"
                        confirm-aria-label="Подтвердить сброс сессий"
                        confirm-label="Ок"
                        tooltip="Подтвердить принудительный сброс всех сессий?"
                        @confirm="forceLogout"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="max-sm:hidden">
                      <LabConfirmActionButton
                        label="Удалить"
                        confirm-label="Подтвердить"
                        tooltip="Подтвердить удаление аккаунта?"
                        icon="ic:round-delete"
                        @confirm="deleteUser"
                      />
                    </span>
                    <span v-if="!targetUserIsAdmin" class="sm:hidden">
                      <LabConfirmActionButton
                        icon-only
                        icon="ic:round-delete"
                        aria-label="Удалить аккаунт"
                        confirm-aria-label="Подтвердить удаление аккаунта"
                        confirm-label="Ок"
                        tooltip="Подтвердить удаление аккаунта?"
                        @confirm="deleteUser"
                      />
                    </span>
                  </div>
                </div>
              </div>
            </section>
          </template>
          <template #panel-payments>
            <LabDataTable
              :columns="paymentHistoryColumns"
              :rows="paymentHistoryRows"
              :loading="ordersLoading"
              empty-text="История транзакций пуста."
            />
          </template>
          <template #panel-activity>
            <LabNavTabs v-model="activityTab" :items="activityTabItems">
              <template #panel-sessions>
                <LabDataTable :columns="sessionColumns" :rows="sessionRows" empty-text="Сессии не найдены.">
                  <template #cell-device="{ row }">
                    <div class="space-y-1">
                      <p class="text-sm">{{ row.device }}</p>
                      <p class="text-xs text-(--lab-text-muted)">{{ row.ip }}</p>
                    </div>
                  </template>
                  <template #cell-activity="{ row }">
                    <div class="space-y-1 text-xs">
                      <p>Сессий: {{ row.count }}</p>
                      <p class="text-(--lab-text-muted)">
                        Последняя активность
                        <LabRelativeTime :datetime="row.lastSeenAt" compact />
                      </p>
                    </div>
                  </template>
                </LabDataTable>
              </template>
              <template #panel-attempts>
                <LabDataTable :columns="loginAttemptColumns" :rows="loginAttemptRows" empty-text="Попыток входа нет.">
                  <template #cell-createdAt="{ row }">
                    <div class="space-y-1 text-xs">
                      <p>{{ row.createdAt }}</p>
                      <p class="text-(--lab-text-muted)">
                        <LabRelativeTime :datetime="row.source.created_at" compact />
                      </p>
                    </div>
                  </template>
                </LabDataTable>
              </template>
              <template #panel-events>
                <LabDataTable
                  :columns="securityEventColumns"
                  :rows="securityEventRows"
                  empty-text="Событий безопасности нет."
                >
                  <template #cell-createdAt="{ row }">
                    <div class="space-y-1 text-xs">
                      <p>{{ row.createdAt }}</p>
                      <p class="text-(--lab-text-muted)">
                        <LabRelativeTime :datetime="row.source.created_at" compact />
                      </p>
                    </div>
                  </template>
                </LabDataTable>
              </template>
            </LabNavTabs>
          </template>
        </LabNavTabs>
      </section>
      <section v-else class="space-y-4 p-4">
        <div class="flex flex-wrap items-start gap-4">
          <LabAvatar version="profile" :user="displayUser" />
          <div class="space-y-2">
            <h1 class="text-2xl font-semibold">{{ displayUser.display_name || 'Пользователь' }}</h1>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>
