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

type Detection struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Label        string     `json:"label"`
	ClassID      int64      `json:"class_id"`
	Confidence   float64    `json:"confidence"`
	X1           float64    `json:"x1"`
	Y1           float64    `json:"y1"`
	X2           float64    `json:"x2"`
	Y2           float64    `json:"y2"`
	MemeID       uuid.UUID  `json:"meme_id"`
	MemeIDObject *Meme      `json:"meme_id_object"`
}

var DetectionTable = "detection"

var DetectionTableWithSchema = fmt.Sprintf("%s.%s", schema, DetectionTable)

var DetectionTableNamespaceID int32 = 1337 + 1

var (
	DetectionTableIDColumn         = "id"
	DetectionTableCreatedAtColumn  = "created_at"
	DetectionTableUpdatedAtColumn  = "updated_at"
	DetectionTableDeletedAtColumn  = "deleted_at"
	DetectionTableLabelColumn      = "label"
	DetectionTableClassIDColumn    = "class_id"
	DetectionTableConfidenceColumn = "confidence"
	DetectionTableX1Column         = "x1"
	DetectionTableY1Column         = "y1"
	DetectionTableX2Column         = "x2"
	DetectionTableY2Column         = "y2"
	DetectionTableMemeIDColumn     = "meme_id"
)

var (
	DetectionTableIDColumnWithTypeCast         = `"id" AS id`
	DetectionTableCreatedAtColumnWithTypeCast  = `"created_at" AS created_at`
	DetectionTableUpdatedAtColumnWithTypeCast  = `"updated_at" AS updated_at`
	DetectionTableDeletedAtColumnWithTypeCast  = `"deleted_at" AS deleted_at`
	DetectionTableLabelColumnWithTypeCast      = `"label" AS label`
	DetectionTableClassIDColumnWithTypeCast    = `"class_id" AS class_id`
	DetectionTableConfidenceColumnWithTypeCast = `"confidence" AS confidence`
	DetectionTableX1ColumnWithTypeCast         = `"x1" AS x1`
	DetectionTableY1ColumnWithTypeCast         = `"y1" AS y1`
	DetectionTableX2ColumnWithTypeCast         = `"x2" AS x2`
	DetectionTableY2ColumnWithTypeCast         = `"y2" AS y2`
	DetectionTableMemeIDColumnWithTypeCast     = `"meme_id" AS meme_id`
)

var DetectionTableColumns = []string{
	DetectionTableIDColumn,
	DetectionTableCreatedAtColumn,
	DetectionTableUpdatedAtColumn,
	DetectionTableDeletedAtColumn,
	DetectionTableLabelColumn,
	DetectionTableClassIDColumn,
	DetectionTableConfidenceColumn,
	DetectionTableX1Column,
	DetectionTableY1Column,
	DetectionTableX2Column,
	DetectionTableY2Column,
	DetectionTableMemeIDColumn,
}

var DetectionTableColumnsWithTypeCasts = []string{
	DetectionTableIDColumnWithTypeCast,
	DetectionTableCreatedAtColumnWithTypeCast,
	DetectionTableUpdatedAtColumnWithTypeCast,
	DetectionTableDeletedAtColumnWithTypeCast,
	DetectionTableLabelColumnWithTypeCast,
	DetectionTableClassIDColumnWithTypeCast,
	DetectionTableConfidenceColumnWithTypeCast,
	DetectionTableX1ColumnWithTypeCast,
	DetectionTableY1ColumnWithTypeCast,
	DetectionTableX2ColumnWithTypeCast,
	DetectionTableY2ColumnWithTypeCast,
	DetectionTableMemeIDColumnWithTypeCast,
}

var DetectionIntrospectedTable *introspect.Table

var DetectionTableColumnLookup map[string]*introspect.Column

var (
	DetectionTablePrimaryKeyColumn = DetectionTableIDColumn
)

func init() {
	DetectionIntrospectedTable = tableByName[DetectionTable]

	/* only needed during templating */
	if DetectionIntrospectedTable == nil {
		DetectionIntrospectedTable = &introspect.Table{}
	}

	DetectionTableColumnLookup = DetectionIntrospectedTable.ColumnByName
}

type DetectionOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type DetectionLoadQueryParams struct {
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

func (m *Detection) GetPrimaryKeyColumn() string {
	return DetectionTablePrimaryKeyColumn
}

func (m *Detection) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Detection) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during DetectionFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during DetectionFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := DetectionTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during DetectionFromItem; item: %#+v",
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

		case "label":
			if v == nil {
				continue
			}

			temp1, err := types.ParseString(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(string)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uulabel.UUID", temp1))
				}
			}

			m.Label = temp2

		case "class_id":
			if v == nil {
				continue
			}

			temp1, err := types.ParseInt(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(int64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuclass_id.UUID", temp1))
				}
			}

			m.ClassID = temp2

		case "confidence":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuconfidence.UUID", temp1))
				}
			}

			m.Confidence = temp2

		case "x1":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uux1.UUID", temp1))
				}
			}

			m.X1 = temp2

		case "y1":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuy1.UUID", temp1))
				}
			}

			m.Y1 = temp2

		case "x2":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uux2.UUID", temp1))
				}
			}

			m.X2 = temp2

		case "y2":
			if v == nil {
				continue
			}

			temp1, err := types.ParseFloat(v)
			if err != nil {
				return wrapError(k, v, err)
			}

			temp2, ok := temp1.(float64)
			if !ok {
				if temp1 != nil {
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuy2.UUID", temp1))
				}
			}

			m.Y2 = temp2

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

		}
	}

	return nil
}

