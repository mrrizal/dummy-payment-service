package ports

import (
	"context"
	"payment-service/internal/core/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	FindByIdempotencyKey(
		ctx context.Context,
		idempotencyKey string,
	) (*domain.Payment, error)
}
