package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/davmontas/exchange-rate-offers/internal/application"
	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/service"
	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/transport"
)

func main() {
	// Load application's config
	cfg := application.Load()

	// Set Up logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Set Up Router
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Load Dependencies
	clients := application.RegisterAPIs(cfg)

	// Exchange service
	exchService := service.NewExchangeService(clients, cfg.Service.Timeout, sugar)

	// Load HTTP Routes
	transport.RegisterRoutes(router, exchService, sugar)

	// Run the server
	srv := &http.Server{
		Addr:    cfg.Server.ListenAddr,
		Handler: router,
	}
	go func() {
		sugar.Infow("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("server crashed", "error", err)
		}
	}()

	// Listen for shutdown signals >:)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Infow("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Errorw("server forced to shutdown", "error", err)
	}
	sugar.Infow("server exited")
}
