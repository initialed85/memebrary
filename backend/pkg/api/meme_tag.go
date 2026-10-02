package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/netip"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/introspect"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/djangolang/pkg/stream"
	"github.com/initialed85/djangolang/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemeTag struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	MemeID       uuid.UUID  `json:"meme_id"`
	MemeIDObject *Meme      `json:"meme_id_object"`
	TagID        uuid.UUID  `json:"tag_id"`
	TagIDObject  *Tag       `json:"tag_id_object"`
}

var MemeTagTable = "meme_tag"

var MemeTagTableWithSchema = fmt.Sprintf("%s.%s", schema, MemeTagTable)

var MemeTagTableNamespaceID int32 = 1337 + 2

var (
	MemeTagTableIDColumn        = "id"
	MemeTagTableCreatedAtColumn = "created_at"
	MemeTagTableUpdatedAtColumn = "updated_at"
	MemeTagTableDeletedAtColumn = "deleted_at"
	MemeTagTableMemeIDColumn    = "meme_id"
	MemeTagTableTagIDColumn     = "tag_id"
)

var (
	MemeTagTableIDColumnWithTypeCast        = `"id" AS id`
	MemeTagTableCreatedAtColumnWithTypeCast = `"created_at" AS created_at`
	MemeTagTableUpdatedAtColumnWithTypeCast = `"updated_at" AS updated_at`
	MemeTagTableDeletedAtColumnWithTypeCast = `"deleted_at" AS deleted_at`
	MemeTagTableMemeIDColumnWithTypeCast    = `"meme_id" AS meme_id`
	MemeTagTableTagIDColumnWithTypeCast     = `"tag_id" AS tag_id`
)

var MemeTagTableColumns = []string{
	MemeTagTableIDColumn,
	MemeTagTableCreatedAtColumn,
	MemeTagTableUpdatedAtColumn,
	MemeTagTableDeletedAtColumn,
	MemeTagTableMemeIDColumn,
	MemeTagTableTagIDColumn,
}

var MemeTagTableColumnsWithTypeCasts = []string{
	MemeTagTableIDColumnWithTypeCast,
	MemeTagTableCreatedAtColumnWithTypeCast,
	MemeTagTableUpdatedAtColumnWithTypeCast,
	MemeTagTableDeletedAtColumnWithTypeCast,
	MemeTagTableMemeIDColumnWithTypeCast,
	MemeTagTableTagIDColumnWithTypeCast,
}

var MemeTagIntrospectedTable *introspect.Table

var MemeTagTableColumnLookup map[string]*introspect.Column

var (
	MemeTagTablePrimaryKeyColumn = MemeTagTableIDColumn
)

func init() {
	MemeTagIntrospectedTable = tableByName[MemeTagTable]

	/* only needed during templating */
	if MemeTagIntrospectedTable == nil {
		MemeTagIntrospectedTable = &introspect.Table{}
	}

	MemeTagTableColumnLookup = MemeTagIntrospectedTable.ColumnByName
}

type MemeTagOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type MemeTagLoadQueryParams struct {
	Depth *int `json:"depth"`
}

/*
TODO: find a way to not need this- there is a piece in the templating logic
that uses goimports but pending where the code is built, it may resolve
the packages to import to the wrong ones (causing odd failures)
these are just here to ensure we don't get unused imports
*/
var _ = []any{
	time.Time{},
	uuid.UUID{},
	pgtype.Hstore{},
	postgis.PointZ{},
	netip.Prefix{},
	errors.Is,
	sql.ErrNoRows,
}

func (m *MemeTag) GetPrimaryKeyColumn() string {
	return MemeTagTablePrimaryKeyColumn
}

func (m *MemeTag) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *MemeTag) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during MemeTagFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during MemeTagFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := MemeTagTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during MemeTagFromItem; item: %#+v",
				k, item,
			)
		}

		switch k {
		case "id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuid.UUID", temp1))
				}
			}

			m.ID = temp2

		case "created_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uucreated_at.UUID", temp1))
				}
			}

			m.CreatedAt = temp2

		case "updated_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuupdated_at.UUID", temp1))
				}
			}

			m.UpdatedAt = temp2

		case "deleted_at":
			if v == nil {
				continue
			}

			temp1, err := types.ParseTime(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(time.Time)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudeleted_at.UUID", temp1))
				}
			}

			m.DeletedAt = &temp2

		case "meme_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumeme_id.UUID", temp1))
				}
			}

			m.MemeID = temp2

		case "tag_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseUUID(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(uuid.UUID)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uutag_id.UUID", temp1))
				}
			}

			m.TagID = temp2

		}
	}

	return nil
}

