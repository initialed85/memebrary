package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	SortOrder            int64    `json:"-"`
	MetadataVersion      int      `json:"-"`
}

type ListResult struct {
	Memes      []Meme `json:"memes"`
	NextCursor string `json:"next_cursor,omitempty"`
	Total      int    `json:"total"`
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
	row := s.db.QueryRow(ctx, `SELECT id, filename, original_name, mime_type, size, description,
		description_status, description_generated, sort_order, metadata_version, created_at
		FROM public.meme WHERE id = $1 AND deleted_at IS NULL`, parsed)
	meme, err := scanMeme(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Meme{}, ErrNotFound
	}
	if err != nil {
		return Meme{}, fmt.Errorf("get meme: %w", err)
	}
	meme.Tags, err = s.tags(ctx, parsed)
	if err != nil {
		return Meme{}, err
	}
	return meme, nil
}

func (s *Store) List(ctx context.Context, limit int, cursor, tag string) (ListResult, error) {
	if limit < 1 || limit > 100 {
		limit = 40
	}
	cursorOrder, cursorTime, cursorID, err := decodeCursor(cursor)
	if err != nil {
		return ListResult{}, err
	}
	where := []string{"m.deleted_at IS NULL"}
	args := make([]any, 0, 7)
	if cursor != "" {
		parsedID, parseErr := uuid.Parse(cursorID)
		if parseErr != nil {
			return ListResult{}, fmt.Errorf("invalid cursor")
		}
		where = append(where, `(m.sort_order < $1 OR (m.sort_order = $2 AND (m.created_at < $3 OR (m.created_at = $4 AND m.id < $5))))`)
		args = append(args, cursorOrder, cursorOrder, cursorTime, cursorTime, parsedID)
	}
	if strings.TrimSpace(tag) != "" {
		where = append(where, fmt.Sprintf(`EXISTS (SELECT 1 FROM public.meme_tag filter_mt JOIN public.tag filter_t ON filter_t.id = filter_mt.tag_id WHERE filter_mt.meme_id = m.id AND filter_mt.deleted_at IS NULL AND filter_t.deleted_at IS NULL AND lower(filter_t.name) LIKE $%d)`, len(args)+1))
		args = append(args, tagPattern(tag))
	}
	// Cursor args above are numbered independently from the query's static SQL.
	// Rewrite the placeholders to make the generated predicate safe when the
	// optional cursor is absent.
	whereSQL := strings.Join(where, " AND ")
	if cursor == "" {
		whereSQL = strings.ReplaceAll(whereSQL, "$1", "$1")
	}
	args = append(args, limit+1)
	query := `SELECT m.id, m.filename, m.original_name, m.mime_type, m.size, m.description,
		m.description_status, m.description_generated, m.sort_order, m.metadata_version, m.created_at,
		COALESCE(string_agg(t.name, ',' ORDER BY t.name), '')
		FROM public.meme m
		LEFT JOIN public.meme_tag mt ON mt.meme_id = m.id AND mt.deleted_at IS NULL
		LEFT JOIN public.tag t ON t.id = mt.tag_id AND t.deleted_at IS NULL
		WHERE ` + whereSQL + `
		GROUP BY m.id ORDER BY m.sort_order DESC, m.created_at DESC, m.id DESC LIMIT $` + strconv.Itoa(len(args))
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return ListResult{}, fmt.Errorf("list memes: %w", err)
	}
	defer rows.Close()
	memes := make([]Meme, 0, limit)
	for rows.Next() {
		var meme Meme
		var id uuid.UUID
		var created time.Time
		var generated int64
		var tagList string
		if err := rows.Scan(&id, &meme.Filename, &meme.OriginalName, &meme.MimeType, &meme.Size, &meme.Description,
			&meme.DescriptionStatus, &generated, &meme.SortOrder, &meme.MetadataVersion, &created, &tagList); err != nil {
			return ListResult{}, fmt.Errorf("scan meme: %w", err)
		}
		meme.ID = id.String()
		meme.CreatedAt = created.UTC().Format(time.RFC3339Nano)
		meme.Description = cleanStoredDescription(meme.Description)
		meme.DescriptionGenerated = generated != 0
		if tagList == "" {
			meme.Tags = []string{}
		} else {
			meme.Tags = strings.Split(tagList, ",")
			sort.Strings(meme.Tags)
		}
		memes = append(memes, meme)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}
	next := ""
	if len(memes) > limit {
		last := memes[limit-1]
		memes = memes[:limit]
		next = encodeCursor(last.SortOrder, last.CreatedAt, last.ID)
	}
	total, err := s.count(ctx, tag)
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{Memes: memes, NextCursor: next, Total: total}, nil
}

