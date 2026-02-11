package sqlite

import (
	"context"

	"payment-service/internal/chaos"
	"payment-service/internal/config"
	"payment-service/internal/core/domain"
	"payment-service/internal/core/ports"
	"payment-service/internal/observability"
)

type PaymentRepositoryChaos struct {
	next ports.PaymentRepository
	cfg  config.ChaosConfig
}

func NewPaymentRepositoryChaos(
	next ports.PaymentRepository,
	cfg config.ChaosConfig,
) ports.PaymentRepository {
	return &PaymentRepositoryChaos{
		next: next,
		cfg:  cfg,
	}
}

func (r *PaymentRepositoryChaos) Create(
	ctx context.Context,
	payment *domain.Payment,
) error {

	if r.cfg.Enabled {
		chaos.MaybeDelay(
			r.cfg.DelayProbability,
			r.cfg.MaxDelay,
		)

		if err := chaos.MaybeError(r.cfg.ErrorProbability); err != nil {
			return err
		}
	}

	return r.next.Create(ctx, payment)
}

func (r *PaymentRepositoryChaos) FindByIdempotencyKey(
	ctx context.Context,
	idempotencyKey string,
) (*domain.Payment, error) {

	if r.cfg.Enabled {
		chaos.MaybeDelay(
			r.cfg.DelayProbability,
			r.cfg.MaxDelay,
		)

		if err := chaos.MaybeError(r.cfg.ErrorProbability); err != nil {
			return nil, err
		}
	}

	return r.next.FindByIdempotencyKey(ctx, idempotencyKey)
}

func (r *PaymentRepositoryChaos) FindbyPublicID(
	ctx context.Context,
	publicID string,
) (*domain.Payment, error) {
	ctx, span := observability.Tracer().Start(ctx, "PaymentRepositoryChaos.FindbyPublicID")
	defer span.End()

	if r.cfg.Enabled {
		chaos.MaybeDelay(
			r.cfg.DelayProbability,
			r.cfg.MaxDelay,
		)

		if err := chaos.MaybeError(r.cfg.ErrorProbability); err != nil {
			return nil, err
		}
	}

	return r.next.FindbyPublicID(ctx, publicID)
}