func (m *MemeTag) ToItem() map[string]any {
	item := make(map[string]any)

	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("%T.ToItem() failed intermediate marshal to JSON: %s", m, err))
	}

	err = json.Unmarshal(b, &item)
	if err != nil {
		panic(fmt.Sprintf("%T.ToItem() failed intermediate unmarshal from JSON: %s", m, err))
	}

	return item
}

func (m *MemeTag) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(MemeTagTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectMemeTag(
		ctx,
		tx,
		fmt.Sprintf("%v = $1%v", m.GetPrimaryKeyColumn(), extraWhere),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return err
	}

	m.ID = o.ID
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt
	m.DeletedAt = o.DeletedAt
	m.MemeID = o.MemeID
	m.MemeIDObject = o.MemeIDObject
	m.TagID = o.TagID
	m.TagIDObject = o.TagIDObject

	return nil
}

func (m *MemeTag) GetColumnsAndValues(setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]string, []any, error) {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, MemeTagTableIDColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableIDColumn)) {
		columns = append(columns, MemeTagTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableCreatedAtColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableCreatedAtColumn) {
		columns = append(columns, MemeTagTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableUpdatedAtColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableUpdatedAtColumn) {
		columns = append(columns, MemeTagTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableDeletedAtColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableDeletedAtColumn) {
		columns = append(columns, MemeTagTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.MemeID) || slices.Contains(forceSetValuesForFields, MemeTagTableMemeIDColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableMemeIDColumn) {
		columns = append(columns, MemeTagTableMemeIDColumn)

		v, err := types.FormatUUID(m.MemeID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.MemeID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TagID) || slices.Contains(forceSetValuesForFields, MemeTagTableTagIDColumn) || isRequired(MemeTagTableColumnLookup, MemeTagTableTagIDColumn) {
		columns = append(columns, MemeTagTableTagIDColumn)

		v, err := types.FormatUUID(m.TagID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.TagID; %v", err)
		}

		values = append(values, v)
	}

	return columns, values, nil
}

func (m *MemeTag) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns, values, err := m.GetColumnsAndValues(setPrimaryKey, setZeroValues, forceSetValuesForFields...)
	if err != nil {
		return fmt.Errorf("failed to get columns and values to insert %#+v; %v", m, err)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	item, err := query.Insert(
		ctx,
		tx,
		MemeTagTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		MemeTagTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[MemeTagTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", MemeTagTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			MemeTagTableIDColumn,
			(*item)[MemeTagTableIDColumn],
			err,
		)
	}

	temp1, err := types.ParseUUID(v)
	if err != nil {
		return wrapError(err)
	}

	temp2, ok := temp1.(uuid.UUID)
	if !ok {
		return wrapError(fmt.Errorf("failed to cast to uuid.UUID"))
	}

	m.ID = temp2

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after insert; %v", err)
	}

	return nil
}

func (m *MemeTag) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableCreatedAtColumn) {
		columns = append(columns, MemeTagTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableUpdatedAtColumn) {
		columns = append(columns, MemeTagTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, MemeTagTableDeletedAtColumn) {
		columns = append(columns, MemeTagTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.MemeID) || slices.Contains(forceSetValuesForFields, MemeTagTableMemeIDColumn) {
		columns = append(columns, MemeTagTableMemeIDColumn)

		v, err := types.FormatUUID(m.MemeID)
		if err != nil {
			return fmt.Errorf("failed to handle m.MemeID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.TagID) || slices.Contains(forceSetValuesForFields, MemeTagTableTagIDColumn) {
		columns = append(columns, MemeTagTableTagIDColumn)

		v, err := types.FormatUUID(m.TagID)
		if err != nil {
			return fmt.Errorf("failed to handle m.TagID; %v", err)
		}

		values = append(values, v)
	}

	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	_, err = query.Update(
		ctx,
		tx,
		MemeTagTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", MemeTagTableIDColumn),
		MemeTagTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update %#+v; %v", m, err)
	}

	err = m.Reload(ctx, tx, slices.Contains(forceSetValuesForFields, "deleted_at"))
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}

	return nil
}

func (m *MemeTag) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(MemeTagTableColumns, "deleted_at") {
		m.DeletedAt = helpers.Ptr(time.Now().UTC())
		err := m.Update(ctx, tx, false, "deleted_at")
		if err != nil {
			return fmt.Errorf("failed to soft-delete (update) %#+v; %v", m, err)
		}
	}

	values := make([]any, 0)
	v, err := types.FormatUUID(m.ID)
	if err != nil {
		return fmt.Errorf("failed to handle m.ID; %v", err)
	}

	values = append(values, v)

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	err = query.Delete(
		ctx,
		tx,
		MemeTagTableWithSchema,
		fmt.Sprintf("%v = $$??", MemeTagTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *MemeTag) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, MemeTagTableWithSchema, timeouts...)
}

func (m *MemeTag) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, MemeTagTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *MemeTag) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, MemeTagTableNamespaceID, key, timeouts...)
}

func (m *MemeTag) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, MemeTagTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectMemeTags(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*MemeTag, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectMemeTags")

		defer func() {
			log.Printf("exited SelectMemeTags in %s", time.Since(before))
		}()
	}
	if slices.Contains(MemeTagTableColumns, "deleted_at") {
		if !strings.Contains(where, "deleted_at") {
			if where != "" {
				where += "\n    AND "
			}

			where += "deleted_at IS null"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	possiblePathValue := query.GetCurrentPathValue(ctx)
	isLoadQuery := possiblePathValue != nil && len(possiblePathValue.VisitedTableNames) > 0

	shouldLoad := query.ShouldLoad(ctx, MemeTagTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", MemeTagTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", MemeTagTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectMemeTag early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*MemeTag{}, 0, 0, 0, 0, nil
	}

	var items *[]map[string]any
	var count int64
	var totalCount int64
	var page int64
	var totalPages int64
	var err error

	useInstead, shouldSkip := query.ShouldSkip[MemeTag](ctx)
	if !shouldSkip {
		items, count, totalCount, page, totalPages, err = query.Select(
			ctx,
			tx,
			MemeTagTableColumnsWithTypeCasts,
			MemeTagTableWithSchema,
			where,
			orderBy,
			limit,
			offset,
			values...,
		)
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectMemeTags; %v", err)
		}
	} else {
		ctx = query.WithoutSkip(ctx)
		count = 1
		totalCount = 1
		page = 1
		totalPages = 1
		items = &[]map[string]any{
			nil,
		}
	}

	objects := make([]*MemeTag, 0)

	for _, item := range *items {
		var object *MemeTag

		if !shouldSkip {
			object = &MemeTag{}
			err = object.FromItem(item)
			if err != nil {
				return nil, 0, 0, 0, 0, err
			}
		} else {
			object = useInstead
		}

		if object == nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("assertion failed: object unexpectedly nil")
		}

		if !types.IsZeroUUID(object.MemeID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", MemeTable, object.MemeID), true)
			shouldLoad := query.ShouldLoad(ctx, MemeTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectMemeTags->SelectMeme for object.MemeIDObject{%s: %v}", MemeTablePrimaryKeyColumn, object.MemeID)
				}

				object.MemeIDObject, _, _, _, _, err = SelectMeme(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", MemeTablePrimaryKeyColumn),
					object.MemeID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectMemeTags->SelectMeme for object.MemeIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		if !types.IsZeroUUID(object.TagID) {
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", TagTable, object.TagID), true)
			shouldLoad := query.ShouldLoad(ctx, TagTable)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectMemeTags->SelectTag for object.TagIDObject{%s: %v}", TagTablePrimaryKeyColumn, object.TagID)
				}

				object.TagIDObject, _, _, _, _, err = SelectTag(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", TagTablePrimaryKeyColumn),
					object.TagID,
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return nil, 0, 0, 0, 0, err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectMemeTags->SelectTag for object.TagIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectMemeTag(ctx context.Context, tx pgx.Tx, where string, values ...any) (*MemeTag, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectMemeTags(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectMemeTag; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectMemeTag returned more than 1 row")
	}

	if len(objects) < 1 {
		return nil, 0, 0, 0, 0, sql.ErrNoRows
	}

	object := objects[0]

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return object, count, totalCount, page, totalPages, nil
}

func InsertMemeTags(ctx context.Context, tx pgx.Tx, objects []*MemeTag, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]*MemeTag, error) {
	var columns []string
	values := make([]any, 0)

	for i, object := range objects {
		thisColumns, thisValues, err := object.GetColumnsAndValues(setPrimaryKey, setZeroValues, forceSetValuesForFields...)
		if err != nil {
			return nil, err
		}

		if columns == nil {
			columns = thisColumns
		} else {
			if len(columns) != len(thisColumns) {
				return nil, fmt.Errorf(
					"assertion failed: call 1 of object.GetColumnsAndValues() gave %d columns but call %d gave %d columns",
					len(columns),
					i+1,
					len(thisColumns),
				)
			}
		}

		values = append(values, thisValues...)
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	items, err := query.BulkInsert(
		ctx,
		tx,
		MemeTagTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		MemeTagTableColumns,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert %d objects; %v", len(objects), err)
	}

	returnedObjects := make([]*MemeTag, 0)

	for _, item := range items {
		v := &MemeTag{}
		err = v.FromItem(*item)
		if err != nil {
			return nil, fmt.Errorf("failed %T.FromItem for %#+v; %v", *item, *item, err)
		}

		err = v.Reload(query.WithSkip(ctx, v), tx)
		if err != nil {
			return nil, fmt.Errorf("failed %T.Reload for %#+v; %v", *item, *item, err)
		}

		returnedObjects = append(returnedObjects, v)
	}

	return returnedObjects, nil
}

func handleGetMemeTags(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*MemeTag, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectMemeTags(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetMemeTag(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*MemeTag, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectMemeTag(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*MemeTag{object}, count, totalCount, page, totalPages, nil
}

func handlePostMemeTag(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*MemeTag, forceSetValuesForFieldsByObjectIndex [][]string) ([]*MemeTag, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	/* TODO: problematic- basically the bulks insert insists all rows have the same schema, which they usually should */
	forceSetValuesForFieldsByObjectIndexMaximal := make(map[string]struct{})
	for _, forceSetforceSetValuesForFields := range forceSetValuesForFieldsByObjectIndex {
		for _, field := range forceSetforceSetValuesForFields {
			forceSetValuesForFieldsByObjectIndexMaximal[field] = struct{}{}
		}
	}

	returnedObjects, err := InsertMemeTags(arguments.Ctx, tx, objects, false, false, slices.Collect(maps.Keys(forceSetValuesForFieldsByObjectIndexMaximal))...)
	if err != nil {
		err = fmt.Errorf("failed to insert %d objects; %v", len(objects), err)
		return nil, 0, 0, 0, 0, err
	}

	copy(objects, returnedObjects)

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, MemeTagTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(len(objects))
	totalCount := count
	page := int64(1)
	totalPages := page

	return objects, count, totalCount, page, totalPages, nil
}

func handlePutMemeTag(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *MemeTag) ([]*MemeTag, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, true)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, MemeTagTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*MemeTag{object}, count, totalCount, page, totalPages, nil
}

func handlePatchMemeTag(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *MemeTag, forceSetValuesForFields []string) ([]*MemeTag, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return nil, 0, 0, 0, 0, err
	}
	_ = xid

	err = object.Update(arguments.Ctx, tx, false, forceSetValuesForFields...)
	if err != nil {
		err = fmt.Errorf("failed to update %#+v; %v", object, err)
		return nil, 0, 0, 0, 0, err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, MemeTagTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return nil, 0, 0, 0, 0, err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return nil, 0, 0, 0, 0, err
	case err = <-errs:
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}
	}

	count := int64(1)
	totalCount := count
	page := int64(1)
	totalPages := page

	return []*MemeTag{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteMemeTag(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *MemeTag) error {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin DB transaction; %v", err)
		return err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	xid, err := query.GetXid(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to get xid; %v", err)
		return err
	}
	_ = xid

	err = object.Delete(arguments.Ctx, tx)
	if err != nil {
		err = fmt.Errorf("failed to delete %#+v; %v", object, err)
		return err
	}

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, MemeTagTable, xid)
		if err != nil {
			err = fmt.Errorf("failed to wait for change; %v", err)
			errs <- err
			return
		}

		errs <- nil
	}()

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		err = fmt.Errorf("failed to commit DB transaction; %v", err)
		return err
	}

	select {
	case <-arguments.Ctx.Done():
		err = fmt.Errorf("context canceled")
		return err
	case err = <-errs:
		if err != nil {
			return err
		}
	}

	return nil
}

func MutateRouterForMemeTag(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/meme-tags",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[MemeTag], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, MemeTagIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[MemeTag]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[MemeTag]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetMemeTags(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[MemeTag]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache missed; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/meme-tags/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeTagOnePathParams,
				queryParams MemeTagLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[MemeTag], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, MemeTagIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[MemeTag]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[MemeTag]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetMemeTag(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[MemeTag]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				/* TODO: it'd be nice to be able to avoid this (i.e. just marshal once, further out) */
				responseAsJSON, err := json.Marshal(response)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[MemeTag]{}, err
				}

				err = server.StoreCachedResponse(arguments.RequestHash, redisConn, responseAsJSON)
				if err != nil {
					log.Printf("warning; %v", err)
				}

				if config.Debug() {
					log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
				}

				return response, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/meme-tags",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams MemeTagLoadQueryParams,
				req []*MemeTag,
				rawReq any,
			) (server.Response[MemeTag], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[MemeTag]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[MemeTag]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range slices.Collect(maps.Keys(item)) {
						if !slices.Contains(MemeTagTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostMemeTag(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[MemeTag]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/meme-tags/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeTagOnePathParams,
				queryParams MemeTagLoadQueryParams,
				req MemeTag,
				rawReq any,
			) (server.Response[MemeTag], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[MemeTag]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutMemeTag(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[MemeTag]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/meme-tags/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeTagOnePathParams,
				queryParams MemeTagLoadQueryParams,
				req MemeTag,
				rawReq any,
			) (server.Response[MemeTag], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[MemeTag]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range slices.Collect(maps.Keys(item)) {
					if !slices.Contains(MemeTagTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchMemeTag(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[MemeTag]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[MemeTag]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    objects,
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/meme-tags/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams MemeTagOnePathParams,
				queryParams MemeTagLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &MemeTag{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteMemeTag(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			MemeTag{},
			MemeTagIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewMemeTagFromItem(item map[string]any) (any, error) {
	object := &MemeTag{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		MemeTagTable,
		MemeTag{},
		NewMemeTagFromItem,
		"/meme-tags",
		MutateRouterForMemeTag,
	)
}
func (m *MemeTag) UpdateField(ctx context.Context, tx pgx.Tx, fieldName string, value any) error {
	var columnName string
	switch fieldName {
	case "id":
		columnName = MemeTagTableIDColumn
	case "created_at":
		columnName = MemeTagTableCreatedAtColumn
	case "updated_at":
		columnName = MemeTagTableUpdatedAtColumn
	case "deleted_at":
		columnName = MemeTagTableDeletedAtColumn
	case "meme_id":
		columnName = MemeTagTableMemeIDColumn
	case "tag_id":
		columnName = MemeTagTableTagIDColumn

	default:
		return fmt.Errorf("unknown field name: %v", fieldName)
	}
	var columnValue any
	var err error
	switch columnName {
	case MemeTagTableIDColumn:
		columnValue, err = types.FormatUUID(value)
	case MemeTagTableCreatedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTagTableUpdatedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTagTableDeletedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTagTableMemeIDColumn:
		columnValue, err = types.FormatUUID(value)
	case MemeTagTableTagIDColumn:
		columnValue, err = types.FormatUUID(value)

	}
	if err != nil {
		return fmt.Errorf("failed to format value for %v; %v", columnName, err)
	}
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()
	ctx = query.WithMaxDepth(ctx, nil)
	_, err = query.Update(
		ctx,
		tx,
		MemeTagTableWithSchema,
		[]string{columnName},
		fmt.Sprintf("%v = $$??", MemeTagTableIDColumn),
		[]string{MemeTagTableIDColumn},
		columnValue,
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update field %v: %v", fieldName, err)
	}
	err = m.Reload(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}
	return nil
}
func (m *MemeTag) UpdateFields(ctx context.Context, tx pgx.Tx, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	fieldNames := make([]string, 0, len(fields))
	for fieldName := range fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	columns := make([]string, 0, len(fields))
	values := make([]any, 0, len(fields)*2)
	for _, fieldName := range fieldNames {
		value := fields[fieldName]
		var columnName string
		switch fieldName {
		case "id":
			columnName = MemeTagTableIDColumn
		case "created_at":
			columnName = MemeTagTableCreatedAtColumn
		case "updated_at":
			columnName = MemeTagTableUpdatedAtColumn
		case "deleted_at":
			columnName = MemeTagTableDeletedAtColumn
		case "meme_id":
			columnName = MemeTagTableMemeIDColumn
		case "tag_id":
			columnName = MemeTagTableTagIDColumn

		default:
			return fmt.Errorf("unknown field name: %v", fieldName)
		}
		var columnValue any
		var err error
		switch columnName {
		case MemeTagTableIDColumn:
			columnValue, err = types.FormatUUID(value)
		case MemeTagTableCreatedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTagTableUpdatedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTagTableDeletedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTagTableMemeIDColumn:
			columnValue, err = types.FormatUUID(value)
		case MemeTagTableTagIDColumn:
			columnValue, err = types.FormatUUID(value)

		}
		if err != nil {
			return fmt.Errorf("failed to format value for %v; %v", columnName, err)
		}
		columns = append(columns, columnName)
		values = append(values, columnValue)
	}
	values = append(values, m.ID)
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()
	ctx = query.WithMaxDepth(ctx, nil)
	_, err := query.Update(
		ctx,
		tx,
		MemeTagTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", MemeTagTableIDColumn),
		[]string{MemeTagTableIDColumn},
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to update fields: %v", err)
	}
	err = m.Reload(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to reload after update")
	}
	return nil
}
