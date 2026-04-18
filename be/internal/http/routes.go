package routes

import (
	"github.com/gofiber/fiber/v3"
	h "sinde.ru/internal/http/handlers"
	authhandlers "sinde.ru/internal/http/handlers/auth"
	paymenthandlers "sinde.ru/internal/http/handlers/payments"
	"sinde.ru/internal/http/middleware"
	"sinde.ru/internal/store"
)

func SetupRoutes(app *fiber.App, authHandler *authhandlers.Handler, paymentHandler *paymenthandlers.Handler, i18nStore *store.I18nStore) {
	registerAPIRoutes(app.Group("/api/v1"), authHandler, paymentHandler, i18nStore)
}

func registerAPIRoutes(api fiber.Router, authHandler *authhandlers.Handler, paymentHandler *paymenthandlers.Handler, i18nStore *store.I18nStore) {
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
	api.Get("/users/:id", authHandler.PublicUserProfile())           // Публичный профиль пользователя.
	api.Get("/i18n/:locale/:namespace", h.GetI18nHandler(i18nStore)) // Переводы интерфейса

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
	api.Get("/media/files/*", h.MediaGetFileHandler()) // Получить медиафайл по ключу хранения.

	mediaAuth := api.Group("/media", middleware.RequireAuth(authHandler.Service()))
	mediaAuth.Post("/upload", h.MediaUploadHandler())

	// Kitchen: публичные маршруты.
	api.Get("/kitchen/catalog", h.KitchenCatalogHandler())               // Каталог, фильтры и ингредиенты.
	api.Get("/kitchen/ingredients", h.KitchenIngredientsHandler())       // Категории ингредиентов.
	api.Get("/kitchen/recipes/latest", h.KitchenLatestRecipesHandler())  // Последние рецепты.
	api.Post("/kitchen/recipes/search", h.KitchenSearchRecipesHandler()) // Поиск рецептов.
	// Kitchen: маршруты авторизованного пользователя.
	kitchenAuth := api.Group("/kitchen", middleware.RequireAuth(authHandler.Service()))
	kitchenAuth.Get("/ingredients/account", h.KitchenAccountIngredientsHandler())
	kitchenAuth.Post("/ingredients/custom", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenCreateCustomIngredientHandler())
	kitchenAuth.Delete("/ingredients/custom/:id", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenDeleteCustomIngredientHandler())
	kitchenAuth.Post("/ingredients/favorites", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenFavoriteIngredientHandler())
	kitchenAuth.Delete("/ingredients/favorites/:id", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenUnfavoriteIngredientHandler())

	kitchenAuth.Get("/recipes/:id/favorite", h.KitchenRecipeFavoriteStatusHandler())
	kitchenAuth.Post("/recipes/:id/favorite", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenFavoriteRecipeHandler())
	kitchenAuth.Delete("/recipes/:id/favorite", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenUnfavoriteRecipeHandler())

	kitchenAuth.Post("/recipes", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenCreateRecipeHandler())
	kitchenAuth.Get("/recipes/mine", h.KitchenMyRecipesHandler())
	kitchenAuth.Get("/recipes/manage/:id", h.KitchenGetManageRecipeHandler())
	kitchenAuth.Patch("/recipes/:id", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenUpdateRecipeHandler())
	kitchenAuth.Delete("/recipes/:id", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenDeleteRecipeHandler())

	api.Get("/kitchen/recipes/:id", h.KitchenGetRecipeHandler()) // Рецепт по ID.

	// Kitchen: административные маршруты.
	kitchenAuth.Get("/admin/recipes/moderation", h.KitchenAdminModerationRecipesHandler())
	kitchenAuth.Post("/admin/recipes/:id/moderate", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenAdminModerateRecipeHandler())
	kitchenAuth.Post("/admin/recipes/:id/owner", middleware.RequireCSRFCookie(authHandler.Service()), h.KitchenAdminChangeRecipeOwnerHandler())

	// Chemistry.
	api.Get("/chemistry/elements", h.ChemistryElementsListHandler()) // Список химических элементов.

	// Minerals.
	api.Get("/minerals", h.MineralsListHandler())                        // Список минералов, поиск и фильтры.
	api.Get("/minerals/:database_id", h.MineralGetByDatabaseIDHandler()) // Карточка минерала.

	// Traits / Sets / Keys: публичные маршруты.
	api.Get("/traits/:uuid", h.MemoryTraitHandler())               // Получить особенность.
	api.Post("/traits", h.MemoryAddTraitHandler())                 // Создать особенность.
	api.Get("/traits/resolve/:uuid", h.MemoryResolveUUIDHandler()) // Разрешить UUID в особенность или набор.

	api.Get("/sets/:uuid", h.StoreSetHandler())              // Получить набор по UUID.
	api.Get("/sets/:uuid/stream", h.StoreSetStreamHandler()) // Потоковая выдача набора по UUID.
	api.Post("/sets", h.CreateSetHandler())                  // Создать набор.
	api.Post("/sets/find", h.FindOrBuildSetHandler())        // Найти или создать набор по списку особенностей.

	api.Post("/keys/meta", h.KeyMetaHandler())                // Meta по одному syn.
	api.Post("/keys/meta/all", h.KeyMetaAllHandler())         // Все meta по списку syn.
	api.Post("/keys/meta/bulk", h.KeyMetaBulkHandler())       // Meta по списку ID.
	api.Post("/keys/enum/options", h.KeyEnumOptionsHandler()) // Enum-опции по syn.
	api.Post("/keys/meta/update", h.KeyMetaUpdateHandler())   // Обновить meta по ID.
	admin.Get("/ip", h.GetIPInfo())
}
