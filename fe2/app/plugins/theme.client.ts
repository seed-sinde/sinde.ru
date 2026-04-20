export default defineNuxtPlugin(() => {
  const { set } = useTheme();

  // Инициализация при загрузке
  const saved = localStorage.getItem("theme") || "system";
  set(saved as any);

  // Слушатель системных изменений
  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", () => {
      if (localStorage.getItem("theme") === "system") set("system");
    });
});
