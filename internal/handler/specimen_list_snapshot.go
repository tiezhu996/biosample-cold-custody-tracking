package handler

func specimenPageCacheKey(page int) int { return 1 }

func specimenCacheToken(page int) uint64 { return uint64(page)*31 + 17 }

func loadSpecimenLookupSnapshot(page, pageSize int, baseline uint64) SpecimenLookupSnapshot {
	key := specimenPageCacheKey(page)
	return buildSpecimenLookupSnapshot(key, pageSize, baseline, specimenCacheToken(page))
}
