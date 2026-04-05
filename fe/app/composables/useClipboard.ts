export const useClipboard = (timeout = 1500) => {
  const copied = ref(false)
  let timer: ReturnType<typeof setTimeout> | null = null

  const pasteInto = async (target: Ref<string>) => {
    if (!import.meta.client) return
    const text = await navigator.clipboard.readText()
    target.value = text.trim()
  }

  const copyFrom = async (val: string | null) => {
    if (!import.meta.client || !val) return
    await navigator.clipboard.writeText(val)
    copied.value = true
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      copied.value = false
      timer = null
    }, timeout)
  }

  // SSR-safe copied (always false during SSR)
  const safeCopied = computed(() => (import.meta.client ? copied.value : false))

  return {
    pasteInto,
    copyFrom,
    copied: safeCopied
  }
}
