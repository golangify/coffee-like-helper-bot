package search

import "strings"

type replacable struct {
	src string
	dst string
}

var replacables = []replacable{
	{"ё", "е"},
	{"й", "и"},
	{"э", "е"},
	{"сс", "с"},
	{"а", "о"},
}

func prepareQueryString(queryString string) string {

	queryRuneSlice := []rune(strings.ToLower(queryString))
	if len(queryRuneSlice) > 100 {
		queryRuneSlice = queryRuneSlice[:100]
	}
	queryString = string(queryRuneSlice)

	for _, r := range replacables {
		queryString = strings.ReplaceAll(queryString, r.src, r.dst)
	}

	return queryString
}
