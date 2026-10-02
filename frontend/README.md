# memebrary frontend

Small Svelte/Vite UI for the anonymous meme library. In development, Vite proxies `/api`, `/media`, and `/healthz` to `http://localhost:8080`.

```sh
npm install
npm run dev
npm run check
npm run build
```

The production image is nginx and serves the compiled app. Its `/api` and `/media` locations proxy to the in-cluster `memebrary-backend` Service.
