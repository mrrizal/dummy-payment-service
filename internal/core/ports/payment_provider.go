package ports

import "context"

type PaymentProvider interface {
	Process(ctx context.Context, method string) error
}
