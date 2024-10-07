package constant

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/volatiletech/null/v8"
)

const (
	TargetingStatusActive   = "active"
	TargetingStatusPaused   = "paused"
	TargetingStatusArchived = "archived"

	TargetingPriceModelCPM      = "CPM"
	TargetingPriceModelRevShare = "RevShare"

	TargetingMinValueCostModelCPM = 0
	TargetingMaxValueCostModelCPM = 50

	TargetingMinValueCostModelRevShare = 0
	TargetingMaxValueCostModelRevShare = 1
)

type Targeting struct {
	Hash          string            `json:"-"`
	RuleID        string            `json:"-"`
	Publisher     string            `json:"publisher" validate:"required"`
	Domain        string            `json:"domain" validate:"required"`
	UnitSize      string            `json:"unit_size" validate:"required"`
	PlacementType string            `json:"placement_type" validate:"required"`
	Country       []string          `json:"country" validate:"countries"`
	DeviceType    []string          `json:"device_type" validate:"devices"`
	Browser       []string          `json:"browser"`
	OS            []string          `json:"os"`
	KV            map[string]string `json:"kv,omitempty"`
	PriceModel    string            `json:"price_model" validate:"targetingPriceModel"`
	Value         float64           `json:"value"`
	DailyCap      int               `json:"daily_cap"`
	Status        string            `json:"status" validate:"targetingStatus"`
}

func (t *Targeting) PrepareData() {
	hash := bcguid.NewFromf(t.Publisher, t.Domain, t.UnitSize, t.PlacementType)
	ruleID := bcguid.NewFromf(hash, t.Country, t.DeviceType, t.Browser, t.OS, t.KV)

	t.Hash = hash
	t.RuleID = ruleID

	slices.SortStableFunc(t.Country, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.DeviceType, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.Browser, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.OS, func(a, b string) int { return cmp.Compare(a, b) })

	if t.Status == "" {
		t.Status = TargetingStatusActive
	}
}

func (t Targeting) ToModel() (*models.Targeting, error) {
	var (
		kv    []byte
		valid bool
		err   error
	)

	if len(t.KV) > 0 {
		kv, err = json.Marshal(t.KV)
		if err != nil {
			return nil, err
		}
		valid = true
	}

	return &models.Targeting{
		Hash:          t.Hash,
		RuleID:        t.RuleID,
		Publisher:     null.StringFrom(t.Publisher),
		Domain:        null.StringFrom(t.Domain),
		UnitSize:      null.StringFrom(t.UnitSize),
		PlacementType: null.StringFrom(t.PlacementType),
		Country:       t.Country,
		DeviceType:    t.DeviceType,
		Browser:       t.Browser,
		Os:            t.OS,
		KV:            null.NewJSON(kv, valid),
		PriceModel:    t.PriceModel,
		Value:         t.Value,
		Status:        TargetingStatusActive,
	}, nil
}

func (t *Targeting) FromModel(mod *models.Targeting) error {
	var kv map[string]string
	if mod.KV.Valid {
		err := json.Unmarshal(mod.KV.JSON, &kv)
		if err != nil {
			return err
		}
	}

	t.Publisher = mod.Publisher.String
	t.Domain = mod.Domain.String
	t.UnitSize = mod.UnitSize.String
	t.PlacementType = mod.PlacementType.String
	t.Country = mod.Country
	t.DeviceType = mod.DeviceType
	t.Browser = mod.Browser
	t.OS = mod.Os
	t.KV = kv
	t.PriceModel = mod.PriceModel
	t.Value = mod.Value
	t.Status = mod.Status

	return nil
}

func GetTargetingRegExp(mod *models.Targeting) (string, error) {
	var kv map[string]string
	if mod.KV.Valid {
		err := json.Unmarshal(mod.KV.JSON, &kv)
		if err != nil {
			return "", err
		}
	}

	p := helpers.GetStringWithDefaultValue(mod.Publisher.String, ".*")
	d := helpers.GetStringWithDefaultValue(mod.Domain.String, ".*")
	s := helpers.GetStringWithDefaultValue(mod.UnitSize.String, ".*")
	pt := helpers.GetStringWithDefaultValue(mod.PlacementType.String, ".*")
	c := helpers.GetStringFromSliceWithDefaultValue(mod.Country, "|", ".*")
	os := helpers.GetStringFromSliceWithDefaultValue(mod.Os, "|", ".*")
	dt := helpers.GetStringFromSliceWithDefaultValue(mod.DeviceType, "|", ".*")
	b := helpers.GetStringFromSliceWithDefaultValue(mod.Browser, "|", ".*")

	baseRegExp := fmt.Sprintf("p=%s__d=%s__s=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, s, c, os, dt, pt, b)

	if len(kv) > 0 {
		kvRegExp := make([]string, 0, len(kv))
		for k, v := range kv {
			kvRegExp = append(kvRegExp, k+"="+v)
		}
		slices.SortFunc(kvRegExp, func(a, b string) int { return cmp.Compare(a, b) })
		return baseRegExp + "__" + strings.Join(kvRegExp, "__"), nil
	}

	return baseRegExp, nil
}

func GetTargetingKey(publisher, domain string) string {
	return utils.JSTagMetaDataKeyPrefix + ":" + publisher + ":" + domain
}
