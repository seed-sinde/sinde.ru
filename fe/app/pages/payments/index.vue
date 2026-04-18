<script setup lang="ts">
const { localeTag, effectiveTheme } = useInterfacePreferences()
const { locale, key, load, t } = useI18nSection('payments')
await useAsyncData(key.value, load, { watch: [locale] })
const title = computed(() => t('index.seo_title'))
const description = computed(() => t('index.seo_description'))

usePageSeo({
  title,
  description
})

const nuxtApp = useNuxtApp()
const route = useRoute()
const router = useRouter()

const { isAuthenticated, isAdmin, ensureLoaded } = useAuth()
const { access, accessLoading, ensureAccessLoaded, createOrder } = usePayments()
const { formatAbsoluteDateTime } = useLocalizedDateTime()

if (import.meta.server || !nuxtApp.isHydrating) {
  await ensureLoaded()
  await ensureAccessLoaded()
}

const proAmount = 39900
const defaultDonationRubles = String(Math.floor(proAmount / 100))
const minDonationRubles = Number.parseInt(defaultDonationRubles, 10)
const donationAnimationDurationMs: number = 300
const donationAnimationFps = 25
const donationAnimationFrames = Math.max(1, Math.round((donationAnimationDurationMs / 1000) * donationAnimationFps))
const donationPresetAdds = [100, 500, 1000] as const

const donationAnimating = ref(false)
const donationEditing = ref(false)
const donationSwitchOn = ref(false)
const donationHasInput = ref(true)
const donationPresetBaseRubles = ref(minDonationRubles)
const donationCurrentRubles = ref(minDonationRubles)
const donationAmountFocused = ref(false)
const creatingPlan = ref<'' | 'pro' | 'donation'>('')
const paymentError = ref('')
const paymentErrorKey = ref(0)
const offerAccepted = ref(false)
const donationInputRef = ref<HTMLElement | null>(null)
let donationAnimationFrameId: number | null = null

const featureKeys = [
  'index.feature_roles',
  'index.feature_deep_workshop',
  'index.feature_chain_links',
  'index.feature_timeline',
  'index.feature_person_views',
  'index.feature_planning',
  'index.feature_collaboration',
  'index.feature_rbac',
  'index.feature_early_access'
] as const

const featureItems = computed(() => featureKeys.map(key => t(key)))
const donationPresetChips = computed(() =>
  donationPresetAdds.map(add => ({
    key: `preset:${add}`,
    add,
    label: `+${add}`
  }))
)
const donationEnabled = computed(() => donationHasInput.value && donationCurrentRubles.value > minDonationRubles)
const donationToggleChecked = computed(() => donationSwitchOn.value || donationEnabled.value)
const donationSelectedPresetAdd = computed<(typeof donationPresetAdds)[number] | null>(
  () => donationPresetAdds.find(add => donationCurrentRubles.value === donationPresetBaseRubles.value + add) ?? null
)
const showDonationChips = computed(() => donationEnabled.value || donationEditing.value)
const paymentActionDisabled = computed(
  () => creatingPlan.value !== '' || (donationEditing.value && !donationHasInput.value) || !offerAccepted.value
)

// Анимация Рейдена активируется, если донат >= 1000
const showElectrifiedEffect = computed(
  () => donationEnabled.value && donationCurrentRubles.value >= donationPresetBaseRubles.value + 1000
)
const electrifiedStyle = computed<Record<string, string>>(() => ({
  '--payments-gold-primary': effectiveTheme.value === 'light' ? '#f59e0b' : '#fbbf24',
  '--payments-gold-glow': effectiveTheme.value === 'light' ? 'rgba(245, 158, 11, 0.4)' : 'rgba(217, 119, 6, 0.6)',
  '--payments-gold-inner': effectiveTheme.value === 'light' ? '#ffffff' : '#fffbeb'
}))

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

const currentPlan = computed<'pro' | 'donation'>(() => (donationEnabled.value ? 'donation' : 'pro'))
const currentAmountKopecks = computed(() => {
  const amountKopecks = readDonationRubles() * 100
  if (isAdmin.value) return amountKopecks
  return donationEnabled.value ? Math.max(proAmount, amountKopecks) : proAmount
})

