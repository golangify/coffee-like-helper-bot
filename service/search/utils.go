package search

import "strings"

func prepareQueryString(queryString string) string {
	tokens := map[string]string{
		"е":  "ё",
		"й":  "и",
		"э":  "ё",
		"сс": "с",
		"а":  "о",
	}

	queryRuneSlice := []rune(strings.ToLower(queryString))
	if len(queryRuneSlice) > 100 {
		queryRuneSlice = queryRuneSlice[:100]
	}
	queryString = string(queryRuneSlice)

	for s, d := range tokens {
		queryString = strings.ReplaceAll(queryString, s, d)
	}

	return queryString
}
