// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testRevenueDailies(t *testing.T) {
	t.Parallel()

	query := RevenueDailies()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testRevenueDailiesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRevenueDailiesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := RevenueDailies().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRevenueDailiesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RevenueDailySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRevenueDailiesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := RevenueDailyExists(ctx, tx, o.Time)
	if err != nil {
		t.Errorf("Unable to check if RevenueDaily exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RevenueDailyExists to return true, but got false.")
	}
}

func testRevenueDailiesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	revenueDailyFound, err := FindRevenueDaily(ctx, tx, o.Time)
	if err != nil {
		t.Error(err)
	}

	if revenueDailyFound == nil {
		t.Error("want a record, got nil")
	}
}

func testRevenueDailiesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = RevenueDailies().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testRevenueDailiesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := RevenueDailies().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRevenueDailiesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	revenueDailyOne := &RevenueDaily{}
	revenueDailyTwo := &RevenueDaily{}
	if err = randomize.Struct(seed, revenueDailyOne, revenueDailyDBTypes, false, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}
	if err = randomize.Struct(seed, revenueDailyTwo, revenueDailyDBTypes, false, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = revenueDailyOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = revenueDailyTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RevenueDailies().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRevenueDailiesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	revenueDailyOne := &RevenueDaily{}
	revenueDailyTwo := &RevenueDaily{}
	if err = randomize.Struct(seed, revenueDailyOne, revenueDailyDBTypes, false, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}
	if err = randomize.Struct(seed, revenueDailyTwo, revenueDailyDBTypes, false, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = revenueDailyOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = revenueDailyTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func revenueDailyBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func revenueDailyAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *RevenueDaily) error {
	*o = RevenueDaily{}
	return nil
}

func testRevenueDailiesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &RevenueDaily{}
	o := &RevenueDaily{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, false); err != nil {
		t.Errorf("Unable to randomize RevenueDaily object: %s", err)
	}

	AddRevenueDailyHook(boil.BeforeInsertHook, revenueDailyBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	revenueDailyBeforeInsertHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.AfterInsertHook, revenueDailyAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	revenueDailyAfterInsertHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.AfterSelectHook, revenueDailyAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	revenueDailyAfterSelectHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.BeforeUpdateHook, revenueDailyBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	revenueDailyBeforeUpdateHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.AfterUpdateHook, revenueDailyAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	revenueDailyAfterUpdateHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.BeforeDeleteHook, revenueDailyBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	revenueDailyBeforeDeleteHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.AfterDeleteHook, revenueDailyAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	revenueDailyAfterDeleteHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.BeforeUpsertHook, revenueDailyBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	revenueDailyBeforeUpsertHooks = []RevenueDailyHook{}

	AddRevenueDailyHook(boil.AfterUpsertHook, revenueDailyAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	revenueDailyAfterUpsertHooks = []RevenueDailyHook{}
}

func testRevenueDailiesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRevenueDailiesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(revenueDailyColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRevenueDailiesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testRevenueDailiesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RevenueDailySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testRevenueDailiesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RevenueDailies().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	revenueDailyDBTypes = map[string]string{`Time`: `timestamp without time zone`, `PublisherImpressions`: `bigint`, `SoldImpressions`: `bigint`, `Cost`: `double precision`, `Revenue`: `double precision`, `DemandPartnerFees`: `double precision`, `MissedOpportunities`: `bigint`, `DataFee`: `double precision`, `DPBidRequests`: `bigint`}
	_                   = bytes.MinRead
)

func testRevenueDailiesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(revenueDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(revenueDailyAllColumns) == len(revenueDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testRevenueDailiesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(revenueDailyAllColumns) == len(revenueDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RevenueDaily{}
	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, revenueDailyDBTypes, true, revenueDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(revenueDailyAllColumns, revenueDailyPrimaryKeyColumns) {
		fields = revenueDailyAllColumns
	} else {
		fields = strmangle.SetComplement(
			revenueDailyAllColumns,
			revenueDailyPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := RevenueDailySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testRevenueDailiesUpsert(t *testing.T) {
	t.Parallel()

	if len(revenueDailyAllColumns) == len(revenueDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := RevenueDaily{}
	if err = randomize.Struct(seed, &o, revenueDailyDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RevenueDaily: %s", err)
	}

	count, err := RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, revenueDailyDBTypes, false, revenueDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RevenueDaily struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RevenueDaily: %s", err)
	}

	count, err = RevenueDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
