# memebrary frontend

Small Svelte/Vite UI for the anonymous meme library. In development, Vite proxies the generated and compatibility `/api` routes to the djangolang service at `http://localhost:7070`.

```sh
npm install
npm run dev
npm run check
npm run build
```

The frontend uses the generated OpenAPI TypeScript definitions with `openapi-fetch`, following the Camry pattern. Read-only timeline data comes from `/api/memes`, `/api/meme-tags`, and `/api/tags`; compatibility routes remain for multipart uploads, ordering, AI actions, and media side effects. The production image is nginx and serves the compiled app. Its `/api` location proxies to the in-cluster `memebrary-backend` Service; compatibility media is served at `/api/custom/media`.

The generated client deliberately loads tags from the direct `/api/meme-tags?tag__load=` endpoint. Asking `/api/memes` for `depth=3` also walks `meme -> meme_tag -> meme`, while the generated query schema cannot express a second-hop load from the meme endpoint without that extra traversal.
