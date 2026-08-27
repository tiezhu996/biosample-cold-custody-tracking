package service

import (
	"context"
	"errors"
	"testing"

	"biosample-cold-custody-tracking/backend/internal/model"
)

func TestProtocolReviewPreservesRequestCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.WithValue(context.Background(), protocolRequestTickKey{}, "request-84"))
	cancel()

	var observed context.Context
	plan := prepareProtocolWritePlan(parent, 84)
	_, _, _, err := runProtocolPersistence(plan.Worker, nil, func(ctx context.Context, _ *model.ProtocolReview) (*model.ProtocolReview, model.Specimen, model.Specimen, error) {
		observed = ctx
		return nil, model.Specimen{}, model.Specimen{}, ctx.Err()
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation did not propagate into persistence: %v", err)
	}
	select {
	case <-observed.Done():
	default:
		t.Fatal("persistence context remained active")
	}
	if plan.RequestTick != "request-84:84" {
		t.Fatalf("persistence plan lost request epoch: %q", plan.RequestTick)
	}
}
