package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/initialed85/djangolang/pkg/helpers"
	"github.com/initialed85/djangolang/pkg/query"
	generated "github.com/initialed85/membrary/backend/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Meme struct {
	ID                   string   `json:"id"`
	Filename             string   `json:"-"`
	OriginalName         string   `json:"original_name"`
	MimeType             string   `json:"mime_type"`
	Size                 int64    `json:"size"`
	Description          string   `json:"description"`
	DescriptionStatus    string   `json:"description_status"`
	DescriptionGenerated bool     `json:"description_generated"`
	Tags                 []string `json:"tags"`
	CreatedAt            string   `json:"created_at"`
	// UpdatedAt is exposed because it is the row's change stamp: djangolang's
	// update trigger rewrites it on every mutation, so the frontend can tell a
	// fresh write response apart from an older poll delta for the same id.
	UpdatedAt string `json:"updated_at"`
	// Deleted marks a soft-deleted row in an `updated_after` delta list, so a
	// client watching for changes can drop a meme it still has on screen.
	Deleted bool `json:"deleted,omitempty"`
	// SortOrder is exposed so a poll delta can be merged into the client's
	// timeline at the position the server actually ordered it.
	SortOrder       int64 `json:"sort_order"`
	MetadataVersion int   `json:"-"` // ForceRegenerate is set only for an explicit user retry. It lets the
	// vision worker ignore existing metadata while retaining it for failure
	// recovery until a replacement result is successfully persisted.
	ForceRegenerate bool `json:"-"`
}

type ListResult struct {
	Memes      []Meme `json:"memes"`
	NextCursor string `json:"next_cursor,omitempty"`
	Total      int    `json:"total"`
	// ServerTime is the database clock at the start of an `updated_after` delta
	// query. The frontend advances its poll watermark to this value instead of
	// trusting a client clock that may be skewed behind the database.
	ServerTime string `json:"server_time,omitempty"`
	// Truncated reports that more changed rows exist behind NextCursor.
	Truncated bool `json:"truncated,omitempty"`
}

type Store struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store { return &Store{db: db} }

func (s *Store) Get(ctx context.Context, id string) (Meme, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return Meme{}, ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Meme{}, fmt.Errorf("begin get meme: %w", err)
	}
	defer tx.Rollback(ctx)
	meme, err := s.getTx(ctx, tx, parsed)
	if err != nil {
		return Meme{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Meme{}, fmt.Errorf("commit get meme: %w", err)
	}
	return meme, nil
}

func (s *Store) getTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (Meme, error) {
	object, _, _, _, _, err := generated.SelectMeme(ctx, tx, fmt.Sprintf("%s = $1", generated.MemeTablePrimaryKeyColumn), id)
	if errors.Is(err, sql.ErrNoRows) {
		return Meme{}, ErrNotFound
	}
	if err != nil {
		return Meme{}, fmt.Errorf("get meme: %w", err)
	}
	meme, err := s.fromGenerated(ctx, tx, object)
	if err != nil {
		return Meme{}, err
	}
	return meme, nil
}

func (s *Store) fromGenerated(ctx context.Context, tx pgx.Tx, object *generated.Meme) (Meme, error) {
	tags, err := s.tagsForMeme(ctx, tx, object.ID)
	if err != nil {
		return Meme{}, err
	}
	return Meme{
		ID:                   object.ID.String(),
		Filename:             object.Filename,
		OriginalName:         object.OriginalName,
		MimeType:             object.MimeType,
		Size:                 object.Size,
		Description:          cleanStoredDescription(object.Description),
		DescriptionStatus:    object.DescriptionStatus,
		DescriptionGenerated: object.DescriptionGenerated != 0,
		Tags:                 tags,
		CreatedAt:            object.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:            object.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Deleted:              object.DeletedAt != nil,
		SortOrder:            object.SortOrder,
		MetadataVersion:      int(object.MetadataVersion),
	}, nil
}

