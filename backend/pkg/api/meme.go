package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
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

type Meme struct {
	ID                               uuid.UUID  `json:"id"`
	CreatedAt                        time.Time  `json:"created_at"`
	UpdatedAt                        time.Time  `json:"updated_at"`
	DeletedAt                        *time.Time `json:"deleted_at"`
	Filename                         string     `json:"filename"`
	OriginalName                     string     `json:"original_name"`
	MimeType                         string     `json:"mime_type"`
	Size                             int64      `json:"size"`
	Description                      string     `json:"description"`
	DescriptionStatus                string     `json:"description_status"`
	DescriptionGenerated             int64      `json:"description_generated"`
	SortOrder                        int64      `json:"sort_order"`
	MetadataVersion                  int64      `json:"metadata_version"`
	AiWorkerClaimedUntil             time.Time  `json:"ai_worker_claimed_until"`
	ReferencedByMemeTagMemeIDObjects []*MemeTag `json:"referenced_by_meme_tag_meme_id_objects"`
}

var MemeTable = "meme"

var MemeTableWithSchema = fmt.Sprintf("%s.%s", schema, MemeTable)

var MemeTableNamespaceID int32 = 1337 + 1

var (
	MemeTableIDColumn                   = "id"
	MemeTableCreatedAtColumn            = "created_at"
	MemeTableUpdatedAtColumn            = "updated_at"
	MemeTableDeletedAtColumn            = "deleted_at"
	MemeTableFilenameColumn             = "filename"
	MemeTableOriginalNameColumn         = "original_name"
	MemeTableMimeTypeColumn             = "mime_type"
	MemeTableSizeColumn                 = "size"
	MemeTableDescriptionColumn          = "description"
	MemeTableDescriptionStatusColumn    = "description_status"
	MemeTableDescriptionGeneratedColumn = "description_generated"
	MemeTableSortOrderColumn            = "sort_order"
	MemeTableMetadataVersionColumn      = "metadata_version"
	MemeTableAiWorkerClaimedUntilColumn = "ai_worker_claimed_until"
)

var (
	MemeTableIDColumnWithTypeCast                   = `"id" AS id`
	MemeTableCreatedAtColumnWithTypeCast            = `"created_at" AS created_at`
	MemeTableUpdatedAtColumnWithTypeCast            = `"updated_at" AS updated_at`
	MemeTableDeletedAtColumnWithTypeCast            = `"deleted_at" AS deleted_at`
	MemeTableFilenameColumnWithTypeCast             = `"filename" AS filename`
	MemeTableOriginalNameColumnWithTypeCast         = `"original_name" AS original_name`
	MemeTableMimeTypeColumnWithTypeCast             = `"mime_type" AS mime_type`
	MemeTableSizeColumnWithTypeCast                 = `"size" AS size`
	MemeTableDescriptionColumnWithTypeCast          = `"description" AS description`
	MemeTableDescriptionStatusColumnWithTypeCast    = `"description_status" AS description_status`
	MemeTableDescriptionGeneratedColumnWithTypeCast = `"description_generated" AS description_generated`
	MemeTableSortOrderColumnWithTypeCast            = `"sort_order" AS sort_order`
	MemeTableMetadataVersionColumnWithTypeCast      = `"metadata_version" AS metadata_version`
	MemeTableAiWorkerClaimedUntilColumnWithTypeCast = `"ai_worker_claimed_until" AS ai_worker_claimed_until`
)

var MemeTableColumns = []string{
	MemeTableIDColumn,
	MemeTableCreatedAtColumn,
	MemeTableUpdatedAtColumn,
	MemeTableDeletedAtColumn,
	MemeTableFilenameColumn,
	MemeTableOriginalNameColumn,
	MemeTableMimeTypeColumn,
	MemeTableSizeColumn,
	MemeTableDescriptionColumn,
	MemeTableDescriptionStatusColumn,
	MemeTableDescriptionGeneratedColumn,
	MemeTableSortOrderColumn,
	MemeTableMetadataVersionColumn,
	MemeTableAiWorkerClaimedUntilColumn,
}

var MemeTableColumnsWithTypeCasts = []string{
	MemeTableIDColumnWithTypeCast,
	MemeTableCreatedAtColumnWithTypeCast,
	MemeTableUpdatedAtColumnWithTypeCast,
	MemeTableDeletedAtColumnWithTypeCast,
	MemeTableFilenameColumnWithTypeCast,
	MemeTableOriginalNameColumnWithTypeCast,
	MemeTableMimeTypeColumnWithTypeCast,
	MemeTableSizeColumnWithTypeCast,
	MemeTableDescriptionColumnWithTypeCast,
	MemeTableDescriptionStatusColumnWithTypeCast,
	MemeTableDescriptionGeneratedColumnWithTypeCast,
	MemeTableSortOrderColumnWithTypeCast,
	MemeTableMetadataVersionColumnWithTypeCast,
	MemeTableAiWorkerClaimedUntilColumnWithTypeCast,
}

var MemeIntrospectedTable *introspect.Table

var MemeTableColumnLookup map[string]*introspect.Column

var (
	MemeTablePrimaryKeyColumn = MemeTableIDColumn
)

