package usecase

import (
	"context"
	"payment-service/internal/core/domain"
	"payment-service/internal/core/ports"
	"payment-service/internal/observability"

	"go.opentelemetry.io/otel/codes"
)

type GetPaymentUsecase struct {
	paymentRepo ports.PaymentRepository
}

func NewGetPaymentUsecase(
	paymentRepo ports.PaymentRepository,
) *GetPaymentUsecase {
	return &GetPaymentUsecase{
		paymentRepo: paymentRepo,
	}
}

func (uc *GetPaymentUsecase) Execute(
	ctx context.Context,
	publicID string,
) (*domain.Payment, error) {
	ctx, span := observability.Tracer().Start(ctx, "GetPaymentUseCase.Execute")
	defer span.End()

	payment, err := uc.paymentRepo.FindbyPublicID(
		ctx,
		publicID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return payment, nil
}
