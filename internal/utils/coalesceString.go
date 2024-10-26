package utils

func CoalesceString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
