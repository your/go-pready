package main

import "regexp"

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

func uniqueSlice(slice *[]string) {
	uniqueSlice := make(map[string]struct{}, len(*slice))
	j := 0
	for _, string := range *slice {
		if _, seen := uniqueSlice[string]; seen {
			continue
		}
		uniqueSlice[string] = struct{}{}
		(*slice)[j] = string
		j++
	}
	*slice = (*slice)[:j]
}

func removeStringFromSlice(str string, slice *[]string) {
	i := stringIndexOfSlice(str, *slice)
	*slice = append((*slice)[:i], (*slice)[i+1:]...)
}

func stringIndexOfSlice(str string, slice []string) int {
	for i, v := range slice {
		if v == str {
			return i
		}
	}
	return -1
}

func stringMatchesRegex(str string, r *regexp.Regexp) bool {
	return r.MatchString(str)
}
