package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DemandPartnerOptimizationRule struct {
	RuleID        string  `json:"rule_id"`
	DemandPartner string  `json:"demand_partners"`
	Publisher     string  `json:"publisher,omitempty"`
	Domain        string  `json:"domain,omitempty"`
	Country       string  `json:"country,omitempty"`
	OS            string  `json:"os,omitempty"`
	DeviceType    string  `json:"device_type,omitempty"`
	PlacementType string  `json:"placement_type,omitempty"`
	Browser       string  `json:"browser,omitempty"`
	Factor        float64 `json:"factor"`
}

type DemandPartnerOptimizationRuleJoined struct {
	RuleID            string  `json:"rule_id"`
	DemandPartnerID   string  `json:"demand_partner_id"`
	Publisher         string  `json:"publisher"`
	Domain            string  `json:"domain"`
	Country           string  `json:"country"`
	OS                string  `json:"os"`
	DeviceType        string  `json:"device_type"`
	PlacementType     string  `json:"placement_type"`
	Browser           string  `json:"browser"`
	Factor            float64 `json:"factor"`
	Name              string  `json:"name"`
	DemandPartnerName string  `json:"demand_partner_name"`
}

type DemandPartnerOptimizationRuleSliceJoined []*DemandPartnerOptimizationRuleJoined

type DemandPartnerOptimizationRuleSlice []*DemandPartnerOptimizationRule

