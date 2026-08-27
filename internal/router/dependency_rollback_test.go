package router

import (
	"context"
	"errors"
	"testing"

	"biosample-cold-custody-tracking/backend/internal/util"
)

func TestRouterDependencyRollbackDoesNotSwallowFailure(t *testing.T) {
	var rollbackErr = errors.New("pool close blocked by active sessions")
	var probeErr = errors.New("redis timeout")

	ready, err := prepareCustodyRuntime(context.Background(),
		util.StartupGuard{Name: "postgres", Probe: func(context.Context) error { return nil }, Rollback: func() error { return rollbackErr }},
		util.StartupGuard{Name: "redis", Probe: func(context.Context) error { return probeErr }, Rollback: func() error { return nil }},
	)

	if ready || !errors.Is(err, probeErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("rollback failure swallowed: ready=%v err=%v", ready, err)
	}
}
