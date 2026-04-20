export const useTheme = () => {
  const mode = useState<"light" | "dark" | "system">("theme", () => "system");

  const apply = (v: typeof mode.value) =>
    document.documentElement.classList.toggle(
      "dark",
      v === "dark" ||
        (v === "system" && matchMedia("(prefers-color-scheme: dark)").matches),
    );

  const set = (v: typeof mode.value) => (
    (mode.value = v),
    localStorage.setItem("theme", v),
    apply(v)
  );

  if (import.meta.client)
    apply(
      (mode.value =
        (localStorage.getItem("theme") as typeof mode.value) ?? "system"),
    );

  return { mode, set };
};
