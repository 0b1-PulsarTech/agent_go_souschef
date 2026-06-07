package goscan

import "strings"

func BestMatch(symbols []string, query string) string {
	query = strings.ToLower(query)
	for _, symbol := range symbols {
		name := strings.ToLower(symbol)
		if name == query || strings.HasSuffix(name, "."+query) {
			return symbol
		}
	}
	for _, symbol := range symbols {
		if strings.Contains(strings.ToLower(symbol), query) {
			return symbol
		}
	}
	return ""
}