func (s *Store) count(ctx context.Context, tag string) (int, error) {
	var count int
	if strings.TrimSpace(tag) == "" {
		err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM public.meme WHERE deleted_at IS NULL`).Scan(&count)
		return count, err
	}
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM public.meme m WHERE m.deleted_at IS NULL AND EXISTS (
		SELECT 1 FROM public.meme_tag mt JOIN public.tag t ON t.id = mt.tag_id
		WHERE mt.meme_id = m.id AND mt.deleted_at IS NULL AND t.deleted_at IS NULL AND lower(t.name) LIKE $1)`, tagPattern(tag)).Scan(&count)
	return count, err
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
	var order int64
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(sort_order), 0) + 1 FROM public.meme WHERE deleted_at IS NULL`).Scan(&order); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO public.meme
		(id, created_at, updated_at, filename, original_name, mime_type, size, description, description_status, description_generated, sort_order, metadata_version)
		VALUES ($1, $2, $2, $3, $4, $5, $6, $7, $8, $9, $10, 0)`,
		id, created, meme.Filename, meme.OriginalName, meme.MimeType, meme.Size, meme.Description,
		defaultString(meme.DescriptionStatus, "none"), boolInt(meme.DescriptionGenerated), order)
	if err != nil {
		return fmt.Errorf("insert meme: %w", err)
	}
	if err := setTagsTx(ctx, tx, id, tags); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit meme: %w", err)
	}
	return nil
}

func (s *Store) AddTags(ctx context.Context, memeID string, tags []string) error {
	id, err := uuid.Parse(memeID)
	if err != nil {
		return ErrNotFound
	}
	if _, err := s.Get(ctx, memeID); err != nil {
		return err
	}
	return setTags(ctx, s.db, id, tags)
}

func (s *Store) RemoveTag(ctx context.Context, memeID, tag string) error {
	id, err := uuid.Parse(memeID)
	if err != nil {
		return ErrNotFound
	}
	result, err := s.db.Exec(ctx, `UPDATE public.meme_tag mt SET deleted_at = now() FROM public.tag t
		WHERE mt.meme_id = $1 AND mt.tag_id = t.id AND mt.deleted_at IS NULL AND t.deleted_at IS NULL AND t.name = $2`, id, tag)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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
	if _, err := tx.Exec(ctx, `UPDATE public.meme_tag SET deleted_at = now() WHERE meme_id = $1 AND deleted_at IS NULL`, id); err != nil {
		return err
	}
	if err := setTagsTx(ctx, tx, id, tags); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateDescription(ctx context.Context, id, description, status string, generated bool) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	result, err := s.db.Exec(ctx, `UPDATE public.meme SET description = $1, description_status = $2, description_generated = $3 WHERE id = $4 AND deleted_at IS NULL`, description, status, boolInt(generated), parsed)
	if err != nil {
		return fmt.Errorf("update description: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) UpdateGeneratedContent(ctx context.Context, id, description, status string, generated bool, tags []string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin generated content update: %w", err)
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `UPDATE public.meme SET description = $1, description_status = $2, description_generated = $3, metadata_version = 2 WHERE id = $4 AND deleted_at IS NULL`, cleanStoredDescription(description), status, boolInt(generated), parsed)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	if err := setTagsTx(ctx, tx, parsed, tags); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit generated content: %w", err)
	}
	return nil
}

