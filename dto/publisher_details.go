package dto

import (
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type PublisherDetailModel struct {
	Publisher       models.Publisher       `boil:"publisher,bind"`
	PublisherDomain models.PublisherDomain `boil:"publisher_domain,bind"`
	User            UserModelCompact       `boil:"user,bind"`
}

// UserModelCompact to support possible null string in first and last names
type UserModelCompact struct {
	ID        int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	FirstName null.String `boil:"first_name" json:"first_name" toml:"first_name" yaml:"first_name"`
	LastName  null.String `boil:"last_name" json:"last_name" toml:"last_name" yaml:"last_name"`
}

type PublisherDetail struct {
	Name                   string  `json:"name"`
	PublisherID            string  `json:"publisher_id"`
	Domain                 string  `json:"domain"`
	AccountManagerID       string  `json:"account_manager_id"`
	AccountManagerFullName string  `json:"account_manager_full_name"`
	Automation             bool    `json:"automation"`
	GPPTarget              float64 `json:"gpp_target"`
	ActivityStatus         string  `json:"activity_status"`
}

func (pd *PublisherDetail) FromModel(mod *PublisherDetailModel, activityStatus map[string]map[string]ActivityStatus) error {
	pd.Name = mod.Publisher.Name
	pd.PublisherID = mod.Publisher.PublisherID
	pd.Domain = mod.PublisherDomain.Domain
	pd.AccountManagerID = mod.Publisher.AccountManagerID.String
	pd.AccountManagerFullName = buildFullName(mod.User)
	pd.Automation = mod.PublisherDomain.Automation
	pd.GPPTarget = mod.PublisherDomain.GPPTarget.Float64
	pd.ActivityStatus = activityStatus[pd.Domain][pd.PublisherID].String()
	return nil
}

type PublisherDetailsSlice []*PublisherDetail

func (pds *PublisherDetailsSlice) FromModel(mods []*PublisherDetailModel, activityStatus map[string]map[string]ActivityStatus) error {
	for _, mod := range mods {
		pd := PublisherDetail{}
		err := pd.FromModel(mod, activityStatus)
		if err != nil {
			return eris.Cause(err)
		}
		*pds = append(*pds, &pd)
	}

	return nil
}

func buildFullName(user UserModelCompact) string {
	if user.FirstName.Valid && user.LastName.Valid {
		return user.FirstName.String + " " + user.LastName.String
	}
	return ""
}
