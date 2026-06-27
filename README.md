# PureBlog v3

> 一个克制的中文博客系统。**辨识度来自排版、留白与唯一的强调色（黛青 `#235A73`），而非装饰。**

PureBlog v3 是一次**断代重写**:与 v2(PHP / ThinkPHP / MySQL)没有代码或数据上的延续,只继承项目名与产品形态(文章 / 多级分类 / 多级评论 / 写作后台)。

## 技术栈

| 层 | 选型 |
|---|---|
| 后端 | **Go + Gin**(JSON API)· **GORM** · **goose** 迁移 · **PostgreSQL** |
| 渲染 | Markdown → HTML 由 **goldmark + bluemonday** 服务端完成 |
| 鉴权 | **JWT**(access + refresh 轮换)· argon2 |
| 前台 | **Nuxt 3**(SSR,SEO/RSS/sitemap) |
| 后台 | **Vue 3 + Vite** SPA · 编辑器 **Vditor** |
| 设计 | 黛 Dài 设计系统(`frontend/packages/ui`) |
| 部署 | **docker-compose** + **Caddy**(自动 HTTPS) |

## 目录结构

```
PureBlog/
├─ backend/      Go · Gin · GORM(六边形:domain/service/store/http/auth/render/config)
├─ frontend/     pnpm 工作区
│  ├─ apps/web        Nuxt 3(前台,SSR)
│  ├─ apps/admin      Vue 3 + Vite(后台 SPA)
│  └─ packages/       ui(黛设计) · api-types(由后端 OpenAPI 生成)
├─ deploy/       docker-compose · Caddyfile · .env.example
└─ docs/
```

## 本地开发

前置:Go ≥ 1.22 · Node ≥ 20 · pnpm ≥ 9 · Docker。

```bash
# 1) 起一个本地 Postgres(或用你自己的)
docker run -d --name pureblog-pg -e POSTGRES_PASSWORD=pureblog \
  -e POSTGRES_DB=pureblog -p 5432:5432 postgres:16-alpine

# 2) 后端
cd backend
cp .env.example .env          # 按需修改;首启会用 ADMIN_USERNAME/ADMIN_PASSWORD 建管理员
go run ./cmd/pureblog         # 启动时自动跑 goose 迁移并播种管理员

# 3) 前端
cd ../frontend
pnpm install
pnpm dev                      # 同时起 web(Nuxt) 与 admin(Vite)
```

- 前台:<http://localhost:3000>
- 后台:<http://localhost:5173>(默认账号见 `backend/.env.example`)
- API:<http://localhost:8080/api/v1>

## 一键起全栈

```bash
cd deploy
cp .env.example .env
docker compose up -d
```

## 生成 API 类型 / OpenAPI

```bash
cd backend && make swag       # 生成 docs/swagger.json
cd ../frontend && pnpm gen:api  # swagger.json → packages/api-types
```

## 许可

[Apache-2.0](./LICENSE)。
