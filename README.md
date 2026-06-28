<div align="center">

# PureBlog v3

**一个克制的中文博客系统。**

辨识度来自排版、留白与唯一的强调色 —— 黛青 `#235A73`,而非装饰。

[![CI](https://github.com/muxik/PureBlog/actions/workflows/ci.yml/badge.svg)](https://github.com/muxik/PureBlog/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-JSON_API-008ECF)](https://gin-gonic.com)
[![Vue 3](https://img.shields.io/badge/Vue-3-42b883?logo=vuedotjs&logoColor=white)](https://vuejs.org)
[![Nuxt 3](https://img.shields.io/badge/Nuxt-3_SSR-00DC82?logo=nuxtdotjs&logoColor=white)](https://nuxt.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-Apache_2.0-D22128)](./LICENSE)

</div>

---

## 这是什么

PureBlog 是一个**为中文写作而生**的个人博客系统:正文用 [霞鹜文楷](https://github.com/lxgw/LxgwWenKai),标题用宋体,35em 的舒适行宽,中英之间留出「盘古之白」。没有花哨的卡片、阴影和渐变 —— 把注意力还给文字本身。

> v3 是一次**断代重写**。v2 是 PHP / ThinkPHP / MySQL;v3 只继承项目名与产品形态,代码与数据完全重来。

## 界面一览

| 首页 · 年表布局 | 文章 · 阅读与评论 |
|:---:|:---:|
| [![首页](docs/screenshots/home.png)](docs/screenshots/home.png) | [![文章页](docs/screenshots/article.png)](docs/screenshots/article.png) |

<p align="center"><i>暖纸底色、唯一的黛青强调、霞鹜文楷正文与宋体标题;首页可在「年表 / 摘要 / 双栏」三种布局间切换,支持深浅色与公历↔农历。</i></p>

<p align="center">
  <a href="docs/screenshots/tags.png"><img src="docs/screenshots/tags.png" width="640" alt="标签页"></a>
</p>

## 为什么值得一看

不只是又一个博客,它是一份**现代全栈工程的范本**:

- 🏛️ **干净的后端架构** —— Go 六边形分层(`domain / service / store / http`),依赖一律指向内层领域。`domain` 既不认识 GORM 也不认识 Gin,换 ORM、加 GraphQL 都不动核心。
- 🔗 **前后端类型贯通** —— 后端用 swaggo 产出 OpenAPI,前端由它**自动生成 TypeScript 类型**。接口改了,前端编译期就报错,告别手抄 DTO。
- ⚡ **SEO 友好的 SSR 前台** —— Nuxt 3 服务端渲染,每篇文章一个真 URL、一份真 HTML;后台是独立的 Vue SPA,各取所长。
- ✍️ **所见即所发** —— Markdown 由后端 goldmark 渲染并净化,后台原生编辑器的实时预览调用**同一个**渲染端点,预览和发布逐字一致。
- 📦 **一条命令起全栈** —— `docker compose up` 拉起 Postgres + Go + Nuxt + 后台 + Caddy(自动 HTTPS)。
- 🈶 **中文优先** —— 内置公历↔农历换算、CJK 排版令牌;字体**自托管 + 按 unicode-range 子集化**,零第三方 CDN。
- 🔍 **SEO 开箱即用** —— 服务端 `sitemap.xml`、RSS `feed.xml`、逐页 Open Graph 与 canonical。

## 设计系统 · 黛 Dài

| | |
|---|---|
| **底色** | 暖纸 `#FAF8F2`,墨色三级 |
| **强调色** | 仅一种 —— 黛青 `#235A73`(深色模式提亮为 `#79B0C9`) |
| **字体** | 正文 霞鹜文楷 · 标题 Noto Serif SC · 拉丁 Newsreader · 等宽 JetBrains Mono(均自托管) |
| **排版** | 35em 行宽 · 行高 1.85 · 盘古之白 · 避头尾 · 无阴影无渐变 |

## 技术栈

| 层 | 选型 |
|---|---|
| 后端 | **Go + Gin** · **GORM** · **goose** 迁移 · **PostgreSQL** |
| 渲染 | Markdown → HTML 由 **goldmark + bluemonday** 服务端完成 |
| 鉴权 | **JWT**(access + refresh 轮换)· argon2 |
| 前台 | **Nuxt 3**(SSR) |
| 后台 | **Vue 3 + Vite** SPA · 原生 Markdown 编辑器(实时预览复用后端渲染端点) |
| 工程 | pnpm 工作区 · OpenAPI → TS 代码生成 · Docker Compose + Caddy |

## 功能

**已实现**
- ✅ 文章:增删改查、草稿/发布、置顶、浏览计数、slug 自动生成、标签关联
- ✅ 评论:访客提交(默认待审)、后台审核、多级回复
- ✅ 站点设置(站名 / 简介 / 社交 / 关于页 / 默认日期格式)
- ✅ JWT 鉴权后台、Markdown 渲染与净化、实时预览端点
- ✅ 前台(SSR):首页三种布局(年表 / 摘要 / 双栏)、文章详情、归档、标签、搜索、关于
- ✅ 深浅色主题、公历↔农历切换、无限滚动
- ✅ SEO:`sitemap.xml`、RSS `feed.xml`、`robots.txt`、逐页 Open Graph / canonical
- ✅ 字体自托管 + unicode-range 子集化(零 CDN)
- ✅ 自动生成的 OpenAPI 文档与前端类型
- ✅ CI:Go 构建 / vet / Postgres 集成测试、前端 typecheck / build、GHCR 镜像构建

**路线图**
- 🚧 全文搜索增强(当前 `ILIKE` → `pg_trgm` / MeiliSearch)
- 🚧 图片上传与媒体管理
- 🚧 评论垃圾过滤与邮件通知

## 快速开始

前置:Go ≥ 1.25 · Node ≥ 22 · pnpm ≥ 11 · PostgreSQL(或 Docker)。

### 本地开发

```bash
# 1) 后端(自动迁移并按 .env 播种管理员)
cd backend
cp .env.example .env        # 改成你的数据库连接
go run ./cmd/pureblog        # API 默认 :8080

# 2) 前端
cd ../frontend
pnpm install
pnpm dev                     # 前台 :3000 · 后台 :5173
```

### 一条命令起全栈

```bash
cd deploy
cp .env.example .env
docker compose up -d         # Postgres + Go + Nuxt + 后台 + Caddy
```

### 前后端类型同步

```bash
cd backend && make swag        # 生成 docs/swagger.json
cd ../frontend && pnpm gen:api  # swagger → packages/api-types
```

## 项目结构

```
PureBlog/
├─ backend/        Go · Gin · GORM(六边形:domain/service/store/http/auth/render/config)
├─ frontend/       pnpm 工作区
│  ├─ apps/web         Nuxt 3 前台(SSR)
│  ├─ apps/admin       Vue 3 + Vite 后台(SPA)
│  └─ packages/        ui(黛 设计系统) · api-types(由 OpenAPI 生成)
├─ deploy/         docker-compose · Caddyfile
└─ docs/           架构说明
```

详见 [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)。

## 许可

[Apache-2.0](./LICENSE) © muxik
