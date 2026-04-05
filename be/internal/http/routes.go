package routes
import (
	"github.com/gofiber/fiber/v3"
	"sinde.ru/internal/http/handlers"
	authhandlers "sinde.ru/internal/http/handlers/auth"
	"sinde.ru/internal/http/middleware"
)
func SetupRoutes(app *fiber.App, authHandler *authhandlers.Handler) {
	registerAPIRoutes(app.Group("/api/v1"), authHandler)
}
func registerAPIRoutes(api fiber.Router, authHandler *authhandlers.Handler) {
	api.Post("/auth/register", authHandler.Register())
	api.Post("/auth/verify-email/request", authHandler.RequestEmailVerification())
	api.Post("/auth/verify-email/confirm", authHandler.VerifyEmail())
	api.Post("/auth/login", authHandler.Login())
	api.Post("/auth/login/2fa", authHandler.Login2FA())
	api.Post("/auth/refresh", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.Refresh())
	api.Post("/auth/logout", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.Logout())
	api.Post("/auth/password/forgot", authHandler.ForgotPassword())
	api.Post("/auth/password/reset", authHandler.ResetPassword())
	authenticated := api.Group("/auth", middleware.RequireAuth(authHandler.Service()))
	authenticated.Get("/me", authHandler.Me())
	authenticated.Patch("/me", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.UpdateMe())
	authenticated.Get("/traits/sets", authHandler.ListSavedTraitSets())                                                              // получить список сохраненных наборов
	authenticated.Post("/traits/sets", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.SaveTraitSet())              // сохранить набор в список
	authenticated.Patch("/traits/sets/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.UpdateSavedTraitSet())  // обновление набора в списке
	authenticated.Delete("/traits/sets/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.DeleteSavedTraitSet()) // удалить набор из списка
	authenticated.Post("/traits/primary", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.SetPrimaryTrait())        // установить основной набор
	authenticated.Get("/sessions", authHandler.ListSessions())
	authenticated.Delete("/sessions/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.RevokeSession())
	authenticated.Get("/login-attempts", authHandler.ListLoginAttempts())
	authenticated.Get("/security-events", authHandler.ListSecurityEvents())
	authenticated.Post("/logout-all", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.LogoutAll())
	authenticated.Post("/password/change", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.ChangePassword())
	authenticated.Post("/2fa/setup", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.SetupTwoFactor())
	authenticated.Post("/2fa/enable", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.EnableTwoFactor())
	authenticated.Post("/2fa/disable", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.DisableTwoFactor())
	authenticated.Get("/summary", authHandler.UserSummary())
	authenticated.Get("/summary/stream", authHandler.UserSummaryStream())
	authenticated.Post("/summary/read", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.UserSummaryRead())
	admin := authenticated.Group("/admin", middleware.RequireRole("admin"))
	admin.Get("/users", authHandler.AdminListUsers())
	admin.Get("/summary", authHandler.AdminSummary())
	admin.Get("/summary/stream", authHandler.AdminSummaryStream())
	admin.Post("/summary/read", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminSummaryRead())
	admin.Patch("/users/:id/role", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminSetUserRole())
	admin.Post("/users/:id/block", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminBlockUser())
	admin.Post("/users/:id/unblock", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminUnblockUser())
	admin.Post("/users/:id/force-logout", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminForceLogoutUser())
	admin.Delete("/users/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminDeleteUser())
	admin.Get("/keys/search", authHandler.AdminSearchTraitKeys())
	admin.Get("/analysis/traits-sets", authHandler.AdminTraitsSetsAnalysis())
	api.Get("/media/files/*", handlers.MediaGetFileHandler())                   // Получить медиафайл по ключу хранения.
	api.Get("/kitchen/catalog", handlers.KitchenCatalogHandler())               // Каталог, фильтры и ингредиенты.
	api.Get("/kitchen/ingredients", handlers.KitchenIngredientsHandler())       // Категории для kitchen.
	api.Get("/kitchen/recipes/latest", handlers.KitchenLatestRecipesHandler())  // Последние рецепты.
	api.Get("/minerals", handlers.MineralsListHandler())                        // Список, поиск и фильтры минералов.
	api.Get("/minerals/:database_id", handlers.MineralGetByDatabaseIDHandler()) // Карточка минерала по slug.
	api.Post("/kitchen/recipes/search", handlers.KitchenSearchRecipesHandler()) // Поиск рецептов.
	kitchenAuth := api.Group("/kitchen", middleware.RequireAuth(authHandler.Service()))
	kitchenAuth.Get("/ingredients/account", handlers.KitchenAccountIngredientsHandler())
	kitchenAuth.Post("/ingredients/custom", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenCreateCustomIngredientHandler())
	kitchenAuth.Delete("/ingredients/custom/:id", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenDeleteCustomIngredientHandler())
	kitchenAuth.Post("/ingredients/favorites", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenFavoriteIngredientHandler())
	kitchenAuth.Delete("/ingredients/favorites/:id", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenUnfavoriteIngredientHandler())
	kitchenAuth.Get("/recipes/:id/favorite", handlers.KitchenRecipeFavoriteStatusHandler())
	kitchenAuth.Post("/recipes/:id/favorite", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenFavoriteRecipeHandler())
	kitchenAuth.Delete("/recipes/:id/favorite", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenUnfavoriteRecipeHandler())
	kitchenAuth.Post("/recipes", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenCreateRecipeHandler())
	kitchenAuth.Get("/recipes/mine", handlers.KitchenMyRecipesHandler())
	kitchenAuth.Get("/recipes/manage/:id", handlers.KitchenGetManageRecipeHandler())
	kitchenAuth.Patch("/recipes/:id", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenUpdateRecipeHandler())
	kitchenAuth.Delete("/recipes/:id", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenDeleteRecipeHandler())
	kitchenAuth.Get("/admin/recipes/moderation", handlers.KitchenAdminModerationRecipesHandler())
	kitchenAuth.Post("/admin/recipes/:id/moderate", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenAdminModerateRecipeHandler())
	kitchenAuth.Post("/admin/recipes/:id/owner", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenAdminChangeRecipeOwnerHandler())
	mediaAuth := api.Group("/media", middleware.RequireAuth(authHandler.Service()))
	mediaAuth.Post("/upload", handlers.MediaUploadHandler())
	api.Get("/kitchen/recipes/:id", handlers.KitchenGetRecipeHandler()) // Получить рецепт по ID.
	api.Get("/sets/:uuid", handlers.StoreSetHandler())              // Вернуть список особенностей по UUID набора.
	api.Get("/sets/:uuid/stream", handlers.StoreSetStreamHandler()) // Стрим особенностей по UUID набора.
	api.Post("/sets", handlers.CreateSetHandler())                  // Создать набор.
	api.Post("/sets/find", handlers.FindOrBuildSetHandler())        // Получить набор по произвольному списку особенностей.
	api.Post("/traits", handlers.MemoryAddTraitHandler())                 // Добавить особенность.
	api.Get("/traits/:uuid", handlers.MemoryTraitHandler())               // Получить особенность.
	api.Get("/traits/resolve/:uuid", handlers.MemoryResolveUUIDHandler()) // Вернуть особенность или набор по UUID.
	api.Post("/keys/meta", handlers.KeyMetaHandler())                // Вернуть meta по syn.
	api.Post("/keys/meta/all", handlers.KeyMetaAllHandler())         // Вернуть все meta по syn.
	api.Post("/keys/meta/bulk", handlers.KeyMetaBulkHandler())       // Вернуть meta по списку ID.
	api.Post("/keys/enum/options", handlers.KeyEnumOptionsHandler()) // Вернуть enum-опции по syn.
	api.Post("/keys/meta/update", handlers.KeyMetaUpdateHandler())   // Обновить meta по ID.
}
