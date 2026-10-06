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
	compatibilityAPI := httpapi.New(dataStore, generator, cfg.MediaDir, cfg.MaxUploadSize, cfg.CORSOrigin, logger)
	compatibilityAPI.Routes(r)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("first argument must be command (one of 'dump-config', 'dump-openapi-json', 'dump-openapi-yaml' or 'serve')")
	}

	command := strings.TrimSpace(strings.ToLower(os.Args[1]))
	switch command {
	case "dump-config":
		config.DumpConfig()
	case "dump-openapi-json":
		api.RunDumpOpenAPIJSON()
	case "dump-openapi-yaml":
		api.RunDumpOpenAPIYAML()
	case "serve":
		api.RunServeWithEnvironment(nil, nil, addCustomHandlers)
	default:
		log.Fatalf("unknown command %q", command)
	}
}