func init() {
	MemeIntrospectedTable = tableByName[MemeTable]

	/* only needed during templating */
	if MemeIntrospectedTable == nil {
		MemeIntrospectedTable = &introspect.Table{}
	}

	MemeTableColumnLookup = MemeIntrospectedTable.ColumnByName
}

type MemeOnePathParams struct {
	PrimaryKey uuid.UUID `json:"primaryKey"`
}

type MemeLoadQueryParams struct {
	Depth *int `json:"depth"`
}

type MemeAiWorkerClaimRequest struct {
	Until          time.Time `json:"until"`
	TimeoutSeconds float64   `json:"timeout_seconds"`
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

func (m *Meme) GetPrimaryKeyColumn() string {
	return MemeTablePrimaryKeyColumn
}

func (m *Meme) GetPrimaryKeyValue() any {
	return m.ID
}

func (m *Meme) FromItem(item map[string]any) error {
	if item == nil {
		return fmt.Errorf(
			"item unexpectedly nil during MemeFromItem",
		)
	}

	if len(item) == 0 {
		return fmt.Errorf(
			"item unexpectedly empty during MemeFromItem",
		)
	}

	wrapError := func(k string, v any, err error) error {
		return fmt.Errorf("%v: %#+v; error; %v", k, v, err)
	}

	for k, v := range item {
		_, ok := MemeTableColumnLookup[k]
		if !ok {
			return fmt.Errorf(
				"item contained unexpected key %#+v during MemeFromItem; item: %#+v",
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

		case "filename":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uufilename.UUID", temp1))
				}
			}

			m.Filename = temp2

		case "original_name":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuoriginal_name.UUID", temp1))
				}
			}

			m.OriginalName = temp2

		case "mime_type":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumime_type.UUID", temp1))
				}
			}

			m.MimeType = temp2

		case "size":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusize.UUID", temp1))
				}
			}

			m.Size = temp2

		case "description":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudescription.UUID", temp1))
				}
			}

			m.Description = temp2

		case "description_status":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudescription_status.UUID", temp1))
				}
			}

			m.DescriptionStatus = temp2

		case "description_generated":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uudescription_generated.UUID", temp1))
				}
			}

			m.DescriptionGenerated = temp2

		case "sort_order":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uusort_order.UUID", temp1))
				}
			}

			m.SortOrder = temp2

		case "metadata_version":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uumetadata_version.UUID", temp1))
				}
			}

			m.MetadataVersion = temp2

		case "ai_worker_claimed_until":
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
					return wrapError(k, v, fmt.Errorf("failed to cast %#+v to uuai_worker_claimed_until.UUID", temp1))
				}
			}

			m.AiWorkerClaimedUntil = temp2

		}
	}

	return nil
}