func (m *Detection) ToItem() map[string]any {
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

func (m *Detection) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(DetectionTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectDetection(
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
	m.Label = o.Label
	m.ClassID = o.ClassID
	m.Confidence = o.Confidence
	m.X1 = o.X1
	m.Y1 = o.Y1
	m.X2 = o.X2
	m.Y2 = o.Y2
	m.MemeID = o.MemeID
	m.MemeIDObject = o.MemeIDObject

	return nil
}

func (m *Detection) GetColumnsAndValues(setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]string, []any, error) {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, DetectionTableIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableIDColumn)) {
		columns = append(columns, DetectionTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableCreatedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableCreatedAtColumn) {
		columns = append(columns, DetectionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableUpdatedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableUpdatedAtColumn) {
		columns = append(columns, DetectionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, DetectionTableDeletedAtColumn) || isRequired(DetectionTableColumnLookup, DetectionTableDeletedAtColumn) {
		columns = append(columns, DetectionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Label) || slices.Contains(forceSetValuesForFields, DetectionTableLabelColumn) || isRequired(DetectionTableColumnLookup, DetectionTableLabelColumn) {
		columns = append(columns, DetectionTableLabelColumn)

		v, err := types.FormatString(m.Label)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Label; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) || slices.Contains(forceSetValuesForFields, DetectionTableClassIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableClassIDColumn) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ClassID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Confidence) || slices.Contains(forceSetValuesForFields, DetectionTableConfidenceColumn) || isRequired(DetectionTableColumnLookup, DetectionTableConfidenceColumn) {
		columns = append(columns, DetectionTableConfidenceColumn)

		v, err := types.FormatFloat(m.Confidence)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Confidence; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.X1) || slices.Contains(forceSetValuesForFields, DetectionTableX1Column) || isRequired(DetectionTableColumnLookup, DetectionTableX1Column) {
		columns = append(columns, DetectionTableX1Column)

		v, err := types.FormatFloat(m.X1)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.X1; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Y1) || slices.Contains(forceSetValuesForFields, DetectionTableY1Column) || isRequired(DetectionTableColumnLookup, DetectionTableY1Column) {
		columns = append(columns, DetectionTableY1Column)

		v, err := types.FormatFloat(m.Y1)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Y1; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.X2) || slices.Contains(forceSetValuesForFields, DetectionTableX2Column) || isRequired(DetectionTableColumnLookup, DetectionTableX2Column) {
		columns = append(columns, DetectionTableX2Column)

		v, err := types.FormatFloat(m.X2)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.X2; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Y2) || slices.Contains(forceSetValuesForFields, DetectionTableY2Column) || isRequired(DetectionTableColumnLookup, DetectionTableY2Column) {
		columns = append(columns, DetectionTableY2Column)

		v, err := types.FormatFloat(m.Y2)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Y2; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.MemeID) || slices.Contains(forceSetValuesForFields, DetectionTableMemeIDColumn) || isRequired(DetectionTableColumnLookup, DetectionTableMemeIDColumn) {
		columns = append(columns, DetectionTableMemeIDColumn)

		v, err := types.FormatUUID(m.MemeID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.MemeID; %v", err)
		}

		values = append(values, v)
	}

	return columns, values, nil
}

func (m *Detection) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
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
		DetectionTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		DetectionTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[DetectionTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", DetectionTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			DetectionTableIDColumn,
			(*item)[DetectionTableIDColumn],
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

func (m *Detection) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableCreatedAtColumn) {
		columns = append(columns, DetectionTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, DetectionTableUpdatedAtColumn) {
		columns = append(columns, DetectionTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, DetectionTableDeletedAtColumn) {
		columns = append(columns, DetectionTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Label) || slices.Contains(forceSetValuesForFields, DetectionTableLabelColumn) {
		columns = append(columns, DetectionTableLabelColumn)

		v, err := types.FormatString(m.Label)
		if err != nil {
			return fmt.Errorf("failed to handle m.Label; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.ClassID) || slices.Contains(forceSetValuesForFields, DetectionTableClassIDColumn) {
		columns = append(columns, DetectionTableClassIDColumn)

		v, err := types.FormatInt(m.ClassID)
		if err != nil {
			return fmt.Errorf("failed to handle m.ClassID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Confidence) || slices.Contains(forceSetValuesForFields, DetectionTableConfidenceColumn) {
		columns = append(columns, DetectionTableConfidenceColumn)

		v, err := types.FormatFloat(m.Confidence)
		if err != nil {
			return fmt.Errorf("failed to handle m.Confidence; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.X1) || slices.Contains(forceSetValuesForFields, DetectionTableX1Column) {
		columns = append(columns, DetectionTableX1Column)

		v, err := types.FormatFloat(m.X1)
		if err != nil {
			return fmt.Errorf("failed to handle m.X1; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Y1) || slices.Contains(forceSetValuesForFields, DetectionTableY1Column) {
		columns = append(columns, DetectionTableY1Column)

		v, err := types.FormatFloat(m.Y1)
		if err != nil {
			return fmt.Errorf("failed to handle m.Y1; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.X2) || slices.Contains(forceSetValuesForFields, DetectionTableX2Column) {
		columns = append(columns, DetectionTableX2Column)

		v, err := types.FormatFloat(m.X2)
		if err != nil {
			return fmt.Errorf("failed to handle m.X2; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroFloat(m.Y2) || slices.Contains(forceSetValuesForFields, DetectionTableY2Column) {
		columns = append(columns, DetectionTableY2Column)

		v, err := types.FormatFloat(m.Y2)
		if err != nil {
			return fmt.Errorf("failed to handle m.Y2; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroUUID(m.MemeID) || slices.Contains(forceSetValuesForFields, DetectionTableMemeIDColumn) {
		columns = append(columns, DetectionTableMemeIDColumn)

		v, err := types.FormatUUID(m.MemeID)
		if err != nil {
			return fmt.Errorf("failed to handle m.MemeID; %v", err)
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
		DetectionTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", DetectionTableIDColumn),
		DetectionTableColumns,
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
func (m *Detection) UpdateField(ctx context.Context, tx pgx.Tx, fieldName string, value any) error {
	var columnName string
	switch fieldName {

	case "id":
		columnName = DetectionTableIDColumn

	case "created_at":
		columnName = DetectionTableCreatedAtColumn

	case "updated_at":
		columnName = DetectionTableUpdatedAtColumn

	case "deleted_at":
		columnName = DetectionTableDeletedAtColumn

	case "label":
		columnName = DetectionTableLabelColumn

	case "class_id":
		columnName = DetectionTableClassIDColumn

	case "confidence":
		columnName = DetectionTableConfidenceColumn

	case "x1":
		columnName = DetectionTableX1Column

	case "y1":
		columnName = DetectionTableY1Column

	case "x2":
		columnName = DetectionTableX2Column

	case "y2":
		columnName = DetectionTableY2Column

	case "meme_id":
		columnName = DetectionTableMemeIDColumn
	default:
		return fmt.Errorf("unknown field name: %v", fieldName)
	}

	var columnValue any
	var err error
	switch columnName {

	case DetectionTableIDColumn:
		columnValue, err = types.FormatUUID(value)

	case DetectionTableCreatedAtColumn:
		columnValue, err = types.FormatTime(value)

	case DetectionTableUpdatedAtColumn:
		columnValue, err = types.FormatTime(value)

	case DetectionTableDeletedAtColumn:
		columnValue, err = types.FormatTime(value)

	case DetectionTableLabelColumn:
		columnValue, err = types.FormatString(value)

	case DetectionTableClassIDColumn:
		columnValue, err = types.FormatInt(value)

	case DetectionTableConfidenceColumn:
		columnValue, err = types.FormatFloat(value)

	case DetectionTableX1Column:
		columnValue, err = types.FormatFloat(value)

	case DetectionTableY1Column:
		columnValue, err = types.FormatFloat(value)

	case DetectionTableX2Column:
		columnValue, err = types.FormatFloat(value)

	case DetectionTableY2Column:
		columnValue, err = types.FormatFloat(value)

	case DetectionTableMemeIDColumn:
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
		DetectionTableWithSchema,
		[]string{columnName},
		fmt.Sprintf("%v = $$??", DetectionTablePrimaryKeyColumn),
		[]string{DetectionTablePrimaryKeyColumn},
		columnValue,
		m.GetPrimaryKeyValue(),
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

func (m *Detection) UpdateFields(ctx context.Context, tx pgx.Tx, fields map[string]any) error {
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
			columnName = DetectionTableIDColumn

		case "created_at":
			columnName = DetectionTableCreatedAtColumn

		case "updated_at":
			columnName = DetectionTableUpdatedAtColumn

		case "deleted_at":
			columnName = DetectionTableDeletedAtColumn

		case "label":
			columnName = DetectionTableLabelColumn

		case "class_id":
			columnName = DetectionTableClassIDColumn

		case "confidence":
			columnName = DetectionTableConfidenceColumn

		case "x1":
			columnName = DetectionTableX1Column

		case "y1":
			columnName = DetectionTableY1Column

		case "x2":
			columnName = DetectionTableX2Column

		case "y2":
			columnName = DetectionTableY2Column

		case "meme_id":
			columnName = DetectionTableMemeIDColumn
		default:
			return fmt.Errorf("unknown field name: %v", fieldName)
		}

		var columnValue any
		var err error
		switch columnName {

		case DetectionTableIDColumn:
			columnValue, err = types.FormatUUID(value)

		case DetectionTableCreatedAtColumn:
			columnValue, err = types.FormatTime(value)

		case DetectionTableUpdatedAtColumn:
			columnValue, err = types.FormatTime(value)

		case DetectionTableDeletedAtColumn:
			columnValue, err = types.FormatTime(value)

		case DetectionTableLabelColumn:
			columnValue, err = types.FormatString(value)

		case DetectionTableClassIDColumn:
			columnValue, err = types.FormatInt(value)

		case DetectionTableConfidenceColumn:
			columnValue, err = types.FormatFloat(value)

		case DetectionTableX1Column:
			columnValue, err = types.FormatFloat(value)

		case DetectionTableY1Column:
			columnValue, err = types.FormatFloat(value)

		case DetectionTableX2Column:
			columnValue, err = types.FormatFloat(value)

		case DetectionTableY2Column:
			columnValue, err = types.FormatFloat(value)

		case DetectionTableMemeIDColumn:
			columnValue, err = types.FormatUUID(value)
		}
		if err != nil {
			return fmt.Errorf("failed to format value for %v; %v", columnName, err)
		}

		columns = append(columns, columnName)
		values = append(values, columnValue)
	}
	values = append(values, m.GetPrimaryKeyValue())

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	_, err := query.Update(
		ctx,
		tx,
		DetectionTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", DetectionTablePrimaryKeyColumn),
		[]string{DetectionTablePrimaryKeyColumn},
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

func (m *Detection) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(DetectionTableColumns, "deleted_at") {
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
		DetectionTableWithSchema,
		fmt.Sprintf("%v = $$??", DetectionTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Detection) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, DetectionTableWithSchema, timeouts...)
}

func (m *Detection) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, DetectionTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *Detection) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, DetectionTableNamespaceID, key, timeouts...)
}

func (m *Detection) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, DetectionTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func SelectDetections(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Detection, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectDetections")

		defer func() {
			log.Printf("exited SelectDetections in %s", time.Since(before))
		}()
	}
	if slices.Contains(DetectionTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, DetectionTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", DetectionTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", DetectionTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectDetection early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Detection{}, 0, 0, 0, 0, nil
	}

	var items *[]map[string]any
	var count int64
	var totalCount int64
	var page int64
	var totalPages int64
	var err error

	useInstead, shouldSkip := query.ShouldSkip[Detection](ctx)
	if !shouldSkip {
		items, count, totalCount, page, totalPages, err = query.Select(
			ctx,
			tx,
			DetectionTableColumnsWithTypeCasts,
			DetectionTableWithSchema,
			where,
			orderBy,
			limit,
			offset,
			values...,
		)
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectDetections; %v", err)
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

	objects := make([]*Detection, 0)

	for _, item := range *items {
		var object *Detection

		if !shouldSkip {
			object = &Detection{}
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
					log.Printf("loading SelectDetections->SelectMeme for object.MemeIDObject{%s: %v}", MemeTablePrimaryKeyColumn, object.MemeID)
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
					log.Printf("loaded SelectDetections->SelectMeme for object.MemeIDObject in %s", time.Since(thisBefore))
				}
			}
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectDetection(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Detection, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectDetections(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectDetection; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectDetection returned more than 1 row")
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

func InsertDetections(ctx context.Context, tx pgx.Tx, objects []*Detection, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]*Detection, error) {
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
		DetectionTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		DetectionTableColumns,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert %d objects; %v", len(objects), err)
	}

	returnedObjects := make([]*Detection, 0)

	for _, item := range items {
		v := &Detection{}
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

func handleGetDetections(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectDetections(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetDetection(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Detection, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectDetection(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handlePostDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Detection, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Detection, int64, int64, int64, int64, error) {
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

	returnedObjects, err := InsertDetections(arguments.Ctx, tx, objects, false, false, slices.Collect(maps.Keys(forceSetValuesForFieldsByObjectIndexMaximal))...)
	if err != nil {
		err = fmt.Errorf("failed to insert %d objects; %v", len(objects), err)
		return nil, 0, 0, 0, 0, err
	}

	copy(objects, returnedObjects)

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, DetectionTable, xid)
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

func handlePutDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection) ([]*Detection, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, DetectionTable, xid)
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

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handlePatchDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection, forceSetValuesForFields []string) ([]*Detection, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, DetectionTable, xid)
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

	return []*Detection{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteDetection(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Detection) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, DetectionTable, xid)
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

func MutateRouterForDetection(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/detections",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Detection], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, DetectionIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Detection]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Detection]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetDetections(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Detection]{
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

					return server.Response[Detection]{}, err
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
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Detection], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, DetectionIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Detection]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Detection]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetDetection(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Detection]{
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

					return server.Response[Detection]{}, err
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
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/detections",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams DetectionLoadQueryParams,
				req []*Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range slices.Collect(maps.Keys(item)) {
						if !slices.Contains(DetectionTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostDetection(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
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
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutDetection(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
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
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/detections/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req Detection,
				rawReq any,
			) (server.Response[Detection], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Detection]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range slices.Collect(maps.Keys(item)) {
					if !slices.Contains(DetectionTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchDetection(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Detection]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Detection]{
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
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/detections/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams DetectionOnePathParams,
				queryParams DetectionLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Detection{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteDetection(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Detection{},
			DetectionIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewDetectionFromItem(item map[string]any) (any, error) {
	object := &Detection{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		DetectionTable,
		Detection{},
		NewDetectionFromItem,
		"/detections",
		MutateRouterForDetection,
	)
}
