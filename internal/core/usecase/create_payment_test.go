package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"payment-service/internal/core/domain"
	"payment-service/internal/observability"
)

// mockPaymentProvider implements ports.PaymentProvider
type mockPaymentProvider struct {
    err       error
    calledWith string
}

func (m *mockPaymentProvider) Process(ctx context.Context, method string) error {
    m.calledWith = method
    return m.err
}

// mockPaymentRepo implements ports.PaymentRepository
type mockPaymentRepo struct {
    createErr                      error
    createdPayment                  *domain.Payment
    findByIdempotencyKeyPayment    *domain.Payment
    findErr                        error
}

func (m *mockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
    // store a copy so tests can inspect independently
    p := *payment
    m.createdPayment = &p
    return m.createErr
}

func (m *mockPaymentRepo) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
    return m.findByIdempotencyKeyPayment, m.findErr
}

func (m *mockPaymentRepo) FindbyPublicID(ctx context.Context, publicID string) (*domain.Payment, error) {
    return nil, errors.New("not implemented")
}

func TestExecute_Success(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    repo := &mockPaymentRepo{}
    provider := &mockPaymentProvider{}

    uc := NewCreatePaymentUsecase(repo, provider)

    input := CreatePaymentInput{
        OrderID:        "order_123",
        PayerID:        42,
        Amount:         1000,
        Currency:       "USD",
        Provider:       "FAKE",
        Method:         "CARD",
        IdempotencyKey: "idem-1",
    }

    out, err := uc.Execute(ctx, input)
    if err != nil {
        t.Fatalf("expected nil error, got %v", err)
    }
    if out == nil {
        t.Fatalf("expected non-nil output")
    }
    if !strings.HasPrefix(out.PaymentID, "pay_") {
        t.Fatalf("expected payment id to start with pay_, got %s", out.PaymentID)
    }
    if out.Status != domain.PaymentStatusPending {
        t.Fatalf("expected status PENDING, got %s", out.Status)
    }
    // repo should have received the created payment
    if repo.createdPayment == nil {
        t.Fatalf("expected repo.Create to be called")
    }
    if repo.createdPayment.IdempotencyKey != input.IdempotencyKey {
        t.Fatalf("idempotency key mismatch: expected %s got %s", input.IdempotencyKey, repo.createdPayment.IdempotencyKey)
    }
}

func TestExecute_InvalidInput(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    repo := &mockPaymentRepo{}
    provider := &mockPaymentProvider{}

    uc := NewCreatePaymentUsecase(repo, provider)

    input := CreatePaymentInput{
        OrderID:        "",
        PayerID:        0,
        Amount:         0, // invalid
        Currency:       "",
        Provider:       "",
        Method:         "",
        IdempotencyKey: "",
    }

    out, err := uc.Execute(ctx, input)
    if err == nil {
        t.Fatalf("expected error for invalid input, got nil and output %v", out)
    }
}

func TestExecute_ProviderError(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    repo := &mockPaymentRepo{}
    provider := &mockPaymentProvider{err: errors.New("provider failed")}

    uc := NewCreatePaymentUsecase(repo, provider)

    input := CreatePaymentInput{
        OrderID:        "order_1",
        PayerID:        1,
        Amount:         10,
        Currency:       "USD",
        Provider:       "FAKE",
        Method:         "CARD",
        IdempotencyKey: "idem-2",
    }

    out, err := uc.Execute(ctx, input)
    if err == nil {
        t.Fatalf("expected provider error, got nil and output %v", out)
    }
}

func TestExecute_IdempotencyConflict(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()

    existing := &domain.Payment{
        PublicID: "pay_existing",
        Status:   domain.PaymentStatusSuccess,
    }

    repo := &mockPaymentRepo{
        createErr:                   errors.New("UNIQUE constraint failed: payments.idempotency_key"),
        findByIdempotencyKeyPayment: existing,
        findErr:                     nil,
    }
    provider := &mockPaymentProvider{}

    uc := NewCreatePaymentUsecase(repo, provider)

    input := CreatePaymentInput{
        OrderID:        "order_x",
        PayerID:        5,
        Amount:         100,
        Currency:       "USD",
        Provider:       "FAKE",
        Method:         "CARD",
        IdempotencyKey: "idem-3",
    }

    out, err := uc.Execute(ctx, input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if out == nil {
        t.Fatalf("expected output, got nil")
    }
    if out.PaymentID != existing.PublicID {
        t.Fatalf("expected payment id %s got %s", existing.PublicID, out.PaymentID)
    }
    if out.Status != existing.Status {
        t.Fatalf("expected status %s got %s", existing.Status, out.Status)
    }
}

// small utility: ensure time-dependent behavior compiles
func TestCreatePayment_TimestampsSet(t *testing.T) {
    observability.InitTracer("test")

    ctx := context.Background()
    repo := &mockPaymentRepo{}
    provider := &mockPaymentProvider{}

    uc := NewCreatePaymentUsecase(repo, provider)

    input := CreatePaymentInput{
        OrderID:        "o",
        PayerID:        2,
        Amount:         1,
        Currency:       "IDR",
        Provider:       "FAKE",
        Method:         "QR",
        IdempotencyKey: "idem-4",
    }

    _, err := uc.Execute(ctx, input)
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if repo.createdPayment == nil {
        t.Fatalf("expected created payment in repo")
    }
    // CreatedAt and UpdatedAt should be close to now
    if time.Since(repo.createdPayment.CreatedAt) < 0 {
        t.Fatalf("CreatedAt in future")
    }
    if repo.createdPayment.UpdatedAt.IsZero() {
        t.Fatalf("UpdatedAt not set")
    }
}
