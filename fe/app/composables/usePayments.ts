type PaymentsApiResult<T> = ApiResponseWithData<T>

let paymentAccessInFlight: Promise<PaymentAccessSummary | null> | null = null

export const usePayments = () => {
  const { json: useApiJson } = useAPI()
  const { isAuthenticated, isAdmin } = useAuth()
  const access = useState<PaymentAccessSummary | null>('payments-access', () => null)
  const accessLoading = useState<boolean>('payments-access-loading', () => false)
  const hasActiveAccess = computed(() => Boolean(access.value?.has_active_access))

  const clearAccess = () => {
    access.value = null
    accessLoading.value = false
    paymentAccessInFlight = null
  }

  const loadAccess = async () => {
    if (!isAuthenticated.value) {
      clearAccess()
      return null
    }
    if (accessLoading.value && paymentAccessInFlight) return await paymentAccessInFlight
    accessLoading.value = true
    paymentAccessInFlight = (async () => {
      try {
        const res = await useApiJson<PaymentsApiResult<PaymentAccessSummary>>('/payments/access', {
          method: 'GET'
        })
        access.value = res.data || null
        return access.value
      } catch {
        access.value = null
        return null
      } finally {
        accessLoading.value = false
        paymentAccessInFlight = null
      }
    })()
    return await paymentAccessInFlight
  }

  const ensureAccessLoaded = async () => {
    if (access.value || !isAuthenticated.value) return access.value
    return await loadAccess()
  }

  const createOrder = async (input: {
    plan_code: 'pro' | 'donation'
    amount?: number
    return_to?: string
  }) => {
    return await useApiJson<PaymentsApiResult<PaymentCreateOrderResult>>('/payments/init', {
      method: 'POST',
      body: {
        plan_code: input.plan_code,
        amount: Number(input.amount || 0),
        return_to: input.return_to || ''
      }
    })
  }

  const lookupPublicOrder = async (input: {
    order_id: string
    token: string
    sync_state?: boolean
  }) => {
    return await useApiJson<PaymentsApiResult<PaymentPublicLookupResult>>('/payments/lookup', {
      method: 'POST',
      auth: {
        allowAutoRefresh: false
      },
      body: {
        order_id: input.order_id,
        token: input.token,
        sync_state: Boolean(input.sync_state)
      }
    })
  }

  const adminListOrders = async (params?: {
    q?: string
    status?: string
    plan?: string
    limit?: number
    offset?: number
  }) => {
    const query = new URLSearchParams()
    if (params?.q) query.set('q', String(params.q))
    if (params?.status) query.set('status', String(params.status))
    if (params?.plan) query.set('plan', String(params.plan))
    if (typeof params?.limit === 'number') query.set('limit', String(params.limit))
    if (typeof params?.offset === 'number') query.set('offset', String(params.offset))
    const suffix = query.toString()
    return await useApiJson<
      PaymentsApiResult<{
        items: PaymentOrderView[]
        total: number
        limit: number
        offset: number
      }>
    >(`/auth/admin/orders${suffix ? `?${suffix}` : ''}`, {
      method: 'GET'
    })
  }

  const adminSummary = async () => {
    return await useApiJson<PaymentsApiResult<PaymentAdminOrdersSummary>>('/auth/admin/payments/summary', {
      method: 'GET'
    })
  }

  watch(
    isAuthenticated,
    authed => {
      if (!authed) {
        clearAccess()
      }
    },
    { immediate: true }
  )

  return {
    access,
    accessLoading,
    hasActiveAccess,
    isAdmin,
    clearAccess,
    loadAccess,
    ensureAccessLoaded,
    createOrder,
    lookupPublicOrder,
    adminListOrders,
    adminSummary
  }
}
