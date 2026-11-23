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

// DerefSlice copies the values pointed to into a slice of values.
func DerefSlice[E any](list []*E) []E {
	out := make([]E, len(list))
	for i, v := range list {
		out[i] = *v
	}
	return out
}
