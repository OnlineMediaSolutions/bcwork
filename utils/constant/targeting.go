package constant

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
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
	ID            int               `json:"id"`
	PublisherID   string            `json:"publisher_id" validate:"required"`
	Domain        string            `json:"domain" validate:"required"`
	UnitSize      string            `json:"unit_size" validate:"required"`
	PlacementType string            `json:"placement_type"`
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
	slices.SortStableFunc(t.Country, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.DeviceType, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.Browser, func(a, b string) int { return cmp.Compare(a, b) })
	slices.SortStableFunc(t.OS, func(a, b string) int { return cmp.Compare(a, b) })

	if t.Status == "" {
		t.Status = TargetingStatusActive
	}
}

func (t Targeting) ToModel() (*models.Targeting, error) {
	modKV, err := GetModelKV(t.KV)
	if err != nil {
		return nil, err
	}

	return &models.Targeting{
		PublisherID:   t.PublisherID,
		Domain:        t.Domain,
		UnitSize:      t.UnitSize,
		PlacementType: null.StringFrom(t.PlacementType),
		Country:       t.Country,
		DeviceType:    t.DeviceType,
		Browser:       t.Browser,
		Os:            t.OS,
		KV:            modKV,
		PriceModel:    t.PriceModel,
		Value:         t.Value,
		Status:        TargetingStatusActive,
		DailyCap:      null.IntFrom(t.DailyCap),
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

	t.ID = mod.ID
	t.PublisherID = mod.PublisherID
	t.Domain = mod.Domain
	t.UnitSize = mod.UnitSize
	t.PlacementType = mod.PlacementType.String
	t.Country = mod.Country
	t.DeviceType = mod.DeviceType
	t.Browser = mod.Browser
	t.OS = mod.Os
	t.KV = kv
	t.PriceModel = mod.PriceModel
	t.Value = mod.Value
	t.Status = mod.Status
	t.DailyCap = mod.DailyCap.Int

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

	p := helpers.GetStringWithDefaultValue(mod.PublisherID, ".*")
	d := helpers.GetStringWithDefaultValue(mod.Domain, ".*")
	s := helpers.GetStringWithDefaultValue(mod.UnitSize, ".*")
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

func GetModelKV(kv map[string]string) (null.JSON, error) {
	var (
		modKV []byte
		valid bool
		err   error
	)

	if len(kv) > 0 {
		modKV, err = json.Marshal(kv)
		if err != nil {
			return null.JSON{}, err
		}
		valid = true
	}

	return null.NewJSON(modKV, valid), nil
}
