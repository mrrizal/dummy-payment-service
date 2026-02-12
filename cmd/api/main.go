package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"payment-service/internal/adapters/sqlite"
	"payment-service/internal/config"
	"payment-service/internal/core/usecase"
	"payment-service/internal/http/handler"
	"payment-service/internal/http/middleware"
	"payment-service/internal/http/router"
	"payment-service/internal/observability"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	fmt.Println("Starting Payment Service...")
	ctx := context.Background()

	// --- load config ---
	cfg := config.LoadConfig()

	// Set up OpenTelemetry.
	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatalf("failed to setup telemetry: %v", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// init tracker
	observability.InitTracer(cfg.App.ServiceName)

	// init metric
	observability.InitMetrics()

	// --- init database (SQLite) ---
	db, err := sqlite.New(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	chaosCfg := config.ChaosConfig{
		Enabled:          true,
		ErrorProbability: 0.2,
		DelayProbability: 0.3,
		MaxDelay:         500 * time.Millisecond,
	}

	// --- init repositories (adapters) ---
	defaultPaymentRepo := sqlite.NewPaymentRepository(db)

	paymentRepoChaos := sqlite.NewPaymentRepositoryChaos(
		defaultPaymentRepo,
		chaosCfg,
	)

	paymentRepo := sqlite.NewPaymentRepositoryMetrics(
		paymentRepoChaos,
	)

	// --- init usecases ---
	createPaymentUC := usecase.NewCreatePaymentUsecase(paymentRepo)
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
		log.Fatalf("failed to start server: %v", err)
	}
}
