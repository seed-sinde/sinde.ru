import { defineEventHandler } from 'h3'

export default defineEventHandler((event) => {
  const { req, res } = event.node
  if (req.url?.startsWith('/.well-known/appspecific/com.chrome.devtools')) {
    res.statusCode = 204
    res.end()
    return
  }
})
