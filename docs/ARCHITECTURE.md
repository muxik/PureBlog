# 架构说明

PureBlog v3 是一个 polyglot monorepo:Go 后端 + pnpm 前端工作区,彻底重写自 v2(PHP/MySQL),仅沿用项目名。

## 后端(`backend/`)—— 六边形

依赖方向始终指向内层 `domain`:

```
http  ──▶  service  ──▶  domain  ◀──  store
(gin)      (用例)       (实体+接口)    (GORM 实现)
                  ▲
        auth · render · config(支撑件)
```

- `domain/` — 纯实体 + 仓储接口,**不认识 GORM,也不认识 Gin**。
- `store/` — GORM model + 仓储实现 + `domain ↔ model` 映射;`mapErr` 把 GORM 错误翻译成 `domain.ErrNotFound` / `ErrConflict`。
- `service/` — 用例(发布时调 `render`、slug 去重、草稿→发布盖时间戳)。
- `http/` — Gin 路由/handler/dto/中间件,handler 带 swag 注解。
- `auth/` — argon2id 口令 + JWT(access + refresh 轮换)。
- `render/` — goldmark + bluemonday,**发布与 `/api/v1/admin/render` 预览共用同一渲染器**,保证"所见=所发"。
- 迁移用 **goose**(`migrations/*.sql`,embed 进二进制),启动自动 `Up`;首启用 `ADMIN_USERNAME/PASSWORD` 播种管理员。

## 前端(`frontend/`)—— pnpm 工作区

```
apps/web      Nuxt 3(SSR,前台,SEO)
apps/admin    Vue 3 + Vite SPA(后台,Vditor 编辑器)
packages/ui   黛 Dài 设计系统(令牌 CSS + lunar.ts),web/admin 共用
packages/api-types  由后端 OpenAPI 生成的 TS 类型
```

- **类型同步(day-1 codegen)**:`make swag`(Swagger 2.0)→ `swagger2openapi`(转 OpenAPI 3)→ `openapi-typescript` → `packages/api-types/src/schema.ts`。`pnpm gen:api` 一键跑完。
- **前台双 API base**:`useApiBase()` 在 SSR 用 `apiBaseServer`(容器内直连 backend),在浏览器用 `public.apiBase`(经 Caddy)。
- 字体仍需自托管 + 子集化(`cn-font-split`)——见 README。

## 部署(`deploy/`)

`docker compose`:`db`(Postgres)· `backend`(Go/distroless)· `web`(Nuxt SSR)· `admin`(静态,Caddy)· `caddy`(反代 + 自动 HTTPS)。
Caddy:`/api/*`→backend,`/`→web,`admin.<域名>`→admin。

## 验证基线(本次脚手架)

| 层 | 验证 |
|---|---|
| 后端 | `go build` ✅ · `go vet` ✅ · `swag init` ✅ |
| 类型生成 | swagger→openapi3→ts ✅(真 schema) |
| 后台 | `vue-tsc --noEmit` ✅ |
| 前台 | `nuxt typecheck` ✅ |
| 部署 | `docker compose config` ✅ |
| 运行时 | 待本地起 Docker/Postgres 后冒烟(本环境 daemon 未启) |
