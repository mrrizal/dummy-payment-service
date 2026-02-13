package usecase

import (
	"context"
	"errors"
	"payment-service/internal/core/domain"
	"payment-service/internal/core/ports"
	"payment-service/internal/observability"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type CreatePaymentInput struct {
	OrderID        string
	PayerID        int
	Amount         int
	Currency       string
	Provider       string
	Method         string
	IdempotencyKey string
}

type CreatePaymentOutput struct {
	PaymentID string
	Status    domain.PaymentStatus
}

type CreatePaymentUsecase struct {
	paymentRepo     ports.PaymentRepository
	paymentProvider ports.PaymentProvider
}

func NewCreatePaymentUsecase(
	paymentRepo ports.PaymentRepository,
	paymentProvider ports.PaymentProvider,
) *CreatePaymentUsecase {
	return &CreatePaymentUsecase{
		paymentRepo:     paymentRepo,
		paymentProvider: paymentProvider,
	}
}

func isValidPaymentInput(input CreatePaymentInput) (bool, error) {
	if input.Amount <= 0 {
		return false, errors.New("amount must be greater than zero")
	}
	if input.Currency == "" {
		return false, errors.New("currency is required")
	}
	if input.Method == "" {
		return false, errors.New("payment method is required")
	}
	if input.Provider == "" {
		return false, errors.New("payment provider is required")
	}
	if input.IdempotencyKey == "" {
		return false, errors.New("idempotency key is required")
	}
	return true, nil
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (uc *CreatePaymentUsecase) Execute(
	ctx context.Context,
	input CreatePaymentInput,
) (*CreatePaymentOutput, error) {
	ctx, span := observability.Tracer().Start(ctx, "CreatePaymentUseCase.Execute")
	defer span.End()

	// --- validate input ---
	if valid, err := isValidPaymentInput(input); !valid {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := uc.paymentProvider.Process(ctx, input.Method); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// --- create domain object ---
	now := time.Now()

	payment := &domain.Payment{
		PublicID:       "pay_" + uuid.NewString(),
		OrderID:        input.OrderID,
		PayerID:        input.PayerID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		Provider:       input.Provider,
		Method:         input.Method,
		IdempotencyKey: input.IdempotencyKey,
		Status:         domain.PaymentStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// --- persist ---
	var paymentOutput CreatePaymentOutput
	err := uc.paymentRepo.Create(ctx, payment)
	if err != nil {
		// --- handle idempotency key conflict ---
		if isUniqueConstraintError(err) {
			existingPayment, findErr := uc.paymentRepo.FindByIdempotencyKey(
				ctx,
				input.IdempotencyKey,
			)
			if findErr != nil {
				return nil, findErr
			}
			paymentOutput = CreatePaymentOutput{
				PaymentID: existingPayment.PublicID,
				Status:    existingPayment.Status,
			}
			return &paymentOutput, nil
		}
	}

	// --- return lightweight response ---
	return &CreatePaymentOutput{
		PaymentID: payment.PublicID,
		Status:    payment.Status,
	}, nil
}
