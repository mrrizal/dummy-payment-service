package usecase

import (
	"context"
	"errors"
	"testing"

	"payment-service/internal/core/domain"
	"payment-service/internal/observability"
)

// mockGetPaymentRepo implements ports.PaymentRepository for get-payment tests
type mockGetPaymentRepo struct {
    returned *domain.Payment
    err      error
}

func (m *mockGetPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
    return nil
}

func (m *mockGetPaymentRepo) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
    return nil, errors.New("not implemented")
}

func (m *mockGetPaymentRepo) FindbyPublicID(ctx context.Context, publicID string) (*domain.Payment, error) {
    return m.returned, m.err
}

func TestGetPayment_Success(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    expected := &domain.Payment{
        PublicID: "pay_123",
        Status:   domain.PaymentStatusSuccess,
    }

    repo := &mockGetPaymentRepo{returned: expected}

    uc := NewGetPaymentUsecase(repo)

    got, err := uc.Execute(ctx, "pay_123")
    if err != nil {
        t.Fatalf("expected nil error, got %v", err)
    }
    if got == nil {
        t.Fatalf("expected non-nil payment")
    }
    if got.PublicID != expected.PublicID {
        t.Fatalf("expected public id %s, got %s", expected.PublicID, got.PublicID)
    }
    if got.Status != expected.Status {
        t.Fatalf("expected status %s, got %s", expected.Status, got.Status)
    }
}

func TestGetPayment_RepoError(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    repo := &mockGetPaymentRepo{err: errors.New("db error")}

    uc := NewGetPaymentUsecase(repo)

    got, err := uc.Execute(ctx, "pay_missing")
    if err == nil {
        t.Fatalf("expected error from repo, got nil and payment %v", got)
    }
}
