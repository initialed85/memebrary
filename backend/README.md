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
- `AI_MAX_TOKENS` (default 2048) — completion budget for one vision pass
- `AI_REASONING_TOKENS` (default 0) — llama.cpp per-request reasoning budget; `0` stops the model thinking before it answers, `-1` omits the field for servers that do not understand it

### Reasoning, the completion budget, and why vision used to fail per input

The classifier is a thinking model: it writes hidden reasoning before the JSON
answer, and llama.cpp strips that reasoning out of `message.content`. When the
reasoning consumed the whole completion budget the endpoint still answered 200, but
with `finish_reason: "length"` and an **empty** `content` — which surfaced as the
unhelpful `unexpected end of JSON input`. Whether a meme succeeded was decided by
how much the model felt like thinking about it, and it was not a context-length
problem: `prompt_tokens` for a 3 MB meme is ~815 against a 32768-token server
context.

Measured over the eight memes that had been failing, on the live classifier:

| reasoning | median time | median completion | usable | visible text tags (capped at 32) |
| --- | --- | --- | --- | --- |
| `AI_REASONING_TOKENS=0` | 4.9s | 164 | 8/8 | 135 |
| budget 128 | 6.1s | 218 | 8/8 | 88 |
| budget 256 | 9.8s | 354 | 8/8 | 97 |
| unlimited thinking | 22.2s | 859 | 8/8 | 147 |

Suppressing reasoning is the default because a meme description does not need
chain-of-thought: full thinking buys 147 tags against 135 for 4.5x the wall time,
and a *partial* budget is worse at text extraction than none at all — the thinking
gets cut off mid-thought and the answer degrades.

`chat_template_kwargs: {"enable_thinking": false}` is still sent for Qwen-style
servers that honour it, but it does nothing for Gemma. `reasoning_budget_tokens`
is the field that works, and llama.cpp only honours it per request from b9982; on
older builds it is silently ignored, which is why the generator also doubles
`AI_MAX_TOKENS` on truncation (ceiling 8192, at most 2 retries) rather than
failing outright.

Probe the real endpoint after a classifier change:

```sh
MEMBRARY_LIVE_AI_URL=https://image-classifier.dev.initialed85.cc/v1 \
MEMBRARY_LIVE_IMAGE=/path/to/meme.png go test ./internal/describe -run Live -v
```

## Container release

From this directory, `./build-tag-and-push.sh` cross-compiles the API, builds the API and frontend images for amd64, and pushes `initialed85/memebrary-backend:latest` and `initialed85/memebrary-frontend:latest`. The API image runs the SQL migrations from `database/migrations` before executing `api serve`.
