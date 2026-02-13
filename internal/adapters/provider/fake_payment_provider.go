package provider

import (
	"context"
	"errors"
	"math/rand/v2"
	"payment-service/internal/core/ports"
	"payment-service/internal/observability"
	"time"

	"go.opentelemetry.io/otel/codes"
)

type FakeProvider struct{}

func NewFakePaymentProvider() ports.PaymentProvider {
	return &FakeProvider{}
}

func (p *FakeProvider) Process(ctx context.Context, method string) error {
	ctx, span := observability.Tracer().
		Start(ctx, "PaymentProvider.Process")
	defer span.End()

	switch method {
	case "credit_card":
		time.Sleep(100 * time.Millisecond)
	case "bank_transfer":
		time.Sleep(200 * time.Millisecond)
	case "ewallet":
		time.Sleep(400 * time.Millisecond)
	}

	if rand.Float64() < 0.15 {
		err := errors.New("provider failure")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
