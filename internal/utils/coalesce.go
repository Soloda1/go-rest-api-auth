package utils

func CoalesceString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func CoalesceSliceStrings(a, b []string) []string {
	if a != nil {
		return a
	}
	return b
}
