package domain

import "time"

type Payment struct {
	ID       int
	PublicID string

	OrderID string
	PayerID int

	Amount   int
	Currency string

	Status   PaymentStatus
	Provider string
	Method   string

	IdempotencyKey string

	CreatedAt time.Time
	UpdatedAt time.Time
	PaidAt    *time.Time
}
