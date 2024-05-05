package order

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Field struct {
	Name string      `json:"name"`
	Desc bool        `json:"desc"`
	Data interface{} `json:"data"`
}

type Sort []Field
type CustomSort map[string]func(Field) (string, error)

func (s Sort) Do(custom CustomSort, primaryKey string) qm.QueryMod {

	if len(s) == 0 {
		return qm.And("TRUE")
	}

	res := make([]string, 0)
	for _, sf := range s {
		if customSortField, found := custom[sf.Name]; found {
			name, err := customSortField(sf)
			if err != nil {
				log.Warn().Err(err).Msg("sorting error")
			} else {
				sf.Name = name
			}
		}
		if sf.Desc {
			res = append(res, fmt.Sprintf(`%s DESC NULLS LAST`, sf.Name))
		} else {
			res = append(res, fmt.Sprintf(`%s ASC NULLS LAST`, sf.Name))
		}
	}

	//always add primary key to the sort so pagination is not broken on equal values
	res = append(res, primaryKey)

	return qm.OrderBy(strings.Join(res, ","))
}
