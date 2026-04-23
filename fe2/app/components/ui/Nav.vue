<template>
  <nav class="flex items-center gap-1 border-b p-1 text-sm">
    <UiButton label="Главная" to="/" />
    <UiButton label="Мастерская" to="/lab" />
    <UiButton label="Особенности" to="/traits" />
    <UiButton label="Вики" to="/edu" />
    <UiButton label="Вид" />
    <UiButton label="Помощь" />
    <ThemeSwitcher />
  </nav>
</template>
