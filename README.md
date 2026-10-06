# meme/brary

A small anonymous image library with a dense Svelte timeline, drag-and-drop uploads, hashtags, and optional vision-generated metadata.

## Development

Prerequisites: Go 1.27+, Docker Compose, Node.js 22+, Java, and `openapi-generator-cli`.

From `backend/`, start the PostgreSQL/Redis generation environment and keep it running:

```sh
./run-env.sh up -d
./build.sh
```

`build.sh` introspects PostgreSQL and regenerates `backend/pkg/api`, dumps the OpenAPI schema, generates `frontend/src/api/api.d.ts`, formats the frontend, and regenerates the Go client. It expects `post-migrate` to have exited successfully. Run the API and frontend in separate terminals:

```sh
# backend/
./run-for-dev.sh

# repository root/
cd frontend && npm run dev -- --host
```

Open <http://localhost:5173>. The frontend uses the generated OpenAPI TypeScript client (`openapi-fetch`) for read-only meme/tag data under `/api`, following the Camry pattern. Compatibility routes under `/api/custom` remain for multipart uploads, ordering, AI actions, and media side effects.

Set `AI_BASE_URL=` to disable vision metadata. `OPENAI_BASE_URL`, `AI_MODEL`, `OPENAI_MODEL`, and the corresponding API key aliases are also accepted.

## Deployment

`backend/build-tag-and-push.sh` builds and pushes the independently deployable images:

```sh
cd backend
./build-tag-and-push.sh
```

The script publishes `initialed85/memebrary-backend:latest` and `initialed85/memebrary-frontend:latest` (override with `BACKEND_IMAGE` and `FRONTEND_IMAGE`). The API image runs PostgreSQL migrations before starting djangolang and stores media under `MEDIA_DIR`.

The manifests for the dev deployment live in `~/Projects/Home/home-ops/applications/memebrary-dev`. Apply them with `kubectl --context home-dev` as described there, then verify both rollouts and the public health endpoint.
