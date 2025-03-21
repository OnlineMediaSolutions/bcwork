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

func testDemandPartnerChildren(t *testing.T) {
	t.Parallel()

	query := DemandPartnerChildren()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testDemandPartnerChildrenDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
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

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerChildrenQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := DemandPartnerChildren().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerChildrenSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DemandPartnerChildSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemandPartnerChildrenExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := DemandPartnerChildExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if DemandPartnerChild exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DemandPartnerChildExists to return true, but got false.")
	}
}

func testDemandPartnerChildrenFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	demandPartnerChildFound, err := FindDemandPartnerChild(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if demandPartnerChildFound == nil {
		t.Error("want a record, got nil")
	}
}

func testDemandPartnerChildrenBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = DemandPartnerChildren().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testDemandPartnerChildrenOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := DemandPartnerChildren().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDemandPartnerChildrenAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demandPartnerChildOne := &DemandPartnerChild{}
	demandPartnerChildTwo := &DemandPartnerChild{}
	if err = randomize.Struct(seed, demandPartnerChildOne, demandPartnerChildDBTypes, false, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}
	if err = randomize.Struct(seed, demandPartnerChildTwo, demandPartnerChildDBTypes, false, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = demandPartnerChildOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = demandPartnerChildTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DemandPartnerChildren().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDemandPartnerChildrenCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	demandPartnerChildOne := &DemandPartnerChild{}
	demandPartnerChildTwo := &DemandPartnerChild{}
	if err = randomize.Struct(seed, demandPartnerChildOne, demandPartnerChildDBTypes, false, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}
	if err = randomize.Struct(seed, demandPartnerChildTwo, demandPartnerChildDBTypes, false, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = demandPartnerChildOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = demandPartnerChildTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func demandPartnerChildBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func demandPartnerChildAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerChild) error {
	*o = DemandPartnerChild{}
	return nil
}

func testDemandPartnerChildrenHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &DemandPartnerChild{}
	o := &DemandPartnerChild{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, false); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild object: %s", err)
	}

	AddDemandPartnerChildHook(boil.BeforeInsertHook, demandPartnerChildBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildBeforeInsertHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.AfterInsertHook, demandPartnerChildAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildAfterInsertHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.AfterSelectHook, demandPartnerChildAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildAfterSelectHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.BeforeUpdateHook, demandPartnerChildBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildBeforeUpdateHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.AfterUpdateHook, demandPartnerChildAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildAfterUpdateHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.BeforeDeleteHook, demandPartnerChildBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildBeforeDeleteHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.AfterDeleteHook, demandPartnerChildAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildAfterDeleteHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.BeforeUpsertHook, demandPartnerChildBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildBeforeUpsertHooks = []DemandPartnerChildHook{}

	AddDemandPartnerChildHook(boil.AfterUpsertHook, demandPartnerChildAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	demandPartnerChildAfterUpsertHooks = []DemandPartnerChildHook{}
}

func testDemandPartnerChildrenInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemandPartnerChildrenInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(demandPartnerChildColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemandPartnerChildToManyAdsTXTS(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DemandPartnerChild
	var b, c AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, adsTXTDBTypes, false, adsTXTColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, adsTXTDBTypes, false, adsTXTColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.DemandPartnerChildID, a.ID)
	queries.Assign(&c.DemandPartnerChildID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.AdsTXTS().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.DemandPartnerChildID, b.DemandPartnerChildID) {
			bFound = true
		}
		if queries.Equal(v.DemandPartnerChildID, c.DemandPartnerChildID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := DemandPartnerChildSlice{&a}
	if err = a.L.LoadAdsTXTS(ctx, tx, false, (*[]*DemandPartnerChild)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AdsTXTS); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.AdsTXTS = nil
	if err = a.L.LoadAdsTXTS(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AdsTXTS); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testDemandPartnerChildToManyAddOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DemandPartnerChild
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demandPartnerChildDBTypes, false, strmangle.SetComplement(demandPartnerChildPrimaryKeyColumns, demandPartnerChildColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*AdsTXT{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, adsTXTDBTypes, false, strmangle.SetComplement(adsTXTPrimaryKeyColumns, adsTXTColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*AdsTXT{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAdsTXTS(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.DemandPartnerChildID) {
			t.Error("foreign key was wrong value", a.ID, first.DemandPartnerChildID)
		}
		if !queries.Equal(a.ID, second.DemandPartnerChildID) {
			t.Error("foreign key was wrong value", a.ID, second.DemandPartnerChildID)
		}

		if first.R.DemandPartnerChild != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.DemandPartnerChild != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.AdsTXTS[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.AdsTXTS[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.AdsTXTS().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testDemandPartnerChildToManySetOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DemandPartnerChild
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demandPartnerChildDBTypes, false, strmangle.SetComplement(demandPartnerChildPrimaryKeyColumns, demandPartnerChildColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*AdsTXT{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, adsTXTDBTypes, false, strmangle.SetComplement(adsTXTPrimaryKeyColumns, adsTXTColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetAdsTXTS(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.AdsTXTS().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetAdsTXTS(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.AdsTXTS().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.DemandPartnerChildID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.DemandPartnerChildID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.DemandPartnerChildID) {
		t.Error("foreign key was wrong value", a.ID, d.DemandPartnerChildID)
	}
	if !queries.Equal(a.ID, e.DemandPartnerChildID) {
		t.Error("foreign key was wrong value", a.ID, e.DemandPartnerChildID)
	}

	if b.R.DemandPartnerChild != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.DemandPartnerChild != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.DemandPartnerChild != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.DemandPartnerChild != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.AdsTXTS[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.AdsTXTS[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testDemandPartnerChildToManyRemoveOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DemandPartnerChild
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demandPartnerChildDBTypes, false, strmangle.SetComplement(demandPartnerChildPrimaryKeyColumns, demandPartnerChildColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*AdsTXT{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, adsTXTDBTypes, false, strmangle.SetComplement(adsTXTPrimaryKeyColumns, adsTXTColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddAdsTXTS(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.AdsTXTS().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveAdsTXTS(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.AdsTXTS().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.DemandPartnerChildID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.DemandPartnerChildID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.DemandPartnerChild != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.DemandPartnerChild != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.DemandPartnerChild != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.DemandPartnerChild != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.AdsTXTS) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.AdsTXTS[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.AdsTXTS[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testDemandPartnerChildToOneDemandPartnerConnectionUsingDPConnection(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local DemandPartnerChild
	var foreign DemandPartnerConnection

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, demandPartnerChildDBTypes, false, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, demandPartnerConnectionDBTypes, false, demandPartnerConnectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerConnection struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.DPConnectionID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.DPConnection().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	ranAfterSelectHook := false
	AddDemandPartnerConnectionHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *DemandPartnerConnection) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := DemandPartnerChildSlice{&local}
	if err = local.L.LoadDPConnection(ctx, tx, false, (*[]*DemandPartnerChild)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.DPConnection == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.DPConnection = nil
	if err = local.L.LoadDPConnection(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.DPConnection == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testDemandPartnerChildToOneSetOpDemandPartnerConnectionUsingDPConnection(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DemandPartnerChild
	var b, c DemandPartnerConnection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demandPartnerChildDBTypes, false, strmangle.SetComplement(demandPartnerChildPrimaryKeyColumns, demandPartnerChildColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, demandPartnerConnectionDBTypes, false, strmangle.SetComplement(demandPartnerConnectionPrimaryKeyColumns, demandPartnerConnectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, demandPartnerConnectionDBTypes, false, strmangle.SetComplement(demandPartnerConnectionPrimaryKeyColumns, demandPartnerConnectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*DemandPartnerConnection{&b, &c} {
		err = a.SetDPConnection(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.DPConnection != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DPConnectionDemandPartnerChildren[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DPConnectionID != x.ID {
			t.Error("foreign key was wrong value", a.DPConnectionID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DPConnectionID))
		reflect.Indirect(reflect.ValueOf(&a.DPConnectionID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DPConnectionID != x.ID {
			t.Error("foreign key was wrong value", a.DPConnectionID, x.ID)
		}
	}
}

func testDemandPartnerChildrenReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
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

func testDemandPartnerChildrenReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DemandPartnerChildSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testDemandPartnerChildrenSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DemandPartnerChildren().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	demandPartnerChildDBTypes = map[string]string{`ID`: `integer`, `DPChildName`: `character varying`, `DPDomain`: `character varying`, `PublisherAccount`: `character varying`, `CertificationAuthorityID`: `character varying`, `IsRequiredForAdsTXT`: `boolean`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`, `IsDirect`: `boolean`, `DPConnectionID`: `integer`}
	_                         = bytes.MinRead
)

func testDemandPartnerChildrenUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(demandPartnerChildPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(demandPartnerChildAllColumns) == len(demandPartnerChildPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testDemandPartnerChildrenSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(demandPartnerChildAllColumns) == len(demandPartnerChildPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DemandPartnerChild{}
	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, demandPartnerChildDBTypes, true, demandPartnerChildPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(demandPartnerChildAllColumns, demandPartnerChildPrimaryKeyColumns) {
		fields = demandPartnerChildAllColumns
	} else {
		fields = strmangle.SetComplement(
			demandPartnerChildAllColumns,
			demandPartnerChildPrimaryKeyColumns,
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

	slice := DemandPartnerChildSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testDemandPartnerChildrenUpsert(t *testing.T) {
	t.Parallel()

	if len(demandPartnerChildAllColumns) == len(demandPartnerChildPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := DemandPartnerChild{}
	if err = randomize.Struct(seed, &o, demandPartnerChildDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DemandPartnerChild: %s", err)
	}

	count, err := DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, demandPartnerChildDBTypes, false, demandPartnerChildPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemandPartnerChild struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DemandPartnerChild: %s", err)
	}

	count, err = DemandPartnerChildren().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
