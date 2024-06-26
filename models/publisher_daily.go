// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PublisherDaily is an object representing the database table.
type PublisherDaily struct {
	Time                 time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string    `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain               string    `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Os                   string    `boil:"os" json:"os" toml:"os" yaml:"os"`
	Country              string    `boil:"country" json:"country" toml:"country" yaml:"country"`
	DeviceType           string    `boil:"device_type" json:"device_type" toml:"device_type" yaml:"device_type"`
	BidRequests          int64     `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	BidResponses         int64     `boil:"bid_responses" json:"bid_responses" toml:"bid_responses" yaml:"bid_responses"`
	BidPriceCount        int64     `boil:"bid_price_count" json:"bid_price_count" toml:"bid_price_count" yaml:"bid_price_count"`
	BidPriceTotal        float64   `boil:"bid_price_total" json:"bid_price_total" toml:"bid_price_total" yaml:"bid_price_total"`
	PublisherImpressions int64     `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	DemandImpressions    int64     `boil:"demand_impressions" json:"demand_impressions" toml:"demand_impressions" yaml:"demand_impressions"`
	MissedOpportunities  int64     `boil:"missed_opportunities" json:"missed_opportunities" toml:"missed_opportunities" yaml:"missed_opportunities"`
	SupplyTotal          float64   `boil:"supply_total" json:"supply_total" toml:"supply_total" yaml:"supply_total"`
	DemandTotal          float64   `boil:"demand_total" json:"demand_total" toml:"demand_total" yaml:"demand_total"`
	DemandPartnerFee     float64   `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`

	R *publisherDailyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L publisherDailyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PublisherDailyColumns = struct {
	Time                 string
	PublisherID          string
	Domain               string
	Os                   string
	Country              string
	DeviceType           string
	BidRequests          string
	BidResponses         string
	BidPriceCount        string
	BidPriceTotal        string
	PublisherImpressions string
	DemandImpressions    string
	MissedOpportunities  string
	SupplyTotal          string
	DemandTotal          string
	DemandPartnerFee     string
}{
	Time:                 "time",
	PublisherID:          "publisher_id",
	Domain:               "domain",
	Os:                   "os",
	Country:              "country",
	DeviceType:           "device_type",
	BidRequests:          "bid_requests",
	BidResponses:         "bid_responses",
	BidPriceCount:        "bid_price_count",
	BidPriceTotal:        "bid_price_total",
	PublisherImpressions: "publisher_impressions",
	DemandImpressions:    "demand_impressions",
	MissedOpportunities:  "missed_opportunities",
	SupplyTotal:          "supply_total",
	DemandTotal:          "demand_total",
	DemandPartnerFee:     "demand_partner_fee",
}

var PublisherDailyTableColumns = struct {
	Time                 string
	PublisherID          string
	Domain               string
	Os                   string
	Country              string
	DeviceType           string
	BidRequests          string
	BidResponses         string
	BidPriceCount        string
	BidPriceTotal        string
	PublisherImpressions string
	DemandImpressions    string
	MissedOpportunities  string
	SupplyTotal          string
	DemandTotal          string
	DemandPartnerFee     string
}{
	Time:                 "publisher_daily.time",
	PublisherID:          "publisher_daily.publisher_id",
	Domain:               "publisher_daily.domain",
	Os:                   "publisher_daily.os",
	Country:              "publisher_daily.country",
	DeviceType:           "publisher_daily.device_type",
	BidRequests:          "publisher_daily.bid_requests",
	BidResponses:         "publisher_daily.bid_responses",
	BidPriceCount:        "publisher_daily.bid_price_count",
	BidPriceTotal:        "publisher_daily.bid_price_total",
	PublisherImpressions: "publisher_daily.publisher_impressions",
	DemandImpressions:    "publisher_daily.demand_impressions",
	MissedOpportunities:  "publisher_daily.missed_opportunities",
	SupplyTotal:          "publisher_daily.supply_total",
	DemandTotal:          "publisher_daily.demand_total",
	DemandPartnerFee:     "publisher_daily.demand_partner_fee",
}

// Generated where

var PublisherDailyWhere = struct {
	Time                 whereHelpertime_Time
	PublisherID          whereHelperstring
	Domain               whereHelperstring
	Os                   whereHelperstring
	Country              whereHelperstring
	DeviceType           whereHelperstring
	BidRequests          whereHelperint64
	BidResponses         whereHelperint64
	BidPriceCount        whereHelperint64
	BidPriceTotal        whereHelperfloat64
	PublisherImpressions whereHelperint64
	DemandImpressions    whereHelperint64
	MissedOpportunities  whereHelperint64
	SupplyTotal          whereHelperfloat64
	DemandTotal          whereHelperfloat64
	DemandPartnerFee     whereHelperfloat64
}{
	Time:                 whereHelpertime_Time{field: "\"publisher_daily\".\"time\""},
	PublisherID:          whereHelperstring{field: "\"publisher_daily\".\"publisher_id\""},
	Domain:               whereHelperstring{field: "\"publisher_daily\".\"domain\""},
	Os:                   whereHelperstring{field: "\"publisher_daily\".\"os\""},
	Country:              whereHelperstring{field: "\"publisher_daily\".\"country\""},
	DeviceType:           whereHelperstring{field: "\"publisher_daily\".\"device_type\""},
	BidRequests:          whereHelperint64{field: "\"publisher_daily\".\"bid_requests\""},
	BidResponses:         whereHelperint64{field: "\"publisher_daily\".\"bid_responses\""},
	BidPriceCount:        whereHelperint64{field: "\"publisher_daily\".\"bid_price_count\""},
	BidPriceTotal:        whereHelperfloat64{field: "\"publisher_daily\".\"bid_price_total\""},
	PublisherImpressions: whereHelperint64{field: "\"publisher_daily\".\"publisher_impressions\""},
	DemandImpressions:    whereHelperint64{field: "\"publisher_daily\".\"demand_impressions\""},
	MissedOpportunities:  whereHelperint64{field: "\"publisher_daily\".\"missed_opportunities\""},
	SupplyTotal:          whereHelperfloat64{field: "\"publisher_daily\".\"supply_total\""},
	DemandTotal:          whereHelperfloat64{field: "\"publisher_daily\".\"demand_total\""},
	DemandPartnerFee:     whereHelperfloat64{field: "\"publisher_daily\".\"demand_partner_fee\""},
}

// PublisherDailyRels is where relationship names are stored.
var PublisherDailyRels = struct {
}{}

// publisherDailyR is where relationships are stored.
type publisherDailyR struct {
}

// NewStruct creates a new relationship struct
func (*publisherDailyR) NewStruct() *publisherDailyR {
	return &publisherDailyR{}
}

// publisherDailyL is where Load methods for each relationship are stored.
type publisherDailyL struct{}

var (
	publisherDailyAllColumns            = []string{"time", "publisher_id", "domain", "os", "country", "device_type", "bid_requests", "bid_responses", "bid_price_count", "bid_price_total", "publisher_impressions", "demand_impressions", "missed_opportunities", "supply_total", "demand_total", "demand_partner_fee"}
	publisherDailyColumnsWithoutDefault = []string{"time", "publisher_id"}
	publisherDailyColumnsWithDefault    = []string{"domain", "os", "country", "device_type", "bid_requests", "bid_responses", "bid_price_count", "bid_price_total", "publisher_impressions", "demand_impressions", "missed_opportunities", "supply_total", "demand_total", "demand_partner_fee"}
	publisherDailyPrimaryKeyColumns     = []string{"time", "publisher_id", "domain", "os", "country", "device_type"}
	publisherDailyGeneratedColumns      = []string{}
)

type (
	// PublisherDailySlice is an alias for a slice of pointers to PublisherDaily.
	// This should almost always be used instead of []PublisherDaily.
	PublisherDailySlice []*PublisherDaily
	// PublisherDailyHook is the signature for custom PublisherDaily hook methods
	PublisherDailyHook func(context.Context, boil.ContextExecutor, *PublisherDaily) error

	publisherDailyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	publisherDailyType                 = reflect.TypeOf(&PublisherDaily{})
	publisherDailyMapping              = queries.MakeStructMapping(publisherDailyType)
	publisherDailyPrimaryKeyMapping, _ = queries.BindMapping(publisherDailyType, publisherDailyMapping, publisherDailyPrimaryKeyColumns)
	publisherDailyInsertCacheMut       sync.RWMutex
	publisherDailyInsertCache          = make(map[string]insertCache)
	publisherDailyUpdateCacheMut       sync.RWMutex
	publisherDailyUpdateCache          = make(map[string]updateCache)
	publisherDailyUpsertCacheMut       sync.RWMutex
	publisherDailyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var publisherDailyAfterSelectMu sync.Mutex
var publisherDailyAfterSelectHooks []PublisherDailyHook

var publisherDailyBeforeInsertMu sync.Mutex
var publisherDailyBeforeInsertHooks []PublisherDailyHook
var publisherDailyAfterInsertMu sync.Mutex
var publisherDailyAfterInsertHooks []PublisherDailyHook

var publisherDailyBeforeUpdateMu sync.Mutex
var publisherDailyBeforeUpdateHooks []PublisherDailyHook
var publisherDailyAfterUpdateMu sync.Mutex
var publisherDailyAfterUpdateHooks []PublisherDailyHook

var publisherDailyBeforeDeleteMu sync.Mutex
var publisherDailyBeforeDeleteHooks []PublisherDailyHook
var publisherDailyAfterDeleteMu sync.Mutex
var publisherDailyAfterDeleteHooks []PublisherDailyHook

var publisherDailyBeforeUpsertMu sync.Mutex
var publisherDailyBeforeUpsertHooks []PublisherDailyHook
var publisherDailyAfterUpsertMu sync.Mutex
var publisherDailyAfterUpsertHooks []PublisherDailyHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *PublisherDaily) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *PublisherDaily) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *PublisherDaily) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *PublisherDaily) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *PublisherDaily) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *PublisherDaily) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *PublisherDaily) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *PublisherDaily) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *PublisherDaily) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range publisherDailyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPublisherDailyHook registers your hook function for all future operations.
func AddPublisherDailyHook(hookPoint boil.HookPoint, publisherDailyHook PublisherDailyHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		publisherDailyAfterSelectMu.Lock()
		publisherDailyAfterSelectHooks = append(publisherDailyAfterSelectHooks, publisherDailyHook)
		publisherDailyAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		publisherDailyBeforeInsertMu.Lock()
		publisherDailyBeforeInsertHooks = append(publisherDailyBeforeInsertHooks, publisherDailyHook)
		publisherDailyBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		publisherDailyAfterInsertMu.Lock()
		publisherDailyAfterInsertHooks = append(publisherDailyAfterInsertHooks, publisherDailyHook)
		publisherDailyAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		publisherDailyBeforeUpdateMu.Lock()
		publisherDailyBeforeUpdateHooks = append(publisherDailyBeforeUpdateHooks, publisherDailyHook)
		publisherDailyBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		publisherDailyAfterUpdateMu.Lock()
		publisherDailyAfterUpdateHooks = append(publisherDailyAfterUpdateHooks, publisherDailyHook)
		publisherDailyAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		publisherDailyBeforeDeleteMu.Lock()
		publisherDailyBeforeDeleteHooks = append(publisherDailyBeforeDeleteHooks, publisherDailyHook)
		publisherDailyBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		publisherDailyAfterDeleteMu.Lock()
		publisherDailyAfterDeleteHooks = append(publisherDailyAfterDeleteHooks, publisherDailyHook)
		publisherDailyAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		publisherDailyBeforeUpsertMu.Lock()
		publisherDailyBeforeUpsertHooks = append(publisherDailyBeforeUpsertHooks, publisherDailyHook)
		publisherDailyBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		publisherDailyAfterUpsertMu.Lock()
		publisherDailyAfterUpsertHooks = append(publisherDailyAfterUpsertHooks, publisherDailyHook)
		publisherDailyAfterUpsertMu.Unlock()
	}
}

