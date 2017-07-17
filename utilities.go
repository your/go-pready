package main

func isStringInSlice(str string, slice []string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func stringOccurencesInSlice(str string, slice []string) int {
	count := 0
	for _, v := range slice {
		if v == str {
			count++
		}
	}
	return count
}

func uniqueSlice(slice []string) []string {
	uniqueSlice := make(map[string]struct{}, len(slice))
	j := 0
	for _, string := range slice {
		if _, seen := uniqueSlice[string]; seen {
			continue
		}
		uniqueSlice[string] = struct{}{}
		slice[j] = string
		j++
	}
	return slice[:j]
}
