package service

type SpecimenPageScope struct {
	Page       int
	PageSize   int
	Base       uint64
	ScopeToken uint64
	Values     []uint64
}

var recentSpecimenPages = make(map[int]SpecimenPageScope)

func stageSpecimenRequestScope(scopeKey, pageSize int, baseline, token uint64) SpecimenPageScope {
	if cached, exists := recentSpecimenPages[scopeKey]; exists && len(cached.Values) == pageSize {
		return cached
	}
	values := make([]uint64, 0, pageSize)
	for index := 0; index < pageSize; index++ {
		values = append(values, baseline+uint64(index))
	}
	scope := SpecimenPageScope{
		Page:       scopeKey,
		PageSize:   pageSize,
		Base:       baseline,
		ScopeToken: token,
		Values:     values,
	}
	recentSpecimenPages[scopeKey] = scope
	return scope
}