// One returns a single publisherDaily record from the query.
func (q publisherDailyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PublisherDaily, error) {
	o := &PublisherDaily{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for publisher_daily")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all PublisherDaily records from the query.
func (q publisherDailyQuery) All(ctx context.Context, exec boil.ContextExecutor) (PublisherDailySlice, error) {
	var o []*PublisherDaily

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PublisherDaily slice")
	}

	if len(publisherDailyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all PublisherDaily records in the query.
func (q publisherDailyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count publisher_daily rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q publisherDailyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if publisher_daily exists")
	}

	return count > 0, nil
}

// PublisherDailies retrieves all the records using an executor.
func PublisherDailies(mods ...qm.QueryMod) publisherDailyQuery {
	mods = append(mods, qm.From("\"publisher_daily\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"publisher_daily\".*"})
	}

	return publisherDailyQuery{q}
}

// FindPublisherDaily retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPublisherDaily(ctx context.Context, exec boil.ContextExecutor, time time.Time, publisherID string, domain string, os string, country string, deviceType string, selectCols ...string) (*PublisherDaily, error) {
	publisherDailyObj := &PublisherDaily{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"publisher_daily\" where \"time\"=$1 AND \"publisher_id\"=$2 AND \"domain\"=$3 AND \"os\"=$4 AND \"country\"=$5 AND \"device_type\"=$6", sel,
	)

	q := queries.Raw(query, time, publisherID, domain, os, country, deviceType)

	err := q.Bind(ctx, exec, publisherDailyObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from publisher_daily")
	}

	if err = publisherDailyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return publisherDailyObj, err
	}

	return publisherDailyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PublisherDaily) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no publisher_daily provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(publisherDailyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	publisherDailyInsertCacheMut.RLock()
	cache, cached := publisherDailyInsertCache[key]
	publisherDailyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			publisherDailyAllColumns,
			publisherDailyColumnsWithDefault,
			publisherDailyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(publisherDailyType, publisherDailyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(publisherDailyType, publisherDailyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"publisher_daily\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"publisher_daily\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into publisher_daily")
	}

	if !cached {
		publisherDailyInsertCacheMut.Lock()
		publisherDailyInsertCache[key] = cache
		publisherDailyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the PublisherDaily.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PublisherDaily) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	publisherDailyUpdateCacheMut.RLock()
	cache, cached := publisherDailyUpdateCache[key]
	publisherDailyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			publisherDailyAllColumns,
			publisherDailyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update publisher_daily, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"publisher_daily\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, publisherDailyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(publisherDailyType, publisherDailyMapping, append(wl, publisherDailyPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update publisher_daily row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for publisher_daily")
	}

	if !cached {
		publisherDailyUpdateCacheMut.Lock()
		publisherDailyUpdateCache[key] = cache
		publisherDailyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q publisherDailyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for publisher_daily")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for publisher_daily")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PublisherDailySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), publisherDailyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"publisher_daily\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, publisherDailyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in publisherDaily slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all publisherDaily")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PublisherDaily) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no publisher_daily provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(publisherDailyColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	publisherDailyUpsertCacheMut.RLock()
	cache, cached := publisherDailyUpsertCache[key]
	publisherDailyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			publisherDailyAllColumns,
			publisherDailyColumnsWithDefault,
			publisherDailyColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			publisherDailyAllColumns,
			publisherDailyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert publisher_daily, could not build update column list")
		}

		ret := strmangle.SetComplement(publisherDailyAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(publisherDailyPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert publisher_daily, could not build conflict column list")
			}

			conflict = make([]string, len(publisherDailyPrimaryKeyColumns))
			copy(conflict, publisherDailyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"publisher_daily\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(publisherDailyType, publisherDailyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(publisherDailyType, publisherDailyMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert publisher_daily")
	}

	if !cached {
		publisherDailyUpsertCacheMut.Lock()
		publisherDailyUpsertCache[key] = cache
		publisherDailyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single PublisherDaily record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PublisherDaily) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PublisherDaily provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), publisherDailyPrimaryKeyMapping)
	sql := "DELETE FROM \"publisher_daily\" WHERE \"time\"=$1 AND \"publisher_id\"=$2 AND \"domain\"=$3 AND \"os\"=$4 AND \"country\"=$5 AND \"device_type\"=$6"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from publisher_daily")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for publisher_daily")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q publisherDailyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no publisherDailyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from publisher_daily")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for publisher_daily")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PublisherDailySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(publisherDailyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), publisherDailyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"publisher_daily\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, publisherDailyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from publisherDaily slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for publisher_daily")
	}

	if len(publisherDailyAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PublisherDaily) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPublisherDaily(ctx, exec, o.Time, o.PublisherID, o.Domain, o.Os, o.Country, o.DeviceType)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PublisherDailySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PublisherDailySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), publisherDailyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"publisher_daily\".* FROM \"publisher_daily\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, publisherDailyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PublisherDailySlice")
	}

	*o = slice

	return nil
}

// PublisherDailyExists checks if the PublisherDaily row exists.
func PublisherDailyExists(ctx context.Context, exec boil.ContextExecutor, time time.Time, publisherID string, domain string, os string, country string, deviceType string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"publisher_daily\" where \"time\"=$1 AND \"publisher_id\"=$2 AND \"domain\"=$3 AND \"os\"=$4 AND \"country\"=$5 AND \"device_type\"=$6 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, time, publisherID, domain, os, country, deviceType)
	}
	row := exec.QueryRowContext(ctx, sql, time, publisherID, domain, os, country, deviceType)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if publisher_daily exists")
	}

	return exists, nil
}

// Exists checks if the PublisherDaily row exists.
func (o *PublisherDaily) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return PublisherDailyExists(ctx, exec, o.Time, o.PublisherID, o.Domain, o.Os, o.Country, o.DeviceType)
}
