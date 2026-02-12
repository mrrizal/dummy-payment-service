package sqlite

import (
	"context"
	"payment-service/internal/core/domain"
	"payment-service/internal/core/ports"
	"payment-service/internal/observability"
	"time"
)

type PaymentRepositoryMetrics struct {
	next ports.PaymentRepository
}

func NewPaymentRepositoryMetrics(next ports.PaymentRepository) ports.PaymentRepository {
	return &PaymentRepositoryMetrics{next: next}
}

func (r *PaymentRepositoryMetrics) Create(ctx context.Context, payment *domain.Payment) error {
	start := time.Now()

	err := r.next.Create(ctx, payment)

	duration := time.Since(start).Seconds()

	observability.DBQueryDuration.WithLabelValues("insert").Observe(duration)

	if err != nil {
		observability.DBErrors.WithLabelValues("insert").Inc()
	}

	return err
}

func (r *PaymentRepositoryMetrics) FindByIdempotencyKey(
	ctx context.Context,
	idempotencyKey string,
) (*domain.Payment, error) {
	start := time.Now()

	payment, err := r.next.FindByIdempotencyKey(ctx, idempotencyKey)

	duration := time.Since(start).Seconds()

	observability.DBQueryDuration.WithLabelValues("select").Observe(duration)

	if err != nil {
		observability.DBErrors.WithLabelValues("select").Inc()
	}

	return payment, err
}

func (r *PaymentRepositoryMetrics) FindbyPublicID(
	ctx context.Context,
	publicID string,
) (*domain.Payment, error) {
	start := time.Now()

	payment, err := r.next.FindbyPublicID(ctx, publicID)

	duration := time.Since(start).Seconds()

	observability.DBQueryDuration.WithLabelValues("select").Observe(duration)

	if err != nil {
		observability.DBErrors.WithLabelValues("select").Inc()
	}

	return payment, err
}
