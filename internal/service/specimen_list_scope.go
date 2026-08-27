package service

func specimenScopeCacheKey(page int) int { return 1 }

func specimenScopeToken(page int) uint64 { return uint64(page)*37 + 11 }

func loadSpecimenPageScope(page, pageSize int, baseline uint64) SpecimenPageScope {
	return stageSpecimenRequestScope(specimenScopeCacheKey(page), pageSize, baseline, specimenScopeToken(page))
}