const accessSummaryText = computed(() => {
  if (accessLoading.value) return t('index.loading')
  if (access.value?.has_active_access && access.value.access_until) {
    return t('index.access_active_description', { date: formatAbsoluteDateTime(access.value.access_until) })
  }
  return ''
})

const latestOrderLine = computed(() => {
  const order = access.value?.latest_order
  if (!order) return ''
  return t('index.latest_order_line', {
    amount: formatPaymentAmount(order.amount || 0, localeTag.value),
    status: paymentStatusLabel(String(order.status || '').trim())
  })
})

function paymentStatusLabel(status: string) {
  switch (status) {
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

const setDonationStateFromText = (value: string) => {
  const digits = String(value).replace(/\D+/g, '')
  donationHasInput.value = digits.length > 0
  donationCurrentRubles.value = Number.parseInt(digits || defaultDonationRubles, 10) || minDonationRubles
}

function setDonationEditableText(value: string) {
  if (!donationInputRef.value) return
  donationInputRef.value.textContent = value
  setDonationStateFromText(value)
}

function handleDonationBeforeInput(event: InputEvent) {
  if (event.inputType === 'insertParagraph' || event.inputType === 'insertLineBreak') {
    event.preventDefault()
    return
  }
  if (!event.inputType.startsWith('insert')) return
  if (event.data == null || /\D/.test(event.data)) {
    event.preventDefault()
  }
}

function handleDonationKeydown(event: KeyboardEvent) {
  if (event.key !== 'Enter') return
  event.preventDefault()
  event.stopPropagation()
  ;(event.currentTarget as HTMLElement | null)?.blur()
}

function handleDonationToggleKeydown(event: KeyboardEvent) {
  if (event.key !== 'Enter' && event.key !== ' ') return
  event.preventDefault()
  toggleDonation(!donationToggleChecked.value)
}

function stopDonationAnimation() {
  if (donationAnimationFrameId == null) return
  cancelAnimationFrame(donationAnimationFrameId)
  donationAnimationFrameId = null
  donationAnimating.value = false
}

const readDonationRubles = () => (donationHasInput.value ? donationCurrentRubles.value : minDonationRubles)

function easeOutCubic(progress: number) {
  return 1 - (1 - progress) ** 3
}

const hasDonationValue = () => donationHasInput.value

function animateDonation(fromRubles: number, targetRubles = minDonationRubles) {
  stopDonationAnimation()
  const from = Math.max(0, fromRubles)
  const target = Math.max(minDonationRubles, targetRubles)
  const delta = target - from
  if (delta <= 0) {
    setDonationEditableText(String(Math.max(target, from)))
    return
  }
  donationAnimating.value = true
  const startAt = performance.now()
  let lastFrameIndex = -1
  const tick = (now: number) => {
    const elapsed = Math.min(now - startAt, donationAnimationDurationMs)
    const progress = donationAnimationDurationMs === 0 ? 1 : elapsed / donationAnimationDurationMs
    const frameIndex = Math.min(donationAnimationFrames, Math.floor(progress * donationAnimationFrames))
    if (frameIndex !== lastFrameIndex) {
      const easedProgress = easeOutCubic(frameIndex / donationAnimationFrames)
      const nextValue = frameIndex >= donationAnimationFrames ? target : Math.round(from + delta * easedProgress)
      setDonationEditableText(String(nextValue))
      lastFrameIndex = frameIndex
    }
    if (frameIndex >= donationAnimationFrames) {
      donationAnimationFrameId = null
      donationAnimating.value = false
      return
    }
    donationAnimationFrameId = requestAnimationFrame(tick)
  }
  donationAnimationFrameId = requestAnimationFrame(tick)
}

function handleDonationInput() {
  setDonationStateFromText(String(donationInputRef.value?.textContent || ''))
}

function handleDonationBlur() {
  donationAmountFocused.value = false
  donationEditing.value = false
  if (!hasDonationValue()) {
    donationPresetBaseRubles.value = minDonationRubles
    donationSwitchOn.value = false
    setDonationEditableText(defaultDonationRubles)
    return
  }
  const currentRubles = readDonationRubles()
  if (!isAdmin.value && currentRubles < minDonationRubles) {
    donationPresetBaseRubles.value = minDonationRubles
    animateDonation(currentRubles, minDonationRubles)
    return
  }
  donationSwitchOn.value = currentRubles > minDonationRubles
  setDonationEditableText(String(currentRubles))
}

function handleDonationFocus() {
  donationAmountFocused.value = true
}

function focusDonationInput(options?: { preservePresetBase?: boolean }) {
  nextTick(() => {
    if (!donationInputRef.value) return
    if (!options?.preservePresetBase) donationPresetBaseRubles.value = Math.max(readDonationRubles(), minDonationRubles)
    donationEditing.value = true
    donationInputRef.value.focus()
  })
}

function toggleDonation(nextValue: boolean) {
  stopDonationAnimation()
  donationEditing.value = false
  donationSwitchOn.value = nextValue
  if (nextValue) {
    donationPresetBaseRubles.value = minDonationRubles
    animateDonation(readDonationRubles(), minDonationRubles + 100)
    focusDonationInput({ preservePresetBase: true })
    return
  }
  donationPresetBaseRubles.value = minDonationRubles
  nextTick(() => setDonationEditableText(defaultDonationRubles))
}

function applyDonationPreset(addRubles: (typeof donationPresetAdds)[number]) {
  stopDonationAnimation()
  donationEditing.value = false
  donationPresetBaseRubles.value = Math.max(donationPresetBaseRubles.value, minDonationRubles)
  setDonationEditableText(String(donationPresetBaseRubles.value + addRubles))
  focusDonationInput({ preservePresetBase: true })
}

function activateDonationEditing() {
  stopDonationAnimation()
  if (isAdmin.value) {
    focusDonationInput()
    return
  }
  if (!donationEnabled.value) {
    toggleDonation(true)
    return
  }
  focusDonationInput()
}

async function submitOrder() {
  paymentError.value = ''
  if (!isAuthenticated.value) {
    await router.push('/auth/login')
    return
  }
  if (!offerAccepted.value) {
    paymentError.value = t('index.offer_required')
    paymentErrorKey.value += 1
    return
  }
  if (paymentActionDisabled.value) return
  creatingPlan.value = currentPlan.value
  try {
    const res = await createOrder({
      plan_code: currentPlan.value,
      amount: currentAmountKopecks.value,
      return_to: currentPath.value
    })
    const paymentURL = String(res?.data?.payment_url || '').trim()
    if (!paymentURL) throw new Error(t('lookup.error_generic'))
    await navigateTo(paymentURL, { external: true })
  } catch (error: any) {
    paymentError.value =
      String(error?.data?.message || error?.message || '').trim() || t('lookup.error_generic')
    paymentErrorKey.value += 1
  } finally {
    creatingPlan.value = ''
  }
}

onBeforeUnmount(() => stopDonationAnimation())
</script>

<template>
  <div class="selection:bg-amber-500/30">
    <LabNavHeader :title="t('index.title')" />
    <main class="relative max-w-6xl p-4">
      <section
        class="relative z-10 max-w-4xl border border-transparent bg-(--lab-bg-surface) transition-all duration-300"
        :class="
          showElectrifiedEffect
            ? 'animate-[payments-gold-flicker_0.1s_infinite_steps(2)] border-(--payments-gold-primary) ring-4 ring-[color-mix(in_srgb,var(--payments-gold-glow)_45%,transparent)]'
            : ''
        "
        :style="electrifiedStyle"
      >
        <div v-if="showElectrifiedEffect" class="pointer-events-none absolute inset-0 z-0 overflow-hidden select-none">
          <div
            class="absolute top-0 z-20 h-0.75 bg-(--payments-gold-inner) opacity-90 drop-shadow-[0_0_15px_var(--payments-gold-primary)]"
            style="
              animation:
                payments-move-h 0.7s cubic-bezier(0.68, -0.55, 0.265, 1.55) infinite,
                payments-size-p-h 0.2s infinite alternate;
            "
          />
          <div
            class="absolute bottom-0 z-20 h-0.75 bg-(--payments-gold-inner) opacity-90 drop-shadow-[0_0_15px_var(--payments-gold-primary)]"
            style="
              animation:
                payments-move-h 1s cubic-bezier(1, 0, 0, 1) infinite reverse,
                payments-size-p-h 0.3s infinite alternate-reverse;
            "
          />
          <div
            class="absolute left-0 z-20 w-0.75 bg-(--payments-gold-inner) opacity-90 drop-shadow-[0_0_15px_var(--payments-gold-primary)]"
            style="
              animation:
                payments-move-v 0.9s cubic-bezier(0.47, 0, 0.745, 0.715) infinite reverse,
                payments-size-p-v 0.25s infinite alternate;
            "
          />
          <div
            class="absolute right-0 z-20 w-0.75 bg-(--payments-gold-inner) opacity-90 drop-shadow-[0_0_15px_var(--payments-gold-primary)]"
            style="
              animation:
                payments-move-v 0.6s cubic-bezier(0.175, 0.885, 0.32, 1.275) infinite,
                payments-size-p-v 0.2s infinite alternate-reverse;
            "
          />
        </div>

        <div class="relative z-10 px-6 py-10 text-center sm:px-12 sm:py-14">
          <div class="mb-10 space-y-4">
            <h1
              class="mx-auto max-w-3xl text-3xl font-black tracking-tight text-(--lab-text-primary) uppercase sm:text-5xl"
            >
              {{ t('index.hero_title') }}
            </h1>
            <p
              v-if="accessSummaryText"
              class="mx-auto max-w-xl text-sm font-bold tracking-wider text-(--lab-text-secondary) uppercase"
            >
              {{ accessSummaryText }}
            </p>
          </div>

          <div class="mb-10 flex flex-col items-center gap-8">
            <div
              class="group relative flex items-baseline justify-center transition-transform"
              :class="{ 'cursor-text': donationEnabled || isAdmin }"
            >
              <span
                id="payments-amount"
                ref="donationInputRef"
                :contenteditable="donationEditing || donationEnabled || isAdmin"
                suppress-contenteditable-warning
                tabindex="0"
                role="textbox"
                class="relative block min-w-[1ch] border-b-2 border-transparent bg-transparent px-4 py-2 text-right text-7xl font-black tracking-tighter text-(--lab-text-primary) tabular-nums outline-none focus-visible:border-(--lab-accent) sm:text-8xl"
                :class="[donationToggleChecked ? 'text-(--lab-accent)' : '']"
                @beforeinput="handleDonationBeforeInput"
                @input="handleDonationInput"
                @keydown="handleDonationKeydown"
                @focus="handleDonationFocus"
                @blur="handleDonationBlur"
                @click="activateDonationEditing"
              >
                {{ defaultDonationRubles }}
              </span>
              <span class="relative ml-2 text-xl font-black text-(--lab-text-muted) uppercase sm:text-2xl">₽</span>
            </div>

            <p v-if="isAdmin" class="text-[10px] font-black tracking-[0.18em] text-(--lab-text-secondary) uppercase">
              {{ t('index.admin_manual_amount_hint') }}
            </p>

            <div
              v-if="showDonationChips"
              class="flex animate-[payments-fade-in_0.3s_ease-out_forwards] flex-wrap items-center justify-center gap-1"
            >
              <button
                v-for="chip in donationPresetChips"
                :key="chip.key"
                type="button"
                class="lab-focus inline-flex h-8 items-center justify-center rounded-full border px-4 text-xs font-black uppercase tabular-nums transition-all"
                :class="
                  donationSelectedPresetAdd === chip.add
                    ? 'border-(--lab-accent) bg-(--lab-accent) text-black'
                    : 'bg-(--lab-bg-surface-subtle) text-(--lab-text-primary) hover:border-(--lab-border-strong)'
                "
                @click="applyDonationPreset(chip.add)"
              >
                {{ chip.label }}
              </button>
            </div>

            <div class="flex flex-col items-center gap-3">
              <label
                tabindex="0"
                role="switch"
                :aria-checked="donationToggleChecked"
                class="group flex cursor-pointer items-center gap-3 rounded-full bg-(--lab-bg-control) p-1 pr-4 transition-colors"
                @keydown="handleDonationToggleKeydown"
              >
                <input
                  type="checkbox"
                  class="peer lab-focus sr-only"
                  tabindex="-1"
                  :checked="donationToggleChecked"
                  @change="event => toggleDonation((event.target as HTMLInputElement).checked)"
                />
                <div
                  class="lab-focus-peer relative h-6 w-11 rounded-full bg-(--lab-border) transition-colors"
                  :class="{ 'bg-(--lab-accent)': donationToggleChecked }"
                >
                  <div
                    class="absolute top-0.5 h-5 w-5 rounded-full bg-white transition-all"
                    :class="donationToggleChecked ? 'left-5.5' : 'left-0.5'"
                  />
                </div>
                <span
                  class="text-[10px] font-black tracking-widest text-(--lab-text-secondary) uppercase group-hover:text-(--lab-text-primary)"
                >
                  {{ t('index.donation_toggle') }}
                </span>
              </label>
            </div>

            <label
              class="flex max-w-2xl items-start gap-3 text-left text-sm leading-6 text-(--lab-text-secondary)"
              for="payments-offer-accept"
            >
              <input
                id="payments-offer-accept"
                v-model="offerAccepted"
                type="checkbox"
                class="lab-focus mt-0.5 h-4 w-4 shrink-0 border border-(--lab-border-strong) bg-transparent"
              />
              <span>
                {{ t('index.offer_accept_prefix') }}
                <NuxtLink to="/docs/offer" class="lab-focus text-(--lab-accent) transition-colors hover:underline">
                  {{ t('index.offer_accept_link') }}
                </NuxtLink>
                {{ t('index.offer_accept_suffix') }}
              </span>
            </label>

            <LabBaseButton
              variant="primary"
              size="xl"
              :loading="creatingPlan !== ''"
              :disabled="paymentActionDisabled"
              :label="t('index.submit')"
              button-class="w-full !rounded-none border-none font-black uppercase tracking-widest sm:w-80"
              @click="submitOrder()"
            />
          </div>

          <ul class="lab-grid-table mt-10 grid-cols-1 sm:grid-cols-3">
            <li v-for="item in featureItems" :key="item" class="lab-grid-table-cell flex items-center gap-3 p-3">
              <Icon name="ic:round-check" class="shrink-0 text-lg text-(--lab-accent)" />
              <span
                class="text-left text-xs leading-tight font-bold wrap-break-word text-(--lab-text-primary) uppercase"
              >
                {{ item }}
              </span>
            </li>
          </ul>

          <footer class="mt-8 flex flex-col items-center gap-4">
            <div
              v-if="latestOrderLine"
              class="inline-flex items-center gap-2 rounded-full bg-(--lab-bg-surface-subtle) px-3 py-1 text-[10px] font-bold text-(--lab-text-secondary) uppercase"
            >
              <span class="relative flex h-1.5 w-1.5">
                <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75" />
                <span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-green-500" />
              </span>
              {{ latestOrderLine }}
            </div>
            <LabNotify
              :key="paymentErrorKey"
              :text="paymentError"
              tone="error"
              size="xs"
              :temporary="true"
              :duration-ms="4200"
              class-name="text-xs font-black uppercase"
            />
          </footer>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
@keyframes move-h {
  0% {
    left: -20%;
    opacity: 0;
  }
  20% {
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    left: 120%;
    opacity: 0;
  }
}

@keyframes move-v {
  0% {
    top: -20%;
    opacity: 0;
  }
  20% {
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    top: 120%;
    opacity: 0;
  }
}

@keyframes size-p-h {
  0% {
    width: 20px;
  }
  100% {
    width: 150px;
  }
}

@keyframes size-p-v {
  0% {
    height: 20px;
  }
  100% {
    height: 130px;
  }
}

@keyframes gold-flicker {
  0%,
  100% {
    border-color: var(--payments-gold-glow);
  }
  50% {
    border-color: var(--payments-gold-primary);
  }
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
