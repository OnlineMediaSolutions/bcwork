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

func testCompetitors(t *testing.T) {
	t.Parallel()

	query := Competitors()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testCompetitorsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
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

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompetitorsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Competitors().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompetitorsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CompetitorSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompetitorsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := CompetitorExists(ctx, tx, o.Name)
	if err != nil {
		t.Errorf("Unable to check if Competitor exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CompetitorExists to return true, but got false.")
	}
}

func testCompetitorsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	competitorFound, err := FindCompetitor(ctx, tx, o.Name)
	if err != nil {
		t.Error(err)
	}

	if competitorFound == nil {
		t.Error("want a record, got nil")
	}
}

func testCompetitorsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Competitors().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testCompetitorsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Competitors().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCompetitorsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	competitorOne := &Competitor{}
	competitorTwo := &Competitor{}
	if err = randomize.Struct(seed, competitorOne, competitorDBTypes, false, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}
	if err = randomize.Struct(seed, competitorTwo, competitorDBTypes, false, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = competitorOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = competitorTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Competitors().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCompetitorsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	competitorOne := &Competitor{}
	competitorTwo := &Competitor{}
	if err = randomize.Struct(seed, competitorOne, competitorDBTypes, false, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}
	if err = randomize.Struct(seed, competitorTwo, competitorDBTypes, false, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = competitorOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = competitorTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func competitorBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func competitorAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Competitor) error {
	*o = Competitor{}
	return nil
}

func testCompetitorsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Competitor{}
	o := &Competitor{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, competitorDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Competitor object: %s", err)
	}

	AddCompetitorHook(boil.BeforeInsertHook, competitorBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	competitorBeforeInsertHooks = []CompetitorHook{}

	AddCompetitorHook(boil.AfterInsertHook, competitorAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	competitorAfterInsertHooks = []CompetitorHook{}

	AddCompetitorHook(boil.AfterSelectHook, competitorAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	competitorAfterSelectHooks = []CompetitorHook{}

	AddCompetitorHook(boil.BeforeUpdateHook, competitorBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	competitorBeforeUpdateHooks = []CompetitorHook{}

	AddCompetitorHook(boil.AfterUpdateHook, competitorAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	competitorAfterUpdateHooks = []CompetitorHook{}

	AddCompetitorHook(boil.BeforeDeleteHook, competitorBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	competitorBeforeDeleteHooks = []CompetitorHook{}

	AddCompetitorHook(boil.AfterDeleteHook, competitorAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	competitorAfterDeleteHooks = []CompetitorHook{}

	AddCompetitorHook(boil.BeforeUpsertHook, competitorBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	competitorBeforeUpsertHooks = []CompetitorHook{}

	AddCompetitorHook(boil.AfterUpsertHook, competitorAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	competitorAfterUpsertHooks = []CompetitorHook{}
}

func testCompetitorsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCompetitorsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(competitorColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCompetitorOneToOneSellersJSONHistoryUsingCompetitorNameSellersJSONHistory(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var foreign SellersJSONHistory
	var local Competitor

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, sellersJSONHistoryDBTypes, true, sellersJSONHistoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SellersJSONHistory struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreign.CompetitorName = local.Name
	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.CompetitorNameSellersJSONHistory().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.CompetitorName != foreign.CompetitorName {
		t.Errorf("want: %v, got %v", foreign.CompetitorName, check.CompetitorName)
	}

	ranAfterSelectHook := false
	AddSellersJSONHistoryHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *SellersJSONHistory) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := CompetitorSlice{&local}
	if err = local.L.LoadCompetitorNameSellersJSONHistory(ctx, tx, false, (*[]*Competitor)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.CompetitorNameSellersJSONHistory == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.CompetitorNameSellersJSONHistory = nil
	if err = local.L.LoadCompetitorNameSellersJSONHistory(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.CompetitorNameSellersJSONHistory == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testCompetitorOneToOneSetOpSellersJSONHistoryUsingCompetitorNameSellersJSONHistory(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Competitor
	var b, c SellersJSONHistory

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, competitorDBTypes, false, strmangle.SetComplement(competitorPrimaryKeyColumns, competitorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, sellersJSONHistoryDBTypes, false, strmangle.SetComplement(sellersJSONHistoryPrimaryKeyColumns, sellersJSONHistoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, sellersJSONHistoryDBTypes, false, strmangle.SetComplement(sellersJSONHistoryPrimaryKeyColumns, sellersJSONHistoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*SellersJSONHistory{&b, &c} {
		err = a.SetCompetitorNameSellersJSONHistory(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.CompetitorNameSellersJSONHistory != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.CompetitorNameCompetitor != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		if a.Name != x.CompetitorName {
			t.Error("foreign key was wrong value", a.Name)
		}

		if exists, err := SellersJSONHistoryExists(ctx, tx, x.CompetitorName); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'x' to exist")
		}

		if a.Name != x.CompetitorName {
			t.Error("foreign key was wrong value", a.Name, x.CompetitorName)
		}

		if _, err = x.Delete(ctx, tx); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}

func testCompetitorsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
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

func testCompetitorsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CompetitorSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testCompetitorsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Competitors().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	competitorDBTypes = map[string]string{`Name`: `character varying`, `URL`: `text`, `Type`: `character varying`, `Position`: `character varying`}
	_                 = bytes.MinRead
)

func testCompetitorsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(competitorPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(competitorAllColumns) == len(competitorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testCompetitorsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(competitorAllColumns) == len(competitorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Competitor{}
	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, competitorDBTypes, true, competitorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(competitorAllColumns, competitorPrimaryKeyColumns) {
		fields = competitorAllColumns
	} else {
		fields = strmangle.SetComplement(
			competitorAllColumns,
			competitorPrimaryKeyColumns,
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

	slice := CompetitorSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testCompetitorsUpsert(t *testing.T) {
	t.Parallel()

	if len(competitorAllColumns) == len(competitorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Competitor{}
	if err = randomize.Struct(seed, &o, competitorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Competitor: %s", err)
	}

	count, err := Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, competitorDBTypes, false, competitorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Competitor struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Competitor: %s", err)
	}

	count, err = Competitors().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
