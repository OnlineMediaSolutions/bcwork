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

func testSeatOwners(t *testing.T) {
	t.Parallel()

	query := SeatOwners()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testSeatOwnersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
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

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSeatOwnersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := SeatOwners().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSeatOwnersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SeatOwnerSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSeatOwnersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := SeatOwnerExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if SeatOwner exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SeatOwnerExists to return true, but got false.")
	}
}

func testSeatOwnersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	seatOwnerFound, err := FindSeatOwner(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if seatOwnerFound == nil {
		t.Error("want a record, got nil")
	}
}

func testSeatOwnersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = SeatOwners().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testSeatOwnersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := SeatOwners().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSeatOwnersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	seatOwnerOne := &SeatOwner{}
	seatOwnerTwo := &SeatOwner{}
	if err = randomize.Struct(seed, seatOwnerOne, seatOwnerDBTypes, false, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}
	if err = randomize.Struct(seed, seatOwnerTwo, seatOwnerDBTypes, false, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = seatOwnerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = seatOwnerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := SeatOwners().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSeatOwnersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	seatOwnerOne := &SeatOwner{}
	seatOwnerTwo := &SeatOwner{}
	if err = randomize.Struct(seed, seatOwnerOne, seatOwnerDBTypes, false, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}
	if err = randomize.Struct(seed, seatOwnerTwo, seatOwnerDBTypes, false, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = seatOwnerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = seatOwnerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func seatOwnerBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func seatOwnerAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *SeatOwner) error {
	*o = SeatOwner{}
	return nil
}

func testSeatOwnersHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &SeatOwner{}
	o := &SeatOwner{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, false); err != nil {
		t.Errorf("Unable to randomize SeatOwner object: %s", err)
	}

	AddSeatOwnerHook(boil.BeforeInsertHook, seatOwnerBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	seatOwnerBeforeInsertHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.AfterInsertHook, seatOwnerAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	seatOwnerAfterInsertHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.AfterSelectHook, seatOwnerAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	seatOwnerAfterSelectHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.BeforeUpdateHook, seatOwnerBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	seatOwnerBeforeUpdateHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.AfterUpdateHook, seatOwnerAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	seatOwnerAfterUpdateHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.BeforeDeleteHook, seatOwnerBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	seatOwnerBeforeDeleteHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.AfterDeleteHook, seatOwnerAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	seatOwnerAfterDeleteHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.BeforeUpsertHook, seatOwnerBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	seatOwnerBeforeUpsertHooks = []SeatOwnerHook{}

	AddSeatOwnerHook(boil.AfterUpsertHook, seatOwnerAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	seatOwnerAfterUpsertHooks = []SeatOwnerHook{}
}

func testSeatOwnersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSeatOwnersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(seatOwnerColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSeatOwnerToManyAdsTXTS(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
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

	queries.Assign(&b.SeatOwnerID, a.ID)
	queries.Assign(&c.SeatOwnerID, a.ID)
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
		if queries.Equal(v.SeatOwnerID, b.SeatOwnerID) {
			bFound = true
		}
		if queries.Equal(v.SeatOwnerID, c.SeatOwnerID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SeatOwnerSlice{&a}
	if err = a.L.LoadAdsTXTS(ctx, tx, false, (*[]*SeatOwner)(&slice), nil); err != nil {
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

func testSeatOwnerToManyDpos(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c Dpo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, dpoDBTypes, false, dpoColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, dpoDBTypes, false, dpoColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.SeatOwnerID, a.ID)
	queries.Assign(&c.SeatOwnerID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Dpos().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.SeatOwnerID, b.SeatOwnerID) {
			bFound = true
		}
		if queries.Equal(v.SeatOwnerID, c.SeatOwnerID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SeatOwnerSlice{&a}
	if err = a.L.LoadDpos(ctx, tx, false, (*[]*SeatOwner)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Dpos); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Dpos = nil
	if err = a.L.LoadDpos(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Dpos); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testSeatOwnerToManyAddOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
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

		if !queries.Equal(a.ID, first.SeatOwnerID) {
			t.Error("foreign key was wrong value", a.ID, first.SeatOwnerID)
		}
		if !queries.Equal(a.ID, second.SeatOwnerID) {
			t.Error("foreign key was wrong value", a.ID, second.SeatOwnerID)
		}

		if first.R.SeatOwner != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.SeatOwner != &a {
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

func testSeatOwnerToManySetOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
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

	if !queries.IsValuerNil(b.SeatOwnerID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.SeatOwnerID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.SeatOwnerID) {
		t.Error("foreign key was wrong value", a.ID, d.SeatOwnerID)
	}
	if !queries.Equal(a.ID, e.SeatOwnerID) {
		t.Error("foreign key was wrong value", a.ID, e.SeatOwnerID)
	}

	if b.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.SeatOwner != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.SeatOwner != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.AdsTXTS[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.AdsTXTS[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testSeatOwnerToManyRemoveOpAdsTXTS(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e AdsTXT

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
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

	if !queries.IsValuerNil(b.SeatOwnerID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.SeatOwnerID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.SeatOwner != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.SeatOwner != &a {
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

func testSeatOwnerToManyAddOpDpos(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e Dpo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Dpo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, dpoDBTypes, false, strmangle.SetComplement(dpoPrimaryKeyColumns, dpoColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Dpo{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDpos(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.SeatOwnerID) {
			t.Error("foreign key was wrong value", a.ID, first.SeatOwnerID)
		}
		if !queries.Equal(a.ID, second.SeatOwnerID) {
			t.Error("foreign key was wrong value", a.ID, second.SeatOwnerID)
		}

		if first.R.SeatOwner != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.SeatOwner != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Dpos[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Dpos[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Dpos().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testSeatOwnerToManySetOpDpos(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e Dpo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Dpo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, dpoDBTypes, false, strmangle.SetComplement(dpoPrimaryKeyColumns, dpoColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDpos(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Dpos().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDpos(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Dpos().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.SeatOwnerID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.SeatOwnerID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.SeatOwnerID) {
		t.Error("foreign key was wrong value", a.ID, d.SeatOwnerID)
	}
	if !queries.Equal(a.ID, e.SeatOwnerID) {
		t.Error("foreign key was wrong value", a.ID, e.SeatOwnerID)
	}

	if b.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.SeatOwner != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.SeatOwner != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Dpos[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Dpos[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testSeatOwnerToManyRemoveOpDpos(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a SeatOwner
	var b, c, d, e Dpo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, seatOwnerDBTypes, false, strmangle.SetComplement(seatOwnerPrimaryKeyColumns, seatOwnerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Dpo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, dpoDBTypes, false, strmangle.SetComplement(dpoPrimaryKeyColumns, dpoColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddDpos(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Dpos().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDpos(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Dpos().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.SeatOwnerID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.SeatOwnerID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.SeatOwner != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.SeatOwner != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.SeatOwner != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Dpos) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Dpos[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Dpos[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testSeatOwnersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
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

func testSeatOwnersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SeatOwnerSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testSeatOwnersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := SeatOwners().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	seatOwnerDBTypes = map[string]string{`ID`: `integer`, `SeatOwnerName`: `character varying`, `SeatOwnerDomain`: `character varying`, `PublisherAccount`: `character varying`, `CertificationAuthorityID`: `character varying`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`}
	_                = bytes.MinRead
)

func testSeatOwnersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(seatOwnerPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(seatOwnerAllColumns) == len(seatOwnerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testSeatOwnersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(seatOwnerAllColumns) == len(seatOwnerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &SeatOwner{}
	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, seatOwnerDBTypes, true, seatOwnerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(seatOwnerAllColumns, seatOwnerPrimaryKeyColumns) {
		fields = seatOwnerAllColumns
	} else {
		fields = strmangle.SetComplement(
			seatOwnerAllColumns,
			seatOwnerPrimaryKeyColumns,
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

	slice := SeatOwnerSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testSeatOwnersUpsert(t *testing.T) {
	t.Parallel()

	if len(seatOwnerAllColumns) == len(seatOwnerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := SeatOwner{}
	if err = randomize.Struct(seed, &o, seatOwnerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert SeatOwner: %s", err)
	}

	count, err := SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, seatOwnerDBTypes, false, seatOwnerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SeatOwner struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert SeatOwner: %s", err)
	}

	count, err = SeatOwners().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
