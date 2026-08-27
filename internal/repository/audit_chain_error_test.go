package repository

import (
	"errors"
	"testing"
	"time"

	"biosample-cold-custody-tracking/backend/internal/model"
)

func sealedEntry(t *testing.T, id uint, previousHash string) model.AuditLog {
	t.Helper()
	entry := model.AuditLog{
		ID:         id,
		CreatedAt:  time.Unix(1760000000+int64(id), 0).UTC(),
		RequestID:  "chain-check",
		ActorID:    id,
		ActorName:  "审计员",
		Action:     "verify",
		EntityType: "Specimen",
		EntityID:   id,
	}
	entry.Seal(previousHash)
	return entry
}

func TestVerifyChainBrokenEntryIsAuditable(t *testing.T) {
	first := sealedEntry(t, 1, "")
	second := sealedEntry(t, 2, first.EntryHash)
	broken := sealedEntry(t, 3, "tampered")
	_, err := verifyAuditChain([]model.AuditLog{first, second, broken})

	if !errors.Is(err, ErrAuditChainBroken) {
		t.Fatalf("broken-chain error = %v, want errors.Is %v", err, ErrAuditChainBroken)
	}
	if got := errors.Unwrap(err); got != ErrAuditChainBroken {
		t.Fatalf("unwrapped = %v", got)
	}
}

func TestVerifyHealthyAuditChainKeepsTailHash(t *testing.T) {
	first := sealedEntry(t, 1, "")
	second := sealedEntry(t, 2, first.EntryHash)
	tail, err := verifyAuditChain([]model.AuditLog{first, second})
	if err != nil {
		t.Fatalf("healthy chain returned %v", err)
	}
	if tail != second.EntryHash {
		t.Fatalf("tail=%q want=%q", tail, second.EntryHash)
	}
}
