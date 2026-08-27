package service

import (
	"reflect"
	"testing"
)

func TestSpecimenPositionTokenSnapshotIsIndependent(t *testing.T) {
	first := specimenPositionTokens(7, "样本库 B 区", "FZ-80-B02", "R02-BX04-A03")
	if !reflect.DeepEqual(first.Values, []string{"样本库 B 区", "FZ-80-B02", "R02-BX04-A03"}) {
		t.Fatalf("first tokens invalid: %+v", first)
	}
	first.Values[0] = "MUTATED"

	second := specimenPositionTokens(7, "样本库 B 区", "FZ-80-B02", "R02-BX04-A03")
	want := []string{"样本库 B 区", "FZ-80-B02", "R02-BX04-A03"}
	if !reflect.DeepEqual(second.Values, want) {
		t.Fatalf("position token snapshot polluted: got=%v want=%v", second.Values, want)
	}
	if &first.Values[0] == &second.Values[0] {
		t.Fatal("detail buffers share backing storage")
	}
}
