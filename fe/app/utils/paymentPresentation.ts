export const formatPaymentAmount = (amountKopecks: number, locale = 'ru-RU') => {
  const amount = Number.isFinite(amountKopecks) ? amountKopecks / 100 : 0
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

export const formatPaymentInputRubles = (amountKopecks: number, locale = 'ru-RU') => {
  const amount = Number.isFinite(amountKopecks) ? amountKopecks / 100 : 0
  return new Intl.NumberFormat(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
    useGrouping: false
  }).format(amount)
}

export const normalizePaymentRublesInput = (value?: string) => {
  const raw = String(value || '').replace(/\s+/g, '').replace(/,/g, '.')
  if (!raw) return ''
  let out = ''
  let hasDot = false
  for (const char of raw) {
    if (char >= '0' && char <= '9') {
      out += char
      continue
    }
    if (char === '.' && !hasDot) {
      out += '.'
      hasDot = true
    }
  }
  const [wholeRaw = '', fractionRaw = ''] = out.split('.', 2)
  const whole = wholeRaw.replace(/^0+(?=\d)/, '')
  const fraction = fractionRaw.slice(0, 2)
  if (!hasDot) return whole
  return `${whole || '0'}.${fraction}`
}

export const paymentRublesInputToKopecks = (value?: string, minimumKopecks = 0) => {
  const normalized = normalizePaymentRublesInput(value)
  if (!normalized) return minimumKopecks
  const [wholeRaw = '0', fractionRaw = ''] = normalized.split('.', 2)
  const whole = Number.parseInt(wholeRaw || '0', 10)
  const fraction = Number.parseInt((fractionRaw + '00').slice(0, 2), 10)
  if (!Number.isFinite(whole) || !Number.isFinite(fraction)) return minimumKopecks
  return Math.max(minimumKopecks, whole * 100 + fraction)
}
