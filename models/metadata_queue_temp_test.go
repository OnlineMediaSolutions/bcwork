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

func testMetadataQueueTemps(t *testing.T) {
	t.Parallel()

	query := MetadataQueueTemps()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testMetadataQueueTempsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
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

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMetadataQueueTempsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := MetadataQueueTemps().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMetadataQueueTempsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MetadataQueueTempSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMetadataQueueTempsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := MetadataQueueTempExists(ctx, tx, o.TransactionID)
	if err != nil {
		t.Errorf("Unable to check if MetadataQueueTemp exists: %s", err)
	}
	if !e {
		t.Errorf("Expected MetadataQueueTempExists to return true, but got false.")
	}
}

func testMetadataQueueTempsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	metadataQueueTempFound, err := FindMetadataQueueTemp(ctx, tx, o.TransactionID)
	if err != nil {
		t.Error(err)
	}

	if metadataQueueTempFound == nil {
		t.Error("want a record, got nil")
	}
}

func testMetadataQueueTempsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = MetadataQueueTemps().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testMetadataQueueTempsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := MetadataQueueTemps().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testMetadataQueueTempsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	metadataQueueTempOne := &MetadataQueueTemp{}
	metadataQueueTempTwo := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, metadataQueueTempOne, metadataQueueTempDBTypes, false, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}
	if err = randomize.Struct(seed, metadataQueueTempTwo, metadataQueueTempDBTypes, false, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = metadataQueueTempOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = metadataQueueTempTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MetadataQueueTemps().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testMetadataQueueTempsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	metadataQueueTempOne := &MetadataQueueTemp{}
	metadataQueueTempTwo := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, metadataQueueTempOne, metadataQueueTempDBTypes, false, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}
	if err = randomize.Struct(seed, metadataQueueTempTwo, metadataQueueTempDBTypes, false, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = metadataQueueTempOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = metadataQueueTempTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func metadataQueueTempBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func metadataQueueTempAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *MetadataQueueTemp) error {
	*o = MetadataQueueTemp{}
	return nil
}

func testMetadataQueueTempsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &MetadataQueueTemp{}
	o := &MetadataQueueTemp{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, false); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp object: %s", err)
	}

	AddMetadataQueueTempHook(boil.BeforeInsertHook, metadataQueueTempBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempBeforeInsertHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.AfterInsertHook, metadataQueueTempAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempAfterInsertHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.AfterSelectHook, metadataQueueTempAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempAfterSelectHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.BeforeUpdateHook, metadataQueueTempBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempBeforeUpdateHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.AfterUpdateHook, metadataQueueTempAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempAfterUpdateHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.BeforeDeleteHook, metadataQueueTempBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempBeforeDeleteHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.AfterDeleteHook, metadataQueueTempAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempAfterDeleteHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.BeforeUpsertHook, metadataQueueTempBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempBeforeUpsertHooks = []MetadataQueueTempHook{}

	AddMetadataQueueTempHook(boil.AfterUpsertHook, metadataQueueTempAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	metadataQueueTempAfterUpsertHooks = []MetadataQueueTempHook{}
}

func testMetadataQueueTempsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMetadataQueueTempsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(metadataQueueTempColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMetadataQueueTempsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
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

func testMetadataQueueTempsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MetadataQueueTempSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMetadataQueueTempsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MetadataQueueTemps().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	metadataQueueTempDBTypes = map[string]string{`TransactionID`: `character varying`, `Key`: `character varying`, `Version`: `character varying`, `Value`: `jsonb`, `CommitedInstances`: `bigint`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`}
	_                        = bytes.MinRead
)

func testMetadataQueueTempsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(metadataQueueTempPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(metadataQueueTempAllColumns) == len(metadataQueueTempPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testMetadataQueueTempsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(metadataQueueTempAllColumns) == len(metadataQueueTempPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MetadataQueueTemp{}
	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, metadataQueueTempDBTypes, true, metadataQueueTempPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(metadataQueueTempAllColumns, metadataQueueTempPrimaryKeyColumns) {
		fields = metadataQueueTempAllColumns
	} else {
		fields = strmangle.SetComplement(
			metadataQueueTempAllColumns,
			metadataQueueTempPrimaryKeyColumns,
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

	slice := MetadataQueueTempSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testMetadataQueueTempsUpsert(t *testing.T) {
	t.Parallel()

	if len(metadataQueueTempAllColumns) == len(metadataQueueTempPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := MetadataQueueTemp{}
	if err = randomize.Struct(seed, &o, metadataQueueTempDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MetadataQueueTemp: %s", err)
	}

	count, err := MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, metadataQueueTempDBTypes, false, metadataQueueTempPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MetadataQueueTemp struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MetadataQueueTemp: %s", err)
	}

	count, err = MetadataQueueTemps().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}