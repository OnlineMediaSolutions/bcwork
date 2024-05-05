package join

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Inner is an helper that generates QueryMod for InnerJoin operation
func Inner(left string, right string, column string) qm.QueryMod {
	return qm.InnerJoin(fmt.Sprintf("%s ON %s.%s = %s.%s", right, left, column, right, column))
}
