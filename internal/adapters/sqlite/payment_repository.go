package sqlite

import (
	"context"
	"database/sql"
	"payment-service/internal/core/domain"
	"payment-service/internal/core/ports"
	"time"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) ports.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	now := time.Now()

	p.CreatedAt = now
	p.UpdatedAt = now

	query := `
	INSERT INTO payments (
	public_id,
	order_id,
	payer_id,
	amount,
	currency,
	status,
	provider,
	method,
	idempotency_key,
	created_at,
	updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		p.PublicID,
		p.OrderID,
		p.PayerID,
		p.Amount,
		p.Currency,
		p.Status,
		p.Provider,
		p.Method,
		p.IdempotencyKey,
		p.CreatedAt,
		p.UpdatedAt,
	)

	return err
}

func (r *paymentRepository) FindByIdempotencyKey(
	ctx context.Context,
	key string,
) (*domain.Payment, error) {

	query := `
	SELECT
		id, public_id, order_id, payer_id,
		amount, currency, status,
		provider, method, idempotency_key,
		created_at, updated_at, paid_at
	FROM payments
	WHERE idempotency_key = ?
	`

	row := r.db.QueryRowContext(ctx, query, key)

	var p domain.Payment
	var paidAt sql.NullTime

	err := row.Scan(
		&p.ID,
		&p.PublicID,
		&p.OrderID,
		&p.PayerID,
		&p.Amount,
		&p.Currency,
		&p.Status,
		&p.Provider,
		&p.Method,
		&p.IdempotencyKey,
		&p.CreatedAt,
		&p.UpdatedAt,
		&paidAt,
	)

	if err != nil {
		return nil, err
	}

	if paidAt.Valid {
		p.PaidAt = &paidAt.Time
	}

	return &p, nil
}
