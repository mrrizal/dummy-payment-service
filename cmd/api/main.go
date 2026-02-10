package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"payment-service/internal/adapters/sqlite"
	"payment-service/internal/config"
	"payment-service/internal/core/usecase"
	"payment-service/internal/http/handler"
	"payment-service/internal/http/router"
)

func main() {
	fmt.Println("Starting Payment Service...")
	// --- load config ---
	cfg := config.LoadConfig()

	// --- init database (SQLite) ---
	db, err := sqlite.New(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	// --- init repositories (adapters) ---
	paymentRepo := sqlite.NewPaymentRepository(db)

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

	// --- register routes ---
	router.Register(r, paymentHandler)

	// --- start server ---
	log.Printf("starting http server on :%s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