func (s *Store) Pending(ctx context.Context) ([]Meme, error) {
	rows, err := s.db.Query(ctx, `SELECT id, filename, mime_type, description, description_status, description_generated, sort_order, metadata_version, created_at
		FROM public.meme WHERE deleted_at IS NULL AND description_status = 'pending' ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list pending descriptions: %w", err)
	}
	defer rows.Close()
	result := []Meme{}
	for rows.Next() {
		var id uuid.UUID
		var meme Meme
		var generated int64
		var created time.Time
		if err := rows.Scan(&id, &meme.Filename, &meme.MimeType, &meme.Description, &meme.DescriptionStatus, &generated, &meme.SortOrder, &meme.MetadataVersion, &created); err != nil {
			return nil, err
		}
		meme.ID = id.String()
		meme.CreatedAt = created.UTC().Format(time.RFC3339Nano)
		meme.Description = cleanStoredDescription(meme.Description)
		meme.DescriptionGenerated = generated != 0
		meme.Tags, err = s.tags(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, meme)
	}
	return result, rows.Err()
}

func (s *Store) AllMetadata(ctx context.Context) ([]Meme, error) {
	rows, err := s.db.Query(ctx, `SELECT id, filename, original_name, mime_type, size, description, description_status, description_generated, sort_order, metadata_version, created_at
		FROM public.meme WHERE deleted_at IS NULL ORDER BY sort_order DESC`)
	if err != nil {
		return nil, fmt.Errorf("list metadata: %w", err)
	}
	defer rows.Close()
	result := []Meme{}
	for rows.Next() {
		var id uuid.UUID
		var meme Meme
		var generated int64
		var created time.Time
		if err := rows.Scan(&id, &meme.Filename, &meme.OriginalName, &meme.MimeType, &meme.Size, &meme.Description, &meme.DescriptionStatus, &generated, &meme.SortOrder, &meme.MetadataVersion, &created); err != nil {
			return nil, err
		}
		meme.ID = id.String()
		meme.CreatedAt = created.UTC().Format(time.RFC3339Nano)
		meme.Description = cleanStoredDescription(meme.Description)
		meme.DescriptionGenerated = generated != 0
		meme.Tags, err = s.tags(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, meme)
	}
	return result, rows.Err()
}

func (s *Store) Filename(ctx context.Context, id string) (string, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return "", ErrNotFound
	}
	var filename string
	err = s.db.QueryRow(ctx, `SELECT filename FROM public.meme WHERE id = $1 AND deleted_at IS NULL`, parsed).Scan(&filename)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return filename, err
}

func (s *Store) Delete(ctx context.Context, id string) (Meme, error) {
	meme, err := s.Get(ctx, id)
	if err != nil {
		return Meme{}, err
	}
	parsed, _ := uuid.Parse(id)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Meme{}, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE public.meme_tag SET deleted_at = now() WHERE meme_id = $1 AND deleted_at IS NULL`, parsed); err != nil {
		return Meme{}, err
	}
	result, err := tx.Exec(ctx, `UPDATE public.meme SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, parsed)
	if err != nil {
		return Meme{}, fmt.Errorf("delete meme: %w", err)
	}
	if result.RowsAffected() == 0 {
		return Meme{}, ErrNotFound
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
	rows, err := tx.Query(ctx, `SELECT id FROM public.meme WHERE deleted_at IS NULL ORDER BY sort_order DESC, created_at DESC, id DESC`)
	if err != nil {
		return err
	}
	ids := []uuid.UUID{}
	for rows.Next() {
		var current uuid.UUID
		if err := rows.Scan(&current); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, current)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	from := -1
	for i, current := range ids {
		if current == parsed {
			from = i
			break
		}
	}
	if from < 0 {
		return ErrNotFound
	}
	ids = append(ids[:from], ids[from+1:]...)
	to := len(ids)
	if beforeID != "" {
		target, parseErr := uuid.Parse(beforeID)
		if parseErr != nil {
			return ErrNotFound
		}
		to = -1
		for i, current := range ids {
			if current == target {
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
		for i, current := range ids {
			if current == target {
				to = i + 1
				break
			}
		}
	}
	if to < 0 {
		return ErrNotFound
	}
	ids = append(ids, uuid.Nil)
	copy(ids[to+1:], ids[to:])
	ids[to] = parsed
	for i, current := range ids {
		if _, err := tx.Exec(ctx, `UPDATE public.meme SET sort_order = $1 WHERE id = $2`, int64(len(ids)-i), current); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) tags(ctx context.Context, id uuid.UUID) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT t.name FROM public.tag t JOIN public.meme_tag mt ON mt.tag_id = t.id
		WHERE mt.meme_id = $1 AND mt.deleted_at IS NULL AND t.deleted_at IS NULL ORDER BY t.name`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	return result, rows.Err()
}

func setTags(ctx context.Context, db interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}, id uuid.UUID, tags []string) error {
	return setTagsTx(ctx, db, id, tags)
}

// setTagsTx accepts both a pool and a transaction; both expose the same pgx
// query methods, which keeps tag writes atomic when called from worker updates.
func setTagsTx(ctx context.Context, db interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}, id uuid.UUID, tags []string) error {
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		var tagID uuid.UUID
		if err := db.QueryRow(ctx, `INSERT INTO public.tag (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET deleted_at = NULL RETURNING id`, tag).Scan(&tagID); err != nil {
			return fmt.Errorf("insert tag: %w", err)
		}
		if _, err := db.Exec(ctx, `INSERT INTO public.meme_tag (meme_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, tagID); err != nil {
			return fmt.Errorf("link tag: %w", err)
		}
	}
	return nil
}

func scanMeme(row pgx.Row) (Meme, error) {
	var id uuid.UUID
	var meme Meme
	var generated int64
	var created time.Time
	err := row.Scan(&id, &meme.Filename, &meme.OriginalName, &meme.MimeType, &meme.Size, &meme.Description,
		&meme.DescriptionStatus, &generated, &meme.SortOrder, &meme.MetadataVersion, &created)
	if err != nil {
		return Meme{}, err
	}
	meme.ID = id.String()
	meme.CreatedAt = created.UTC().Format(time.RFC3339Nano)
	meme.Description = cleanStoredDescription(meme.Description)
	meme.DescriptionGenerated = generated != 0
	return meme, nil
}

func boolInt(value bool) int {
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
