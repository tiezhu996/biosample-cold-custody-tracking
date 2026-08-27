package service

import (
	"reflect"
	"testing"
)

func TestSpecimenListPageScopeIsIndependentAcrossRequests(t *testing.T) {
	first := loadSpecimenPageScope(1, 4, 100)
	if first.ScopeToken != uint64(48) || !reflect.DeepEqual(first.Values, []uint64{100, 101, 102, 103}) {
		t.Fatalf("first scope invalid: %+v", first)
	}
	first.Values[0] = 999999

	samePage := loadSpecimenPageScope(1, 4, 100)
	if !reflect.DeepEqual(samePage.Values, []uint64{100, 101, 102, 103}) || samePage.ScopeToken != uint64(48) {
		t.Fatalf("same-page scope polluted: %+v", samePage)
	}

	otherPage := loadSpecimenPageScope(2, 4, 200)
	wantOther := []uint64{200, 201, 202, 203}
	if otherPage.Page != 2 || otherPage.ScopeToken != uint64(85) || !reflect.DeepEqual(otherPage.Values, wantOther) {
		t.Fatalf("cross-page scope key polluted: %+v", otherPage)
	}
}
