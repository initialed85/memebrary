# memebrary backend

The backend uses djangolang-generated PostgreSQL models/routes plus a compatibility API for the existing frontend. PostgreSQL and Redis are required at runtime because djangolang uses PostgreSQL logical replication for its change stream and Redis for coordination.

## Local workflow

```sh
# shell 1: preserve volumes and run the local dependencies
./run-env.sh up -d

# shell 2: introspect Postgres, regenerate Go/OpenAPI/frontend types
./build.sh

# shell 3: run the API
./run-for-dev.sh
```

The API listens on `:7070` by default. The generated djangolang endpoints are under `/api`; the frontend-compatible upload/timeline/media endpoints are under `/api/custom`:

- `GET /api/custom/memes?limit=36&cursor=...&tag=cats`
- `POST /api/custom/memes` multipart form (`file`, optional `tags`)
- `POST /api/custom/memes/:id/describe`
- `PATCH /api/custom/memes/:id/order`
- `POST /api/custom/memes/:id/tags`
- `DELETE /api/custom/memes/:id/tags/:tag`
- `DELETE /api/custom/memes/:id`
- `GET /api/custom/media/:id`

`build.sh` must run with healthy Compose dependencies. It regenerates `pkg/api/` from `database/schema.yaml`, dumps `schema/openapi.json` by briefly serving the custom entrypoint, generates `frontend/src/api/api.d.ts`, formats the Svelte source, and regenerates `pkg/api_client/`.

## Configuration

Djangolang connection settings use `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_SSLMODE`, `REDIS_URL`, `DJANGOLANG_API_ROOT`, and `PORT`. The compatibility layer additionally accepts:

- `DATA_DIR` / `MEDIA_DIR` (default `./data` / `./data/media`)
- `MAX_UPLOAD_BYTES` (default 20 MiB)
- `AI_BASE_URL` or `OPENAI_BASE_URL` (empty disables vision metadata)
- `AI_MODEL` / `OPENAI_MODEL`, `AI_API_KEY` / `OPENAI_API_KEY`
- `AI_WORKERS` and one-shot `AI_REPROCESS_EXISTING`

## Container release

From this directory, `./build-tag-and-push.sh` cross-compiles the API, builds the API and frontend images for amd64, and pushes `initialed85/memebrary-backend:latest` and `initialed85/memebrary-frontend:latest`. The API image runs the SQL migrations from `database/migrations` before executing `api serve`.
