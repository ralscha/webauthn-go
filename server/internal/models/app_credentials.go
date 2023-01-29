// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// AppCredential is an object representing the database table.
type AppCredential struct {
	ID         []byte `boil:"id" json:"id" toml:"id" yaml:"id"`
	AppUserID  int64  `boil:"app_user_id" json:"app_user_id" toml:"app_user_id" yaml:"app_user_id"`
	Credential []byte `boil:"credential" json:"credential" toml:"credential" yaml:"credential"`

	R *appCredentialR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L appCredentialL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AppCredentialColumns = struct {
	ID         string
	AppUserID  string
	Credential string
}{
	ID:         "id",
	AppUserID:  "app_user_id",
	Credential: "credential",
}

var AppCredentialTableColumns = struct {
	ID         string
	AppUserID  string
	Credential string
}{
	ID:         "app_credentials.id",
	AppUserID:  "app_credentials.app_user_id",
	Credential: "app_credentials.credential",
}

// Generated where

type whereHelper__byte struct{ field string }

func (w whereHelper__byte) EQ(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelper__byte) NEQ(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelper__byte) LT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelper__byte) LTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelper__byte) GT(x []byte) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelper__byte) GTE(x []byte) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var AppCredentialWhere = struct {
	ID         whereHelper__byte
	AppUserID  whereHelperint64
	Credential whereHelper__byte
}{
	ID:         whereHelper__byte{field: "\"app_credentials\".\"id\""},
	AppUserID:  whereHelperint64{field: "\"app_credentials\".\"app_user_id\""},
	Credential: whereHelper__byte{field: "\"app_credentials\".\"credential\""},
}

// AppCredentialRels is where relationship names are stored.
var AppCredentialRels = struct {
	AppUser string
}{
	AppUser: "AppUser",
}

// appCredentialR is where relationships are stored.
type appCredentialR struct {
	AppUser *AppUser `boil:"AppUser" json:"AppUser" toml:"AppUser" yaml:"AppUser"`
}

// NewStruct creates a new relationship struct
func (*appCredentialR) NewStruct() *appCredentialR {
	return &appCredentialR{}
}

func (r *appCredentialR) GetAppUser() *AppUser {
	if r == nil {
		return nil
	}
	return r.AppUser
}

// appCredentialL is where Load methods for each relationship are stored.
type appCredentialL struct{}

var (
	appCredentialAllColumns            = []string{"id", "app_user_id", "credential"}
	appCredentialColumnsWithoutDefault = []string{"id", "app_user_id", "credential"}
	appCredentialColumnsWithDefault    = []string{}
	appCredentialPrimaryKeyColumns     = []string{"id", "app_user_id"}
	appCredentialGeneratedColumns      = []string{}
)

type (
	// AppCredentialSlice is an alias for a slice of pointers to AppCredential.
	// This should almost always be used instead of []AppCredential.
	AppCredentialSlice []*AppCredential

	appCredentialQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	appCredentialType                 = reflect.TypeOf(&AppCredential{})
	appCredentialMapping              = queries.MakeStructMapping(appCredentialType)
	appCredentialPrimaryKeyMapping, _ = queries.BindMapping(appCredentialType, appCredentialMapping, appCredentialPrimaryKeyColumns)
	appCredentialInsertCacheMut       sync.RWMutex
	appCredentialInsertCache          = make(map[string]insertCache)
	appCredentialUpdateCacheMut       sync.RWMutex
	appCredentialUpdateCache          = make(map[string]updateCache)
	appCredentialUpsertCacheMut       sync.RWMutex
	appCredentialUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single appCredential record from the query.
func (q appCredentialQuery) One(ctx context.Context, exec boil.ContextExecutor) (*AppCredential, error) {
	o := &AppCredential{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for app_credentials")
	}

	return o, nil
}

// All returns all AppCredential records from the query.
func (q appCredentialQuery) All(ctx context.Context, exec boil.ContextExecutor) (AppCredentialSlice, error) {
	var o []*AppCredential

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to AppCredential slice")
	}

	return o, nil
}

// Count returns the count of all AppCredential records in the query.
func (q appCredentialQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count app_credentials rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q appCredentialQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if app_credentials exists")
	}

	return count > 0, nil
}

// AppUser pointed to by the foreign key.
func (o *AppCredential) AppUser(mods ...qm.QueryMod) appUserQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AppUserID),
	}

	queryMods = append(queryMods, mods...)

	return AppUsers(queryMods...)
}

// LoadAppUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (appCredentialL) LoadAppUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeAppCredential interface{}, mods queries.Applicator) error {
	var slice []*AppCredential
	var object *AppCredential

	if singular {
		var ok bool
		object, ok = maybeAppCredential.(*AppCredential)
		if !ok {
			object = new(AppCredential)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeAppCredential)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeAppCredential))
			}
		}
	} else {
		s, ok := maybeAppCredential.(*[]*AppCredential)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeAppCredential)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeAppCredential))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &appCredentialR{}
		}
		args = append(args, object.AppUserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &appCredentialR{}
			}

			for _, a := range args {
				if a == obj.AppUserID {
					continue Outer
				}
			}

			args = append(args, obj.AppUserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`app_user`),
		qm.WhereIn(`app_user.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AppUser")
	}

	var resultSlice []*AppUser
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AppUser")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for app_user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for app_user")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.AppUser = foreign
		if foreign.R == nil {
			foreign.R = &appUserR{}
		}
		foreign.R.AppCredentials = append(foreign.R.AppCredentials, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AppUserID == foreign.ID {
				local.R.AppUser = foreign
				if foreign.R == nil {
					foreign.R = &appUserR{}
				}
				foreign.R.AppCredentials = append(foreign.R.AppCredentials, local)
				break
			}
		}
	}

	return nil
}

// SetAppUser of the appCredential to the related item.
// Sets o.R.AppUser to related.
// Adds o to related.R.AppCredentials.
func (o *AppCredential) SetAppUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *AppUser) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"app_credentials\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"app_user_id"}),
		strmangle.WhereClause("\"", "\"", 2, appCredentialPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID, o.AppUserID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AppUserID = related.ID
	if o.R == nil {
		o.R = &appCredentialR{
			AppUser: related,
		}
	} else {
		o.R.AppUser = related
	}

	if related.R == nil {
		related.R = &appUserR{
			AppCredentials: AppCredentialSlice{o},
		}
	} else {
		related.R.AppCredentials = append(related.R.AppCredentials, o)
	}

	return nil
}

// AppCredentials retrieves all the records using an executor.
func AppCredentials(mods ...qm.QueryMod) appCredentialQuery {
	mods = append(mods, qm.From("\"app_credentials\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"app_credentials\".*"})
	}

	return appCredentialQuery{q}
}

// FindAppCredential retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAppCredential(ctx context.Context, exec boil.ContextExecutor, iD []byte, appUserID int64, selectCols ...string) (*AppCredential, error) {
	appCredentialObj := &AppCredential{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"app_credentials\" where \"id\"=$1 AND \"app_user_id\"=$2", sel,
	)

	q := queries.Raw(query, iD, appUserID)

	err := q.Bind(ctx, exec, appCredentialObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from app_credentials")
	}

	return appCredentialObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AppCredential) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no app_credentials provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(appCredentialColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	appCredentialInsertCacheMut.RLock()
	cache, cached := appCredentialInsertCache[key]
	appCredentialInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			appCredentialAllColumns,
			appCredentialColumnsWithDefault,
			appCredentialColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(appCredentialType, appCredentialMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(appCredentialType, appCredentialMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"app_credentials\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"app_credentials\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into app_credentials")
	}

	if !cached {
		appCredentialInsertCacheMut.Lock()
		appCredentialInsertCache[key] = cache
		appCredentialInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the AppCredential.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AppCredential) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	appCredentialUpdateCacheMut.RLock()
	cache, cached := appCredentialUpdateCache[key]
	appCredentialUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			appCredentialAllColumns,
			appCredentialPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return errors.New("models: unable to update app_credentials, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"app_credentials\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, appCredentialPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(appCredentialType, appCredentialMapping, append(wl, appCredentialPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	_, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update app_credentials row")
	}

	if !cached {
		appCredentialUpdateCacheMut.Lock()
		appCredentialUpdateCache[key] = cache
		appCredentialUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q appCredentialQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for app_credentials")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AppCredentialSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appCredentialPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"app_credentials\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, appCredentialPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in appCredential slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AppCredential) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no app_credentials provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(appCredentialColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	appCredentialUpsertCacheMut.RLock()
	cache, cached := appCredentialUpsertCache[key]
	appCredentialUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			appCredentialAllColumns,
			appCredentialColumnsWithDefault,
			appCredentialColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			appCredentialAllColumns,
			appCredentialPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert app_credentials, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(appCredentialPrimaryKeyColumns))
			copy(conflict, appCredentialPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"app_credentials\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(appCredentialType, appCredentialMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(appCredentialType, appCredentialMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert app_credentials")
	}

	if !cached {
		appCredentialUpsertCacheMut.Lock()
		appCredentialUpsertCache[key] = cache
		appCredentialUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single AppCredential record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AppCredential) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no AppCredential provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), appCredentialPrimaryKeyMapping)
	sql := "DELETE FROM \"app_credentials\" WHERE \"id\"=$1 AND \"app_user_id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from app_credentials")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q appCredentialQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no appCredentialQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from app_credentials")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AppCredentialSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appCredentialPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"app_credentials\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, appCredentialPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from appCredential slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AppCredential) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAppCredential(ctx, exec, o.ID, o.AppUserID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AppCredentialSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AppCredentialSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), appCredentialPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"app_credentials\".* FROM \"app_credentials\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, appCredentialPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AppCredentialSlice")
	}

	*o = slice

	return nil
}

// AppCredentialExists checks if the AppCredential row exists.
func AppCredentialExists(ctx context.Context, exec boil.ContextExecutor, iD []byte, appUserID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"app_credentials\" where \"id\"=$1 AND \"app_user_id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD, appUserID)
	}
	row := exec.QueryRowContext(ctx, sql, iD, appUserID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if app_credentials exists")
	}

	return exists, nil
}

// Exists checks if the AppCredential row exists.
func (o *AppCredential) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return AppCredentialExists(ctx, exec, o.ID, o.AppUserID)
}