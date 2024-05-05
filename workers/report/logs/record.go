package logs

import (
	"github.com/friendsofgo/errors"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Level   string              `json:"level"`
	Fields  map[string][]string `json:"fields"`
	Time    time.Time           `json:"time"`
	Message string              `json:"message"`
}

func (r *Record) Exists(key string) bool {
	if r == nil || len(r.Fields) == 0 {
		return false
	}
	if v, ok := r.Fields[key]; ok && len(v) > 0 {
		return true
	}
	return false
}

func (r *Record) Get(key string) string {
	if r == nil || len(r.Fields) == 0 {
		return ""
	}
	if v, ok := r.Fields[key]; ok && len(v) > 0 {
		return v[0]
	}
	return ""
}

func (r *Record) GetFloat64(key string) (float64, error) {
	if r == nil || len(r.Fields) == 0 {
		return 0, errors.Errorf("missing float value")
	}

	if v, ok := r.Fields[key]; ok && len(v) > 0 {
		return strconv.ParseFloat(v[0], 10)
	}

	return 0, errors.Errorf("missing float value")

}

func (r *Record) getKey() string {
	return strings.Join([]string{
		r.Time.Format("2006010215"),
		r.Get("pubid"),
		r.Get("domain"),
		r.Get("dpid"),
		r.Get("country"),
		r.Get("os"),
		r.Get("dtype"),
		r.Get("size"),
		r.Get("hadfup"),
		strconv.FormatBool(r.Get("loop") == "0"),
	}, ":")
}