type DPOFactorOptions struct {
	Filter     DPORuleFilter          `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type DPORuleFilter struct {
	RuleId          filter.StringArrayFilter `json:"rule_id,omitempty"`
	DemandPartnerId filter.StringArrayFilter `json:"demand_partner_id,omitempty"`
	Publisher       filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain          filter.StringArrayFilter `json:"domain,omitempty"`
	Country         filter.StringArrayFilter `json:"country,omitempty"`
	Device          filter.StringArrayFilter `json:"device_type,omitempty"`
	Factor          filter.StringArrayFilter `json:"factor,omitempty"`
	Active          filter.StringArrayFilter `json:"active,omitempty"`
}

func (d *DPOService) GetJoinedDPORule(ctx context.Context, ops *DPOFactorOptions) (DemandPartnerOptimizationRuleSliceJoined, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.DpoRuleColumns.RuleID).
		AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.DpoRuleColumns.RuleID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		qmods = qmods.Add(qm.Load(models.DpoRuleRels.DemandPartner))
		qmods = qmods.Add(qm.Load(models.DpoRuleRels.DpoRulePublisher))

	}
	mods, err := models.DpoRules(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve Dpo rule")
	}

	res := make(DemandPartnerOptimizationRuleSliceJoined, 0)
	res.FromJoinedModel(mods)

	return res, nil
}

type JoinedDpo struct {
	DemandPartnerId string      `json:"demand_partner_id"`
	IsInclude       bool        `json:"is_include"`
	CreatedAt       null.Time   `json:"created_at"`
	UpdatedAt       null.Time   `json:"updated_at"`
	RuleId          string      `json:"rule_id"`
	Publisher       string      `json:"publisher"`
	Domain          string      `json:"domain"`
	Country         string      `json:"country"`
	Browser         null.String `json:"browser"`
	OS              null.String `json:"os,omitempty"`
	DeviceType      string      `json:"device_type"`
	PlacementType   null.String `json:"placement_type"`
	Factor          float64     `json:"factor"`
	Name            string      `json:"name"`
}

func (dpo *DemandPartnerOptimizationRuleJoined) FromJoinedModel(mod *models.DpoRule) {
	dpo.RuleID = mod.RuleID
	dpo.DemandPartnerID = mod.DemandPartnerID
	dpo.Factor = mod.Factor

	if mod.R.DemandPartner.DemandPartnerName.Valid {
		dpo.DemandPartnerName = mod.R.DemandPartner.DemandPartnerName.String
	}
	if mod.R.DpoRulePublisher != nil {
		dpo.Name = mod.R.DpoRulePublisher.Name
	}
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

func (dpos *DemandPartnerOptimizationRuleSliceJoined) FromJoinedModel(slice models.DpoRuleSlice) {
	for _, mod := range slice {
		dpo := DemandPartnerOptimizationRuleJoined{}
		dpo.FromJoinedModel(mod)
		*dpos = append(*dpos, &dpo)
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
	if len(dpo.RuleID) > 0 {
		return dpo.RuleID
	} else {
		return bcguid.NewFrom(dpo.GetFormula())
	}
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

func (d *DPOService) SetDPORule(ctx context.Context, data *DPOUpdateRequest) (string, error) {
	dpoRule := &DemandPartnerOptimizationRule{
		DemandPartner: data.DemandPartner,
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		OS:            data.OS,
		DeviceType:    data.DeviceType,
		PlacementType: data.PlacementType,
		Browser:       data.Browser,
		Factor:        data.Factor,
	}

	ruleID, err := d.saveDPORule(ctx, dpoRule)
	if err != nil {
		return "", err
	}

	go func() {
		err := sendToRT(context.Background(), data.DemandPartner)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
		}
	}()

	return ruleID, nil
}

func (d *DPOService) UpdateDPORule(ctx context.Context, ruleId string, factor float64) error {
	rule, err := models.DpoRules(models.DpoRuleWhere.RuleID.EQ(ruleId)).One(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to get dpo rule: %w", err)
	}

	oldRule := *rule

	rule.Factor = factor
	rule.Active = true

	updated, err := rule.Update(ctx, bcdb.DB(), boil.Whitelist(models.DpoRuleColumns.Factor, models.DpoRuleColumns.Active))
	if err != nil {
		return fmt.Errorf("failed to update dpo rule: %w", err)
	}

	if updated > 0 {
		go func() {
			err := sendToRT(context.Background(), rule.DemandPartnerID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
			}
		}()
	}

	d.historyModule.SaveAction(ctx, &oldRule, rule, nil)

	return nil
}

func (d *DPOService) DeleteDPORule(ctx context.Context, dpoRules []string) error {
	mods, err := models.DpoRules(models.DpoRuleWhere.RuleID.IN(dpoRules)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed getting dpo rules for soft deleting: %w", err)
	}

	oldMods := make([]any, 0, len(mods))
	newMods := make([]any, 0, len(mods))

	for i := range mods {
		oldMods = append(oldMods, mods[i])
		newMods = append(newMods, nil)
	}

	deleteQuery := createDeleteQuery(dpoRules)

	_, err = queries.Raw(deleteQuery).Exec(bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed soft deleting dpo rules: %w", err)
	}

	d.historyModule.SaveAction(ctx, oldMods, newMods, nil)

	return nil
}

func (filter *DPORuleFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.DpoRuleColumns.Publisher))
	}

	if len(filter.RuleId) > 0 {
		mods = append(mods, filter.RuleId.AndIn(models.DpoRuleColumns.RuleID))
	}

	if len(filter.DemandPartnerId) > 0 {
		mods = append(mods, filter.DemandPartnerId.AndIn(models.DpoRuleColumns.DemandPartnerID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.DpoRuleColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.DpoRuleColumns.Country))
	}

	if len(filter.Factor) > 0 {
		mods = append(mods, filter.Factor.AndIn(models.DpoRuleColumns.Factor))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.DpoRuleColumns.DeviceType))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.DpoRuleColumns.Active))
	}
	return mods
}

func (d *DPOService) saveDPORule(ctx context.Context, dpo *DemandPartnerOptimizationRule) (string, error) {
	mod := dpo.ToModel()

	var old any
	oldMod, err := models.DpoRules(models.DpoRuleWhere.RuleID.EQ(mod.RuleID)).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return "", eris.Wrapf(err, "Failed to get dpo rule(rule=%s)", dpo.GetFormula())
	}

	if oldMod != nil {
		old = oldMod
	}

	err = mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.DpoRuleColumns.RuleID},
		boil.Blacklist(models.DpoRuleColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return "", eris.Wrapf(err, "Failed to upsert dpo rule(rule=%s)", dpo.GetFormula())
	}

	d.historyModule.SaveAction(ctx, old, mod, nil)

	return mod.RuleID, nil
}

func sendToRT(ctx context.Context, demandPartnerID string) error {
	modDpos, err := models.DpoRules(models.DpoRuleWhere.DemandPartnerID.EQ(demandPartnerID)).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "failed to fetch dpo rules(dpid:%s)", demandPartnerID)
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
