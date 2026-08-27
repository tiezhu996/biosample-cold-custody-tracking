package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"biosample-cold-custody-tracking/backend/internal/repository"
	"biosample-cold-custody-tracking/backend/internal/util"
)

type auditVerifyRepo struct {
	repository.AuditRepository
	err error
}

func (r auditVerifyRepo) VerifyChain(context.Context) error { return r.err }

func TestAuditServiceReportsBrokenChainWithStatus(t *testing.T) {
	broken := repository.ErrAuditChainBroken
	svc := NewAuditService(auditVerifyRepo{err: broken})
	err := svc.Verify(context.Background())

	var apiErr *util.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusServiceUnavailable || apiErr.Code != "AUDIT_CHAIN_BROKEN" {
		t.Fatalf("broken chain did not map through sentinel: %v", err)
	}
}
