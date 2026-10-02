package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/initialed85/djangolang/pkg/config"
	"github.com/initialed85/djangolang/pkg/query"
	"github.com/initialed85/djangolang/pkg/server"
	"github.com/initialed85/membrary/backend/pkg/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = api.ThisLogger()

func addCustomHandlers(r chi.Router, db *pgxpool.Pool, redisPool *redis.Pool) error {
	postHandler, err := server.GetHTTPHandler(
		http.MethodPost,
		"/memes/{primaryKey}/do-some-custom-action",
		http.StatusCreated,
		func(
			ctx context.Context,
			pathParams api.MemeOnePathParams,
			queryParams server.EmptyQueryParams,
			req server.EmptyRequest,
			rawReq any,
		) (server.Response[api.Meme], error) {
			tx, err := db.Begin(ctx)
			if err != nil {
				return server.Response[api.Meme]{}, fmt.Errorf("failed to begin DB transaction; %v", err)
			}

			defer func() {
				_ = tx.Rollback(ctx)
			}()

			meme, _, _, _, _, err := api.SelectMeme(
				query.WithLoad(ctx, "referenced_by_meme_tag"),
				tx,
				fmt.Sprintf("%s = $$??", api.MemeTablePrimaryKeyColumn),
				pathParams.PrimaryKey,
			)
			if err != nil {
				return server.Response[api.Meme]{}, fmt.Errorf(
					"failed to get job for job name %#+v; %v",
					pathParams.PrimaryKey, err,
				)
			}

			for i, memeTag := range meme.ReferencedByMemeTagMemeIDObjects {
				memeTag.Reload(
					query.WithLoad(ctx, "tag"),
					tx,
				)

				meme.ReferencedByMemeTagMemeIDObjects[i] = memeTag
			}

			err = tx.Commit(ctx)
			if err != nil {
				return server.Response[api.Meme]{}, fmt.Errorf(
					"failed to commit DB transaction; %v", err,
				)
			}

			return server.Response[api.Meme]{
				Status:     http.StatusOK,
				Success:    true,
				Error:      nil,
				Objects:    []*api.Meme{meme},
				Count:      1,
				TotalCount: 1,
				Limit:      1,
				Offset:     0,
			}, nil
		},
	)
	if err != nil {
		return err
	}

	r.Post(postHandler.FullPath, postHandler.ServeHTTP)

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
	}
}
