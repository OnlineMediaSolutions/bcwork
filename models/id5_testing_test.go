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

func testID5Testings(t *testing.T) {
	t.Parallel()

	query := ID5Testings()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testID5TestingsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
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

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testID5TestingsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ID5Testings().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testID5TestingsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ID5TestingSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testID5TestingsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ID5TestingExists(ctx, tx, o.Time, o.DemandPartnerID)
	if err != nil {
		t.Errorf("Unable to check if ID5Testing exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ID5TestingExists to return true, but got false.")
	}
}

func testID5TestingsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	id5TestingFound, err := FindID5Testing(ctx, tx, o.Time, o.DemandPartnerID)
	if err != nil {
		t.Error(err)
	}

	if id5TestingFound == nil {
		t.Error("want a record, got nil")
	}
}

func testID5TestingsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ID5Testings().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testID5TestingsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ID5Testings().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testID5TestingsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	id5TestingOne := &ID5Testing{}
	id5TestingTwo := &ID5Testing{}
	if err = randomize.Struct(seed, id5TestingOne, id5TestingDBTypes, false, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}
	if err = randomize.Struct(seed, id5TestingTwo, id5TestingDBTypes, false, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = id5TestingOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = id5TestingTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ID5Testings().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testID5TestingsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	id5TestingOne := &ID5Testing{}
	id5TestingTwo := &ID5Testing{}
	if err = randomize.Struct(seed, id5TestingOne, id5TestingDBTypes, false, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}
	if err = randomize.Struct(seed, id5TestingTwo, id5TestingDBTypes, false, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = id5TestingOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = id5TestingTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func id5TestingBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func id5TestingAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ID5Testing) error {
	*o = ID5Testing{}
	return nil
}

func testID5TestingsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ID5Testing{}
	o := &ID5Testing{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, id5TestingDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ID5Testing object: %s", err)
	}

	AddID5TestingHook(boil.BeforeInsertHook, id5TestingBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	id5TestingBeforeInsertHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.AfterInsertHook, id5TestingAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	id5TestingAfterInsertHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.AfterSelectHook, id5TestingAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	id5TestingAfterSelectHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.BeforeUpdateHook, id5TestingBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	id5TestingBeforeUpdateHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.AfterUpdateHook, id5TestingAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	id5TestingAfterUpdateHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.BeforeDeleteHook, id5TestingBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	id5TestingBeforeDeleteHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.AfterDeleteHook, id5TestingAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	id5TestingAfterDeleteHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.BeforeUpsertHook, id5TestingBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	id5TestingBeforeUpsertHooks = []ID5TestingHook{}

	AddID5TestingHook(boil.AfterUpsertHook, id5TestingAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	id5TestingAfterUpsertHooks = []ID5TestingHook{}
}

func testID5TestingsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testID5TestingsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(id5TestingColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testID5TestingsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
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

func testID5TestingsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ID5TestingSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testID5TestingsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ID5Testings().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	id5TestingDBTypes = map[string]string{`Time`: `timestamp without time zone`, `DemandPartnerID`: `character varying`, `ID5Requests`: `bigint`, `NonID5Requests`: `bigint`, `ID5Impressions`: `bigint`, `NonID5Impressions`: `bigint`}
	_                 = bytes.MinRead
)

func testID5TestingsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(id5TestingPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(id5TestingAllColumns) == len(id5TestingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testID5TestingsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(id5TestingAllColumns) == len(id5TestingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ID5Testing{}
	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, id5TestingDBTypes, true, id5TestingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(id5TestingAllColumns, id5TestingPrimaryKeyColumns) {
		fields = id5TestingAllColumns
	} else {
		fields = strmangle.SetComplement(
			id5TestingAllColumns,
			id5TestingPrimaryKeyColumns,
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

	slice := ID5TestingSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testID5TestingsUpsert(t *testing.T) {
	t.Parallel()

	if len(id5TestingAllColumns) == len(id5TestingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ID5Testing{}
	if err = randomize.Struct(seed, &o, id5TestingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ID5Testing: %s", err)
	}

	count, err := ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, id5TestingDBTypes, false, id5TestingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ID5Testing struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ID5Testing: %s", err)
	}

	count, err = ID5Testings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
