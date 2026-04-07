<script setup lang="ts">
  const title = 'Оплата'

  usePageSeo({
    title,
    description: 'Получение доступа к платным данным и инструментам на один месяц.'
  })

  const route = useRoute()
  const router = useRouter()

  const { isAuthenticated, ensureLoaded } = useAuth()
  const { access, accessLoading, ensureAccessLoaded, createOrder } = usePayments()
  const { formatAbsoluteDateTime } = useLocalizedDateTime()

  await ensureLoaded()
  await ensureAccessLoaded()

  const proAmount = 39900
  const donationAmount = ref<number>(proAmount)

  const creatingPlan = ref<'' | 'pro' | 'donation'>('')
  const formError = ref('')

  const currentPath = computed(() => {
    const query = new URLSearchParams()
    for (const [key, value] of Object.entries(route.query)) {
      if (Array.isArray(value)) {
        for (const item of value) query.append(key, String(item))
        continue
      }
      if (value != null) query.set(key, String(value))
    }
    const suffix = query.toString()
    return `${route.path}${suffix ? `?${suffix}` : ''}`
  })

  const activeUntilText = computed(() => {
    if (!access.value?.access_until) return ''
    return formatAbsoluteDateTime(access.value.access_until)
  })

  const latestOrder = computed(() => access.value?.latest_order || null)

  const latestOrderStatusLabel = computed(() => {
    const status = String(latestOrder.value?.status || '').trim()
    switch (status) {
      case 'success':
        return 'Оплачен'
      case 'pending':
        return 'Ожидает подтверждения'
      case 'failed':
        return 'Неуспешен'
      case 'canceled':
        return 'Отменён'
      case 'refunded':
        return 'Возвращён'
      default:
        return '—'
    }
  })

  const latestOrderCreatedAtText = computed(() => {
    if (!latestOrder.value?.created_at) return ''
    return formatAbsoluteDateTime(latestOrder.value.created_at)
  })

  const latestOrderAmountText = computed(() => formatPrice(latestOrder.value?.amount || 0))
  const donationHintText = computed(() => {
    const tip = Math.max(0, Number(donationAmount.value || 0) - proAmount)
    return tip > 0 ? `Поддержка проекта: ${formatPrice(tip)}` : 'Без дополнительной поддержки'
  })

  function formatPrice(amount: number) {
    const rub = Math.floor(Number(amount || 0) / 100)
    return new Intl.NumberFormat('ru-RU').format(rub) + ' ₽'
  }

  function normalizeDonationAmount(value: number) {
    if (!Number.isFinite(value)) return proAmount
    return Math.max(proAmount, Math.round(value))
  }

  async function submitOrder(plan: 'pro' | 'donation') {
    formError.value = ''

    if (!isAuthenticated.value) {
      await router.push('/auth/login')
      return
    }

    creatingPlan.value = plan

    try {
      const amount = plan === 'donation' ? normalizeDonationAmount(donationAmount.value) : proAmount

      const res = await createOrder({
        plan_code: plan,
        amount,
        return_to: currentPath.value
      })

      const paymentURL = String(res?.data?.payment_url || '').trim()
      if (!paymentURL) {
        throw new Error('Платёжная ссылка не была получена')
      }

      await navigateTo(paymentURL, { external: true })
    } catch (error: any) {
      formError.value =
        String(error?.data?.message || error?.message || '').trim() || 'Не удалось создать платёжный заказ'
    } finally {
      creatingPlan.value = ''
    }
  }
</script>

