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

func testConfigurations(t *testing.T) {
	t.Parallel()

	query := Configurations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testConfigurationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
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

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConfigurationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Configurations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConfigurationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ConfigurationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConfigurationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ConfigurationExists(ctx, tx, o.Key)
	if err != nil {
		t.Errorf("Unable to check if Configuration exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ConfigurationExists to return true, but got false.")
	}
}

func testConfigurationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	configurationFound, err := FindConfiguration(ctx, tx, o.Key)
	if err != nil {
		t.Error(err)
	}

	if configurationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testConfigurationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Configurations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testConfigurationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Configurations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testConfigurationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	configurationOne := &Configuration{}
	configurationTwo := &Configuration{}
	if err = randomize.Struct(seed, configurationOne, configurationDBTypes, false, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}
	if err = randomize.Struct(seed, configurationTwo, configurationDBTypes, false, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = configurationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = configurationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Configurations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testConfigurationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	configurationOne := &Configuration{}
	configurationTwo := &Configuration{}
	if err = randomize.Struct(seed, configurationOne, configurationDBTypes, false, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}
	if err = randomize.Struct(seed, configurationTwo, configurationDBTypes, false, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = configurationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = configurationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func configurationBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func configurationAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Configuration) error {
	*o = Configuration{}
	return nil
}

func testConfigurationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Configuration{}
	o := &Configuration{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, configurationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Configuration object: %s", err)
	}

	AddConfigurationHook(boil.BeforeInsertHook, configurationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	configurationBeforeInsertHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.AfterInsertHook, configurationAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	configurationAfterInsertHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.AfterSelectHook, configurationAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	configurationAfterSelectHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.BeforeUpdateHook, configurationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	configurationBeforeUpdateHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.AfterUpdateHook, configurationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	configurationAfterUpdateHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.BeforeDeleteHook, configurationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	configurationBeforeDeleteHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.AfterDeleteHook, configurationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	configurationAfterDeleteHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.BeforeUpsertHook, configurationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	configurationBeforeUpsertHooks = []ConfigurationHook{}

	AddConfigurationHook(boil.AfterUpsertHook, configurationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	configurationAfterUpsertHooks = []ConfigurationHook{}
}

func testConfigurationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testConfigurationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(configurationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testConfigurationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
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

func testConfigurationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ConfigurationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testConfigurationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Configurations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	configurationDBTypes = map[string]string{`Key`: `character varying`, `Value`: `text`, `Description`: `text`, `UpdatedAt`: `timestamp without time zone`, `CreatedAt`: `timestamp without time zone`}
	_                    = bytes.MinRead
)

func testConfigurationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(configurationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(configurationAllColumns) == len(configurationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testConfigurationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(configurationAllColumns) == len(configurationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Configuration{}
	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, configurationDBTypes, true, configurationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(configurationAllColumns, configurationPrimaryKeyColumns) {
		fields = configurationAllColumns
	} else {
		fields = strmangle.SetComplement(
			configurationAllColumns,
			configurationPrimaryKeyColumns,
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

	slice := ConfigurationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testConfigurationsUpsert(t *testing.T) {
	t.Parallel()

	if len(configurationAllColumns) == len(configurationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Configuration{}
	if err = randomize.Struct(seed, &o, configurationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Configuration: %s", err)
	}

	count, err := Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, configurationDBTypes, false, configurationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Configuration struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Configuration: %s", err)
	}

	count, err = Configurations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
