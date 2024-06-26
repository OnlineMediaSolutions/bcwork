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

func testDemandPartnerDailies(t *testing.T) {
	t.Parallel()

	query := DemandPartnerDailies()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testDemandPartnerDailiesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
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

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerDailiesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := DemandPartnerDailies().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerDailiesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DemandPartnerDailySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerDailiesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := DemandPartnerDailyExists(ctx, tx, o.Time, o.DemandPartnerID, o.Domain)
	if err != nil {
		t.Errorf("Unable to check if DemandPartnerDaily exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DemandPartnerDailyExists to return true, but got false.")
	}
}

func testDemandPartnerDailiesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	demandPartnerDailyFound, err := FindDemandPartnerDaily(ctx, tx, o.Time, o.DemandPartnerID, o.Domain)
	if err != nil {
		t.Error(err)
	}

	if demandPartnerDailyFound == nil {
		t.Error("want a record, got nil")
	}
}

func testDemandPartnerDailiesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = DemandPartnerDailies().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testDemandPartnerDailiesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := DemandPartnerDailies().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDemandPartnerDailiesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demandPartnerDailyOne := &DemandPartnerDaily{}
	demandPartnerDailyTwo := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, demandPartnerDailyOne, demandPartnerDailyDBTypes, false, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}
	if err = randomize.Struct(seed, demandPartnerDailyTwo, demandPartnerDailyDBTypes, false, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = demandPartnerDailyOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = demandPartnerDailyTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DemandPartnerDailies().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDemandPartnerDailiesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	demandPartnerDailyOne := &DemandPartnerDaily{}
	demandPartnerDailyTwo := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, demandPartnerDailyOne, demandPartnerDailyDBTypes, false, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}
	if err = randomize.Struct(seed, demandPartnerDailyTwo, demandPartnerDailyDBTypes, false, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = demandPartnerDailyOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = demandPartnerDailyTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func demandPartnerDailyBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func demandPartnerDailyAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerDaily) error {
	*o = DemandPartnerDaily{}
	return nil
}

func testDemandPartnerDailiesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &DemandPartnerDaily{}
	o := &DemandPartnerDaily{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, false); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily object: %s", err)
	}

	AddDemandPartnerDailyHook(boil.BeforeInsertHook, demandPartnerDailyBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyBeforeInsertHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.AfterInsertHook, demandPartnerDailyAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyAfterInsertHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.AfterSelectHook, demandPartnerDailyAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyAfterSelectHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.BeforeUpdateHook, demandPartnerDailyBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyBeforeUpdateHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.AfterUpdateHook, demandPartnerDailyAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyAfterUpdateHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.BeforeDeleteHook, demandPartnerDailyBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyBeforeDeleteHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.AfterDeleteHook, demandPartnerDailyAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyAfterDeleteHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.BeforeUpsertHook, demandPartnerDailyBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyBeforeUpsertHooks = []DemandPartnerDailyHook{}

	AddDemandPartnerDailyHook(boil.AfterUpsertHook, demandPartnerDailyAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerDailyAfterUpsertHooks = []DemandPartnerDailyHook{}
}

func testDemandPartnerDailiesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemandPartnerDailiesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(demandPartnerDailyColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemandPartnerDailiesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
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

func testDemandPartnerDailiesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DemandPartnerDailySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testDemandPartnerDailiesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DemandPartnerDailies().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	demandPartnerDailyDBTypes = map[string]string{`Time`: `timestamp without time zone`, `DemandPartnerID`: `character varying`, `Domain`: `character varying`, `Impression`: `bigint`, `Revenue`: `double precision`}
	_                         = bytes.MinRead
)

func testDemandPartnerDailiesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(demandPartnerDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(demandPartnerDailyAllColumns) == len(demandPartnerDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testDemandPartnerDailiesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(demandPartnerDailyAllColumns) == len(demandPartnerDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerDaily{}
	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, demandPartnerDailyDBTypes, true, demandPartnerDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(demandPartnerDailyAllColumns, demandPartnerDailyPrimaryKeyColumns) {
		fields = demandPartnerDailyAllColumns
	} else {
		fields = strmangle.SetComplement(
			demandPartnerDailyAllColumns,
			demandPartnerDailyPrimaryKeyColumns,
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

	slice := DemandPartnerDailySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testDemandPartnerDailiesUpsert(t *testing.T) {
	t.Parallel()

	if len(demandPartnerDailyAllColumns) == len(demandPartnerDailyPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := DemandPartnerDaily{}
	if err = randomize.Struct(seed, &o, demandPartnerDailyDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DemandPartnerDaily: %s", err)
	}

	count, err := DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, demandPartnerDailyDBTypes, false, demandPartnerDailyPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerDaily struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DemandPartnerDaily: %s", err)
	}

	count, err = DemandPartnerDailies().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
