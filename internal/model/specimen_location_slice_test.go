package model

import (
	"testing"

	"biosample-cold-custody-tracking/backend/internal/constants"
)

func TestSpecimenLocationLabelStableAcrossInventorySnapshot(t *testing.T) {
	containerID := uint(11)
	item := Specimen{
		State:              constants.SpecimenStateStored,
		StorageContainerID: &containerID,
		StorageContainer: &StorageContainer{
			Code:     "FZ-80-B02",
			Location: "样本库 B 区",
		},
		Position: "R02-BX04-A03",
	}

	want := "样本库 B 区 / FZ-80-B02 / R02-BX04-A03"
	for round := 0; round < 10; round++ {
		if got := item.LocationLabel(); got != want {
			t.Fatalf("round %d label = %q, want %q", round, got, want)
		}
	}
}

func TestSpecimenLocationLabelOmitsEmptyParts(t *testing.T) {
	item := Specimen{State: constants.SpecimenStateAliquoted}
	if got := item.LocationLabel(); got != "intake" {
		t.Fatalf("empty-part label = %q", got)
	}
}
