package core

import (
	"strings"
)

func ExtractTags(key string) map[string]string {
	tokens := extractTags(key)
	res := make(map[string]string)

	for _, tok := range tokens {
		tokVals := extractTagsValues(tok)
		for k, v := range tokVals {
			if _, found := res[k]; !found {
				res[k] = v
			}
		}
	}

	return res
}

func ExtractKey(key string) string {
	return strings.TrimSpace(strings.Split(key, "(")[0])
}

func extractTags(text string) []string {
	res := []string{}

	in := false
	start := 0
	nested := 0
	for i, s := range text {
		if !in {
			if s == '(' {
				in = true
				start = i + 1
			}
		} else {
			if s == '(' {
				nested++
				continue
			}
			if s == ')' {
				if nested > 0 {
					nested--
					continue
				}

				res = append(res, text[start:i])
				in = false
			}
		}
	}

	return res
}

func extractTagsValues(text string) map[string]string {
	res := map[string]string{}
	tokens := strings.Split(text, ",")
	for _, tok := range tokens {
		if tok[0] == '$' || tok[0] == '~' {
			continue
		}
		i := strings.Index(tok, "=")
		if i > 0 {
			res[tok[:i]] = tok[i+1:]
		}

	}
	return res
}
