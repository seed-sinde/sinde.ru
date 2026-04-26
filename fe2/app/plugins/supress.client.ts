export default defineNuxtPlugin(() => {
  if (import.meta.env.DEV) {
    // eslint-disable-next-line no-console
    const originalInfo = console.info
    // eslint-disable-next-line no-console
    console.info = (...args) => {
      if (
        typeof args[0] === "string" &&
        args[0].includes("<Suspense> is an experimental feature")
      ) {
        return // terminate msg
      }
      originalInfo(...args)
    }
  }
})
