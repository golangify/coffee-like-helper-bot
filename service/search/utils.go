package search

import "strings"

const maxQueryLength = 100

type replacement struct {
	src string
	dst string
}

var replacements = []replacement{
	{"ё", "е"},
	{"й", "и"},
	{"э", "е"},
	{"сс", "с"},
	{"а", "о"},
}

func prepareQueryString(query string) string {
	if len(query) == 0 {
		return query
	}

	query = strings.ToLower(query)

	if len(query) > maxQueryLength {
		if idx := strings.LastIndex(query[:maxQueryLength], " "); idx > 0 {
			query = query[:idx]
		} else {
			query = query[:maxQueryLength]
		}
	}

	for _, r := range replacements {
		query = strings.ReplaceAll(query, r.src, r.dst)
	}

	return query
}
