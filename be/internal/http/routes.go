package routes

import (
	"github.com/gofiber/fiber/v3"
	"sinde.ru/internal/http/handlers"
	authhandlers "sinde.ru/internal/http/handlers/auth"
	paymenthandlers "sinde.ru/internal/http/handlers/payments"
	"sinde.ru/internal/http/middleware"
)

func SetupRoutes(app *fiber.App, authHandler *authhandlers.Handler, paymentHandler *paymenthandlers.Handler) {
	registerAPIRoutes(app.Group("/api/v1"), authHandler, paymentHandler)
}

func registerAPIRoutes(api fiber.Router, authHandler *authhandlers.Handler, paymentHandler *paymenthandlers.Handler) {
	// Auth: публичные маршруты.
	api.Post("/auth/register", authHandler.Register())
	api.Post("/auth/verify-email/request", authHandler.RequestEmailVerification())
	api.Post("/auth/verify-email/confirm", authHandler.VerifyEmail())
	api.Post("/auth/login", authHandler.Login())
	api.Post("/auth/login/2fa", authHandler.Login2FA())
	api.Post("/auth/refresh", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.Refresh())
	api.Post("/auth/logout", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.Logout())
	api.Post("/auth/password/forgot", authHandler.ForgotPassword())
	api.Post("/auth/password/reset", authHandler.ResetPassword())
	api.Get("/users/:id", authHandler.PublicUserProfile()) // Публичный профиль пользователя.

	// Auth: маршруты авторизованного пользователя.
	authenticated := api.Group("/auth", middleware.RequireAuth(authHandler.Service()))
	authenticated.Get("/me", authHandler.Me())
	authenticated.Patch("/me", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.UpdateMe())
	authenticated.Post("/email/change/request", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.RequestEmailChange())

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

	// Auth: сохранённые наборы особенностей пользователя.
	authenticated.Get("/traits/sets", authHandler.ListSavedTraitSets())                                                              // Список сохранённых наборов.
	authenticated.Post("/traits/sets", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.SaveTraitSet())              // Сохранить набор.
	authenticated.Patch("/traits/sets/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.UpdateSavedTraitSet())  // Обновить сохранённый набор.
	authenticated.Delete("/traits/sets/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.DeleteSavedTraitSet()) // Удалить сохранённый набор.
	authenticated.Post("/traits/primary", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.SetPrimaryTrait())        // Назначить основной набор.

	// Auth: административные маршруты.
	admin := authenticated.Group("/admin", middleware.RequireRole("admin"))
	admin.Get("/users", authHandler.AdminListUsers())
	admin.Get("/users/:id", authHandler.AdminUserDetail())
	admin.Patch("/users/:id/role", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminSetUserRole())
	admin.Post("/users/:id/block", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminBlockUser())
	admin.Post("/users/:id/unblock", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminUnblockUser())
	admin.Post("/users/:id/force-logout", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminForceLogoutUser())
	admin.Delete("/users/:id", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminDeleteUser())

	admin.Get("/summary", authHandler.AdminSummary())
	admin.Get("/summary/stream", authHandler.AdminSummaryStream())
	admin.Post("/summary/read", middleware.RequireCSRFCookie(authHandler.Service()), authHandler.AdminSummaryRead())

	admin.Get("/keys/search", authHandler.AdminSearchTraitKeys())
	admin.Get("/analysis/traits-sets", authHandler.AdminTraitsSetsAnalysis())

	// Payments: публичные маршруты.
	api.Post("/payments/lookup", paymentHandler.PublicLookup())            // Публичная проверка заказа по order_id и token.
	api.Post("/payments/tbank/notify", paymentHandler.TBankNotification()) // Входящий webhook T-Bank.

	// Payments: маршруты авторизованного пользователя.
	paymentsAuth := api.Group("/payments", middleware.RequireAuth(authHandler.Service()))
	paymentsAuth.Get("/access", paymentHandler.Access()) // Текущий доступ пользователя.
	paymentsAuth.Get("/history", paymentHandler.History())
	paymentsAuth.Post("/init", middleware.RequireCSRFCookie(authHandler.Service()), paymentHandler.CreateOrder())
	paymentsAuth.Post("/:orderId/refund", middleware.RequireCSRFCookie(authHandler.Service()), paymentHandler.Refund())

	// Payments: административные маршруты.
	admin.Get("/orders", paymentHandler.AdminListOrders())
	admin.Get("/payments/summary", paymentHandler.AdminSummary())
	admin.Get("/users/:id/access", paymentHandler.AdminUserAccess())
	admin.Get("/users/:id/orders", paymentHandler.AdminUserOrders())

	// Media.
	api.Get("/media/files/*", handlers.MediaGetFileHandler()) // Получить медиафайл по ключу хранения.

	mediaAuth := api.Group("/media", middleware.RequireAuth(authHandler.Service()))
	mediaAuth.Post("/upload", handlers.MediaUploadHandler())

	// Kitchen: публичные маршруты.
	api.Get("/kitchen/catalog", handlers.KitchenCatalogHandler())               // Каталог, фильтры и ингредиенты.
	api.Get("/kitchen/ingredients", handlers.KitchenIngredientsHandler())       // Категории ингредиентов.
	api.Get("/kitchen/recipes/latest", handlers.KitchenLatestRecipesHandler())  // Последние рецепты.
	api.Post("/kitchen/recipes/search", handlers.KitchenSearchRecipesHandler()) // Поиск рецептов.
	api.Get("/kitchen/recipes/:id", handlers.KitchenGetRecipeHandler())         // Рецепт по ID.

	// Kitchen: маршруты авторизованного пользователя.
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

	// Kitchen: административные маршруты.
	kitchenAuth.Get("/admin/recipes/moderation", handlers.KitchenAdminModerationRecipesHandler())
	kitchenAuth.Post("/admin/recipes/:id/moderate", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenAdminModerateRecipeHandler())
	kitchenAuth.Post("/admin/recipes/:id/owner", middleware.RequireCSRFCookie(authHandler.Service()), handlers.KitchenAdminChangeRecipeOwnerHandler())

	// Chemistry.
	api.Get("/chemistry/elements", handlers.ChemistryElementsListHandler()) // Список химических элементов.

	// Minerals.
	api.Get("/minerals", handlers.MineralsListHandler())                        // Список минералов, поиск и фильтры.
	api.Get("/minerals/:database_id", handlers.MineralGetByDatabaseIDHandler()) // Карточка минерала.

	// Traits / Sets / Keys: публичные маршруты.
	api.Get("/traits/:uuid", handlers.MemoryTraitHandler())               // Получить особенность.
	api.Post("/traits", handlers.MemoryAddTraitHandler())                 // Создать особенность.
	api.Get("/traits/resolve/:uuid", handlers.MemoryResolveUUIDHandler()) // Разрешить UUID в особенность или набор.

	api.Get("/sets/:uuid", handlers.StoreSetHandler())              // Получить набор по UUID.
	api.Get("/sets/:uuid/stream", handlers.StoreSetStreamHandler()) // Потоковая выдача набора по UUID.
	api.Post("/sets", handlers.CreateSetHandler())                  // Создать набор.
	api.Post("/sets/find", handlers.FindOrBuildSetHandler())        // Найти или создать набор по списку особенностей.

	api.Post("/keys/meta", handlers.KeyMetaHandler())                // Meta по одному syn.
	api.Post("/keys/meta/all", handlers.KeyMetaAllHandler())         // Все meta по списку syn.
	api.Post("/keys/meta/bulk", handlers.KeyMetaBulkHandler())       // Meta по списку ID.
	api.Post("/keys/enum/options", handlers.KeyEnumOptionsHandler()) // Enum-опции по syn.
	api.Post("/keys/meta/update", handlers.KeyMetaUpdateHandler())   // Обновить meta по ID.

	// Traits / Sets / Keys: публичные маршруты. v2

}
