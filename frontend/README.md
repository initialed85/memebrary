# memebrary frontend

Small Svelte/Vite UI for the anonymous meme library. In development, Vite proxies the generated and compatibility `/api` routes to the djangolang service at `http://localhost:7070`.

```sh
npm install
npm run dev
npm run check
npm run build
```

The frontend uses the generated OpenAPI TypeScript definitions with `openapi-fetch`, following the Camry pattern. Read-only timeline data comes from `/api/memes?depth=3` and `/api/tags`/`/api/meme-tags` for tag filtering; compatibility routes remain for multipart uploads, ordering, AI actions, and media side effects. The production image is nginx and serves the compiled app. Its `/api` location proxies to the in-cluster `memebrary-backend` Service; compatibility media is served at `/api/custom/media`.

Cycle-aware djangolang loading now makes `depth=3` safe for the timeline: meme-tag children expose `tag_id_object`, while the repeated root meme is suppressed. The frontend maps those generated nested objects directly instead of issuing a separate tag-hydration request.
