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

func testBidCachings(t *testing.T) {
	t.Parallel()

	query := BidCachings()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testBidCachingsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
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

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBidCachingsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := BidCachings().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBidCachingsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BidCachingSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBidCachingsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := BidCachingExists(ctx, tx, o.RuleID)
	if err != nil {
		t.Errorf("Unable to check if BidCaching exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BidCachingExists to return true, but got false.")
	}
}

func testBidCachingsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	bidCachingFound, err := FindBidCaching(ctx, tx, o.RuleID)
	if err != nil {
		t.Error(err)
	}

	if bidCachingFound == nil {
		t.Error("want a record, got nil")
	}
}

func testBidCachingsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = BidCachings().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testBidCachingsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := BidCachings().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBidCachingsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	bidCachingOne := &BidCaching{}
	bidCachingTwo := &BidCaching{}
	if err = randomize.Struct(seed, bidCachingOne, bidCachingDBTypes, false, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}
	if err = randomize.Struct(seed, bidCachingTwo, bidCachingDBTypes, false, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = bidCachingOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = bidCachingTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BidCachings().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBidCachingsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	bidCachingOne := &BidCaching{}
	bidCachingTwo := &BidCaching{}
	if err = randomize.Struct(seed, bidCachingOne, bidCachingDBTypes, false, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}
	if err = randomize.Struct(seed, bidCachingTwo, bidCachingDBTypes, false, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = bidCachingOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = bidCachingTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func bidCachingBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func bidCachingAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *BidCaching) error {
	*o = BidCaching{}
	return nil
}

func testBidCachingsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &BidCaching{}
	o := &BidCaching{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, bidCachingDBTypes, false); err != nil {
		t.Errorf("Unable to randomize BidCaching object: %s", err)
	}

	AddBidCachingHook(boil.BeforeInsertHook, bidCachingBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	bidCachingBeforeInsertHooks = []BidCachingHook{}

	AddBidCachingHook(boil.AfterInsertHook, bidCachingAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	bidCachingAfterInsertHooks = []BidCachingHook{}

	AddBidCachingHook(boil.AfterSelectHook, bidCachingAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	bidCachingAfterSelectHooks = []BidCachingHook{}

	AddBidCachingHook(boil.BeforeUpdateHook, bidCachingBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	bidCachingBeforeUpdateHooks = []BidCachingHook{}

	AddBidCachingHook(boil.AfterUpdateHook, bidCachingAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	bidCachingAfterUpdateHooks = []BidCachingHook{}

	AddBidCachingHook(boil.BeforeDeleteHook, bidCachingBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	bidCachingBeforeDeleteHooks = []BidCachingHook{}

	AddBidCachingHook(boil.AfterDeleteHook, bidCachingAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	bidCachingAfterDeleteHooks = []BidCachingHook{}

	AddBidCachingHook(boil.BeforeUpsertHook, bidCachingBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	bidCachingBeforeUpsertHooks = []BidCachingHook{}

	AddBidCachingHook(boil.AfterUpsertHook, bidCachingAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	bidCachingAfterUpsertHooks = []BidCachingHook{}
}

func testBidCachingsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBidCachingsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(bidCachingColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBidCachingToOnePublisherUsingBidCachingPublisher(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local BidCaching
	var foreign Publisher

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, bidCachingDBTypes, false, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, publisherDBTypes, false, publisherColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Publisher struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.Publisher = foreign.PublisherID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.BidCachingPublisher().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.PublisherID != foreign.PublisherID {
		t.Errorf("want: %v, got %v", foreign.PublisherID, check.PublisherID)
	}

	ranAfterSelectHook := false
	AddPublisherHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *Publisher) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := BidCachingSlice{&local}
	if err = local.L.LoadBidCachingPublisher(ctx, tx, false, (*[]*BidCaching)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.BidCachingPublisher == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.BidCachingPublisher = nil
	if err = local.L.LoadBidCachingPublisher(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.BidCachingPublisher == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testBidCachingToOneSetOpPublisherUsingBidCachingPublisher(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a BidCaching
	var b, c Publisher

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bidCachingDBTypes, false, strmangle.SetComplement(bidCachingPrimaryKeyColumns, bidCachingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, publisherDBTypes, false, strmangle.SetComplement(publisherPrimaryKeyColumns, publisherColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, publisherDBTypes, false, strmangle.SetComplement(publisherPrimaryKeyColumns, publisherColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Publisher{&b, &c} {
		err = a.SetBidCachingPublisher(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.BidCachingPublisher != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.BidCachings[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.Publisher != x.PublisherID {
			t.Error("foreign key was wrong value", a.Publisher)
		}

		zero := reflect.Zero(reflect.TypeOf(a.Publisher))
		reflect.Indirect(reflect.ValueOf(&a.Publisher)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.Publisher != x.PublisherID {
			t.Error("foreign key was wrong value", a.Publisher, x.PublisherID)
		}
	}
}

func testBidCachingsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
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

func testBidCachingsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BidCachingSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBidCachingsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := BidCachings().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	bidCachingDBTypes = map[string]string{`Publisher`: `character varying`, `Domain`: `character varying`, `Country`: `character varying`, `Device`: `character varying`, `BidCaching`: `smallint`, `CreatedAt`: `timestamp without time zone`, `UpdatedAt`: `timestamp without time zone`, `RuleID`: `character varying`, `DemandPartnerID`: `character varying`, `Browser`: `character varying`, `Os`: `character varying`, `PlacementType`: `character varying`, `Active`: `boolean`}
	_                 = bytes.MinRead
)

func testBidCachingsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(bidCachingPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(bidCachingAllColumns) == len(bidCachingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testBidCachingsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(bidCachingAllColumns) == len(bidCachingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &BidCaching{}
	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, bidCachingDBTypes, true, bidCachingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(bidCachingAllColumns, bidCachingPrimaryKeyColumns) {
		fields = bidCachingAllColumns
	} else {
		fields = strmangle.SetComplement(
			bidCachingAllColumns,
			bidCachingPrimaryKeyColumns,
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

	slice := BidCachingSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testBidCachingsUpsert(t *testing.T) {
	t.Parallel()

	if len(bidCachingAllColumns) == len(bidCachingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := BidCaching{}
	if err = randomize.Struct(seed, &o, bidCachingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BidCaching: %s", err)
	}

	count, err := BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, bidCachingDBTypes, false, bidCachingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize BidCaching struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert BidCaching: %s", err)
	}

	count, err = BidCachings().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}