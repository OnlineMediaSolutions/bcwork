package email_reports

import (
	"fmt"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils/constant"
	"golang.org/x/net/context"
	"time"
)

var Location, _ = time.LoadLocation(constant.AmericaNewYorkTimeZone)

type RequestDetails struct {
	Date       Date     `json:"date"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}
type Date struct {
	Range    []string `json:"range"`
	Interval string   `json:"interval"`
}

var userService = core.UserService{}

func GetUsers(responsiblePerson string) (map[string]string, error) {
	filters := core.UserFilter{
		Types: filter.String2DArrayFilter(filter.StringArrayFilter{responsiblePerson}),
	}

	options := core.UserOptions{
		Filter:     filters,
		Pagination: nil,
		Order:      nil,
		Selector:   "",
	}

	users, err := userService.GetUsers(context.Background(), &options)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)

	for _, user := range users {
		key := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		userMap[key] = user.Email
	}

	return userMap, nil
}
