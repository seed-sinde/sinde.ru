package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"sinde.ru/db"
	"sinde.ru/db/services"
	authsvc "sinde.ru/internal/auth"
	routes "sinde.ru/internal/http"
	authhandlers "sinde.ru/internal/http/handlers/auth"
	paymenthandlers "sinde.ru/internal/http/handlers/payments"
	"sinde.ru/internal/http/middleware"
	paymentsvc "sinde.ru/internal/payments"
	"sinde.ru/internal/store"
	"sinde.ru/utils"
)

const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
)

func main() {
	if err := runServer(); err != nil {
		log.SetFlags(0)
		log.SetPrefix("")
		log.Fatalf("%v", err)
	}
}
func runServer() error {
	utils.LoadEnv()
	log.SetFlags(0)   // disable default prefix
	log.SetPrefix("") // disable "[INFO]" & date
	if err := db.Init(); err != nil {
		return fmt.Errorf("DB initilization error: %w", err)
	}
	defer db.PDB.Close()
	err := services.PdbLoadAllData()
	if err != nil {
		return fmt.Errorf("Postgres data load error: %w", err)
	}
	authConfig, err := authsvc.LoadConfig()
	if err != nil {
		return fmt.Errorf("Auth config error: %w", err)
	}
	authService, err := authsvc.NewService(authsvc.NewRepository(db.PDB), authConfig, authsvc.Dependencies{})
	if err != nil {
		return fmt.Errorf("Auth service init error: %w", err)
	}
	defer authService.Close()
	authHandler := authhandlers.New(authService)
	paymentConfig, err := paymentsvc.LoadConfig()
	if err != nil {
		return fmt.Errorf("Payments config error: %w", err)
	}
	paymentService := paymentsvc.NewService(paymentsvc.NewRepository(db.PDB), paymentConfig, paymentsvc.Dependencies{})
	paymentHandler := paymenthandlers.New(paymentService)
	i18nStore := store.NewI18nStore()
	app := fiber.New(fiber.Config{
		CaseSensitive:      true,
		StrictRouting:      false,
		ServerHeader:       "",
		AppName:            "",
		JSONEncoder:        json.Marshal,
		JSONDecoder:        json.Unmarshal,
		DisableDefaultDate: true,
		ProxyHeader:        "X-Forwarded-For",
		TrustProxy:         true,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback:  true,
			LinkLocal: true,
			Private:   true,
		},
		EnableIPValidation: true,
	})
	app.Use(middleware.AccessLog())
	app.Use(etag.New(etag.Config{
		Next: func(c fiber.Ctx) bool {
			path := strings.TrimSpace(c.Path())
			if strings.HasSuffix(path, "/stream") {
				return true
			}
			return strings.Contains(strings.ToLower(strings.TrimSpace(c.Get("Accept"))), "text/event-stream")
		},
	}))
	app.Use(favicon.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     authConfig.AllowedOrigins,
		AllowMethods:     []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
		AllowHeaders:     []string{"Origin", "Referer", "Content-Type", "Accept", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Etag", "Vary", "Date", "Set-Cookie"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(app, authHandler, paymentHandler, i18nStore)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		// Очищаем текущую строку терминала, чтобы отображаемое "^C" визуально не мешало.
		fmt.Print("\r\033[2K")
		fmt.Printf("%sShutdown signal received. Stopping server...%s\n", ansiRed, ansiReset)
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error during app shutdown:", err)
		}
	}()
	if err := app.Listen(":3001", fiber.ListenConfig{
		EnablePrefork:         false,
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
	}); err != nil {
		return fmt.Errorf("Server start error: %w", err)
	}
	fmt.Println("Server stopped.")
	return nil
}
