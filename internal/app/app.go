package app

import (
	"log"
	"log/slog"
	"net/http"
	_ "testapp/docs"
	app_config "testapp/internal/configs"
	logger "testapp/internal/logs"
	"testapp/internal/middleware"
	"testapp/internal/providers"
	"testapp/internal/routers"
	"testapp/internal/storage"
	"testapp/internal/web"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Wallet API
// @version 1.0
// @description API для управления балансом кошельков
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func Run() {
	err := logger.Init("app.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	c, err := app_config.LoadConfig("..")
	if err != nil {
		logger.Log.Debug("Конфиг не читается", slog.Any("error", err))
	}

	mux := web.NewMux()
	DB := storage.DatabaseInit(c)
	p := providers.New(DB)
	r := routers.New(&p)

	// Роутер для API

	go mux.HandleFunc("GET /api/v1/wallet/{WALLET_UUID}", r.GetBalance)
	go mux.HandleFunc("POST /api/v1/wallet/update_balance", r.UpdateBalance)
	go mux.HandleFunc("GET /api/v1/wallets", r.Wallets)

	// Роутер для Swagger UI
	go mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	loggedMux := middleware.LoggingMiddleware(mux)

	logger.Log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		logger.Log.Error("Server failed",
			"error", err.Error(),
		)
	}
}
