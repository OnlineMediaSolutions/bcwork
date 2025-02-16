package email_reports

import (
	"fmt"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/core"
	"golang.org/x/net/context"
)

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