func (s *Store) tagsForMeme(ctx context.Context, tx pgx.Tx, memeID uuid.UUID) ([]string, error) {
	objects, _, _, _, _, err := generated.SelectMemeTags(
		query.WithLoad(ctx, "tag"),
		tx,
		fmt.Sprintf("%s = $1", generated.MemeTagTableMemeIDColumn),
		nil,
		nil,
		nil,
		memeID,
	)
	if err != nil {
		return nil, fmt.Errorf("list meme tags: %w", err)
	}

	result := make([]string, 0, len(objects))
	for _, object := range objects {
		if object.TagIDObject == nil {
			object, err = reloadMemeTagWithTag(ctx, tx, object)
			if err != nil {
				return nil, err
			}
		}
		if object.TagIDObject != nil {
			result = append(result, object.TagIDObject.Name)
		}
	}
	sort.Strings(result)
	return result, nil
}

func reloadMemeTagWithTag(ctx context.Context, tx pgx.Tx, object *generated.MemeTag) (*generated.MemeTag, error) {
	if err := object.Reload(query.WithLoad(ctx, "tag"), tx); err != nil {
		return nil, fmt.Errorf("load meme tag: %w", err)
	}
	return object, nil
}

// touchMeme bumps the meme row so the change is visible in that row's
// updated_at. Tag edits live on meme_tag rows, and djangolang's CDC stream
// still invalidates cached meme queries through the object graph, but the
// frontend's poll ticks watch meme.updated_at, so the parent has to record the
// touch as well.
func touchMeme(ctx context.Context, tx pgx.Tx, memeID uuid.UUID) error {
	object, _, _, _, _, err := generated.SelectMeme(ctx, tx, fmt.Sprintf("%s = $1", generated.MemeTablePrimaryKeyColumn), memeID)
	if err != nil {
		return fmt.Errorf("select meme to touch: %w", err)
	}
	if err := object.UpdateFields(ctx, tx, map[string]any{"updated_at": time.Now().UTC()}); err != nil {
		return fmt.Errorf("touch meme: %w", err)
	}
	return nil
}

