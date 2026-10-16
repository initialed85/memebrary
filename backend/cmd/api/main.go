package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/config"
	internalconfig "github.com/initialed85/membrary/backend/internal/config"
	"github.com/initialed85/membrary/backend/internal/describe"
	"github.com/initialed85/membrary/backend/internal/detect"
	"github.com/initialed85/membrary/backend/internal/httpapi"
	"github.com/initialed85/membrary/backend/internal/store"
	"github.com/initialed85/membrary/backend/pkg/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = api.ThisLogger()

func addCustomHandlers(r chi.Router, db *pgxpool.Pool, _ *redis.Pool) error {
	cfg := internalconfig.Load()
	if err := os.MkdirAll(cfg.MediaDir, 0o755); err != nil {
		return fmt.Errorf("create media directory: %w", err)
	}
	dataStore := store.New(db)
	logger := slog.Default()
	generator := describe.New(cfg.AIBaseURL, cfg.AIModel, cfg.AIAPIKey, dataStore, cfg.ReprocessExisting, cfg.AIWorkers, logger)
	generator.Start(context.Background())
	worker := detect.New(cfg.DetectionBaseURL, cfg.DetectionConf, dataStore, cfg.DetectionWorkers, logger)
	worker.Start(context.Background())
	compatibilityAPI := httpapi.New(dataStore, generator, worker, cfg.MediaDir, cfg.MaxUploadSize, cfg.CORSOrigin, logger)
	compatibilityAPI.Routes(r)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'dump-config', 'dump-openapi-json', 'dump-openapi-yaml', 'reprocess-detections' or 'serve')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))
	switch command {
	case "dump-config":
		config.DumpConfig()
	case "dump-openapi-json":
		api.RunDumpOpenAPIJSON()
	case "dump-openapi-yaml":
		api.RunDumpOpenAPIYAML()
	case "reprocess-detections":
		runReprocessDetections()
	case "serve":
		api.RunServeWithEnvironment(nil, nil, addCustomHandlers)
	default:
		log.Fatalf("unknown command %q", command)
	}
}

// runReprocessDetections walks every meme once and re-runs detection for it.
func runReprocessDetections() {
	cfg := internalconfig.Load()
	if cfg.DetectionBaseURL == "" {
		log.Fatal("DETECTION_BASE_URL is required")
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.GetDSN())
	if err != nil {
		log.Fatalf("connect to postgres failed: %v", err)
	}
	defer db.Close()

	dataStore := store.New(db)
	worker := detect.New(cfg.DetectionBaseURL, cfg.DetectionConf, dataStore, 1, slog.Default())

	cursor := ""
	seen := 0
	for {
		result, err := dataStore.List(ctx, 100, cursor, "")
		if err != nil {
			log.Fatalf("list memes failed: %v", err)
		}
		for _, meme := range result.Memes {
			seen++
			if meme.Deleted {
				continue
			}
			worker.ProcessOne(ctx, meme)
		}
		if result.NextCursor == "" {
			log.Printf("reprocess detections complete: seen=%d", seen)
			return
		}
		cursor = result.NextCursor
	}
}
