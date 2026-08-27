package model

import (
	"strings"
	"testing"
	"time"

	"biosample-cold-custody-tracking/backend/internal/constants"
)

func TestPreparedTransferRejectsIdempotentResolution(t *testing.T) {
	prepared := CustodyTransfer{State: constants.TransferStatePrepared}
	if prepared.CanResolveTo(constants.TransferStatePrepared) {
		t.Fatal("prepared transfer incorrectly permits prepared-to-prepared migration")
	}
}

func TestResolvedTransferCannotBeResolvedAgain(t *testing.T) {
	accepted := CustodyTransfer{State: constants.TransferStateAccepted}
	if accepted.CanResolveTo(constants.TransferStateAccepted) {
		t.Fatal("accepted transfer incorrectly permits repeat resolution")
	}
	if accepted.CanResolveTo(constants.TransferStateCancelled) {
		t.Fatal("accepted transfer incorrectly switches to cancelled")
	}
}

func TestPreparedTransferAllowsOnlyTerminalResolutions(t *testing.T) {
	prepared := CustodyTransfer{State: constants.TransferStatePrepared}
	for _, state := range []constants.TransferState{
		constants.TransferStateAccepted,
		constants.TransferStateRejected,
		constants.TransferStateCancelled,
	} {
		if !prepared.CanResolveTo(state) {
			t.Fatalf("prepared cannot resolve to %s", state)
		}
	}
}

func TestPreparedTransferCannotCarryResolutionMetadata(t *testing.T) {
	resolvedByID := uint(9)
	resolvedAt := time.Now()
	preparedAt := resolvedAt.Add(-time.Minute)
	item := CustodyTransfer{
		SpecimenID:     5,
		TransferNo:     "CT-2026-PREP",
		FromCustodian:  "样本接收员",
		ToCustodian:    "冻存保管员",
		FromLocation:   "intake",
		ToLocation:     "样本库 B 区",
		PreparedByID:   1,
		PreparedByName: "样本接收员",
		PreparedAt:     preparedAt,
		State:          constants.TransferStatePrepared,
		AcceptedByID:   &resolvedByID,
		AcceptedByName: "提前受理人",
		ResolvedAt:     &resolvedAt,
	}

	err := item.Validate()
	if err == nil || !strings.Contains(err.Error(), "resolution metadata") {
		t.Fatalf("prepared transfer accepted resolution metadata: err=%v", err)
	}
}
