import type { Ref } from 'vue'
import { compareDateText, datePlaceholder, isValidDateText, maskDateText } from './useDateText'
type DatePairModel = Record<string, string>
type DatePairOptions = {
  leftKey: string
  rightKey: string
}
export const useDatePairModel = (model: Ref<DatePairModel>, options: DatePairOptions) => {
  const mode = 'datetime'
  const placeholder = computed(() => datePlaceholder(mode))
  const left = ref(String(model.value?.[options.leftKey] ?? ''))
  const right = ref(String(model.value?.[options.rightKey] ?? ''))
  const leftValid = computed(() => !left.value || isValidDateText(left.value, mode))
  const rightValid = computed(() => !right.value || isValidDateText(right.value, mode))
  const isRangeInvalid = computed(() => {
    if (!left.value || !right.value || !leftValid.value || !rightValid.value) return false
    const cmp = compareDateText(left.value, right.value, mode)
    return cmp !== null && cmp > 0
  })
  const onLeftInput = (raw: string) => {
    left.value = maskDateText(raw, mode)
  }
  const onRightInput = (raw: string) => {
    right.value = maskDateText(raw, mode)
  }
  watch([left, right], ([leftValue, rightValue]) => {
    const current = model.value ?? {}
    if (
      String(current?.[options.leftKey] ?? '') === leftValue &&
      String(current?.[options.rightKey] ?? '') === rightValue
    )
      return
    model.value = {
      ...current,
      [options.leftKey]: leftValue,
      [options.rightKey]: rightValue
    }
  })
  watch(
    (): [string, string] => [
      String(model.value?.[options.leftKey] ?? ''),
      String(model.value?.[options.rightKey] ?? '')
    ],
    ([nextLeft, nextRight]) => {
      if (left.value !== nextLeft) left.value = nextLeft
      if (right.value !== nextRight) right.value = nextRight
    }
  )
  return {
    mode,
    placeholder,
    left,
    right,
    leftValid,
    rightValid,
    isRangeInvalid,
    onLeftInput,
    onRightInput
  }
}
