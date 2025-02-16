package omserr

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type ErrorContext struct {
	Err     error             `json:"err"`
	Context map[string]string `json:"context"`
}

func (errctx ErrorContext) GetValue(key string) string {
	return errctx.Context[key]
}

func (errctx ErrorContext) GetValueDefault(key string, def string) string {
	val, found := errctx.Context[key]
	if !found {
		return def
	}

	return val
}

func (errctx ErrorContext) GetNumericValue(key string, def int) int {
	val, found := errctx.Context[key]
	if !found || val == "" {
		return def
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		log.Info().Err(err).Str("val", val).Msg("failed to convert ErrorContext numeric value to an int")

		return def
	}

	return n
}

func ExtractErrContext(err error) ErrorContext {
	tokens := extractContext(err.Error())
	res := ErrorContext{
		Err:     err,
		Context: map[string]string{},
	}

	for _, tok := range tokens {
		tokVals := extractContextValues(tok)
		for k, v := range tokVals {
			if _, found := res.Context[k]; !found {
				res.Context[k] = v
			}
		}
	}

	return res
}

func extractContext(text string) []string {
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

func extractContextValues(text string) map[string]string {
	res := map[string]string{}
	tokens := strings.Split(text, ",")
	for _, tok := range tokens {
		i := strings.Index(tok, ":")
		if i > 0 {
			res[tok[:i]] = tok[i+1:]
		}
	}

	return res
}
