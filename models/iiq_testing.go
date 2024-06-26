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

// IiqTesting is an object representing the database table.
type IiqTesting struct {
	Time              time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`
	DemandPartnerID   string    `boil:"demand_partner_id" json:"demand_partner_id" toml:"demand_partner_id" yaml:"demand_partner_id"`
	IiqRequests       int64     `boil:"iiq_requests" json:"iiq_requests" toml:"iiq_requests" yaml:"iiq_requests"`
	NonIiqRequests    int64     `boil:"non_iiq_requests" json:"non_iiq_requests" toml:"non_iiq_requests" yaml:"non_iiq_requests"`
	IiqImpressions    int64     `boil:"iiq_impressions" json:"iiq_impressions" toml:"iiq_impressions" yaml:"iiq_impressions"`
	NonIiqImpressions int64     `boil:"non_iiq_impressions" json:"non_iiq_impressions" toml:"non_iiq_impressions" yaml:"non_iiq_impressions"`

	R *iiqTestingR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L iiqTestingL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var IiqTestingColumns = struct {
	Time              string
	DemandPartnerID   string
	IiqRequests       string
	NonIiqRequests    string
	IiqImpressions    string
	NonIiqImpressions string
}{
	Time:              "time",
	DemandPartnerID:   "demand_partner_id",
	IiqRequests:       "iiq_requests",
	NonIiqRequests:    "non_iiq_requests",
	IiqImpressions:    "iiq_impressions",
	NonIiqImpressions: "non_iiq_impressions",
}

var IiqTestingTableColumns = struct {
	Time              string
	DemandPartnerID   string
	IiqRequests       string
	NonIiqRequests    string
	IiqImpressions    string
	NonIiqImpressions string
}{
	Time:              "iiq_testing.time",
	DemandPartnerID:   "iiq_testing.demand_partner_id",
	IiqRequests:       "iiq_testing.iiq_requests",
	NonIiqRequests:    "iiq_testing.non_iiq_requests",
	IiqImpressions:    "iiq_testing.iiq_impressions",
	NonIiqImpressions: "iiq_testing.non_iiq_impressions",
}

// Generated where

var IiqTestingWhere = struct {
	Time              whereHelpertime_Time
	DemandPartnerID   whereHelperstring
	IiqRequests       whereHelperint64
	NonIiqRequests    whereHelperint64
	IiqImpressions    whereHelperint64
	NonIiqImpressions whereHelperint64
}{
	Time:              whereHelpertime_Time{field: "\"iiq_testing\".\"time\""},
	DemandPartnerID:   whereHelperstring{field: "\"iiq_testing\".\"demand_partner_id\""},
	IiqRequests:       whereHelperint64{field: "\"iiq_testing\".\"iiq_requests\""},
	NonIiqRequests:    whereHelperint64{field: "\"iiq_testing\".\"non_iiq_requests\""},
	IiqImpressions:    whereHelperint64{field: "\"iiq_testing\".\"iiq_impressions\""},
	NonIiqImpressions: whereHelperint64{field: "\"iiq_testing\".\"non_iiq_impressions\""},
}

// IiqTestingRels is where relationship names are stored.
var IiqTestingRels = struct {
}{}

// iiqTestingR is where relationships are stored.
type iiqTestingR struct {
}

// NewStruct creates a new relationship struct
func (*iiqTestingR) NewStruct() *iiqTestingR {
	return &iiqTestingR{}
}

// iiqTestingL is where Load methods for each relationship are stored.
type iiqTestingL struct{}

var (
	iiqTestingAllColumns            = []string{"time", "demand_partner_id", "iiq_requests", "non_iiq_requests", "iiq_impressions", "non_iiq_impressions"}
	iiqTestingColumnsWithoutDefault = []string{"time", "demand_partner_id", "iiq_requests", "non_iiq_requests", "iiq_impressions", "non_iiq_impressions"}
	iiqTestingColumnsWithDefault    = []string{}
	iiqTestingPrimaryKeyColumns     = []string{"time", "demand_partner_id"}
	iiqTestingGeneratedColumns      = []string{}
)

type (
	// IiqTestingSlice is an alias for a slice of pointers to IiqTesting.
	// This should almost always be used instead of []IiqTesting.
	IiqTestingSlice []*IiqTesting
	// IiqTestingHook is the signature for custom IiqTesting hook methods
	IiqTestingHook func(context.Context, boil.ContextExecutor, *IiqTesting) error

	iiqTestingQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	iiqTestingType                 = reflect.TypeOf(&IiqTesting{})
	iiqTestingMapping              = queries.MakeStructMapping(iiqTestingType)
	iiqTestingPrimaryKeyMapping, _ = queries.BindMapping(iiqTestingType, iiqTestingMapping, iiqTestingPrimaryKeyColumns)
	iiqTestingInsertCacheMut       sync.RWMutex
	iiqTestingInsertCache          = make(map[string]insertCache)
	iiqTestingUpdateCacheMut       sync.RWMutex
	iiqTestingUpdateCache          = make(map[string]updateCache)
	iiqTestingUpsertCacheMut       sync.RWMutex
	iiqTestingUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var iiqTestingAfterSelectMu sync.Mutex
var iiqTestingAfterSelectHooks []IiqTestingHook

var iiqTestingBeforeInsertMu sync.Mutex
var iiqTestingBeforeInsertHooks []IiqTestingHook
var iiqTestingAfterInsertMu sync.Mutex
var iiqTestingAfterInsertHooks []IiqTestingHook

var iiqTestingBeforeUpdateMu sync.Mutex
var iiqTestingBeforeUpdateHooks []IiqTestingHook
var iiqTestingAfterUpdateMu sync.Mutex
var iiqTestingAfterUpdateHooks []IiqTestingHook

var iiqTestingBeforeDeleteMu sync.Mutex
var iiqTestingBeforeDeleteHooks []IiqTestingHook
var iiqTestingAfterDeleteMu sync.Mutex
var iiqTestingAfterDeleteHooks []IiqTestingHook

var iiqTestingBeforeUpsertMu sync.Mutex
var iiqTestingBeforeUpsertHooks []IiqTestingHook
var iiqTestingAfterUpsertMu sync.Mutex
var iiqTestingAfterUpsertHooks []IiqTestingHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *IiqTesting) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *IiqTesting) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *IiqTesting) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *IiqTesting) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *IiqTesting) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *IiqTesting) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *IiqTesting) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *IiqTesting) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *IiqTesting) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range iiqTestingAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddIiqTestingHook registers your hook function for all future operations.
func AddIiqTestingHook(hookPoint boil.HookPoint, iiqTestingHook IiqTestingHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		iiqTestingAfterSelectMu.Lock()
		iiqTestingAfterSelectHooks = append(iiqTestingAfterSelectHooks, iiqTestingHook)
		iiqTestingAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		iiqTestingBeforeInsertMu.Lock()
		iiqTestingBeforeInsertHooks = append(iiqTestingBeforeInsertHooks, iiqTestingHook)
		iiqTestingBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		iiqTestingAfterInsertMu.Lock()
		iiqTestingAfterInsertHooks = append(iiqTestingAfterInsertHooks, iiqTestingHook)
		iiqTestingAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		iiqTestingBeforeUpdateMu.Lock()
		iiqTestingBeforeUpdateHooks = append(iiqTestingBeforeUpdateHooks, iiqTestingHook)
		iiqTestingBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		iiqTestingAfterUpdateMu.Lock()
		iiqTestingAfterUpdateHooks = append(iiqTestingAfterUpdateHooks, iiqTestingHook)
		iiqTestingAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		iiqTestingBeforeDeleteMu.Lock()
		iiqTestingBeforeDeleteHooks = append(iiqTestingBeforeDeleteHooks, iiqTestingHook)
		iiqTestingBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		iiqTestingAfterDeleteMu.Lock()
		iiqTestingAfterDeleteHooks = append(iiqTestingAfterDeleteHooks, iiqTestingHook)
		iiqTestingAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		iiqTestingBeforeUpsertMu.Lock()
		iiqTestingBeforeUpsertHooks = append(iiqTestingBeforeUpsertHooks, iiqTestingHook)
		iiqTestingBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		iiqTestingAfterUpsertMu.Lock()
		iiqTestingAfterUpsertHooks = append(iiqTestingAfterUpsertHooks, iiqTestingHook)
		iiqTestingAfterUpsertMu.Unlock()
	}
}

// One returns a single iiqTesting record from the query.
func (q iiqTestingQuery) One(ctx context.Context, exec boil.ContextExecutor) (*IiqTesting, error) {
	o := &IiqTesting{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for iiq_testing")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all IiqTesting records from the query.
func (q iiqTestingQuery) All(ctx context.Context, exec boil.ContextExecutor) (IiqTestingSlice, error) {
	var o []*IiqTesting

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to IiqTesting slice")
	}

	if len(iiqTestingAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all IiqTesting records in the query.
func (q iiqTestingQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count iiq_testing rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q iiqTestingQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if iiq_testing exists")
	}

	return count > 0, nil
}

// IiqTestings retrieves all the records using an executor.
func IiqTestings(mods ...qm.QueryMod) iiqTestingQuery {
	mods = append(mods, qm.From("\"iiq_testing\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"iiq_testing\".*"})
	}

	return iiqTestingQuery{q}
}

// FindIiqTesting retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindIiqTesting(ctx context.Context, exec boil.ContextExecutor, time time.Time, demandPartnerID string, selectCols ...string) (*IiqTesting, error) {
	iiqTestingObj := &IiqTesting{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"iiq_testing\" where \"time\"=$1 AND \"demand_partner_id\"=$2", sel,
	)

	q := queries.Raw(query, time, demandPartnerID)

	err := q.Bind(ctx, exec, iiqTestingObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from iiq_testing")
	}

	if err = iiqTestingObj.doAfterSelectHooks(ctx, exec); err != nil {
		return iiqTestingObj, err
	}

	return iiqTestingObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *IiqTesting) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no iiq_testing provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(iiqTestingColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	iiqTestingInsertCacheMut.RLock()
	cache, cached := iiqTestingInsertCache[key]
	iiqTestingInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			iiqTestingAllColumns,
			iiqTestingColumnsWithDefault,
			iiqTestingColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(iiqTestingType, iiqTestingMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(iiqTestingType, iiqTestingMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"iiq_testing\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"iiq_testing\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into iiq_testing")
	}

	if !cached {
		iiqTestingInsertCacheMut.Lock()
		iiqTestingInsertCache[key] = cache
		iiqTestingInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the IiqTesting.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *IiqTesting) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	iiqTestingUpdateCacheMut.RLock()
	cache, cached := iiqTestingUpdateCache[key]
	iiqTestingUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			iiqTestingAllColumns,
			iiqTestingPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update iiq_testing, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"iiq_testing\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, iiqTestingPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(iiqTestingType, iiqTestingMapping, append(wl, iiqTestingPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update iiq_testing row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for iiq_testing")
	}

	if !cached {
		iiqTestingUpdateCacheMut.Lock()
		iiqTestingUpdateCache[key] = cache
		iiqTestingUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q iiqTestingQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for iiq_testing")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for iiq_testing")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o IiqTestingSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), iiqTestingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"iiq_testing\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, iiqTestingPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in iiqTesting slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all iiqTesting")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *IiqTesting) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no iiq_testing provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(iiqTestingColumnsWithDefault, o)

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

	iiqTestingUpsertCacheMut.RLock()
	cache, cached := iiqTestingUpsertCache[key]
	iiqTestingUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			iiqTestingAllColumns,
			iiqTestingColumnsWithDefault,
			iiqTestingColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			iiqTestingAllColumns,
			iiqTestingPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert iiq_testing, could not build update column list")
		}

		ret := strmangle.SetComplement(iiqTestingAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(iiqTestingPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert iiq_testing, could not build conflict column list")
			}

			conflict = make([]string, len(iiqTestingPrimaryKeyColumns))
			copy(conflict, iiqTestingPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"iiq_testing\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(iiqTestingType, iiqTestingMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(iiqTestingType, iiqTestingMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert iiq_testing")
	}

	if !cached {
		iiqTestingUpsertCacheMut.Lock()
		iiqTestingUpsertCache[key] = cache
		iiqTestingUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single IiqTesting record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *IiqTesting) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no IiqTesting provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), iiqTestingPrimaryKeyMapping)
	sql := "DELETE FROM \"iiq_testing\" WHERE \"time\"=$1 AND \"demand_partner_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from iiq_testing")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for iiq_testing")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q iiqTestingQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no iiqTestingQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from iiq_testing")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for iiq_testing")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o IiqTestingSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(iiqTestingBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), iiqTestingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"iiq_testing\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, iiqTestingPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from iiqTesting slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for iiq_testing")
	}

	if len(iiqTestingAfterDeleteHooks) != 0 {
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
func (o *IiqTesting) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindIiqTesting(ctx, exec, o.Time, o.DemandPartnerID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *IiqTestingSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := IiqTestingSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), iiqTestingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"iiq_testing\".* FROM \"iiq_testing\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, iiqTestingPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in IiqTestingSlice")
	}

	*o = slice

	return nil
}

// IiqTestingExists checks if the IiqTesting row exists.
func IiqTestingExists(ctx context.Context, exec boil.ContextExecutor, time time.Time, demandPartnerID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"iiq_testing\" where \"time\"=$1 AND \"demand_partner_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, time, demandPartnerID)
	}
	row := exec.QueryRowContext(ctx, sql, time, demandPartnerID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if iiq_testing exists")
	}

	return exists, nil
}

// Exists checks if the IiqTesting row exists.
func (o *IiqTesting) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return IiqTestingExists(ctx, exec, o.Time, o.DemandPartnerID)
}
