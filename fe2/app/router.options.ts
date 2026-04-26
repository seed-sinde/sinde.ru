import type {RouterConfig} from "nuxt/schema"
const tabAlias = (path: string) => (path === "/" ? "/:tab([^/]+)" : `${path}/:tab(.*)*`)
const tabRouteName = (name: unknown) => (typeof name === "string" ? `${name}-tabs` : undefined)
export default {
  routes: routes =>
    routes.flatMap(route =>
      route.meta?.tabs
        ? [
            route,
            {
              ...route,
              name: tabRouteName(route.name),
              path: tabAlias(route.path)
            }
          ]
        : [route]
    )
} satisfies RouterConfig
