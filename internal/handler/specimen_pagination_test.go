package handler

import (
	"reflect"
	"testing"
)

func TestSpecimenListPageLookupSnapshotDistinctBaseline(t *testing.T) {
	first := loadSpecimenLookupSnapshot(1, 4, 100)
	if first.CacheToken != uint64(48) || !reflect.DeepEqual(first.Values, []uint64{100, 101, 102, 103}) {
		t.Fatalf("first snapshot invalid: %+v", first)
	}
	first.Values[0] = 999999

	samePage := loadSpecimenLookupSnapshot(1, 4, 100)
	if samePage.CacheToken != uint64(48) || !reflect.DeepEqual(samePage.Values, []uint64{100, 101, 102, 103}) {
		t.Fatalf("same-page snapshot polluted: %+v", samePage)
	}

	otherPage := loadSpecimenLookupSnapshot(2, 4, 200)
	if otherPage.Page != 2 || otherPage.CacheToken != uint64(79) || !reflect.DeepEqual(otherPage.Values, []uint64{200, 201, 202, 203}) {
		t.Fatalf("cross-page cache key polluted: %+v", otherPage)
	}
}
