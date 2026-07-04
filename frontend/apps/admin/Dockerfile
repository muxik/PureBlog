# Vue 3 + Vite admin SPA, built to static and served by Nginx.
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

FROM nginx:1.27-alpine AS run
COPY --from=build /app/apps/admin/dist /usr/share/nginx/html
RUN printf 'server {\n\tlisten 80;\n\tserver_name _;\n\troot /usr/share/nginx/html;\n\tindex index.html;\n\tgzip on;\n\tgzip_types text/plain text/css application/json application/javascript image/svg+xml;\n\n\tlocation / {\n\t\ttry_files $uri $uri/ /index.html;\n\t}\n}\n' > /etc/nginx/conf.d/default.conf
EXPOSE 80
