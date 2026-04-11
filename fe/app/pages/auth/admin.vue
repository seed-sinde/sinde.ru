<script setup lang="ts">
  definePageMeta({
    middleware: ['admin-only']
  })
  const title = 'Админка'
  usePageSeo({
    title,
    description: 'Управление пользователями, ролями и базовая аналитика.'
  })
  const {
    user,
    adminListUsers,
    sharedAdminSummary,
    loadSharedAdminSummary,
    adminMarkSummaryRead,
    adminSetUserRole,
    adminBlockUser,
    adminUnblockUser,
    adminForceLogoutUser,
    adminDeleteUser,
    adminSearchKeys,
    adminTraitsSetsAnalysis
  } = useAuth()
  const { adminSummary: loadPaymentsAdminSummary } = usePayments()
  const { formatAbsoluteDateTime } = useLocalizedDateTime()
  const { localeTag } = useInterfacePreferences()
  const users = ref<AdminUserView[]>([])
  const usersTotal = ref(0)
  const usersLimit = ref(50)
  const usersOffset = ref(0)
  const usersLoading = ref(false)
  const usersError = ref('')
  const usersInfo = ref('')
  const userSearch = ref('')
  const userStatusFilter = ref('')
  const userRoleFilter = ref('')
  const userStatusOptions: SelectOptionInput[] = [
    { value: '', label: 'все' },
    { value: 'active', label: 'Активен' },
    { value: 'pending_verification', label: 'Ожидает подтверждения' },
    { value: 'blocked', label: 'Заблокирован' }
  ]
  const userRoleOptions: SelectOptionInput[] = [
    { value: '', label: 'все' },
    { value: 'admin', label: 'admin' },
    { value: 'user', label: 'user' }
  ]
  const editableRoleOptions: SelectOptionInput[] = [
    { value: 'user', label: 'user' },
    { value: 'admin', label: 'admin' }
  ]
  const roleDrafts = reactive<Record<string, 'admin' | 'user'>>({})
  const roleSavingUserId = ref<string | null>(null)
  const keysLoading = ref(false)
  const keysError = ref('')
  const keysQuery = ref('')
  const keys = ref<AdminTraitKeySearchItem[]>([])
  let keysRequestId = 0
  const analysisLoading = ref(false)
  const analysisError = ref('')
  const analysis = ref<AdminTraitsSetsAnalysis | null>(null)
  const summaryLoading = ref(false)
  const summaryError = ref('')
  const summaryInfo = ref('')
  const paymentsSummaryLoading = ref(false)
  const paymentsSummaryError = ref('')
  const paymentsSummary = ref<PaymentAdminOrdersSummary | null>(null)
  const summary = computed(() => sharedAdminSummary.value)
  const summaryReadPending = ref(false)
  const adminTab = ref<AdminTab>('users')
  const adminBreadcrumbItems = computed<BreadcrumbItem[]>(() => {
    if (adminTab.value === 'users') {
      return [{ label: 'Админка', current: true }]
    }
    const labels: Record<Exclude<AdminTab, 'users'>, string> = {
      moderation: 'Модерация рецептов',
      keys: 'Ключи',
      analysis: 'Аналитика'
    }
    return [
      { label: 'Админка', to: '/auth/admin' },
      { label: labels[adminTab.value as Exclude<AdminTab, 'users'>], current: true, kind: 'tab' }
    ]
  })
  const adminTabItems = computed<LabTabItem[]>(() => {
    const items: LabTabItem[] = [
      { value: 'users', label: 'Пользователи' },
      { value: 'moderation', label: 'Модерация рецептов' },
      { value: 'keys', label: 'Ключи' },
      { value: 'analysis', label: 'Аналитика' }
    ]
    if (usersTotal.value) items[0]!.badge = usersTotal.value
    if (moderationStatusTotals.value.pending) items[1]!.badge = moderationStatusTotals.value.pending
    if (keys.value.length) items[2]!.badge = keys.value.length
    if (summary.value?.has_unread) {
      items[3]!.badge =
        Number(summary.value.new_users_since_last_login || 0) +
        Number(summary.value.new_pending_recipes_since_last_login || 0)
    }
    return items
  })
  const adminSectionClass = 'space-y-4 py-1'
  const adminSubsectionClass = 'space-y-3 py-1'
  const adminStatCardClass = 'space-y-1'
  const adminMetaTextClass = 'text-xs text-zinc-500'
  const adminConfirmButtonClass = 'h-8 min-w-28'
  const adminCompactButtonClass = ''
  const adminInlineHintClass = 'text-xs text-zinc-500'
  const moderationLoading = ref(false)
  const moderationError = ref('')
  const moderationInfo = ref('')
  const moderationStatus = ref<AdminModerationStatus | 'all'>('pending')
  const moderationItems = ref<KitchenRecipe[]>([])
  const moderationTotal = ref(0)
  const moderationLimit = ref(30)
  const moderationOffset = ref(0)
  const moderationStatusTotals = ref<Record<string, number>>({
    all: 0,
    pending: 0,
    approved: 0,
    rejected: 0,
    draft: 0
  })
  const moderationOwnerEditorRecipeId = ref<string | null>(null)
  const moderationOwnerQuery = ref('')
  const moderationOwnerCandidates = ref<AdminUserView[]>([])
  const moderationOwnerCandidatesLoading = ref(false)
  const moderationOwnerCandidatesError = ref('')
  const moderationOwnerSavingRecipeId = ref<string | null>(null)
  const moderationOwnerDirectory = reactive<Record<string, { email: string; display_name: string }>>({})
  const moderationNoteDrafts = reactive<Record<string, string>>({})
  let moderationOwnerSearchRequestId = 0
  const moderationHasPrevPage = computed(() => moderationOffset.value > 0)
  const moderationHasNextPage = computed(() => moderationOffset.value + moderationLimit.value < moderationTotal.value)
  const normalizeModerationStatusTotals = (raw?: Record<string, number> | null) => {
    const base = {
      all: 0,
      pending: 0,
      approved: 0,
      rejected: 0,
      draft: 0
    }
    if (!raw || typeof raw !== 'object') return base
    for (const key of Object.keys(base)) {
      const value = Number(raw[key] || 0)
      base[key as keyof typeof base] = Number.isFinite(value) && value > 0 ? value : 0
    }
    return base
  }
  const hasPrevPage = computed(() => usersOffset.value > 0)
  const hasNextPage = computed(() => usersOffset.value + usersLimit.value < usersTotal.value)
  const formatPercent = (value?: number) => `${Math.round(Math.max(0, Number(value || 0)) * 100)}%`
  const displayNumber = (value?: number | null) => (value === null || value === undefined ? '—' : value)
  const displayMoney = (value?: number | null) => formatPaymentAmount(Number(value || 0), localeTag.value)
  const displayDateTime = (value?: string | null) => (value ? formatAbsoluteDateTime(value) : '—')
  const paymentsStats = computed(() => [
    { key: 'orders_total', label: 'Всего заказов', value: displayNumber(paymentsSummary.value?.orders_total) },
    { key: 'orders_success', label: 'Успешных', value: displayNumber(paymentsSummary.value?.orders_success) },
    { key: 'orders_pending', label: 'В обработке', value: displayNumber(paymentsSummary.value?.orders_pending) },
    { key: 'orders_failed', label: 'Неуспешных', value: displayNumber(paymentsSummary.value?.orders_failed) },
    { key: 'orders_refunded', label: 'Возвратов', value: displayNumber(paymentsSummary.value?.orders_refunded) },
    {
      key: 'paid_users_total',
      label: 'Платящих пользователей',
      value: displayNumber(paymentsSummary.value?.paid_users_total)
    },
    {
      key: 'patron_users_total',
      label: 'Поддержавших проект',
      value: displayNumber(paymentsSummary.value?.patron_users_total)
    },
    {
      key: 'active_access_users',
      label: 'Активный доступ',
      value: displayNumber(paymentsSummary.value?.active_access_users)
    },
    { key: 'gross_revenue', label: 'Оборот', value: displayMoney(paymentsSummary.value?.gross_revenue) },
    { key: 'net_revenue', label: 'Чистая выручка', value: displayMoney(paymentsSummary.value?.net_revenue) },
    { key: 'tip_revenue', label: 'Донаты', value: displayMoney(paymentsSummary.value?.tip_revenue) },
    { key: 'mrr', label: 'MRR', value: displayMoney(paymentsSummary.value?.mrr) },
    { key: 'churn_rate', label: 'Churn', value: formatPercent(paymentsSummary.value?.churn_rate) },
    { key: 'patron_share', label: 'Доля поддержавших', value: formatPercent(paymentsSummary.value?.patron_share) }
  ])
  const currentRole = (item: AdminUserView): 'admin' | 'user' => (item.roles.includes('admin') ? 'admin' : 'user')
  const isAdminAccount = (item: AdminUserView) => currentRole(item) === 'admin'
  const userStatusLabel = (value?: string) => {
    const key = String(value || '').trim()
    if (key === 'active') return 'Активен'
    if (key === 'pending_verification') return 'Ожидает подтверждения'
    if (key === 'blocked') return 'Заблокирован'
    return key || '—'
  }
  const syncRoleDrafts = (items: AdminUserView[]) => {
    for (const item of items) {
      roleDrafts[item.user_id] = currentRole(item)
    }
  }
  const adminListUserToAuthUser = (item: AdminUserView): AuthUser => ({
    user_id: item.user_id,
    email: item.email,
    status: item.status,
    display_name: item.display_name,
    locale: item.locale,
    timezone: item.timezone,
    roles: item.roles,
    is_two_factor_enabled: item.is_two_factor_enabled,
    last_login_at: item.last_login_at || null,
    blocked_reason: item.blocked_reason || '',
    blocked_at: item.blocked_at || null,
    created_at: item.created_at,
    profile: item.profile || {},
    settings: item.settings || {}
  })
  const usersTableColumns: LabDataTableColumn[] = [
    { key: 'avatar', label: '', nowrap: true, widthClass: 'w-12 min-w-12' },
    { key: 'user', label: 'Пользователь' },
    { key: 'status', label: 'Статус' },
    { key: 'role', label: 'Роль' },
    { key: 'login', label: 'Вход', nowrap: true },
    { key: 'actions', label: 'Действия', nowrap: true, widthClass: 'min-w-44' }
  ]
  const keysTableColumns: LabDataTableColumn[] = [
    { key: 'id', label: 'ID', nowrap: true },
    { key: 'key', label: 'Ключ' },
    { key: 'meta', label: 'Meta' },
    { key: 'traits', label: 'Traits', nowrap: true }
  ]
  const moderationTableColumns: LabDataTableColumn[] = [
    { key: 'recipe', label: 'Рецепт', widthClass: 'w-full min-w-96' },
    { key: 'owner', label: 'Владелец', widthClass: 'min-w-32' },
    { key: 'status', label: 'Статус', nowrap: true, widthClass: 'min-w-32' },
    { key: 'created', label: 'Создан', nowrap: true, widthClass: 'min-w-36' },
    { key: 'actions', label: 'Действия', nowrap: true, widthClass: 'min-w-44' }
  ]
  const loadUsers = async () => {
    usersLoading.value = true
    usersError.value = ''
    usersInfo.value = ''
    try {
      const res = await adminListUsers({
        q: userSearch.value.trim(),
        status: userStatusFilter.value.trim(),
        role: userRoleFilter.value.trim(),
        limit: usersLimit.value,
        offset: usersOffset.value
      })
      users.value = res.data.items || []
      usersTotal.value = Number(res.data.total || 0)
      syncRoleDrafts(users.value)
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось загрузить пользователей.'
      users.value = []
      usersTotal.value = 0
    } finally {
      usersLoading.value = false
    }
  }
  const setRole = async (item: AdminUserView) => {
    const role = roleDrafts[item.user_id]
    if (!role) return
    if (item.user_id === user.value?.user_id && role !== 'admin') {
      roleDrafts[item.user_id] = 'admin'
      usersError.value = 'Нельзя изменить собственную роль администратора.'
      return
    }
    if (role === currentRole(item)) return
    if (roleSavingUserId.value === item.user_id) return
    roleSavingUserId.value = item.user_id
    usersError.value = ''
    usersInfo.value = ''
    try {
      await adminSetUserRole(item.user_id, role)
      usersInfo.value = `Роль обновлена: ${item.email} -> ${role}.`
      await loadUsers()
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось обновить роль.'
    } finally {
      if (roleSavingUserId.value === item.user_id) {
        roleSavingUserId.value = null
      }
    }
  }
  const blockUser = async (item: AdminUserView) => {
    if (isAdminAccount(item)) {
      usersError.value = 'Нельзя блокировать пользователя с ролью admin.'
      return
    }
    usersError.value = ''
    usersInfo.value = ''
    try {
      await adminBlockUser(item.user_id, '')
      usersInfo.value = `Пользователь ${item.email} заблокирован.`
      await loadUsers()
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось заблокировать пользователя.'
    }
  }
  const unblockUser = async (item: AdminUserView) => {
    usersError.value = ''
    usersInfo.value = ''
    try {
      await adminUnblockUser(item.user_id)
      usersInfo.value = `Пользователь ${item.email} разблокирован.`
      await loadUsers()
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось разблокировать пользователя.'
    }
  }
  const forceLogoutUser = async (item: AdminUserView) => {
    if (isAdminAccount(item)) {
      usersError.value = 'Нельзя сбрасывать сессии пользователя с ролью admin.'
      return
    }
    usersError.value = ''
    usersInfo.value = ''
    try {
      await adminForceLogoutUser(item.user_id)
      usersInfo.value = `Все сессии пользователя ${item.email} завершены.`
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось принудительно завершить сессии.'
    }
  }
  const deleteUser = async (item: AdminUserView) => {
    if (isAdminAccount(item)) {
      usersError.value = 'Нельзя удалять пользователя с ролью admin.'
      return
    }
    usersError.value = ''
    usersInfo.value = ''
    try {
      await adminDeleteUser(item.user_id)
      usersInfo.value = `Аккаунт ${item.email} удалён.`
      await loadUsers()
    } catch (err: any) {
      usersError.value = err?.data?.message || err?.message || 'Не удалось удалить пользователя.'
    }
  }
  const loadKeys = async (query?: string) => {
    const requestId = ++keysRequestId
    keysLoading.value = true
    keysError.value = ''
    try {
      const res = await adminSearchKeys(String(query ?? keysQuery.value).trim(), 50)
      if (requestId !== keysRequestId) return
      keys.value = res.data.items || []
    } catch (err: any) {
      if (requestId !== keysRequestId) return
      keysError.value = err?.data?.message || err?.message || 'Не удалось выполнить поиск ключей.'
      keys.value = []
    } finally {
      if (requestId !== keysRequestId) return
      keysLoading.value = false
    }
  }
  const debouncedLoadKeys = debounce((query: string) => {
    void loadKeys(query)
  }, 300)
  const keyMetaRecord = (meta?: Record<string, any>) => {
    if (!meta || typeof meta !== 'object' || Array.isArray(meta)) return {}
    return meta
  }
  const keyMetaSummaryItems = (meta?: Record<string, any>) => {
    const safe = keyMetaRecord(meta)
    const items: string[] = []
    const dataType = String(safe.dataType || '').trim()
    const unitCategory = String(safe.unitCategory || '').trim()
    const unit = String(safe.unit || '').trim()
    const mode = String(safe.mode || '').trim()
    const optionType = String(safe.optionType || '').trim()
    const optionsCount = Array.isArray(safe.options) ? safe.options.length : 0
    if (dataType) items.push(`тип: ${dataType}`)
    if (unitCategory) items.push(`категория ед.: ${unitCategory}`)
    if (unit) items.push(`ед. изм.: ${unit}`)
    if (mode) items.push(`режим: ${mode}`)
    if (optionType) items.push(`тип опций: ${optionType}`)
    if (optionsCount > 0) items.push(`опций: ${optionsCount}`)
    return items
  }
  const keyMetaJson = (meta?: Record<string, any>) => {
    const safe = keyMetaRecord(meta)
    if (!Object.keys(safe).length) return '—'
    return JSON.stringify(safe, null, 2)
  }
  const moderationStatusLabel = (value: string) => {
    if (value === 'all') return 'Все'
    if (value === 'pending') return 'На модерации'
    if (value === 'approved') return 'Одобрен'
    if (value === 'rejected') return 'Отклонён'
    if (value === 'draft') return 'Черновик'
    return value || '—'
  }
  const userStatusTone = (value?: string): 'success' | 'warning' | 'danger' | 'muted' => {
    if (value === 'active') return 'success'
    if (value === 'pending_verification') return 'warning'
    if (value === 'blocked') return 'danger'
    return 'muted'
  }
  const moderationStatusTone = (value?: string): 'success' | 'warning' | 'danger' | 'muted' => {
    if (value === 'approved') return 'success'
    if (value === 'pending') return 'warning'
    if (value === 'rejected') return 'danger'
    return 'muted'
  }
  const moderationStatusStatItems = computed(
    () =>
      [
        { key: 'all', label: 'Все' },
        { key: 'pending', label: 'На модерации' },
        { key: 'approved', label: 'Одобрено' },
        { key: 'rejected', label: 'Отклонено' },
        { key: 'draft', label: 'Черновики' }
      ] as const
  )
  const moderationCanApprove = (item: KitchenRecipe) => item.moderation_status !== 'approved'
  const moderationCanReject = (item: KitchenRecipe) =>
    item.moderation_status === 'pending' || item.moderation_status === 'approved'
  const moderationRecipeLink = (item: KitchenRecipe) => `/kitchen/${item.id}`
  const moderationNoteDraft = (item: KitchenRecipe) => {
    const recipeId = String(item.id || '').trim()
    if (!recipeId) return ''
    if (!(recipeId in moderationNoteDrafts)) {
      moderationNoteDrafts[recipeId] = String(item.moderation_note || '').trim()
    }
    return moderationNoteDrafts[recipeId]
  }
  const cacheOwnerUsers = (items: AdminUserView[]) => {
    for (const item of items) {
      const id = String(item.user_id || '').trim()
      if (!id) continue
      moderationOwnerDirectory[id] = {
        email: String(item.email || '').trim(),
        display_name: String(item.display_name || '').trim()
      }
    }
  }
  const moderationOwnerDisplay = (ownerUserId?: string | null) => {
    const id = String(ownerUserId || '').trim()
    if (!id) return '—'
    const cached = moderationOwnerDirectory[id]
    if (!cached) return id
    if (cached.display_name) return `${cached.display_name} · ${cached.email}`
    if (cached.email) return cached.email
    return id
  }
  const moderationOwnerSubLabel = (ownerUserId?: string | null) => {
    const id = String(ownerUserId || '').trim()
    if (!id) return 'без владельца'
    const cached = moderationOwnerDirectory[id]
    if (!cached) return id
    return id
  }
  const isModerationOwnerEditorOpen = (recipeId?: string) =>
    String(recipeId || '').trim() !== '' && moderationOwnerEditorRecipeId.value === recipeId
  const closeModerationOwnerEditor = () => {
    moderationOwnerEditorRecipeId.value = null
    moderationOwnerQuery.value = ''
    moderationOwnerCandidates.value = []
    moderationOwnerCandidatesError.value = ''
  }
  const searchModerationOwnerCandidates = async (queryRaw: string) => {
    if (!moderationOwnerEditorRecipeId.value) return
    const query = String(queryRaw || '').trim()
    if (query.length < 2) {
      moderationOwnerCandidates.value = []
      moderationOwnerCandidatesError.value = ''
      moderationOwnerCandidatesLoading.value = false
      return
    }
    const requestId = ++moderationOwnerSearchRequestId
    moderationOwnerCandidatesLoading.value = true
    moderationOwnerCandidatesError.value = ''
    try {
      const res = await adminListUsers({ q: query, limit: 12, offset: 0 })
      if (requestId !== moderationOwnerSearchRequestId) return
      const items = res.data.items || []
      moderationOwnerCandidates.value = items
      cacheOwnerUsers(items)
    } catch (err: any) {
      if (requestId !== moderationOwnerSearchRequestId) return
      moderationOwnerCandidates.value = []
      moderationOwnerCandidatesError.value = err?.data?.message || err?.message || 'Не удалось найти пользователей.'
    } finally {
      if (requestId !== moderationOwnerSearchRequestId) return
      moderationOwnerCandidatesLoading.value = false
    }
  }
  const debouncedSearchModerationOwnerCandidates = debounce((query: string) => {
    void searchModerationOwnerCandidates(query)
  }, 250)
  const openModerationOwnerEditor = (item: KitchenRecipe) => {
    const recipeId = String(item.id || '').trim()
    if (!recipeId) return
    if (moderationOwnerSavingRecipeId.value === recipeId) return
    if (moderationOwnerEditorRecipeId.value === recipeId) {
      closeModerationOwnerEditor()
      return
    }
    moderationOwnerEditorRecipeId.value = recipeId
    moderationOwnerCandidates.value = []
    moderationOwnerCandidatesError.value = ''
    const ownerId = String(item.owner_user_id || '').trim()
    const cached = ownerId ? moderationOwnerDirectory[ownerId] : null
    moderationOwnerQuery.value = String(cached?.email || cached?.display_name || '').trim()
    if (moderationOwnerQuery.value.length >= 2) {
      debouncedSearchModerationOwnerCandidates(moderationOwnerQuery.value)
    }
  }
  const changeModerationRecipeOwner = async (item: KitchenRecipe, nextOwner: AdminUserView) => {
    const recipeId = String(item.id || '').trim()
    const nextOwnerId = String(nextOwner.user_id || '').trim()
    if (!recipeId || !nextOwnerId) return
    if (moderationOwnerSavingRecipeId.value === recipeId) return
    if (String(item.owner_user_id || '').trim() === nextOwnerId) {
      closeModerationOwnerEditor()
      return
    }
    moderationOwnerSavingRecipeId.value = recipeId
    moderationError.value = ''
    moderationInfo.value = ''
    try {
      const res = await adminChangeKitchenRecipeOwner(recipeId, nextOwnerId)
      moderationItems.value = moderationItems.value.map(row => (row.id === recipeId ? res.data : row))
      cacheOwnerUsers([nextOwner])
      moderationInfo.value = `Владелец рецепта «${item.title}» изменён на ${nextOwner.email}.`
      closeModerationOwnerEditor()
    } catch (err: any) {
      moderationError.value = err?.data?.message || err?.message || 'Не удалось изменить владельца рецепта.'
    } finally {
      if (moderationOwnerSavingRecipeId.value === recipeId) {
        moderationOwnerSavingRecipeId.value = null
      }
    }
  }
  const loadModeration = async () => {
    moderationLoading.value = true
    moderationError.value = ''
    try {
      const res = await adminListKitchenModerationRecipes({
        status: moderationStatus.value,
        limit: moderationLimit.value,
        offset: moderationOffset.value
      })
      moderationItems.value = res.data.items || []
      for (const item of moderationItems.value) {
        const recipeId = String(item.id || '').trim()
        if (!recipeId) continue
        moderationNoteDrafts[recipeId] = String(item.moderation_note || '').trim()
      }
      moderationTotal.value = Number(res.data.total || 0)
      moderationStatusTotals.value = normalizeModerationStatusTotals(res.data.status_totals)
    } catch (err: any) {
      moderationError.value = err?.data?.message || err?.message || 'Не удалось загрузить рецепты для модерации.'
      moderationItems.value = []
      moderationTotal.value = 0
    } finally {
      moderationLoading.value = false
    }
    if (moderationOwnerEditorRecipeId.value) {
      const exists = moderationItems.value.some(item => item.id === moderationOwnerEditorRecipeId.value)
      if (!exists) closeModerationOwnerEditor()
    }
  }
  const moderationPrevPage = async () => {
    if (!moderationHasPrevPage.value) return
    moderationOffset.value = Math.max(0, moderationOffset.value - moderationLimit.value)
    await loadModeration()
  }
  const moderationNextPage = async () => {
    if (!moderationHasNextPage.value) return
    moderationOffset.value = moderationOffset.value + moderationLimit.value
    await loadModeration()
  }
  const moderateRecipe = async (item: KitchenRecipe, approve: boolean) => {
    moderationError.value = ''
    moderationInfo.value = ''
    const note = String(moderationNoteDraft(item) || '').trim()
    if (!approve && !note) {
      moderationError.value = 'Укажите причину отклонения рецепта.'
      return
    }
    try {
      await adminModerateKitchenRecipe(item.id, approve, note)
      moderationInfo.value =
        approve ? `Рецепт «${item.title}» одобрен и опубликован.` : `Рецепт «${item.title}» отклонён.`
      await loadModeration()
    } catch (err: any) {
      moderationError.value = err?.data?.message || err?.message || 'Не удалось изменить статус модерации.'
    }
  }
  const loadAnalysis = async () => {
    analysisLoading.value = true
    analysisError.value = ''
    try {
      const res = await adminTraitsSetsAnalysis()
      analysis.value = res.data || null
    } catch (err: any) {
      analysisError.value = err?.data?.message || err?.message || 'Не удалось загрузить аналитику.'
      analysis.value = null
    } finally {
      analysisLoading.value = false
    }
  }
  const loadAdminSummary = async () => {
    summaryLoading.value = true
    summaryError.value = ''
    try {
      const nextSummary = await loadSharedAdminSummary()
      if (nextSummary?.recipe_status_totals) {
        moderationStatusTotals.value = normalizeModerationStatusTotals(nextSummary.recipe_status_totals)
      }
    } catch (err: any) {
      summaryError.value = err?.data?.message || err?.message || 'Не удалось загрузить сводку.'
    } finally {
      summaryLoading.value = false
    }
  }
  const loadPaymentsSummary = async () => {
    paymentsSummaryLoading.value = true
    paymentsSummaryError.value = ''
    try {
      const res = await loadPaymentsAdminSummary()
      paymentsSummary.value = res.data || null
    } catch (err: any) {
      paymentsSummaryError.value = err?.data?.message || err?.message || 'Не удалось загрузить payments summary.'
      paymentsSummary.value = null
    } finally {
      paymentsSummaryLoading.value = false
    }
  }
  const markSummaryRead = async () => {
    if (summaryReadPending.value) return
    summaryReadPending.value = true
    summaryError.value = ''
    summaryInfo.value = ''
    try {
      await adminMarkSummaryRead()
      summaryInfo.value = 'Уведомления отмечены как прочитанные.'
    } catch (err: any) {
      summaryError.value = err?.data?.message || err?.message || 'Не удалось отметить уведомления как прочитанные.'
    } finally {
      summaryReadPending.value = false
    }
  }
  const prevPage = async () => {
    if (!hasPrevPage.value) return
    usersOffset.value = Math.max(0, usersOffset.value - usersLimit.value)
    await loadUsers()
  }
  const nextPage = async () => {
    if (!hasNextPage.value) return
    usersOffset.value = usersOffset.value + usersLimit.value
    await loadUsers()
  }
  const applyUserFilters = async () => {
    usersOffset.value = 0
    await loadUsers()
  }
  const debouncedApplyUserFilters = debounce(() => {
    void applyUserFilters()
  }, 300)
  watch(
    adminTab,
    async tab => {
      if (tab !== 'moderation') {
        closeModerationOwnerEditor()
      }
      if (tab === 'users') {
        if (!users.value.length && !usersLoading.value) await loadUsers()
        return
      }
      if (tab === 'moderation') {
        if (!moderationItems.value.length && !moderationLoading.value) await loadModeration()
        return
      }
      if (tab === 'keys') {
        if (!keys.value.length && !keysLoading.value) await loadKeys()
        return
      }
      if (!analysis.value && !analysisLoading.value) {
        await loadAnalysis()
      }
      if (!summary.value && !summaryLoading.value) {
        await loadAdminSummary()
      }
      if (!paymentsSummary.value && !paymentsSummaryLoading.value) {
        await loadPaymentsSummary()
      }
    },
    { immediate: true }
  )
  watch(moderationStatus, async (next, prev) => {
    if (next === prev) return
    moderationOffset.value = 0
    closeModerationOwnerEditor()
    if (adminTab.value === 'moderation') {
      await loadModeration()
    }
  })
  watch(moderationOwnerQuery, next => {
    if (!moderationOwnerEditorRecipeId.value) return
    debouncedSearchModerationOwnerCandidates(String(next || ''))
  })
  watch(keysQuery, next => {
    if (adminTab.value !== 'keys') return
    debouncedLoadKeys(String(next || ''))
  })
  watch([userSearch, userStatusFilter, userRoleFilter], () => {
    if (adminTab.value !== 'users') return
    debouncedApplyUserFilters()
  })
  watch(users, items => {
    cacheOwnerUsers(items || [])
  })
  watch(
    sharedAdminSummary,
    nextSummary => {
      if (!nextSummary?.recipe_status_totals) return
      moderationStatusTotals.value = normalizeModerationStatusTotals(nextSummary.recipe_status_totals)
    },
    { immediate: true }
  )
  onMounted(() => {
    if (!summary.value && !summaryLoading.value) {
      void loadAdminSummary()
    }
  })
