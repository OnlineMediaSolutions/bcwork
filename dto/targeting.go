package dto

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/volatiletech/null/v8"
)

const (
	TargetingStatusActive   = "Active"
	TargetingStatusPaused   = "Paused"
	TargetingStatusArchived = "Archived"

	TargetingPriceModelCPM      = "CPM"
	TargetingPriceModelRevShare = "Rev Share"

	TargetingMinValueCostModelCPM = 0
	TargetingMaxValueCostModelCPM = 50

	TargetingMinValueCostModelRevShare = 0
	TargetingMaxValueCostModelRevShare = 1

	JSTagHeaderTemplate = "<!-- HTML Tag for publisher='{{ .PublisherName }}', domain='{{ .Domain }}', size='{{ .UnitSize }}', {{ if .KV }}{{ range $key, $value := .KV }}{{ if ne $value \".*\" }}{{ $key }}='{{ $value }}', {{ end }}{{ end }}{{ end }}exported='{{ .DateOfExport }}' -->\n"
	JSTagBodyTemplate   = "<script src=\"https://rt.marphezis.com/js?pid={{ .PublisherID }}&size={{ .UnitSize }}&dom={{ .Domain }}{{ if .KV }}{{ range $key, $value := .KV }}{{ if ne $value \".*\" }}&{{ $key }}={{ $value }}{{ end }}{{ end }}{{ end }}{{ if .AddGDPR }}&gdpr=${GDPR}&gdpr_concent=${GDPR_CONSENT_883}{{ end }}\"></script>"
)

var (
	tmpl = template.Must(
		template.New("JSTag").
			Parse(JSTagHeaderTemplate + JSTagBodyTemplate),
	)
	ErrFoundDuplicate = errors.New("found duplicate")
)

type Tags struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

type Targeting struct {
	ID            int               `json:"id"`
	RuleID        string            `json:"rule_id"`
	PublisherID   string            `json:"publisher_id" validate:"required"`
	PublisherName string            `json:"publisher_name"`
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
	DailyCap      *int              `json:"daily_cap"`
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

	targeting := &models.Targeting{
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
		Status:        t.Status,
		DailyCap:      null.IntFromPtr(t.DailyCap),
	}

	ruleID, err := CalculateTargetingRuleID(targeting)
	if err != nil {
		return nil, err
	}

	targeting.RuleID = ruleID

	return targeting, nil
}

func (t *Targeting) FromModel(mod *models.Targeting) error {
	var kv map[string]string
	if mod.KV.Valid {
		err := json.Unmarshal(mod.KV.JSON, &kv)
		if err != nil {
			return err
		}
	}

	var (
		country    = []string{}
		deviceType = []string{}
		browser    = []string{}
		os         = []string{}
	)
	if mod.Country != nil {
		country = mod.Country
	}
	if mod.DeviceType != nil {
		deviceType = mod.DeviceType
	}
	if mod.Browser != nil {
		browser = mod.Browser
	}
	if mod.Os != nil {
		os = mod.Os
	}

	if mod.R != nil && mod.R.Publisher != nil {
		t.PublisherName = mod.R.Publisher.Name
	}
	t.ID = mod.ID
	t.RuleID = mod.RuleID
	t.PublisherID = mod.PublisherID
	t.Domain = mod.Domain
	t.UnitSize = mod.UnitSize
	t.PlacementType = mod.PlacementType.String
	t.Country = country
	t.DeviceType = deviceType
	t.Browser = browser
	t.OS = os
	t.KV = kv
	t.PriceModel = mod.PriceModel
	t.Value = mod.Value
	t.Status = mod.Status
	t.DailyCap = mod.DailyCap.Ptr()

	return nil
}

func CalculateTargetingRuleID(mod *models.Targeting) (string, error) {
	formula, err := GetTargetingRegExp(mod)
	if err != nil {
		return "", err
	}

	return bcguid.NewFrom(formula), nil
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

	// now KV limit to one key named "oms"
	// if no KV was provided, "oms" will be set to "all"
	defaultKV := "__oms=.*"

	return baseRegExp + defaultKV, nil
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

func GetJSTagString(mod *models.Targeting, addGDPR bool) (string, error) {
	type tag struct {
		PublisherID   string
		PublisherName string
		UnitSize      string
		Domain        string
		KV            map[string]string
		AddGDPR       bool
		DateOfExport  string
	}

	var kv map[string]string
	if mod.KV.Valid {
		err := json.Unmarshal(mod.KV.JSON, &kv)
		if err != nil {
			return "", err
		}
	}

	var publisherName string
	if mod.R.GetPublisher() != nil {
		publisherName = mod.R.GetPublisher().Name
	}

	data := tag{
		PublisherID:   mod.PublisherID,
		PublisherName: publisherName,
		UnitSize:      mod.UnitSize,
		Domain:        mod.Domain,
		KV:            kv,
		AddGDPR:       addGDPR,
		DateOfExport:  time.Now().Format(time.DateOnly),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
