package dto

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/m6yf/bcwork/models"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
)

type History struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	UserID       int       `json:"user_id"`
	UserFullName string    `json:"user_full_name"`
	Action       string    `json:"action"`
	Subject      string    `json:"subject"`
	Item         string    `json:"item"`
	Changes      []Changes `json:"children"`
}

type Changes struct {
	ID       string `json:"id"`
	Property string `json:"property"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

func (h *History) FromModel(mod *models.History, usersMap map[int]string) error {
	h.ID = mod.ID
	h.Date = mod.Date
	h.UserID = mod.UserID
	h.UserFullName = func() string {
		if mod.UserID == supertokens_module.WorkerUserID {
			return supertokens_module.WorkerUserName
		}
		return usersMap[mod.UserID]
	}()
	h.Action = mod.Action
	h.Subject = mod.Subject
	h.Item = mod.Item

	if mod.Changes.Valid {
		err := json.Unmarshal(mod.Changes.JSON, &h.Changes)
		if err != nil {
			return err
		}

		counter := 1
		for i := range h.Changes {
			h.Changes[i].ID = strconv.Itoa(h.ID) + "-" + strconv.Itoa(counter)
			counter++
		}
	}

	return nil
}
