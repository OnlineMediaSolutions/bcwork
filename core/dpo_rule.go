package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"sort"
	"strings"
	"time"
)

type DemandPartnerOptimizationRule struct {
	RuleID        string `json:"rule_id"`
	DemandPartner string `json:"demand_partners"`
	Publisher     string `json:"publisher,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Country       string `json:"country,omitempty"`
	OS            string `json:"os,omitempty"`
	DeviceType    string `json:"device_type,omitempty"`
	PlacementType string `json:"placement_type,omitempty"`
	Browser       string `json:"browser,omitempty"`

	Factor float64 `json:"factor"`
}

type DemandPartnerOptimizationRuleSlice []*DemandPartnerOptimizationRule

func (dpo *DemandPartnerOptimizationRule) Save(ctx context.Context) (string, error) {

	mod := dpo.ToModel()
	err := mod.Upsert(ctx, bcdb.DB(), true, []string{models.DpoRuleColumns.RuleID}, boil.Infer(), boil.Infer())
	if err != nil {
		return "", eris.Wrapf(err, "failed to updsert dpo rule(rule=%s)", dpo.GetFormula())
	}

	return mod.RuleID, nil

}

func (dpo *DemandPartnerOptimizationRule) FromModel(mod *models.DpoRule) {

	dpo.RuleID = mod.RuleID
	dpo.DemandPartner = mod.DemandPartnerID
	dpo.Factor = mod.Factor

	if mod.Publisher.Valid {
		dpo.Publisher = mod.Publisher.String
	}

	if mod.Domain.Valid {
		dpo.Domain = mod.Domain.String
	}

	if mod.Country.Valid {
		dpo.Country = mod.Country.String
	}

	if mod.Os.Valid {
		dpo.OS = mod.Os.String
	}

	if mod.DeviceType.Valid {
		dpo.DeviceType = mod.DeviceType.String
	}

	if mod.PlacementType.Valid {
		dpo.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		dpo.Browser = mod.Browser.String
	}
}

func (dpos *DemandPartnerOptimizationRuleSlice) FromModel(slice models.DpoRuleSlice) {

	for _, mod := range slice {
		dpo := DemandPartnerOptimizationRule{}
		dpo.FromModel(mod)
		*dpos = append(*dpos, &dpo)
	}

}

func (dpo *DemandPartnerOptimizationRule) ToModel() *models.DpoRule {

	mod := models.DpoRule{
		RuleID:          dpo.GetRuleID(),
		DemandPartnerID: dpo.DemandPartner,
		Factor:          dpo.Factor,
	}

	if dpo.Publisher != "" {
		mod.Publisher = null.StringFrom(dpo.Publisher)
	}

	if dpo.Domain != "" {
		mod.Domain = null.StringFrom(dpo.Domain)
	}

	if dpo.Country != "" {
		mod.Country = null.StringFrom(dpo.Country)
	}

	if dpo.OS != "" {
		mod.Os = null.StringFrom(dpo.OS)
	}

	if dpo.DeviceType != "" {
		mod.DeviceType = null.StringFrom(dpo.DeviceType)
	}

	if dpo.PlacementType != "" {
		mod.PlacementType = null.StringFrom(dpo.PlacementType)
	}

	if dpo.Browser != "" {
		mod.Browser = null.StringFrom(dpo.Browser)
	}

	return &mod

}

func (dpo *DemandPartnerOptimizationRule) GetRuleID() string {
	return bcguid.NewFrom(dpo.GetFormula())
}

func (dpo *DemandPartnerOptimizationRule) GetFormula() string {
	dp := dpo.DemandPartner
	if dp == "" {
		dp = "*"
	}

	p := dpo.Publisher
	if p == "" {
		p = "*"
	}

	d := dpo.Domain
	if d == "" {
		d = "*"
	}

	c := dpo.Country
	if c == "" {
		c = "*"
	}

	os := dpo.OS
	if os == "" {
		os = "*"
	}

	dt := dpo.DeviceType
	if dt == "" {
		dt = "*"
	}

	pt := dpo.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := dpo.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("dp=%s__p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", dp, p, d, c, os, dt, pt, b)
}

func (dpo *DemandPartnerOptimizationRule) GetFormulaRegex() string {

	p := dpo.Publisher
	if p == "" {
		p = ".*"
	}

	d := dpo.Domain
	if d == "" {
		d = ".*"
	}

	c := dpo.Country
	if c == "" {
		c = ".*"
	}

	os := dpo.OS
	if os == "" {
		os = ".*"
	}

	dt := dpo.DeviceType
	if dt == "" {
		dt = ".*"
	}

	pt := dpo.PlacementType
	if pt == "" {
		pt = ".*"
	}

	b := dpo.Browser
	if b == "" {
		b = ".*"
	}

	return fmt.Sprintf("(p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s)", p, d, c, os, dt, pt, b)
}

type DpoRealtimeRecord struct {
	RuleID string  `json:"rule_id"`
	Rule   string  `json:"rule"`
	Factor float64 `json:"factor"`
}

type DpoRealtimeRecordSlice []*DpoRealtimeRecord

type DpoRT struct {
	DemandPartnerID string                 `json:"demand_partner_id"`
	IsInclude       bool                   `json:"is_include"`
	Rules           DpoRealtimeRecordSlice `json:"rules"`
}

func (dpo *DemandPartnerOptimizationRule) ToRtRule() *DpoRealtimeRecord {
	return &DpoRealtimeRecord{
		RuleID: dpo.RuleID,
		Rule:   dpo.GetFormulaRegex(),
		Factor: dpo.Factor,
	}
}

func (dpos DpoRealtimeRecordSlice) Sort() {
	sort.Slice(dpos, func(i, j int) bool {
		return strings.Count(dpos[i].Rule, "*") < strings.Count(dpos[j].Rule, "*")
	})
}

func SendToRT(ctx context.Context, demandPartnerID string) error {
	modDpos, err := models.DpoRules(models.DpoRuleWhere.DemandPartnerID.EQ(demandPartnerID)).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "failed to fetch dpo rules(dpid:%)", demandPartnerID)
	}

	if len(modDpos) == 0 {
		return nil
	}

	dpos := make(DemandPartnerOptimizationRuleSlice, 0, 0)
	dpos.FromModel(modDpos)

	dposRT := DpoRT{
		DemandPartnerID: demandPartnerID,
		IsInclude:       false,
	}
	for _, dpo := range dpos {
		dposRT.Rules = append(dposRT.Rules, dpo.ToRtRule())
	}
	dposRT.Rules.Sort()

	b, err := json.Marshal(dposRT)
	if err != nil && err != sql.ErrNoRows {
		return eris.Cause(err)
	}

	modMeta := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(time.Now()),
		Key:           "dpo:" + demandPartnerID,
		Value:         b,
	}

	err = modMeta.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrapf(err, "failed to insert metadata record")
	}

	return nil
}
