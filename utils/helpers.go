package utils

// Contais returns true if the given string is in the given slice
func Contains(arr []string, entry string) bool {
	for _, a := range arr {
		if a == entry {
			return true
		}
	}
	return false
}
