package domain

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusExpired    PaymentStatus = "EXPIRED"
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending,
		PaymentStatusProcessing,
		PaymentStatusSuccess,
		PaymentStatusFailed,
		PaymentStatusExpired:
		return true
	default:
		return false
	}
}

func (s PaymentStatus) IsFinal() bool {
	return s == PaymentStatusSuccess ||
		s == PaymentStatusFailed ||
		s == PaymentStatusExpired
}

func (p *Payment) CanTransitionTo(next PaymentStatus) bool {
	if p.Status.IsFinal() {
		return false
	}

	switch p.Status {
	case PaymentStatusPending:
		return next == PaymentStatusProcessing ||
			next == PaymentStatusFailed ||
			next == PaymentStatusExpired

	case PaymentStatusProcessing:
		return next == PaymentStatusSuccess ||
			next == PaymentStatusFailed

	default:
		return false
	}
}
