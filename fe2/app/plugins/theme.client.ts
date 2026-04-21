export default defineNuxtPlugin(() => {
  const {apply, mode} = useTheme()
  apply()
  if (mode.value === "system") {
    window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", apply)
  }
})