func (s *Store) List(ctx context.Context, limit int, cursor, tag string) (ListResult, error) {
	if limit < 1 || limit > 100 {
		limit = 40
	}
	cursorOrder, cursorTime, cursorID, err := decodeCursor(cursor)
	if err != nil {
		return ListResult{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return ListResult{}, fmt.Errorf("begin list memes: %w", err)
	}
	defer tx.Rollback(ctx)

	whereParts := make([]string, 0, 3)
	values := make([]any, 0, 8)
	if cursor != "" {
		parsedID, parseErr := uuid.Parse(cursorID)
		if parseErr != nil {
			return ListResult{}, fmt.Errorf("invalid cursor")
		}
		whereParts = append(whereParts, "(sort_order < $1 OR (sort_order = $2 AND (created_at < $3 OR (created_at = $4 AND id < $5))))")
		values = append(values, cursorOrder, cursorOrder, cursorTime, cursorTime, parsedID)
	}

	if strings.TrimSpace(tag) != "" {
		matchingIDs, matchErr := matchingMemeIDs(ctx, tx, tag, false)
		if matchErr != nil {
			return ListResult{}, matchErr
		}
		if len(matchingIDs) == 0 {
			return ListResult{Memes: []Meme{}, Total: 0}, tx.Commit(ctx)
		}
		placeholders := make([]string, 0, len(matchingIDs))
		for _, id := range matchingIDs {
			placeholders = append(placeholders, fmt.Sprintf("id = $%d", len(values)+1))
			values = append(values, id)
		}
		whereParts = append(whereParts, "("+strings.Join(placeholders, " OR ")+")")
	}

	where := strings.Join(whereParts, " AND ")
	orderBy := "sort_order DESC, created_at DESC, id DESC"
	pageLimit := limit + 1
	objects, _, totalCount, _, _, err := generated.SelectMemes(
		ctx,
		tx,
		where,
		&orderBy,
		&pageLimit,
		helperInt(0),
		values...,
	)
	if err != nil {
		return ListResult{}, fmt.Errorf("list memes: %w", err)
	}

	memes := make([]Meme, 0, len(objects))
	for _, object := range objects {
		meme, convertErr := s.fromGenerated(ctx, tx, object)
		if convertErr != nil {
			return ListResult{}, convertErr
		}
		memes = append(memes, meme)
	}

	next := ""
	if len(memes) > limit {
		last := memes[limit-1]
		memes = memes[:limit]
		next = encodeCursor(last.SortOrder, last.CreatedAt, last.ID)
	}
	if err := tx.Commit(ctx); err != nil {
		return ListResult{}, fmt.Errorf("commit list memes: %w", err)
	}
	// djangolang intentionally uses PostgreSQL's planner row estimate for
	// totalCount. It can be stale (including on an empty table), so preserve
	// the compatibility API's exact empty result while retaining the estimate
	// for populated pages.
	if len(objects) == 0 {
		totalCount = 0
	} else if totalCount < int64(len(objects)) {
		// A stale planner estimate must never claim fewer records than the
		// rows returned in this page.
		totalCount = int64(len(objects))
	}
	return ListResult{Memes: memes, NextCursor: next, Total: int(totalCount)}, nil
}

// ListChanged returns memes whose metadata moved since `updatedAfter`, newest
// change first, including rows that were soft-deleted in the same window. The
// client's poll ticks use this as a delta feed: watching meme.updated_at costs
// a handful of rows per tick instead of re-reading the whole first page, and
// tombstones make deletions in another tab visible here.
func (s *Store) ListChanged(ctx context.Context, limit int, updatedAfter, cursor, tag string) (ListResult, error) {
	if limit < 1 || limit > 100 {
		limit = 40
	}
	since, err := time.Parse(time.RFC3339Nano, updatedAfter)
	if err != nil {
		return ListResult{}, fmt.Errorf("invalid updated_after")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return ListResult{}, fmt.Errorf("begin list changed memes: %w", err)
	}
	defer tx.Rollback(ctx)

	// The watermark the client resumes from must come from the same clock the
	// updated_at trigger writes with, so read the database time first.
	var dbNow time.Time
	if err := tx.QueryRow(ctx, "SELECT now()").Scan(&dbNow); err != nil {
		return ListResult{}, fmt.Errorf("read database time: %w", err)
	}

	whereParts := []string{"updated_at >= $1"}
	values := []any{since}
	if strings.TrimSpace(cursor) != "" {
		cursorTime, cursorOrder, cursorID, err := decodeChangedCursor(cursor)
		if err != nil {
			return ListResult{}, err
		}
		// Keyset continuation for the delta ordering (updated_at, sort_order,
		// id, all ascending). Row comparison is needed because a single reorder
		// transaction stamps every meme with the same updated_at.
		whereParts = append(whereParts, fmt.Sprintf(
			"(updated_at, %s, id) > ($%d::timestamptz, $%d::bigint, $%d::uuid)",
			generated.MemeTableSortOrderColumn, len(values)+1, len(values)+2, len(values)+3,
		))
		values = append(values, cursorTime, cursorOrder, cursorID)
	}
	if strings.TrimSpace(tag) != "" {
		matchingIDs, matchErr := matchingMemeIDs(ctx, tx, tag, true)
		if matchErr != nil {
			return ListResult{}, matchErr
		}
		if len(matchingIDs) == 0 {
			// Still report the server clock: the client advances its watermark to
			// it, and a filter matching nothing today can match tomorrow.
			return ListResult{Memes: []Meme{}, Total: 0, ServerTime: dbNow.UTC().Format(time.RFC3339Nano)}, tx.Commit(ctx)
		}
		placeholders := make([]string, 0, len(matchingIDs))
		for _, id := range matchingIDs {
			placeholders = append(placeholders, fmt.Sprintf("id = $%d", len(values)+1))
			values = append(values, id)
		}
		whereParts = append(whereParts, "("+strings.Join(placeholders, " OR ")+")")
	}

	// The delta must include rows deleted since the watermark, so opt out of the
	// generated selector's implicit `deleted_at IS null` filter.
	whereParts = append(whereParts, "(deleted_at IS null OR deleted_at IS NOT null)")
	where := strings.Join(whereParts, "\n    AND ")
	orderBy := "updated_at ASC, sort_order ASC, id ASC"
	pageLimit := limit + 1
	objects, _, _, _, _, err := generated.SelectMemes(ctx, tx, where, &orderBy, &pageLimit, helperInt(0), values...)
	if err != nil {
		return ListResult{}, fmt.Errorf("list changed memes: %w", err)
	}

	memes := make([]Meme, 0, len(objects))
	for _, object := range objects {
		meme, convertErr := s.fromGenerated(ctx, tx, object)
		if convertErr != nil {
			return ListResult{}, convertErr
		}
		memes = append(memes, meme)
	}
	// The cursor is the exact last row of the page in the delta's total ordering,
	// so a page that ends inside a tie group (a reorder transaction stamps every
	// meme it touched with one updated_at) simply continues on the next tick.
	next := ""
	truncated := len(memes) > limit
	if truncated {
		last := memes[limit-1]
		memes = memes[:limit]
		next = encodeChangedCursor(last.UpdatedAt, last.SortOrder, last.ID)
	}
	serverTime := dbNow.UTC().Format(time.RFC3339Nano)
	if err := tx.Commit(ctx); err != nil {
		return ListResult{}, fmt.Errorf("commit list changed memes: %w", err)
	}
	return ListResult{Memes: memes, NextCursor: next, ServerTime: serverTime, Truncated: truncated}, nil
}

// matchingMemeIDs resolves a tag pattern to the memes carrying it.
//
// includeDeletedLinks matters for the poll delta: a meme deleted out of band
// reaches other tabs as a tombstone, but Delete soft-deletes its meme_tag links
// too, so skipping those links would hide the very rows a filtered view has to
// drop. The plain list passes false, because SelectMemes excludes deleted memes
// anyway and an untagged-but-live row must never leak into a filtered view.
func matchingMemeIDs(ctx context.Context, tx pgx.Tx, tag string, includeDeletedLinks bool) ([]uuid.UUID, error) {
	pattern := tagPattern(tag)
	tagWhere := "lower(name) LIKE $1"
	linkWhere := fmt.Sprintf("%s = $1", generated.MemeTagTableTagIDColumn)
	if includeDeletedLinks {
		tagWhere += " AND (deleted_at IS NULL OR deleted_at IS NOT NULL)"
		linkWhere += " AND (deleted_at IS NULL OR deleted_at IS NOT NULL)"
	}
	tags, _, _, _, _, err := generated.SelectTags(ctx, tx, tagWhere, nil, nil, nil, pattern)
	if err != nil {
		return nil, fmt.Errorf("find tags: %w", err)
	}
	seen := map[uuid.UUID]struct{}{}
	result := make([]uuid.UUID, 0)
	for _, tagObject := range tags {
		links, _, _, _, _, selectErr := generated.SelectMemeTags(ctx, tx, linkWhere, nil, nil, nil, tagObject.ID)
		if selectErr != nil {
			return nil, fmt.Errorf("find tagged memes: %w", selectErr)
		}
		for _, link := range links {
			if _, ok := seen[link.MemeID]; ok {
				continue
			}
			seen[link.MemeID] = struct{}{}
			result = append(result, link.MemeID)
		}
	}
	return result, nil
}

func (s *Store) Create(ctx context.Context, meme Meme, tags []string) error {
	id := uuid.New()
	if meme.ID != "" {
		parsed, err := uuid.Parse(meme.ID)
		if err != nil {
			return fmt.Errorf("invalid meme id: %w", err)
		}
		id = parsed
	}
	created := time.Now().UTC()
	if meme.CreatedAt != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, meme.CreatedAt); err == nil {
			created = parsed
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin create meme: %w", err)
	}
	defer tx.Rollback(ctx)
	order, err := nextSortOrder(ctx, tx)
	if err != nil {
		return err
	}
	object := &generated.Meme{
		ID:                   id,
		CreatedAt:            created,
		UpdatedAt:            created,
		Filename:             meme.Filename,
		OriginalName:         meme.OriginalName,
		MimeType:             meme.MimeType,
		Size:                 meme.Size,
		Description:          meme.Description,
		DescriptionStatus:    defaultString(meme.DescriptionStatus, "none"),
		DescriptionGenerated: boolInt(meme.DescriptionGenerated),
		SortOrder:            order,
	}
	if err := object.Insert(ctx, tx, true, false); err != nil {
		return fmt.Errorf("insert meme: %w", err)
	}
	if err := setTagsTx(ctx, tx, id, tags); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func nextSortOrder(ctx context.Context, tx pgx.Tx) (int64, error) {
	orderBy := "sort_order DESC"
	objects, _, _, _, _, err := generated.SelectMemes(ctx, tx, "", &orderBy, helperInt(1), helperInt(0))
	if err != nil {
		return 0, fmt.Errorf("allocate meme sort order: %w", err)
	}
	if len(objects) == 0 {
		return 1, nil
	}
	return objects[0].SortOrder + 1, nil
}

func (s *Store) AddTags(ctx context.Context, memeID string, tags []string) error {
	id, err := uuid.Parse(memeID)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := s.getTx(ctx, tx, id); err != nil {
		return err
	}
	if err := setTagsTx(ctx, tx, id, tags); err != nil {
		return err
	}
	if err := touchMeme(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) RemoveTag(ctx context.Context, memeID, tag string) error {
	id, err := uuid.Parse(memeID)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := s.getTx(ctx, tx, id); err != nil {
		return err
	}
	tags, _, _, _, _, err := generated.SelectTags(ctx, tx, "name = $1", nil, nil, nil, tag)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return ErrNotFound
	}
	links, _, _, _, _, err := generated.SelectMemeTags(ctx, tx, "meme_id = $1 AND tag_id = $2", nil, nil, nil, id, tags[0].ID)
	if err != nil {
		return err
	}
	if len(links) == 0 {
		return ErrNotFound
	}
	links[0].DeletedAt = helpers.Ptr(time.Now().UTC())
	if err := links[0].Update(ctx, tx, false, "deleted_at"); err != nil {
		return err
	}
	if err := touchMeme(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) ReplaceTags(ctx context.Context, memeID string, tags []string) error {
	id, err := uuid.Parse(memeID)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := s.getTx(ctx, tx, id); err != nil {
		return err
	}
	links, _, _, _, _, err := generated.SelectMemeTags(ctx, tx, "meme_id = $1", nil, nil, nil, id)
	if err != nil {
		return err
	}
	for _, link := range links {
		if err := link.Delete(ctx, tx); err != nil {
			return err
		}
	}
	if err := setTagsTx(ctx, tx, id, tags); err != nil {
		return err
	}
	if err := touchMeme(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateDescription(ctx context.Context, id, description, status string, generated bool) error {
	return s.updateMemeFields(ctx, id, map[string]any{
		"description":           description,
		"description_status":    status,
		"description_generated": boolInt(generated),
	})
}

func (s *Store) updateMemeFields(ctx context.Context, id string, fields map[string]any) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	object, _, _, _, _, err := generated.SelectMeme(ctx, tx, "id = $1", parsed)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if err := object.UpdateFields(ctx, tx, fields); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateGeneratedContent(ctx context.Context, id, description, status string, generatedFlag bool, tags []string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin generated content update: %w", err)
	}
	defer tx.Rollback(ctx)
	object, _, _, _, _, err := generated.SelectMeme(ctx, tx, "id = $1", parsed)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if err := object.UpdateFields(ctx, tx, map[string]any{
		"description":           cleanStoredDescription(description),
		"description_status":    status,
		"description_generated": boolInt(generatedFlag),
		"metadata_version":      int64(2),
	}); err != nil {
		return err
	}
	if err := setTagsTx(ctx, tx, parsed, tags); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) Pending(ctx context.Context) ([]Meme, error) {
	return s.listMetadata(ctx, "description_status = $1", "pending")
}

func (s *Store) AllMetadata(ctx context.Context) ([]Meme, error) {
	return s.listMetadata(ctx, "", nil)
}

func (s *Store) listMetadata(ctx context.Context, where string, values ...any) ([]Meme, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	orderBy := "sort_order DESC"
	objects, _, _, _, _, err := generated.SelectMemes(ctx, tx, where, &orderBy, nil, nil, values...)
	if err != nil {
		return nil, err
	}
	result := make([]Meme, 0, len(objects))
	for _, object := range objects {
		meme, err := s.fromGenerated(ctx, tx, object)
		if err != nil {
			return nil, err
		}
		result = append(result, meme)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) Filename(ctx context.Context, id string) (string, error) {
	meme, err := s.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return meme.Filename, nil
}

func (s *Store) Delete(ctx context.Context, id string) (Meme, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return Meme{}, ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Meme{}, err
	}
	defer tx.Rollback(ctx)
	meme, err := s.getTx(ctx, tx, parsed)
	if err != nil {
		return Meme{}, err
	}
	links, _, _, _, _, err := generated.SelectMemeTags(ctx, tx, "meme_id = $1", nil, nil, nil, parsed)
	if err != nil {
		return Meme{}, err
	}
	for _, link := range links {
		if err := link.Delete(ctx, tx); err != nil {
			return Meme{}, err
		}
	}
	object, _, _, _, _, err := generated.SelectMeme(ctx, tx, "id = $1", parsed)
	if err != nil {
		return Meme{}, err
	}
	if err := object.Delete(ctx, tx); err != nil {
		return Meme{}, fmt.Errorf("delete meme: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return Meme{}, err
	}
	return meme, nil
}

func (s *Store) Move(ctx context.Context, id, beforeID string) error {
	return s.move(ctx, id, beforeID, "")
}

func (s *Store) MoveAfter(ctx context.Context, id, afterID string) error {
	return s.move(ctx, id, "", afterID)
}

func (s *Store) move(ctx context.Context, id, beforeID, afterID string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	orderBy := "sort_order DESC, created_at DESC, id DESC"
	objects, _, _, _, _, err := generated.SelectMemes(ctx, tx, "", &orderBy, nil, nil)
	if err != nil {
		return err
	}
	from := -1
	for i, object := range objects {
		if object.ID == parsed {
			from = i
			break
		}
	}
	if from < 0 {
		return ErrNotFound
	}
	objects = append(objects[:from], objects[from+1:]...)
	to := len(objects)
	if beforeID != "" {
		target, parseErr := uuid.Parse(beforeID)
		if parseErr != nil {
			return ErrNotFound
		}
		to = -1
		for i, object := range objects {
			if object.ID == target {
				to = i
				break
			}
		}
	} else if afterID != "" {
		target, parseErr := uuid.Parse(afterID)
		if parseErr != nil {
			return ErrNotFound
		}
		to = -1
		for i, object := range objects {
			if object.ID == target {
				to = i + 1
				break
			}
		}
	}
	if to < 0 {
		return ErrNotFound
	}
	objects = append(objects, nil)
	copy(objects[to+1:], objects[to:])
	objects[to] = &generated.Meme{ID: parsed}
	for i, object := range objects {
		if object.ID == parsed {
			// The selected object is reloaded so UpdateFields has a complete model
			// and the generated API can apply its normal update/reload path.
			object, _, _, _, _, err = generated.SelectMeme(ctx, tx, "id = $1", parsed)
			if err != nil {
				return err
			}
		}
		if err := object.UpdateFields(ctx, tx, map[string]any{"sort_order": int64(len(objects) - i)}); err != nil {
			return fmt.Errorf("update meme order: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func setTagsTx(ctx context.Context, tx pgx.Tx, memeID uuid.UUID, tags []string) error {
	for _, name := range tags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tag, err := ensureTag(ctx, tx, name)
		if err != nil {
			return err
		}
		links, _, _, _, _, err := generated.SelectMemeTags(ctx, tx, "meme_id = $1 AND tag_id = $2 AND (deleted_at IS NULL OR deleted_at IS NOT NULL)", nil, nil, nil, memeID, tag.ID)
		if err != nil {
			return fmt.Errorf("find meme tag: %w", err)
		}
		if len(links) > 0 {
			if links[0].DeletedAt != nil {
				links[0].DeletedAt = nil
				if err := links[0].Update(ctx, tx, false, "deleted_at"); err != nil {
					return fmt.Errorf("restore meme tag: %w", err)
				}
			}
			continue
		}
		link := &generated.MemeTag{MemeID: memeID, TagID: tag.ID}
		if err := link.Insert(ctx, tx, false, false); err != nil {
			return fmt.Errorf("insert meme tag: %w", err)
		}
	}
	return nil
}

func ensureTag(ctx context.Context, tx pgx.Tx, name string) (*generated.Tag, error) {
	tags, _, _, _, _, err := generated.SelectTags(ctx, tx, "name = $1 AND (deleted_at IS NULL OR deleted_at IS NOT NULL)", nil, nil, nil, name)
	if err != nil {
		return nil, fmt.Errorf("find tag: %w", err)
	}
	if len(tags) > 0 {
		tag := tags[0]
		if tag.DeletedAt != nil {
			tag.DeletedAt = nil
			if err := tag.Update(ctx, tx, false, "deleted_at"); err != nil {
				return nil, fmt.Errorf("restore tag: %w", err)
			}
		}
		return tag, nil
	}
	tag := &generated.Tag{Name: name}
	if err := tag.Insert(ctx, tx, false, false); err != nil {
		return nil, fmt.Errorf("insert tag: %w", err)
	}
	return tag, nil
}

func helperInt(value int) *int { return &value }

func boolInt(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

var fillerWords = map[string]struct{}{"a": {}, "an": {}, "and": {}, "are": {}, "as": {}, "at": {}, "be": {}, "by": {}, "for": {}, "from": {}, "has": {}, "he": {}, "in": {}, "is": {}, "it": {}, "of": {}, "on": {}, "or": {}, "she": {}, "that": {}, "the": {}, "this": {}, "to": {}, "was": {}, "we": {}, "were": {}, "what": {}, "when": {}, "where": {}, "which": {}, "who": {}, "with": {}, "you": {}, "your": {}}
var placeholderTags = map[string]struct{}{"na": {}, "n-a": {}, "none": {}, "null": {}, "unknown": {}, "placeholder": {}}

func IsFillerWord(value string) bool {
	_, ok := fillerWords[strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "#"))]
	return ok
}

func IsNoisyTag(value string) bool {
	value = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "#"))
	if value == "" || IsFillerWord(value) {
		return true
	}
	if _, ok := placeholderTags[value]; ok {
		return true
	}
	allDigits := true
	for _, r := range value {
		if !unicode.IsDigit(r) {
			allDigits = false
			break
		}
	}
	return allDigits
}

func cleanStoredDescription(value string) string {
	value = strings.TrimSpace(value)
	var payload struct {
		Description string `json:"description"`
	}
	if strings.HasPrefix(value, "{") && json.Unmarshal([]byte(value), &payload) == nil && strings.TrimSpace(payload.Description) != "" {
		return strings.TrimSpace(payload.Description)
	}
	marker := strings.Index(value, `"description"`)
	if marker >= 0 {
		colon := strings.Index(value[marker+len(`"description"`):], ":")
		if colon >= 0 {
			start := marker + len(`"description"`) + colon + 1
			for start < len(value) && (value[start] == ' ' || value[start] == '\t') {
				start++
			}
			if start < len(value) && value[start] == '"' {
				start++
				for end := start; end < len(value); end++ {
					if value[end] == '"' && (end == start || value[end-1] != '\\') {
						var description string
						if json.Unmarshal([]byte(`"`+value[start:end]+`"`), &description) == nil && strings.TrimSpace(description) != "" {
							return strings.TrimSpace(description)
						}
					}
				}
			}
		}
	}
	return value
}

func tagPattern(tag string) string {
	value := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(tag), "#"))
	return "%" + value + "%"
}

// encodeChangedCursor is a keyset cursor over the delta ordering
// (updated_at ASC, sort_order ASC, id ASC).
func encodeChangedCursor(updatedAt string, sortOrder int64, id string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%s\x00%d\x00%s", updatedAt, sortOrder, id)))
}

func decodeChangedCursor(value string) (string, int64, uuid.UUID, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", 0, uuid.Nil, fmt.Errorf("invalid cursor")
	}
	parts := strings.SplitN(string(decoded), "\x00", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", 0, uuid.Nil, fmt.Errorf("invalid cursor")
	}
	parsed, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return "", 0, uuid.Nil, fmt.Errorf("invalid cursor")
	}
	var order int64
	if _, err := fmt.Sscan(parts[1], &order); err != nil {
		return "", 0, uuid.Nil, fmt.Errorf("invalid cursor")
	}
	id, err := uuid.Parse(parts[2])
	if err != nil {
		return "", 0, uuid.Nil, fmt.Errorf("invalid cursor")
	}
	return parsed.UTC().Format(time.RFC3339Nano), order, id, nil
}

func encodeCursor(sortOrder int64, createdAt, id string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d\x00%s\x00%s", sortOrder, createdAt, id)))
}

func decodeCursor(value string) (int64, string, string, error) {
	if value == "" {
		return 0, "", "", nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid cursor")
	}
	parts := strings.SplitN(string(decoded), "\x00", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return 0, "", "", fmt.Errorf("invalid cursor")
	}
	var order int64
	if _, err := fmt.Sscan(parts[0], &order); err != nil {
		return 0, "", "", fmt.Errorf("invalid cursor")
	}
	return order, parts[1], parts[2], nil
}

var ErrNotFound = errors.New("meme not found")
