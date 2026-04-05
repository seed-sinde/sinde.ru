export const useInputRejectFeedback = <T extends string | number>() => {
  const rejectedTarget = ref<T | null>(null)
  const rejectedHint = ref('')
  let timer: ReturnType<typeof setTimeout> | null = null
  const clearRejectFeedback = () => {
    rejectedTarget.value = null
    rejectedHint.value = ''
    if (!timer) return
    clearTimeout(timer)
    timer = null
  }
  const markRejected = (target: T, hint = '', durationMs = 1000) => {
    rejectedTarget.value = target
    rejectedHint.value = hint
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      rejectedTarget.value = null
      rejectedHint.value = ''
      timer = null
    }, durationMs)
  }
  onBeforeUnmount(() => {
    if (timer) clearTimeout(timer)
  })
  return {
    rejectedTarget,
    rejectedHint,
    markRejected,
    clearRejectFeedback,
  }
}
