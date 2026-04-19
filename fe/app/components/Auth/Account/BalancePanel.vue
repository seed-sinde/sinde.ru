<script setup lang="ts">
type PaymentHistoryRow = LabDataTableRow & {
  id: string
  orderId: string
  createdAt: string
  plan: string
  amount: string
  status: string
  access: string
  refund: string
  canRefund: boolean
}

const { locale, key, load, t } = useI18nSection('payments')
await useAsyncData(key.value, load, { watch: [locale] })
const { formatAbsoluteDateTime } = useLocalizedDateTime()
const { access, accessLoading, ensureAccessLoaded, history, historyLoading, loadHistory, refundOrder, loadAccess } =
  usePayments()

const refundPendingOrderId = ref('')
const paymentHistoryError = ref('')
const paymentHistoryInfo = ref('')

const paymentAccessUntilText = computed(() => {
  if (!access.value?.access_until) return ''
  return formatAbsoluteDateTime(access.value.access_until)
})
const balanceSubscriptionStatusLabel = computed(() =>
  access.value?.has_active_access ? `Подписка активна до ${paymentAccessUntilText.value || '—'}` : 'Подписка не активна'
)
const paymentHistoryColumns = computed<LabDataTableColumn[]>(() => [
  { key: 'createdAt', label: 'Дата', nowrap: true },
  { key: 'plan', label: 'План' },
  { key: 'amount', label: 'Сумма', nowrap: true },
  { key: 'status', label: 'Статус', nowrap: true },
  { key: 'access', label: 'Доступ' },
  { key: 'refund', label: 'Возврат', nowrap: true }
])
const formatDateTime = (value?: string | null) =>
  formatAbsoluteDateTime(value, { dateStyle: 'medium', timeStyle: 'short' })
const paymentStatusLabel = (status?: string | null) => {
  switch (String(status || '').trim()) {
    case 'success':
      return t('status.success')
    case 'pending':
      return t('status.pending')
    case 'failed':
      return t('status.failed')
    case 'canceled':
      return t('status.canceled')
    case 'refunded':
      return t('status.refunded')
    default:
      return t('status.unknown')
  }
}
const paymentPlanLabel = (planCode?: string | null) =>
  String(planCode || '').trim() === 'donation' ? t('plan.donation') : t('plan.pro')
const paymentHistoryRows = computed<{
  id: string
  orderId: string
  createdAt: string
  plan: string
  amount: string
  status: string
  access: string
  refund: string
  canRefund: boolean
}[]>(() =>
  history.value.map(item => ({
    id: item.order_id,
    orderId: item.order_id,
    createdAt: formatDateTime(item.created_at),
    plan: paymentPlanLabel(item.plan_code),
    amount: formatPaymentWholeRubles(item.amount),
    status: paymentStatusLabel(item.status),
    access: item.access_until
      ? `До ${formatDateTime(item.access_until)}`
      : item.access_from
        ? `С ${formatDateTime(item.access_from)}`
        : '—',
    refund: item.can_refund ? 'Доступен' : item.refunded_at ? 'Выполнен' : 'Недоступен',
    canRefund: item.can_refund
  }))
)

const loadPaymentHistoryState = async (force = false) => {
  paymentHistoryError.value = ''
  try {
    await loadHistory(force)
  } catch (err) {
    paymentHistoryError.value = extractApiErrorMessage(err, 'Не удалось загрузить историю оплат.')
  }
}
const submitRefund = async (orderID: string) => {
  if (!orderID || refundPendingOrderId.value) return
  paymentHistoryError.value = ''
  paymentHistoryInfo.value = ''
  refundPendingOrderId.value = orderID
  try {
    await refundOrder(orderID)
    await Promise.all([loadAccess(true), loadPaymentHistoryState(true)])
    paymentHistoryInfo.value = 'Возврат отправлен в платёжную систему.'
  } catch (err) {
    paymentHistoryError.value = extractApiErrorMessage(err, 'Не удалось выполнить возврат.')
  } finally {
    refundPendingOrderId.value = ''
  }
}

await ensureAccessLoaded()
await loadPaymentHistoryState()
</script>

<template>
  <section class="space-y-3 text-(--lab-text-primary)">
    <div class="space-y-1">
      <h2 class="text-base font-medium text-(--lab-text-primary)">Статус подписки</h2>
      <p class="text-sm text-(--lab-text-muted)">
        {{ accessLoading ? 'Загрузка сведений о доступе…' : balanceSubscriptionStatusLabel }}
      </p>
    </div>
    <div>
      <NuxtLink
        to="/payments"
        class="inline-flex min-h-11 items-center justify-center border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary)"
      >
        Пополнить баланс
      </NuxtLink>
    </div>
    <section class="space-y-4 pt-3">
      <h2 class="text-base font-medium text-(--lab-text-primary)">История транзакций</h2>

      <LabNotify :text="paymentHistoryError" tone="error" size="xs" />
      <LabNotify :text="paymentHistoryInfo" tone="success" size="xs" />

      <LabDataTable
        :columns="paymentHistoryColumns"
        :rows="paymentHistoryRows"
        :loading="historyLoading"
        empty-text="У вас пока нет завершённых или созданных платёжных операций."
        row-key="orderId"
      >
        <template #cell-refund="{ row }">
          <LabConfirmActionButton
            v-if="(row as PaymentHistoryRow).canRefund"
            icon="ic:round-undo"
            confirm-icon="ic:round-check"
            label="Оформить возврат"
            confirm-label="Подтвердить"
            tooltip="Вернуть этот платёж?"
            :disabled="refundPendingOrderId !== ''"
            @confirm="submitRefund((row as PaymentHistoryRow).orderId)"
          />
          <span v-else>{{ (row as PaymentHistoryRow).refund }}</span>
        </template>
      </LabDataTable>
    </section>
  </section>
</template>
