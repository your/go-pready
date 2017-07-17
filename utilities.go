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
