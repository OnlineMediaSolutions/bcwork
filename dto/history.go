package dto

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/m6yf/bcwork/models"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/volatiletech/null/v8"
)

type HistoryModelExtended struct {
	models.History    `boil:",bind"`
	FirstName         null.String `boil:"first_name" json:"first_name"`
	LastName          null.String `boil:"last_name" json:"last_name"`
	DemandPartnerName null.String `boil:"demand_partner_name" json:"demand_partner_name"`
}

type History struct {
	ID                int       `json:"id"`
	Date              time.Time `json:"date"`
	UserFullName      string    `json:"user_full_name"`
	Action            string    `json:"action"`
	Subject           string    `json:"subject"`
	Item              string    `json:"item"`
	Changes           []Changes `json:"children"`
	DemandPartnerName *string   `json:"demand_partner_name"`
}

type Changes struct {
	ID       string `json:"id"`
	Property string `json:"property"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

func (h *History) FromModel(mod *HistoryModelExtended) error {
	h.ID = mod.ID
	h.Date = mod.Date
	h.UserFullName = func() string {
		switch mod.UserID {
		case supertokens_module.WorkerUserID:
			return supertokens_module.WorkerUserName
		case supertokens_module.AutomationUserID:
			return supertokens_module.AutomationUserName
		default:
			return mod.FirstName.String + " " + mod.LastName.String
		}
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

	if mod.DemandPartnerName.Valid {
		h.DemandPartnerName = &mod.DemandPartnerName.String
	}

	return nil
}
