# memebrary frontend

Small Svelte/Vite UI for the anonymous meme library. In development, Vite proxies the generated and compatibility `/api` routes to the djangolang service at `http://localhost:7070`.

```sh
npm install
npm run dev
npm run check
npm run build
```

The production image is nginx and serves the compiled app. Its `/api` location proxies to the in-cluster `memebrary-backend` Service; compatibility media is served at `/api/custom/media`.
