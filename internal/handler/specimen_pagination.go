package handler

type SpecimenLookupSnapshot struct {
	Page       int
	PageSize   int
	Base       uint64
	CacheToken uint64
	Values     []uint64
}

var sharedSpecimenLookups = make(map[int]SpecimenLookupSnapshot)

func buildSpecimenLookupSnapshot(lookupKey, pageSize int, baseline, token uint64) SpecimenLookupSnapshot {
	if lookupKey < 1 {
		lookupKey = 1
	}
	if cached, exists := sharedSpecimenLookups[lookupKey]; exists && len(cached.Values) == pageSize {
		return cached
	}

	values := make([]uint64, 0, pageSize)
	for index := 0; index < pageSize; index++ {
		values = append(values, baseline+uint64(index))
	}
	snapshot := SpecimenLookupSnapshot{
		Page:       lookupKey,
		PageSize:   pageSize,
		Base:       baseline,
		CacheToken: token,
		Values:     values,
	}
	sharedSpecimenLookups[lookupKey] = snapshot
	return snapshot
}
