package core

func SliceIsSubset[S ~[]E, E comparable](subset S, superset S) bool {
	checkset := make(map[E]bool)
	for _, subVal := range subset {
		checkset[subVal] = true
	}
	for _, supVal := range superset {
		if checkset[supVal] {
			delete(checkset, supVal)
		}
	}
	return len(checkset) == 0
}