<template>
  <div class="space-y-8">
    <LabNavHeader :title />

    <section class="space-y-4 rounded-3xl border border-(--lab-border) bg-(--lab-bg-elevated) p-5">
      <div class="space-y-2">
        <h1 class="text-lg font-medium text-(--lab-text-primary)">Доступ к платным данным и инструментам</h1>
        <p class="text-sm text-(--lab-text-muted)">
          Один платёж открывает доступ на один месяц. Можно оплатить обычный тариф или отправить сумму выше как
          поддержку проекта.
        </p>
      </div>

      <div
        v-if="accessLoading"
        class="rounded-2xl border border-(--lab-border) bg-(--lab-bg-soft) p-4 text-sm text-(--lab-text-muted)">
        Загрузка сведений о доступе…
      </div>

      <div
        v-else-if="access?.has_active_access"
        class="rounded-2xl border border-(--lab-success)/30 bg-(--lab-success)/8 p-4">
        <div class="text-sm font-medium text-(--lab-text-primary)">Доступ уже активен</div>
        <p class="mt-1 text-sm text-(--lab-text-muted)">Текущий доступ действует до {{ activeUntilText || '—' }}.</p>
      </div>

      <div v-else class="rounded-2xl border border-(--lab-border) bg-(--lab-bg-soft) p-4">
        <div class="text-sm font-medium text-(--lab-text-primary)">Активного доступа нет</div>
        <p class="mt-1 text-sm text-(--lab-text-muted)">После успешной оплаты доступ откроется автоматически.</p>
      </div>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <article class="space-y-4 rounded-3xl border border-(--lab-border) bg-(--lab-bg-elevated) p-5">
        <div class="space-y-2">
          <h2 class="text-base font-medium text-(--lab-text-primary)">Pro на 1 месяц</h2>
          <p class="text-sm text-(--lab-text-muted)">Базовый платёж за доступ.</p>
        </div>

        <div class="text-2xl font-semibold text-(--lab-text-primary)">
          {{ formatPrice(proAmount) }}
        </div>

        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-2xl border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary) transition hover:border-(--lab-accent) hover:text-(--lab-text-primary) disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="creatingPlan !== ''"
          @click="submitOrder('pro')">
          {{ creatingPlan === 'pro' ? 'Переход к оплате…' : 'Оплатить Pro' }}
        </button>
      </article>

      <article class="space-y-4 rounded-3xl border border-(--lab-border) bg-(--lab-bg-elevated) p-5">
        <div class="space-y-2">
          <h2 class="text-base font-medium text-(--lab-text-primary)">Поддержка проекта</h2>
          <p class="text-sm text-(--lab-text-muted)">Можно отправить сумму выше базовой стоимости.</p>
        </div>

        <label class="space-y-2">
          <span class="text-sm text-(--lab-text-muted)">Сумма в копейках</span>
          <input
            v-model.number="donationAmount"
            type="number"
            inputmode="numeric"
            min="39900"
            step="100"
            class="min-h-11 w-full rounded-2xl border border-(--lab-border) bg-transparent px-4 text-sm text-(--lab-text-primary) outline-hidden" />
        </label>

        <div class="text-sm text-(--lab-text-muted)">
          {{ donationHintText }}
        </div>

        <div class="text-2xl font-semibold text-(--lab-text-primary)">
          {{ formatPrice(normalizeDonationAmount(donationAmount)) }}
        </div>

        <button
          type="button"
          class="inline-flex min-h-11 items-center justify-center rounded-2xl border border-(--lab-border-strong) px-4 text-sm font-medium text-(--lab-text-primary) transition hover:border-(--lab-accent) hover:text-(--lab-text-primary) disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="creatingPlan !== ''"
          @click="submitOrder('donation')">
          {{ creatingPlan === 'donation' ? 'Переход к оплате…' : 'Оплатить и поддержать' }}
        </button>
      </article>
    </section>

    <section
      v-if="formError"
      class="rounded-3xl border border-(--lab-danger)/30 bg-(--lab-danger)/8 p-4 text-sm text-(--lab-text-primary)">
      {{ formError }}
    </section>

    <section v-if="latestOrder" class="space-y-3 rounded-3xl border border-(--lab-border) bg-(--lab-bg-elevated) p-5">
      <h2 class="text-base font-medium text-(--lab-text-primary)">Последний заказ</h2>

      <dl class="grid gap-3 sm:grid-cols-2">
        <div class="space-y-1">
          <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Статус</dt>
          <dd class="text-sm text-(--lab-text-primary)">
            {{ latestOrderStatusLabel }}
          </dd>
        </div>

        <div class="space-y-1">
          <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Сумма</dt>
          <dd class="text-sm text-(--lab-text-primary)">
            {{ latestOrderAmountText }}
          </dd>
        </div>

        <div class="space-y-1">
          <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">План</dt>
          <dd class="text-sm text-(--lab-text-primary)">
            {{ latestOrder.plan_code === 'donation' ? 'Поддержка проекта' : 'Pro' }}
          </dd>
        </div>

        <div class="space-y-1">
          <dt class="text-xs uppercase tracking-wide text-(--lab-text-muted)">Создан</dt>
          <dd class="text-sm text-(--lab-text-primary)">
            {{ latestOrderCreatedAtText || '—' }}
          </dd>
        </div>
      </dl>
    </section>
  </div>
</template>
