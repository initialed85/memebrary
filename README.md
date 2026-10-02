# meme/brary

A small, anonymous image library with a dense timeline, drag-and-drop uploads, hashtags, and optional image descriptions generated through an OpenAI-compatible vision endpoint.

## Local development

Prerequisites: Go 1.26+, Node.js 22+, and npm.

Run the API and web app in separate terminals:

```sh
make dev-backend
make dev-frontend
```

Open <http://localhost:5173>. Vite proxies `/api`, `/media`, and `/healthz` to the Go service on port 8080. Uploaded images and the SQLite database live under `backend/data/`.

The backend uses the local llama-server by default:

```sh
AI_BASE_URL=http://192.168.137.111:8088/v1 \
AI_MODEL='unsloth/Qwen3.6-35B-A3B-MTP-GGUF:UD-Q6_K_XL' \
make dev-backend
```

Set `AI_BASE_URL=` to disable generated descriptions. `OPENAI_BASE_URL`, `OPENAI_MODEL`, and `OPENAI_API_KEY` are accepted aliases. The AI worker sends the image as a base64 `image_url` in `/chat/completions`, so the endpoint must support OpenAI-compatible vision messages. An unavailable AI service never prevents the original image from being saved; the card shows a retry action.

Or run both pieces through Docker Compose:

```sh
docker compose up --build
```

Then open <http://localhost:8081>.

## API

- `GET /healthz`
- `GET /api/memes?limit=36&cursor=...&tag=cats`
- `POST /api/memes` multipart form: `file`, optional `tags`, optional `description`
- `POST /api/memes/:id/describe` queues (or retries) AI metadata generation
- `PATCH /api/memes/:id/order` moves a meme before `{ "before_id": "..." }` (omit it to move to the end)
- `DELETE /api/memes/:id` removes a meme and its stored image
- `GET /media/:id`

Images are validated as JPEG, PNG, GIF, or WebP and stored under a UUID filename. SQLite and media are kept together under the configured data directory.

## Container images

Build and publish the independently deployable images from the repository root (after `docker login`):

```sh
docker build --platform linux/amd64 -t initialed85/memebrary-backend:latest ./backend
docker build --platform linux/amd64 -t initialed85/memebrary-frontend:latest ./frontend
docker push initialed85/memebrary-backend:latest
docker push initialed85/memebrary-frontend:latest
```

Then apply the manifests in `~/Projects/Home/home-ops/applications/memebrary` as described in that directory's README. The frontend nginx config proxies `/api` and `/media` to the Kubernetes `memebrary-backend` service. Change that service name in `frontend/nginx.conf` if the deployment naming changes.
