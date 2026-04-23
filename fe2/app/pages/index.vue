<script setup lang="ts">
import IconMoreVert from "~icons/ic/round-more-vert"
const details = ref(false)
const lorem = `Lorem, ipsum dolor sit amet consectetur adipisicing elit. Vel doloremque itaque numquam ad.`
const email = ref("")
const hasEmailError = ref(true)
const hasPassError = ref(true)
const email_error = ref("Неверный формат почты")
const password = ref("")
const password_error = ref("Пароль не верный")
const {data, error} = await useAsyncData("test-api", () => useAPI().json("/i18n/ru/ui"))
const onSubmit = () => {
  return
}
useSeoMeta({
  title: "Главная",
  description: "Платформа прикладны"
})
</script>
<template>
  <UiStack>
    <UiSpoiler v-model="details" label="Скрытый блок">
      <p class="my-1">Содержимое скрытого блока!</p>
    </UiSpoiler>
    <UiButton label="Кнопка с иконкой слева" :icon="IconMoreVert" icon-class="text-lg" />
    <UiButton label="Кнопка с иконкой справа" :icon="IconMoreVert" icon-class="text-lg" icon-position="right" />
    <UiTooltip :text="lorem" side="right">
      <span class="border-b border-dotted">Показать всплывающую подсказку на тексте</span>
    </UiTooltip>
    <form class="flex w-fit flex-col" @submit.prevent="onSubmit">
      <UiField label="Почта" hint="Введите пароль не менее 12 символов" :error="email_error" required>
        <UiInput
          v-model="email"
          type="email"
          placeholder="seed@sinde.ru"
          :invalid="hasEmailError"
          autocomplete="username"
        />
      </UiField>
      <UiField label="Пароль" hint="Введите пароль не менее 12 символов" :error="password_error" required>
        <UiInput
          v-model="password"
          type="password"
          placeholder="Пароль"
          :invalid="hasPassError"
          autocomplete="current-password"
        />
      </UiField>
    </form>
    <div v-if="data">
      <MessageError :text="error?.message" />
      {{ data }}
    </div>
  </UiStack>
</template>
