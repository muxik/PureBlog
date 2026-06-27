# Vue 3 + Vite admin SPA, built to static and served by Caddy.
# Build context = frontend/ (the pnpm workspace root).
FROM node:22-alpine AS build
RUN corepack enable
WORKDIR /app
# VITE_API_BASE is inlined at build time (browser → public API URL via the proxy).
ARG VITE_API_BASE=http://localhost/api/v1
ENV VITE_API_BASE=$VITE_API_BASE
COPY . .
RUN pnpm install --frozen-lockfile
RUN pnpm --filter @pureblog/admin build

FROM caddy:2-alpine AS run
COPY --from=build /app/apps/admin/dist /srv
RUN printf ':80 {\n\troot * /srv\n\tencode gzip\n\ttry_files {path} /index.html\n\tfile_server\n}\n' > /etc/caddy/Caddyfile
EXPOSE 80
