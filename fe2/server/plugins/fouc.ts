export default defineNitroPlugin(nitroApp => {
  nitroApp.hooks.hook("render:html", (html, {event}) => {
    const theme = getCookie(event, "theme") ?? "system"
    html.head.unshift(
      `<script>(m=>{(m==='dark'||(m==='system'&&matchMedia('(prefers-color-scheme: dark)').matches))&&document.documentElement.classList.add('dark')})('${theme}')</script>`
    )
  })
})
