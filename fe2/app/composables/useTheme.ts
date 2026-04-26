type ThemeMode = "light" | "dark" | "system"
const isThemeMode = (v: unknown): v is ThemeMode => v === "light" || v === "dark" || v === "system"
export const useTheme = () => {
  const theme = useCookie<ThemeMode>("theme", {
    default: () => "system",
    path: "/"
  })
  const mode = useState<ThemeMode>("theme", () =>
    isThemeMode(theme.value) ? theme.value : "system"
  )
  const isDark = () =>
    mode.value === "dark" ||
    (mode.value === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches)
  const apply = () => document.documentElement.classList.toggle("dark", isDark())
  const set = (v: ThemeMode) => {
    mode.value = v
    theme.value = v
    apply()
  }
  const toggle = () =>
    set(mode.value === "light" ? "dark" : mode.value === "dark" ? "system" : "light")
  return {mode, set, toggle, apply, isDark}
}
