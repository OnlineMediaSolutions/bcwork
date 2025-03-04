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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// SearchView is an object representing the database table.
type SearchView struct {
	SectionType   null.String `boil:"section_type" json:"section_type,omitempty" toml:"section_type" yaml:"section_type,omitempty"`
	PublisherID   null.String `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id,omitempty"`
	PublisherName null.String `boil:"publisher_name" json:"publisher_name,omitempty" toml:"publisher_name" yaml:"publisher_name,omitempty"`
	Domain        null.String `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Query         null.String `boil:"query" json:"query,omitempty" toml:"query" yaml:"query,omitempty"`
}

var SearchViewColumns = struct {
	SectionType   string
	PublisherID   string
	PublisherName string
	Domain        string
	Query         string
}{
	SectionType:   "section_type",
	PublisherID:   "publisher_id",
	PublisherName: "publisher_name",
	Domain:        "domain",
	Query:         "query",
}

var SearchViewTableColumns = struct {
	SectionType   string
	PublisherID   string
	PublisherName string
	Domain        string
	Query         string
}{
	SectionType:   "search_view.section_type",
	PublisherID:   "search_view.publisher_id",
	PublisherName: "search_view.publisher_name",
	Domain:        "search_view.domain",
	Query:         "search_view.query",
}

// Generated where

var SearchViewWhere = struct {
	SectionType   whereHelpernull_String
	PublisherID   whereHelpernull_String
	PublisherName whereHelpernull_String
	Domain        whereHelpernull_String
	Query         whereHelpernull_String
}{
	SectionType:   whereHelpernull_String{field: "\"search_view\".\"section_type\""},
	PublisherID:   whereHelpernull_String{field: "\"search_view\".\"publisher_id\""},
	PublisherName: whereHelpernull_String{field: "\"search_view\".\"publisher_name\""},
	Domain:        whereHelpernull_String{field: "\"search_view\".\"domain\""},
	Query:         whereHelpernull_String{field: "\"search_view\".\"query\""},
}

var (
	searchViewAllColumns            = []string{"section_type", "publisher_id", "publisher_name", "domain", "query"}
	searchViewColumnsWithoutDefault = []string{}
	searchViewColumnsWithDefault    = []string{"section_type", "publisher_id", "publisher_name", "domain", "query"}
	searchViewPrimaryKeyColumns     = []string{}
	searchViewGeneratedColumns      = []string{}
)

type (
	// SearchViewSlice is an alias for a slice of pointers to SearchView.
	// This should almost always be used instead of []SearchView.
	SearchViewSlice []*SearchView
	// SearchViewHook is the signature for custom SearchView hook methods
	SearchViewHook func(context.Context, boil.ContextExecutor, *SearchView) error

	searchViewQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	searchViewType           = reflect.TypeOf(&SearchView{})
	searchViewMapping        = queries.MakeStructMapping(searchViewType)
	searchViewInsertCacheMut sync.RWMutex
	searchViewInsertCache    = make(map[string]insertCache)
	searchViewUpdateCacheMut sync.RWMutex
	searchViewUpdateCache    = make(map[string]updateCache)
	searchViewUpsertCacheMut sync.RWMutex
	searchViewUpsertCache    = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
	// These are used in some views
	_ = fmt.Sprintln("")
	_ = reflect.Int
	_ = strings.Builder{}
	_ = sync.Mutex{}
	_ = strmangle.Plural("")
	_ = strconv.IntSize
)

var searchViewAfterSelectMu sync.Mutex
var searchViewAfterSelectHooks []SearchViewHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SearchView) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range searchViewAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSearchViewHook registers your hook function for all future operations.
func AddSearchViewHook(hookPoint boil.HookPoint, searchViewHook SearchViewHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		searchViewAfterSelectMu.Lock()
		searchViewAfterSelectHooks = append(searchViewAfterSelectHooks, searchViewHook)
		searchViewAfterSelectMu.Unlock()
	}
}

// One returns a single searchView record from the query.
func (q searchViewQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SearchView, error) {
	o := &SearchView{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for search_view")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SearchView records from the query.
func (q searchViewQuery) All(ctx context.Context, exec boil.ContextExecutor) (SearchViewSlice, error) {
	var o []*SearchView

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SearchView slice")
	}

	if len(searchViewAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SearchView records in the query.
func (q searchViewQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count search_view rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q searchViewQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if search_view exists")
	}

	return count > 0, nil
}

// SearchViews retrieves all the records using an executor.
func SearchViews(mods ...qm.QueryMod) searchViewQuery {
	mods = append(mods, qm.From("\"search_view\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"search_view\".*"})
	}

	return searchViewQuery{q}
}
