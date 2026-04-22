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
useSeoMeta({
  title: "Главная",
  description: "Платформа прикладны"
})
</script>
<template>
  <Stack>
    <Spoiler v-model="details" label="Скрытый блок">
      <p class="my-1">Содержимое скрытого блока!</p>
    </Spoiler>
    <Button label="Кнопка с иконкой слева" :icon="IconMoreVert" icon-class="text-lg" />
    <Button label="Кнопка с иконкой справа" :icon="IconMoreVert" icon-class="text-lg" icon-position="right" />
    <Tooltip :text="lorem" side="right">
      <span class="border-b border-dotted">Показать всплывающую подсказку на тексте</span>
    </Tooltip>
    <div class="flex w-fit flex-col">
      <Field label="Почта" hint="Введите пароль не менее 12 символов" :error="email_error" required>
        <Input v-model="email" type="email" placeholder="seed@sinde.ru" :invalid="hasEmailError" />
      </Field>
      <Field label="Пароль" hint="Введите пароль не менее 12 символов" :error="password_error" required>
        <Input v-model="password" type="password" placeholder="Пароль" :invalid="hasPassError" />
      </Field>
    </div>
    <div v-if="data">
      <MessageError />
      {{ data }}</div>
  </Stack>
</template>