</script>
<template>
  <div>
    <LabNavHeader :title :breadcrumb-items="adminBreadcrumbItems">
      <template #actions>
        <LabBaseBadge v-if="summary?.has_unread" variant="warning" size="xs" :rounded="false">
          новое
          {{
            Number(summary?.new_users_since_last_login || 0) +
            Number(summary?.new_pending_recipes_since_last_login || 0)
          }}
        </LabBaseBadge>
      </template>
    </LabNavHeader>
    <LabNavTabs
      v-model="adminTab"
      :items="adminTabItems"
      :no-select="true"
      :render-panels="false"
      route-query-key="tab"
      route-default-value="users" />
    <div class="mt-4 space-y-4 px-3 md:px-4">
      <section v-show="adminTab === 'users'" :class="adminSectionClass">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="space-y-1">
            <h2 class="text-base font-semibold text-zinc-100">Пользователи и роли</h2>
            <p :class="adminMetaTextClass">Поиск, блокировки, смена роли и принудительный выход из сессий.</p>
          </div>
        </div>
        <div class="grid gap-3 md:grid-cols-3 xl:grid-cols-4">
          <LabField label="Поиск" for-id="admin_user_search">
            <LabBaseInput
              id="admin_user_search"
              v-model="userSearch"
              name="admin_user_search"
              type="text"
              placeholder="email или имя" />
          </LabField>
          <LabField label="Статус" for-id="admin_user_status">
            <LabBaseSelect
              id="admin_user_status"
              v-model="userStatusFilter"
              name="admin_user_status"
              :options="userStatusOptions" />
          </LabField>
          <LabField label="Роль" for-id="admin_user_role_filter">
            <LabBaseSelect
              id="admin_user_role_filter"
              v-model="userRoleFilter"
              name="admin_user_role_filter"
              :options="userRoleOptions" />
          </LabField>
        </div>
        <LabNotify :text="usersError" tone="error" size="xs" />
        <LabNotify :text="usersInfo" tone="success" size="xs" />
        <LabDataTable
          :columns="usersTableColumns"
          :rows="users"
          :loading="usersLoading"
          empty-text="Пользователи не найдены."
          :row-key="(row: AdminUserView) => row.user_id">
          <template #cell-avatar="{ row }">
            <NuxtLink :to="`/users/${row.user_id}`" class="inline-flex">
              <LabAvatar version="preview" :user="adminListUserToAuthUser(row)" :show-label="false" />
            </NuxtLink>
          </template>
          <template #cell-user="{ row }">
            <div class="space-y-1">
              <NuxtLink :to="`/users/${row.user_id}`" class="text-(--lab-text-primary) hover:text-(--lab-accent)">
                {{ row.email }}
              </NuxtLink>
              <p class="text-zinc-500">{{ row.display_name || '—' }}</p>
            </div>
          </template>
          <template #cell-status="{ row }">
            <div class="space-y-1">
              <LabBaseBadge :variant="userStatusTone(row.status)" size="xs" :rounded="false">
                {{ userStatusLabel(row.status) }}
              </LabBaseBadge>
              <p v-if="row.blocked_reason" :class="adminMetaTextClass">причина: {{ row.blocked_reason }}</p>
            </div>
          </template>
          <template #cell-role="{ row }">
            <div class="flex items-center gap-2">
              <LabBaseSelect
                v-model="roleDrafts[row.user_id]"
                :name="`admin_user_role_${row.user_id}`"
                :options="editableRoleOptions"
                class="min-w-28"
                :disabled="roleSavingUserId === row.user_id || row.user_id === user?.user_id"
                @change="setRole(row)" />
              <span v-if="roleSavingUserId === row.user_id" :class="adminInlineHintClass">сохранение...</span>
              <span v-else-if="row.user_id === user?.user_id" :class="adminInlineHintClass">текущий админ</span>
            </div>
          </template>
          <template #cell-login="{ row }">
            <LabRelativeTime :datetime="row.last_login_at" compact />
          </template>
          <template #cell-actions="{ row }">
            <div class="flex min-w-max items-center gap-1 whitespace-nowrap">
              <span v-if="row.status !== 'blocked' && !isAdminAccount(row)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Блок"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить блокировку аккаунта?"
                  icon="ic:round-block"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_12%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_18%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="blockUser(row)" />
              </span>
              <span v-if="row.status !== 'blocked' && !isAdminAccount(row)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-block"
                  aria-label="Заблокировать"
                  confirm-aria-label="Подтвердить блокировку"
                  confirm-label="Ок"
                  tooltip="Подтвердить блокировку аккаунта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_12%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_18%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="blockUser(row)" />
              </span>
              <LabBaseButton
                v-if="row.status === 'blocked'"
                class="max-sm:hidden"
                variant="secondary"
                size="xs"
                icon="ic:round-lock-open"
                label="Разблок"
                :button-class="`${adminCompactButtonClass} border-[color-mix(in_srgb,var(--lab-info)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-info)_10%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-info)_16%,var(--lab-bg-surface-hover))]`"
                @click="unblockUser(row)" />
              <LabBaseButton
                v-if="row.status === 'blocked'"
                class="sm:hidden"
                variant="secondary"
                size="xs"
                icon="ic:round-lock-open"
                icon-only
                aria-label="Разблокировать"
                :button-class="`${adminCompactButtonClass} border-[color-mix(in_srgb,var(--lab-info)_42%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-info)_10%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-info)_16%,var(--lab-bg-surface-hover))]`"
                @click="unblockUser(row)" />
              <span v-if="!isAdminAccount(row)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Сброс сессий"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить принудительный сброс всех сессий?"
                  icon="ic:round-logout"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-primary) hover:bg-(--lab-bg-surface-hover)"
                  confirm-class="border-(--lab-warning) bg-(--lab-warning) text-white hover:bg-[color-mix(in_srgb,var(--lab-warning)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-warning)_30%,transparent)]"
                  tooltip-class="border-[color-mix(in_srgb,var(--lab-warning)_52%,var(--lab-border))] text-(--lab-text-primary)"
                  @confirm="forceLogoutUser(row)" />
              </span>
              <span v-if="!isAdminAccount(row)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-logout"
                  aria-label="Сбросить сессии"
                  confirm-aria-label="Подтвердить сброс сессий"
                  confirm-label="Ок"
                  tooltip="Подтвердить принудительный сброс всех сессий?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-(--lab-border) bg-(--lab-bg-surface) text-(--lab-text-primary) hover:bg-(--lab-bg-surface-hover)"
                  confirm-class="border-(--lab-warning) bg-(--lab-warning) text-white hover:bg-[color-mix(in_srgb,var(--lab-warning)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-warning)_30%,transparent)]"
                  tooltip-class="border-[color-mix(in_srgb,var(--lab-warning)_52%,var(--lab-border))] text-(--lab-text-primary)"
                  @confirm="forceLogoutUser(row)" />
              </span>
              <span v-if="!isAdminAccount(row)" class="max-sm:hidden">
                <LabConfirmActionButton
                  label="Удалить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить удаление аккаунта?"
                  icon="ic:round-delete"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_46%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_20%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="deleteUser(row)" />
              </span>
              <span v-if="!isAdminAccount(row)" class="sm:hidden">
                <LabConfirmActionButton
                  icon-only
                  icon="ic:round-delete"
                  aria-label="Удалить аккаунт"
                  confirm-aria-label="Подтвердить удаление аккаунта"
                  confirm-label="Ок"
                  tooltip="Подтвердить удаление аккаунта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border-[color-mix(in_srgb,var(--lab-danger)_46%,var(--lab-border))] bg-[color-mix(in_srgb,var(--lab-danger)_14%,var(--lab-bg-surface))] text-(--lab-text-primary) hover:bg-[color-mix(in_srgb,var(--lab-danger)_20%,var(--lab-bg-surface-hover))]"
                  confirm-class="border-(--lab-danger) bg-(--lab-danger) text-white hover:bg-[color-mix(in_srgb,var(--lab-danger)_88%,black)]"
                  progress-class="bg-[color-mix(in_srgb,var(--lab-danger)_30%,transparent)]"
                  @confirm="deleteUser(row)" />
              </span>
            </div>
          </template>
        </LabDataTable>
        <div class="flex items-center gap-2">
          <LabBaseButton
            label="Назад"
            variant="secondary"
            size="xs"
            :disabled="!hasPrevPage || usersLoading"
            @click="prevPage" />
          <LabBaseButton
            label="Вперёд"
            variant="secondary"
            size="xs"
            :disabled="!hasNextPage || usersLoading"
            @click="nextPage" />
        </div>
      </section>
      <section v-show="adminTab === 'moderation'" :class="adminSectionClass">
        <div class="space-y-1">
          <div class="space-y-1">
            <h2 class="text-base font-semibold text-zinc-100">Модерация рецептов</h2>
            <p :class="adminMetaTextClass">Фильтрация по статусу, смена владельца и решение по публикации.</p>
          </div>
        </div>
        <div class="flex flex-wrap gap-2">
          <LabBaseButton
            v-for="item in moderationStatusStatItems"
            :key="`moderation-stat:${item.key}`"
            :variant="moderationStatus === item.key ? 'primary' : 'secondary'"
            size="xs"
            @click="moderationStatus = item.key">
            {{ item.label }} {{ displayNumber(moderationStatusTotals[item.key]) }}
          </LabBaseButton>
        </div>
        <LabNotify :text="moderationError" tone="error" size="xs" />
        <LabNotify :text="moderationInfo" tone="success" size="xs" />
        <LabDataTable
          :columns="moderationTableColumns"
          :rows="moderationItems"
          :loading="moderationLoading"
          empty-text="Рецепты для выбранного статуса не найдены."
          :row-key="(row: KitchenRecipe) => row.id">
          <template #cell-recipe="{ row }">
            <div class="min-w-0">
              <NuxtLink
                :to="moderationRecipeLink(row)"
                class="inline-flex max-w-full items-center gap-1 text-zinc-200 hover:text-zinc-100 hover:underline">
                <span class="truncate">{{ row.title }}</span>
                <Icon name="ic:round-link" class="h-3.5 w-3.5 shrink-0 text-zinc-500" />
              </NuxtLink>
              <p class="text-zinc-500 line-clamp-2">{{ row.description || '—' }}</p>
            </div>
          </template>
          <template #cell-owner="{ row }">
            <div class="min-w-0 space-y-1">
              <LabBaseButton
                variant="ghost"
                size="none"
                button-class="truncate inline-flex max-w-full items-start justify-start px-0 text-left text-zinc-300 hover:text-zinc-100"
                :disabled="moderationOwnerSavingRecipeId === row.id"
                @click="openModerationOwnerEditor(row)"
                :label="moderationOwnerDisplay(row.owner_user_id)" />
              <p class="truncate text-xs text-zinc-500">
                {{ moderationOwnerSubLabel(row.owner_user_id) }}
              </p>
              <div v-if="isModerationOwnerEditorOpen(row.id)" :class="adminSubsectionClass">
                <LabBaseInput
                  v-model="moderationOwnerQuery"
                  :name="`moderation_owner_query_${row.id}`"
                  type="text"
                  placeholder="Поиск: email или имя" />
                <LabNotify :text="moderationOwnerCandidatesError" tone="error" size="xs" />
                <p v-if="moderationOwnerCandidatesLoading" :class="adminInlineHintClass">Поиск...</p>
                <div v-else-if="moderationOwnerCandidates.length" class="max-h-36 space-y-1 overflow-y-auto pr-1">
                  <LabBaseButton
                    v-for="candidate in moderationOwnerCandidates"
                    :key="`owner-candidate:${row.id}:${candidate.user_id}`"
                    variant="secondary"
                    size="none"
                    button-class="w-full justify-start px-2 py-1 text-left text-xs truncate"
                    :disabled="moderationOwnerSavingRecipeId === row.id"
                    @click="changeModerationRecipeOwner(row, candidate)"
                    :label="candidate.display_name || candidate.email" />
                </div>
                <p v-else-if="moderationOwnerQuery.trim().length >= 2" :class="adminInlineHintClass">
                  Пользователи не найдены.
                </p>
                <p v-else :class="adminInlineHintClass">Введите минимум 2 символа для поиска.</p>
                <LabBaseButton
                  variant="ghost"
                  size="xs"
                  button-class="px-0 text-zinc-400 hover:text-zinc-200"
                  @click="closeModerationOwnerEditor">
                  Отмена
                </LabBaseButton>
              </div>
            </div>
          </template>
          <template #cell-status="{ row }">
            <LabBaseBadge :variant="moderationStatusTone(row.moderation_status)" size="xs" :rounded="false">
              {{ moderationStatusLabel(row.moderation_status) }}
            </LabBaseBadge>
          </template>
          <template #cell-created="{ row }">
            <LabRelativeTime :datetime="row.created_at" compact />
          </template>
          <template #cell-actions="{ row }">
            <div class="min-w-64 space-y-2">
              <LabBaseTextarea
                v-if="moderationCanReject(row)"
                :model-value="moderationNoteDraft(row)"
                :name="`moderation_note_${row.id}`"
                rows="2"
                placeholder="Причина отклонения для автора"
                class="w-full text-xs"
                @update:model-value="moderationNoteDrafts[row.id] = String($event || '')" />
              <p v-else-if="row.moderation_note" :class="adminInlineHintClass">
                Последняя причина: {{ row.moderation_note }}
              </p>
              <div class="flex flex-wrap gap-1">
                <LabConfirmActionButton
                  v-if="moderationCanApprove(row)"
                  label="Одобрить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить публикацию рецепта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border border-emerald-500/40 bg-emerald-500/10 text-emerald-100 hover:bg-emerald-500/20"
                  confirm-class="border border-emerald-300/90 bg-emerald-600 text-white hover:bg-emerald-500"
                  progress-class="bg-emerald-300/45"
                  @confirm="moderateRecipe(row, true)" />
                <LabConfirmActionButton
                  v-if="moderationCanReject(row)"
                  label="Отклонить"
                  confirm-label="Подтвердить"
                  tooltip="Подтвердить отклонение рецепта?"
                  :button-class="adminConfirmButtonClass"
                  idle-class="border border-rose-500/50 bg-rose-500/10 text-rose-200 hover:bg-rose-500/20"
                  confirm-class="border border-rose-300/90 bg-rose-600 text-white hover:bg-rose-500"
                  progress-class="bg-rose-300/45"
                  @confirm="moderateRecipe(row, false)" />
              </div>
            </div>
          </template>
        </LabDataTable>
        <div class="flex items-center gap-2">
          <LabBaseButton
            variant="secondary"
            size="xs"
            label="Назад"
            :disabled="!moderationHasPrevPage || moderationLoading"
            @click="moderationPrevPage" />
          <LabBaseButton
            variant="secondary"
            size="xs"
            label="Вперёд"
            :disabled="!moderationHasNextPage || moderationLoading"
            @click="moderationNextPage" />
        </div>
      </section>
      <section v-show="adminTab === 'keys'" :class="adminSectionClass">
        <div class="space-y-1">
          <h2 class="text-base font-semibold text-zinc-100">Поиск по ключам особенностей</h2>
          <p :class="adminMetaTextClass">Проверка syn, meta и количества traits для ключа.</p>
        </div>
        <LabField label="Запрос" for-id="admin_keys_query" class="max-w-xl">
          <LabBaseInput
            id="admin_keys_query"
            v-model="keysQuery"
            name="admin_keys_query"
            type="text"
            placeholder="часть названия ключа" />
        </LabField>
        <LabNotify :text="keysError" tone="error" size="xs" />
        <LabDataTable
          :columns="keysTableColumns"
          :rows="keys"
          :loading="keysLoading"
          empty-text="Ничего не найдено."
          :row-key="(row: AdminTraitKeySearchItem) => row.key_id">
          <template #cell-id="{ row }">
            <span class="text-zinc-400">{{ row.key_id }}</span>
          </template>
          <template #cell-key="{ row }">
            <span class="text-zinc-200">{{ row.syn }}</span>
          </template>
          <template #cell-meta="{ row }">
            <div class="space-y-1">
              <div class="flex flex-wrap gap-1">
                <LabBaseBadge
                  v-for="label in keyMetaSummaryItems(row.meta)"
                  :key="`${row.key_id}:${label}`"
                  variant="muted"
                  size="xs"
                  :rounded="false">
                  {{ label }}
                </LabBaseBadge>
                <LabBaseBadge v-if="!keyMetaSummaryItems(row.meta).length" variant="muted" size="xs" :rounded="false">
                  —
                </LabBaseBadge>
              </div>
              <pre class="max-w-80 overflow-x-auto p-2 text-xs leading-5 text-zinc-400">{{
                keyMetaJson(row.meta)
              }}</pre>
            </div>
          </template>
          <template #cell-traits="{ row }">
            <span class="text-zinc-300">{{ row.trait_count }}</span>
          </template>
        </LabDataTable>
      </section>
      <section v-show="adminTab === 'analysis'" :class="adminSectionClass">
        <div class="space-y-1">
          <h2 class="text-base font-semibold text-zinc-100">Аналитика</h2>
          <p :class="adminMetaTextClass">Сводка по уведомлениям, coverage traits и ключам с максимальной нагрузкой.</p>
        </div>
        <LabNotify :text="analysisError" tone="error" size="xs" />
        <LabNotify :text="summaryError" tone="error" size="xs" />
        <LabNotify :text="paymentsSummaryError" tone="error" size="xs" />
        <LabNotify :text="summaryInfo" tone="success" size="xs" />
        <p v-if="analysisLoading || summaryLoading || paymentsSummaryLoading" class="text-xs text-zinc-400">
          Обновление данных...
        </p>
        <div :class="adminSubsectionClass">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <p class="text-zinc-200">Сводка уведомлений</p>
            <LabBaseButton
              variant="secondary"
              size="xs"
              label="Прочитано"
              :disabled="summaryReadPending || !summary?.has_unread"
              @click="markSummaryRead" />
          </div>
          <div class="grid gap-2 sm:grid-cols-2">
            <div class="space-y-1 p-2.5">
              <p class="text-zinc-500">Новых регистраций</p>
              <p class="text-zinc-100">{{ displayNumber(summary?.new_users_since_last_login) }}</p>
            </div>
            <div class="space-y-1 p-2.5">
              <p class="text-zinc-500">Новых рецептов на модерации</p>
              <p class="text-zinc-100">{{ displayNumber(summary?.new_pending_recipes_since_last_login) }}</p>
            </div>
          </div>
        </div>
        <div :class="adminSubsectionClass">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="space-y-1">
              <p class="text-zinc-200">Payments</p>
              <p :class="adminMetaTextClass">Сводка заказов, выручки, активного доступа и доли поддержки проекта.</p>
            </div>
            <p :class="adminMetaTextClass">
              Последний успешный платёж:
              <span class="text-zinc-300">{{ displayDateTime(paymentsSummary?.last_successful_paid) }}</span>
            </p>
          </div>
          <div class="grid gap-3 text-xs sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            <div v-for="item in paymentsStats" :key="item.key" :class="adminStatCardClass">
              <p class="text-zinc-500">{{ item.label }}</p>
              <p class="text-zinc-100">{{ item.value }}</p>
            </div>
          </div>
        </div>
        <div class="grid gap-3 text-xs sm:grid-cols-2 lg:grid-cols-3">
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Traits</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.total_traits) }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Уникальные пары (key+value)</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.unique_trait_pairs) }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Уникальные ключи</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.unique_trait_keys) }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Sets</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.total_sets) }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Уникальность наборов</p>
            <p class="text-zinc-100">{{ analysis ? formatPercent(analysis.set_uniqueness_rate) : '—' }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Исхожесть (доля производных наборов)</p>
            <p class="text-zinc-100">{{ analysis ? formatPercent(analysis.derived_set_rate) : '—' }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Traits в наборах</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.traits_referenced_in_sets) }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Покрытие traits наборами</p>
            <p class="text-zinc-100">{{ analysis ? formatPercent(analysis.trait_coverage_in_sets_rate) : '—' }}</p>
          </div>
          <div :class="adminStatCardClass">
            <p class="text-zinc-500">Traits вне наборов</p>
            <p class="text-zinc-100">{{ displayNumber(analysis?.orphan_traits) }}</p>
          </div>
        </div>
        <div :class="adminSubsectionClass">
          <p class="text-zinc-400">Топ ключей по количеству traits:</p>
          <ul v-if="analysis?.top_keys?.length" class="space-y-1">
            <li
              v-for="item in analysis?.top_keys || []"
              :key="item.key_id"
              class="flex items-center justify-between gap-2 px-2 py-1.5 text-zinc-300">
              <span class="truncate">{{ item.syn }} (id {{ item.key_id }})</span>
              <span class="text-zinc-100">{{ item.trait_count }}</span>
            </li>
          </ul>
          <p v-else class="text-zinc-500">Данные пока не получены.</p>
        </div>
      </section>
    </div>
  </div>
</template>