func (m *Meme) ToItem() map[string]any {
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

func (m *Meme) Reload(ctx context.Context, tx pgx.Tx, includeDeleteds ...bool) error {
	extraWhere := ""
	if len(includeDeleteds) > 0 && includeDeleteds[0] {
		if slices.Contains(MemeTableColumns, "deleted_at") {
			extraWhere = "\n    AND (deleted_at IS null OR deleted_at IS NOT null)"
		}
	}

	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	o, _, _, _, _, err := SelectMeme(
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
	m.Filename = o.Filename
	m.OriginalName = o.OriginalName
	m.MimeType = o.MimeType
	m.Size = o.Size
	m.Description = o.Description
	m.DescriptionStatus = o.DescriptionStatus
	m.DescriptionGenerated = o.DescriptionGenerated
	m.SortOrder = o.SortOrder
	m.MetadataVersion = o.MetadataVersion
	m.AiWorkerClaimedUntil = o.AiWorkerClaimedUntil
	m.ReferencedByMemeTagMemeIDObjects = o.ReferencedByMemeTagMemeIDObjects

	return nil
}

func (m *Meme) GetColumnsAndValues(setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]string, []any, error) {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setPrimaryKey && (setZeroValues || !types.IsZeroUUID(m.ID) || slices.Contains(forceSetValuesForFields, MemeTableIDColumn) || isRequired(MemeTableColumnLookup, MemeTableIDColumn)) {
		columns = append(columns, MemeTableIDColumn)

		v, err := types.FormatUUID(m.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.ID; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, MemeTableCreatedAtColumn) || isRequired(MemeTableColumnLookup, MemeTableCreatedAtColumn) {
		columns = append(columns, MemeTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, MemeTableUpdatedAtColumn) || isRequired(MemeTableColumnLookup, MemeTableUpdatedAtColumn) {
		columns = append(columns, MemeTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, MemeTableDeletedAtColumn) || isRequired(MemeTableColumnLookup, MemeTableDeletedAtColumn) {
		columns = append(columns, MemeTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Filename) || slices.Contains(forceSetValuesForFields, MemeTableFilenameColumn) || isRequired(MemeTableColumnLookup, MemeTableFilenameColumn) {
		columns = append(columns, MemeTableFilenameColumn)

		v, err := types.FormatString(m.Filename)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Filename; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.OriginalName) || slices.Contains(forceSetValuesForFields, MemeTableOriginalNameColumn) || isRequired(MemeTableColumnLookup, MemeTableOriginalNameColumn) {
		columns = append(columns, MemeTableOriginalNameColumn)

		v, err := types.FormatString(m.OriginalName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.OriginalName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.MimeType) || slices.Contains(forceSetValuesForFields, MemeTableMimeTypeColumn) || isRequired(MemeTableColumnLookup, MemeTableMimeTypeColumn) {
		columns = append(columns, MemeTableMimeTypeColumn)

		v, err := types.FormatString(m.MimeType)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.MimeType; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Size) || slices.Contains(forceSetValuesForFields, MemeTableSizeColumn) || isRequired(MemeTableColumnLookup, MemeTableSizeColumn) {
		columns = append(columns, MemeTableSizeColumn)

		v, err := types.FormatInt(m.Size)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Size; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Description) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionColumn) || isRequired(MemeTableColumnLookup, MemeTableDescriptionColumn) {
		columns = append(columns, MemeTableDescriptionColumn)

		v, err := types.FormatString(m.Description)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.Description; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.DescriptionStatus) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionStatusColumn) || isRequired(MemeTableColumnLookup, MemeTableDescriptionStatusColumn) {
		columns = append(columns, MemeTableDescriptionStatusColumn)

		v, err := types.FormatString(m.DescriptionStatus)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DescriptionStatus; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.DescriptionGenerated) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionGeneratedColumn) || isRequired(MemeTableColumnLookup, MemeTableDescriptionGeneratedColumn) {
		columns = append(columns, MemeTableDescriptionGeneratedColumn)

		v, err := types.FormatInt(m.DescriptionGenerated)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.DescriptionGenerated; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SortOrder) || slices.Contains(forceSetValuesForFields, MemeTableSortOrderColumn) || isRequired(MemeTableColumnLookup, MemeTableSortOrderColumn) {
		columns = append(columns, MemeTableSortOrderColumn)

		v, err := types.FormatInt(m.SortOrder)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.SortOrder; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.MetadataVersion) || slices.Contains(forceSetValuesForFields, MemeTableMetadataVersionColumn) || isRequired(MemeTableColumnLookup, MemeTableMetadataVersionColumn) {
		columns = append(columns, MemeTableMetadataVersionColumn)

		v, err := types.FormatInt(m.MetadataVersion)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.MetadataVersion; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.AiWorkerClaimedUntil) || slices.Contains(forceSetValuesForFields, MemeTableAiWorkerClaimedUntilColumn) || isRequired(MemeTableColumnLookup, MemeTableAiWorkerClaimedUntilColumn) {
		columns = append(columns, MemeTableAiWorkerClaimedUntilColumn)

		v, err := types.FormatTime(m.AiWorkerClaimedUntil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to handle m.AiWorkerClaimedUntil; %v", err)
		}

		values = append(values, v)
	}

	return columns, values, nil
}

func (m *Meme) Insert(ctx context.Context, tx pgx.Tx, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) error {
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
		MemeTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		MemeTableColumns,
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %#+v; %v", m, err)
	}
	v := (*item)[MemeTableIDColumn]

	if v == nil {
		return fmt.Errorf("failed to find %v in %#+v", MemeTableIDColumn, item)
	}

	wrapError := func(err error) error {
		return fmt.Errorf(
			"failed to treat %v: %#+v as uuid.UUID: %v",
			MemeTableIDColumn,
			(*item)[MemeTableIDColumn],
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

func (m *Meme) Update(ctx context.Context, tx pgx.Tx, setZeroValues bool, forceSetValuesForFields ...string) error {
	columns := make([]string, 0)
	values := make([]any, 0)

	if setZeroValues || !types.IsZeroTime(m.CreatedAt) || slices.Contains(forceSetValuesForFields, MemeTableCreatedAtColumn) {
		columns = append(columns, MemeTableCreatedAtColumn)

		v, err := types.FormatTime(m.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.CreatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.UpdatedAt) || slices.Contains(forceSetValuesForFields, MemeTableUpdatedAtColumn) {
		columns = append(columns, MemeTableUpdatedAtColumn)

		v, err := types.FormatTime(m.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.UpdatedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.DeletedAt) || slices.Contains(forceSetValuesForFields, MemeTableDeletedAtColumn) {
		columns = append(columns, MemeTableDeletedAtColumn)

		v, err := types.FormatTime(m.DeletedAt)
		if err != nil {
			return fmt.Errorf("failed to handle m.DeletedAt; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Filename) || slices.Contains(forceSetValuesForFields, MemeTableFilenameColumn) {
		columns = append(columns, MemeTableFilenameColumn)

		v, err := types.FormatString(m.Filename)
		if err != nil {
			return fmt.Errorf("failed to handle m.Filename; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.OriginalName) || slices.Contains(forceSetValuesForFields, MemeTableOriginalNameColumn) {
		columns = append(columns, MemeTableOriginalNameColumn)

		v, err := types.FormatString(m.OriginalName)
		if err != nil {
			return fmt.Errorf("failed to handle m.OriginalName; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.MimeType) || slices.Contains(forceSetValuesForFields, MemeTableMimeTypeColumn) {
		columns = append(columns, MemeTableMimeTypeColumn)

		v, err := types.FormatString(m.MimeType)
		if err != nil {
			return fmt.Errorf("failed to handle m.MimeType; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.Size) || slices.Contains(forceSetValuesForFields, MemeTableSizeColumn) {
		columns = append(columns, MemeTableSizeColumn)

		v, err := types.FormatInt(m.Size)
		if err != nil {
			return fmt.Errorf("failed to handle m.Size; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.Description) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionColumn) {
		columns = append(columns, MemeTableDescriptionColumn)

		v, err := types.FormatString(m.Description)
		if err != nil {
			return fmt.Errorf("failed to handle m.Description; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroString(m.DescriptionStatus) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionStatusColumn) {
		columns = append(columns, MemeTableDescriptionStatusColumn)

		v, err := types.FormatString(m.DescriptionStatus)
		if err != nil {
			return fmt.Errorf("failed to handle m.DescriptionStatus; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.DescriptionGenerated) || slices.Contains(forceSetValuesForFields, MemeTableDescriptionGeneratedColumn) {
		columns = append(columns, MemeTableDescriptionGeneratedColumn)

		v, err := types.FormatInt(m.DescriptionGenerated)
		if err != nil {
			return fmt.Errorf("failed to handle m.DescriptionGenerated; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.SortOrder) || slices.Contains(forceSetValuesForFields, MemeTableSortOrderColumn) {
		columns = append(columns, MemeTableSortOrderColumn)

		v, err := types.FormatInt(m.SortOrder)
		if err != nil {
			return fmt.Errorf("failed to handle m.SortOrder; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroInt(m.MetadataVersion) || slices.Contains(forceSetValuesForFields, MemeTableMetadataVersionColumn) {
		columns = append(columns, MemeTableMetadataVersionColumn)

		v, err := types.FormatInt(m.MetadataVersion)
		if err != nil {
			return fmt.Errorf("failed to handle m.MetadataVersion; %v", err)
		}

		values = append(values, v)
	}

	if setZeroValues || !types.IsZeroTime(m.AiWorkerClaimedUntil) || slices.Contains(forceSetValuesForFields, MemeTableAiWorkerClaimedUntilColumn) {
		columns = append(columns, MemeTableAiWorkerClaimedUntilColumn)

		v, err := types.FormatTime(m.AiWorkerClaimedUntil)
		if err != nil {
			return fmt.Errorf("failed to handle m.AiWorkerClaimedUntil; %v", err)
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
		MemeTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", MemeTableIDColumn),
		MemeTableColumns,
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

func (m *Meme) Delete(ctx context.Context, tx pgx.Tx, hardDeletes ...bool) error {
	hardDelete := false
	if len(hardDeletes) > 0 {
		hardDelete = hardDeletes[0]
	}

	if !hardDelete && slices.Contains(MemeTableColumns, "deleted_at") {
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
		MemeTableWithSchema,
		fmt.Sprintf("%v = $$??", MemeTableIDColumn),
		values...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete %#+v; %v", m, err)
	}

	_ = m.Reload(ctx, tx, true)

	return nil
}

func (m *Meme) LockTable(ctx context.Context, tx pgx.Tx, timeouts ...time.Duration) error {
	return query.LockTable(ctx, tx, MemeTableWithSchema, timeouts...)
}

func (m *Meme) LockTableWithRetries(ctx context.Context, tx pgx.Tx, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.LockTableWithRetries(ctx, tx, MemeTableWithSchema, overallTimeout, individualAttempttimeout)
}

func (m *Meme) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int32, timeouts ...time.Duration) error {
	return query.AdvisoryLock(ctx, tx, MemeTableNamespaceID, key, timeouts...)
}

func (m *Meme) AdvisoryLockWithRetries(ctx context.Context, tx pgx.Tx, key int32, overallTimeout time.Duration, individualAttempttimeout time.Duration) error {
	return query.AdvisoryLockWithRetries(ctx, tx, MemeTableNamespaceID, key, overallTimeout, individualAttempttimeout)
}

func (m *Meme) AiWorkerClaim(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration) error {
	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to claim (advisory lock): %s", err.Error())
	}

	_, _, _, _, _, err = SelectMeme(
		ctx,
		tx,
		fmt.Sprintf(
			"%s = $$?? AND (ai_worker_claimed_until IS null OR ai_worker_claimed_until < now())",
			MemeTablePrimaryKeyColumn,
		),
		m.GetPrimaryKeyValue(),
	)
	if err != nil {
		return fmt.Errorf("failed to claim (select): %s", err.Error())
	}

	m.AiWorkerClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("failed to claim (update): %s", err.Error())
	}

	return nil
}

func SelectMemes(ctx context.Context, tx pgx.Tx, where string, orderBy *string, limit *int, offset *int, values ...any) ([]*Meme, int64, int64, int64, int64, error) {
	before := time.Now()

	if config.Debug() {
		log.Printf("entered SelectMemes")

		defer func() {
			log.Printf("exited SelectMemes in %s", time.Since(before))
		}()
	}
	if slices.Contains(MemeTableColumns, "deleted_at") {
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

	shouldLoad := query.ShouldLoad(ctx, MemeTable) || query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", MemeTable))

	var ok bool
	ctx, ok = query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("%s{%v}", MemeTable, nil), !isLoadQuery)
	if !ok && !shouldLoad {
		if config.Debug() {
			log.Printf("skipping SelectMeme early (query.ShouldLoad(): %v, query.HandleQueryPathGraphCycles(): %v)", shouldLoad, ok)
		}
		return []*Meme{}, 0, 0, 0, 0, nil
	}

	var items *[]map[string]any
	var count int64
	var totalCount int64
	var page int64
	var totalPages int64
	var err error

	useInstead, shouldSkip := query.ShouldSkip[Meme](ctx)
	if !shouldSkip {
		items, count, totalCount, page, totalPages, err = query.Select(
			ctx,
			tx,
			MemeTableColumnsWithTypeCasts,
			MemeTableWithSchema,
			where,
			orderBy,
			limit,
			offset,
			values...,
		)
		if err != nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectMemes; %v", err)
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

	objects := make([]*Meme, 0)

	for _, item := range *items {
		var object *Meme

		if !shouldSkip {
			object = &Meme{}
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

		err = func() error {
			shouldLoad := query.ShouldLoad(ctx, fmt.Sprintf("referenced_by_%s", MemeTagTable))
			ctx, ok := query.HandleQueryPathGraphCycles(ctx, fmt.Sprintf("__ReferencedBy__%s{%v}", MemeTagTable, object.GetPrimaryKeyValue()), true)
			if ok || shouldLoad {
				thisBefore := time.Now()

				if config.Debug() {
					log.Printf("loading SelectMemes->SelectMemeTags for object.ReferencedByMemeTagMemeIDObjects")
				}

				object.ReferencedByMemeTagMemeIDObjects, _, _, _, _, err = SelectMemeTags(
					ctx,
					tx,
					fmt.Sprintf("%v = $1", MemeTagTableMemeIDColumn),
					nil,
					nil,
					nil,
					object.GetPrimaryKeyValue(),
				)
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
				}

				if config.Debug() {
					log.Printf("loaded SelectMemes->SelectMemeTags for object.ReferencedByMemeTagMemeIDObjects in %s", time.Since(thisBefore))
				}

			}

			return nil
		}()
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		objects = append(objects, object)
	}

	return objects, count, totalCount, page, totalPages, nil
}

func SelectMeme(ctx context.Context, tx pgx.Tx, where string, values ...any) (*Meme, int64, int64, int64, int64, error) {
	ctx, cleanup := query.WithQueryID(ctx)
	defer cleanup()

	ctx = query.WithMaxDepth(ctx, nil)

	objects, _, _, _, _, err := SelectMemes(
		ctx,
		tx,
		where,
		nil,
		helpers.Ptr(2),
		helpers.Ptr(0),
		values...,
	)
	if err != nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("failed to call SelectMeme; %v", err)
	}

	if len(objects) > 1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("attempt to call SelectMeme returned more than 1 row")
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

func InsertMemes(ctx context.Context, tx pgx.Tx, objects []*Meme, setPrimaryKey bool, setZeroValues bool, forceSetValuesForFields ...string) ([]*Meme, error) {
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
		MemeTableWithSchema,
		columns,
		nil,
		nil,
		nil,
		false,
		false,
		MemeTableColumns,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert %d objects; %v", len(objects), err)
	}

	returnedObjects := make([]*Meme, 0)

	for _, item := range items {
		v := &Meme{}
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

func AiWorkerClaimMeme(ctx context.Context, tx pgx.Tx, until time.Time, timeout time.Duration, where string, orderBy *string, values ...any) (*Meme, error) {
	m := &Meme{}

	err := m.AdvisoryLockWithRetries(ctx, tx, math.MinInt32, timeout, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if strings.TrimSpace(where) != "" {
		where += " AND\n"
	}

	where += "    (ai_worker_claimed_until IS null OR ai_worker_claimed_until < now())"

	if orderBy == nil {
		orderBy = helpers.Ptr("ai_worker_claimed_until ASC, ID ASC")
	}

	ms, _, _, _, _, err := SelectMemes(
		ctx,
		tx,
		where,
		orderBy,
		helpers.Ptr(1),
		nil,
		values...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	if len(ms) == 0 {
		return nil, nil
	}

	m = ms[0]

	m.AiWorkerClaimedUntil = until

	err = m.Update(ctx, tx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to claim: %s", err.Error())
	}

	return m, nil
}

func handleGetMemes(arguments *server.SelectManyArguments, db *pgxpool.Pool) ([]*Meme, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	objects, count, totalCount, page, totalPages, err := SelectMemes(arguments.Ctx, tx, arguments.Where, arguments.OrderBy, arguments.Limit, arguments.Offset, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return objects, count, totalCount, page, totalPages, nil
}

func handleGetMeme(arguments *server.SelectOneArguments, db *pgxpool.Pool, primaryKey uuid.UUID) ([]*Meme, int64, int64, int64, int64, error) {
	tx, err := db.Begin(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	defer func() {
		_ = tx.Rollback(arguments.Ctx)
	}()

	object, count, totalCount, page, totalPages, err := SelectMeme(arguments.Ctx, tx, arguments.Where, arguments.Values...)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = tx.Commit(arguments.Ctx)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	return []*Meme{object}, count, totalCount, page, totalPages, nil
}

func handlePostMeme(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, objects []*Meme, forceSetValuesForFieldsByObjectIndex [][]string) ([]*Meme, int64, int64, int64, int64, error) {
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

	returnedObjects, err := InsertMemes(arguments.Ctx, tx, objects, false, false, slices.Collect(maps.Keys(forceSetValuesForFieldsByObjectIndexMaximal))...)
	if err != nil {
		err = fmt.Errorf("failed to insert %d objects; %v", len(objects), err)
		return nil, 0, 0, 0, 0, err
	}

	copy(objects, returnedObjects)

	errs := make(chan error, 1)
	go func() {
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.INSERT}, MemeTable, xid)
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

func handlePutMeme(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Meme) ([]*Meme, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, MemeTable, xid)
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

	return []*Meme{object}, count, totalCount, page, totalPages, nil
}

func handlePatchMeme(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Meme, forceSetValuesForFields []string) ([]*Meme, int64, int64, int64, int64, error) {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.UPDATE, stream.SOFT_DELETE, stream.SOFT_RESTORE, stream.SOFT_UPDATE}, MemeTable, xid)
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

	return []*Meme{object}, count, totalCount, page, totalPages, nil
}

func handleDeleteMeme(arguments *server.LoadArguments, db *pgxpool.Pool, waitForChange server.WaitForChange, object *Meme) error {
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
		_, err := waitForChange(arguments.Ctx, []stream.Action{stream.DELETE, stream.SOFT_DELETE}, MemeTable, xid)
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

func MutateRouterForMeme(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool, objectMiddlewares []server.ObjectMiddleware, waitForChange server.WaitForChange) {

	func() {
		postHandlerForAiWorkerClaim, err := getHTTPHandler(
			http.MethodPost,
			"/ai-worker-claim-meme",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req MemeAiWorkerClaimRequest,
				rawReq any,
			) (server.Response[Meme], error) {
				tx, err := db.Begin(ctx)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				defer func() {
					_ = tx.Rollback(ctx)
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, MemeIntrospectedTable, nil, nil)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				object, err := AiWorkerClaimMeme(ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000), arguments.Where, arguments.OrderBy, arguments.Values...)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				count := int64(0)

				totalCount := int64(0)

				limit := int64(0)

				offset := int64(0)

				if object == nil {
					return server.Response[Meme]{
						Status:     http.StatusOK,
						Success:    true,
						Error:      nil,
						Objects:    []*Meme{},
						Count:      count,
						TotalCount: totalCount,
						Limit:      limit,
						Offset:     offset,
					}, nil
				}

				err = tx.Commit(ctx)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				return server.Response[Meme]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Meme{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}, nil
			},
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForAiWorkerClaim.FullPath, postHandlerForAiWorkerClaim.ServeHTTP)

		postHandlerForAiWorkerClaimOne, err := getHTTPHandler(
			http.MethodPost,
			"/memes/{primaryKey}/ai-worker-claim",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeOnePathParams,
				queryParams MemeLoadQueryParams,
				req MemeAiWorkerClaimRequest,
				rawReq any,
			) (server.Response[Meme], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, MemeIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				/* note: deliberately no attempt at a cache hit */

				var object *Meme
				var count int64
				var totalCount int64

				err = func() error {
					tx, err := db.Begin(arguments.Ctx)
					if err != nil {
						return err
					}

					defer func() {
						_ = tx.Rollback(arguments.Ctx)
					}()

					object, count, totalCount, _, _, err = SelectMeme(arguments.Ctx, tx, arguments.Where, arguments.Values...)
					if err != nil {
						return fmt.Errorf("failed to select object to claim: %s", err.Error())
					}

					err = object.AiWorkerClaim(arguments.Ctx, tx, req.Until, time.Millisecond*time.Duration(req.TimeoutSeconds*1000))
					if err != nil {
						return err
					}

					err = tx.Commit(arguments.Ctx)
					if err != nil {
						return err
					}

					return nil
				}()
				if err != nil {
					if config.Debug() {
						log.Printf("request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Meme]{
					Status:     http.StatusOK,
					Success:    true,
					Error:      nil,
					Objects:    []*Meme{object},
					Count:      count,
					TotalCount: totalCount,
					Limit:      limit,
					Offset:     offset,
				}

				return response, nil
			},
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandlerForAiWorkerClaimOne.FullPath, postHandlerForAiWorkerClaimOne.ServeHTTP)
	}()

	func() {
		getManyHandler, err := getHTTPHandler(
			http.MethodGet,
			"/memes",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams map[string]any,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Meme], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectManyArguments(ctx, queryParams, MemeIntrospectedTable, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Meme]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Meme]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetMemes(arguments, db)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				limit := int64(0)
				if arguments.Limit != nil {
					limit = int64(*arguments.Limit)
				}

				offset := int64(0)
				if arguments.Offset != nil {
					offset = int64(*arguments.Offset)
				}

				response := server.Response[Meme]{
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

					return server.Response[Meme]{}, err
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
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getManyHandler.FullPath, getManyHandler.ServeHTTP)
	}()

	func() {
		getOneHandler, err := getHTTPHandler(
			http.MethodGet,
			"/memes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeOnePathParams,
				queryParams MemeLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.Response[Meme], error) {
				before := time.Now()

				redisConn := redisPool.Get()
				defer func() {
					_ = redisConn.Close()
				}()

				arguments, err := server.GetSelectOneArguments(ctx, queryParams.Depth, MemeIntrospectedTable, pathParams.PrimaryKey, nil, nil)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache not yet reached; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				cachedResponseAsJSON, cacheHit, err := server.GetCachedResponseAsJSON(arguments.RequestHash, redisConn)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache failed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				if cacheHit {
					var cachedResponse server.Response[Meme]

					/* TODO: it'd be nice to be able to avoid this (i.e. just pass straight through) */
					err = json.Unmarshal(cachedResponseAsJSON, &cachedResponse)
					if err != nil {
						if config.Debug() {
							log.Printf("request cache hit but failed unmarshal; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
						}

						return server.Response[Meme]{}, err
					}

					if config.Debug() {
						log.Printf("request cache hit; request succeeded in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return cachedResponse, nil
				}

				objects, count, totalCount, _, _, err := handleGetMeme(arguments, db, pathParams.PrimaryKey)
				if err != nil {
					if config.Debug() {
						log.Printf("request cache missed; request failed in %s %s path: %#+v query: %#+v req: %#+v", time.Since(before), http.MethodGet, pathParams, queryParams, req)
					}

					return server.Response[Meme]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				response := server.Response[Meme]{
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

					return server.Response[Meme]{}, err
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
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Get(getOneHandler.FullPath, getOneHandler.ServeHTTP)
	}()

	func() {
		postHandler, err := getHTTPHandler(
			http.MethodPost,
			"/memes",
			http.StatusCreated,
			func(
				ctx context.Context,
				pathParams server.EmptyPathParams,
				queryParams MemeLoadQueryParams,
				req []*Meme,
				rawReq any,
			) (server.Response[Meme], error) {
				allRawItems, ok := rawReq.([]any)
				if !ok {
					return server.Response[Meme]{}, fmt.Errorf("failed to cast %#+v to []map[string]any", rawReq)
				}

				allItems := make([]map[string]any, 0)
				for _, rawItem := range allRawItems {
					item, ok := rawItem.(map[string]any)
					if !ok {
						return server.Response[Meme]{}, fmt.Errorf("failed to cast %#+v to map[string]any", rawItem)
					}

					allItems = append(allItems, item)
				}

				forceSetValuesForFieldsByObjectIndex := make([][]string, 0)
				for _, item := range allItems {
					forceSetValuesForFields := make([]string, 0)
					for _, possibleField := range slices.Collect(maps.Keys(item)) {
						if !slices.Contains(MemeTableColumns, possibleField) {
							continue
						}

						forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
					}
					forceSetValuesForFieldsByObjectIndex = append(forceSetValuesForFieldsByObjectIndex, forceSetValuesForFields)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				objects, count, totalCount, _, _, err := handlePostMeme(arguments, db, waitForChange, req, forceSetValuesForFieldsByObjectIndex)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Meme]{
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
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Post(postHandler.FullPath, postHandler.ServeHTTP)
	}()

	func() {
		putHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/memes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeOnePathParams,
				queryParams MemeLoadQueryParams,
				req Meme,
				rawReq any,
			) (server.Response[Meme], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Meme]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePutMeme(arguments, db, waitForChange, object)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Meme]{
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
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Put(putHandler.FullPath, putHandler.ServeHTTP)
	}()

	func() {
		patchHandler, err := getHTTPHandler(
			http.MethodPatch,
			"/memes/{primaryKey}",
			http.StatusOK,
			func(
				ctx context.Context,
				pathParams MemeOnePathParams,
				queryParams MemeLoadQueryParams,
				req Meme,
				rawReq any,
			) (server.Response[Meme], error) {
				item, ok := rawReq.(map[string]any)
				if !ok {
					return server.Response[Meme]{}, fmt.Errorf("failed to cast %#+v to map[string]any", item)
				}

				forceSetValuesForFields := make([]string, 0)
				for _, possibleField := range slices.Collect(maps.Keys(item)) {
					if !slices.Contains(MemeTableColumns, possibleField) {
						continue
					}

					forceSetValuesForFields = append(forceSetValuesForFields, possibleField)
				}

				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				object := &req
				object.ID = pathParams.PrimaryKey

				objects, count, totalCount, _, _, err := handlePatchMeme(arguments, db, waitForChange, object, forceSetValuesForFields)
				if err != nil {
					return server.Response[Meme]{}, err
				}

				limit := int64(0)

				offset := int64(0)

				return server.Response[Meme]{
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
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Patch(patchHandler.FullPath, patchHandler.ServeHTTP)
	}()

	func() {
		deleteHandler, err := getHTTPHandler(
			http.MethodDelete,
			"/memes/{primaryKey}",
			http.StatusNoContent,
			func(
				ctx context.Context,
				pathParams MemeOnePathParams,
				queryParams MemeLoadQueryParams,
				req server.EmptyRequest,
				rawReq any,
			) (server.EmptyResponse, error) {
				arguments, err := server.GetLoadArguments(ctx, queryParams.Depth)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				object := &Meme{}
				object.ID = pathParams.PrimaryKey

				err = handleDeleteMeme(arguments, db, waitForChange, object)
				if err != nil {
					return server.EmptyResponse{}, err
				}

				return server.EmptyResponse{}, nil
			},
			Meme{},
			MemeIntrospectedTable,
		)
		if err != nil {
			panic(err)
		}
		r.Delete(deleteHandler.FullPath, deleteHandler.ServeHTTP)
	}()
}

func NewMemeFromItem(item map[string]any) (any, error) {
	object := &Meme{}

	err := object.FromItem(item)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func init() {
	register(
		MemeTable,
		Meme{},
		NewMemeFromItem,
		"/memes",
		MutateRouterForMeme,
	)
}
func (m *Meme) UpdateField(ctx context.Context, tx pgx.Tx, fieldName string, value any) error {
	var columnName string
	switch fieldName {
	case "id":
		columnName = MemeTableIDColumn
	case "created_at":
		columnName = MemeTableCreatedAtColumn
	case "updated_at":
		columnName = MemeTableUpdatedAtColumn
	case "deleted_at":
		columnName = MemeTableDeletedAtColumn
	case "filename":
		columnName = MemeTableFilenameColumn
	case "original_name":
		columnName = MemeTableOriginalNameColumn
	case "mime_type":
		columnName = MemeTableMimeTypeColumn
	case "size":
		columnName = MemeTableSizeColumn
	case "description":
		columnName = MemeTableDescriptionColumn
	case "description_status":
		columnName = MemeTableDescriptionStatusColumn
	case "description_generated":
		columnName = MemeTableDescriptionGeneratedColumn
	case "sort_order":
		columnName = MemeTableSortOrderColumn
	case "metadata_version":
		columnName = MemeTableMetadataVersionColumn
	case "ai_worker_claimed_until":
		columnName = MemeTableAiWorkerClaimedUntilColumn

	default:
		return fmt.Errorf("unknown field name: %v", fieldName)
	}
	var columnValue any
	var err error
	switch columnName {
	case MemeTableIDColumn:
		columnValue, err = types.FormatUUID(value)
	case MemeTableCreatedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTableUpdatedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTableDeletedAtColumn:
		columnValue, err = types.FormatTime(value)
	case MemeTableFilenameColumn:
		columnValue, err = types.FormatString(value)
	case MemeTableOriginalNameColumn:
		columnValue, err = types.FormatString(value)
	case MemeTableMimeTypeColumn:
		columnValue, err = types.FormatString(value)
	case MemeTableSizeColumn:
		columnValue, err = types.FormatInt(value)
	case MemeTableDescriptionColumn:
		columnValue, err = types.FormatString(value)
	case MemeTableDescriptionStatusColumn:
		columnValue, err = types.FormatString(value)
	case MemeTableDescriptionGeneratedColumn:
		columnValue, err = types.FormatInt(value)
	case MemeTableSortOrderColumn:
		columnValue, err = types.FormatInt(value)
	case MemeTableMetadataVersionColumn:
		columnValue, err = types.FormatInt(value)
	case MemeTableAiWorkerClaimedUntilColumn:
		columnValue, err = types.FormatTime(value)

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
		MemeTableWithSchema,
		[]string{columnName},
		fmt.Sprintf("%v = $$??", MemeTableIDColumn),
		[]string{MemeTableIDColumn},
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
func (m *Meme) UpdateFields(ctx context.Context, tx pgx.Tx, fields map[string]any) error {
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
			columnName = MemeTableIDColumn
		case "created_at":
			columnName = MemeTableCreatedAtColumn
		case "updated_at":
			columnName = MemeTableUpdatedAtColumn
		case "deleted_at":
			columnName = MemeTableDeletedAtColumn
		case "filename":
			columnName = MemeTableFilenameColumn
		case "original_name":
			columnName = MemeTableOriginalNameColumn
		case "mime_type":
			columnName = MemeTableMimeTypeColumn
		case "size":
			columnName = MemeTableSizeColumn
		case "description":
			columnName = MemeTableDescriptionColumn
		case "description_status":
			columnName = MemeTableDescriptionStatusColumn
		case "description_generated":
			columnName = MemeTableDescriptionGeneratedColumn
		case "sort_order":
			columnName = MemeTableSortOrderColumn
		case "metadata_version":
			columnName = MemeTableMetadataVersionColumn
		case "ai_worker_claimed_until":
			columnName = MemeTableAiWorkerClaimedUntilColumn

		default:
			return fmt.Errorf("unknown field name: %v", fieldName)
		}
		var columnValue any
		var err error
		switch columnName {
		case MemeTableIDColumn:
			columnValue, err = types.FormatUUID(value)
		case MemeTableCreatedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTableUpdatedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTableDeletedAtColumn:
			columnValue, err = types.FormatTime(value)
		case MemeTableFilenameColumn:
			columnValue, err = types.FormatString(value)
		case MemeTableOriginalNameColumn:
			columnValue, err = types.FormatString(value)
		case MemeTableMimeTypeColumn:
			columnValue, err = types.FormatString(value)
		case MemeTableSizeColumn:
			columnValue, err = types.FormatInt(value)
		case MemeTableDescriptionColumn:
			columnValue, err = types.FormatString(value)
		case MemeTableDescriptionStatusColumn:
			columnValue, err = types.FormatString(value)
		case MemeTableDescriptionGeneratedColumn:
			columnValue, err = types.FormatInt(value)
		case MemeTableSortOrderColumn:
			columnValue, err = types.FormatInt(value)
		case MemeTableMetadataVersionColumn:
			columnValue, err = types.FormatInt(value)
		case MemeTableAiWorkerClaimedUntilColumn:
			columnValue, err = types.FormatTime(value)

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
		MemeTableWithSchema,
		columns,
		fmt.Sprintf("%v = $$??", MemeTableIDColumn),
		[]string{MemeTableIDColumn},
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
