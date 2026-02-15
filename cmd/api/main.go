package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"payment-service/internal/adapters/provider"
	"payment-service/internal/adapters/sqlite"
	"payment-service/internal/config"
	"payment-service/internal/core/usecase"
	"payment-service/internal/http/handler"
	"payment-service/internal/http/middleware"
	"payment-service/internal/http/router"
	"payment-service/internal/observability"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func run() error {
	fmt.Println("Starting Payment Service...")
	ctx := context.Background()

	// --- load config ---
	cfg := config.LoadConfig()

	// Set up OpenTelemetry.
	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup telemetry: %w", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		_ = errors.Join(err, otelShutdown(context.Background()))
	}()

	// init tracker
	observability.InitTracer(cfg.App.ServiceName)

	// init metric
	observability.InitMetrics()

	// --- init database (SQLite) ---
	db, err := sqlite.New(cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	chaosCfg := config.ChaosConfig{
		Enabled:          true,
		ErrorProbability: 0.2,
		DelayProbability: 0.3,
		MaxDelay:         700 * time.Millisecond,
	}

	// --- init repositories (adapters) ---
	defaultPaymentRepo := sqlite.NewPaymentRepository(db)
	paymentRepoWithMetrics := sqlite.NewPaymentRepositoryMetrics(
		defaultPaymentRepo,
	)
	paymentRepo := sqlite.NewPaymentRepositoryChaos(
		paymentRepoWithMetrics,
		chaosCfg,
	)

	// --- payment provider
	paymentProvider := provider.NewFakePaymentProvider()

	// --- init usecases ---
	createPaymentUC := usecase.NewCreatePaymentUsecase(
		paymentRepo,
		paymentProvider,
	)
	getPaymentUC := usecase.NewGetPaymentUsecase(paymentRepo)

	// --- init handlers ---
	paymentHandler := handler.NewPaymentHandler(
		createPaymentUC,
		getPaymentUC,
	)

	// --- init gin ---
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(cfg.App.ServiceName))
	r.Use(middleware.MetricsMiddleware())

	// --- register routes ---
	router.Register(r, paymentHandler)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// --- start server ---
	log.Printf("starting http server on :%s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
