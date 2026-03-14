# forge-site — Decisions

Architectural decisions for forge-cms.dev. Append only — revisions require an Amendment.

| # | Title | Status | Date |
|---|-------|--------|------|
| D1 | Content types | Locked | 2026-03-12 |
| D2 | URL structure | Locked | 2026-03-12 |
| D3 | Forge features per module | Locked | 2026-03-12 |
| D4 | Storage: SQLite | Locked | 2026-03-12 |
| D5 | Deployment stack | Locked | 2026-03-12 |
| D6 | No frontend dependencies | Locked | 2026-03-12 |
| D7 | Home page as static route | Locked | 2026-03-12 |
| D8 | Analytics: third-party interim | Locked | 2026-03-12 |
| D9 | Health endpoint | Locked | 2026-03-12 |

---

## D1 — Content types

Two content types: `Post` (devlog) and `DocPage` (documentation). Both embed `forge.Node`.

```go
type Post struct {
    forge.Node
    Title  string   `forge:"required,min=3"`
    Body   string   `forge:"required"`
    Tags   []string
}

type DocPage struct {
    forge.Node
    Title   string `forge:"required,min=3"`
    Body    string `forge:"required"`
    Section string
    Order   int
}
```

Both implement `Headable`, `Markdownable`, and `AIDocSummary`.
Adding a field requires an Amendment. A third content type requires a new Decision.

---

## D2 — URL structure

| Content | List | Detail | Feed |
|---------|------|--------|------|
| Post | `/devlog` | `/devlog/{slug}` | `/devlog/feed.xml` |
| DocPage | `/docs` | `/docs/{slug}` | — |
| Home | `/` | — | — |

Global (auto): `/sitemap.xml`, `/llms.txt`, `/llms-full.txt`, `/robots.txt`, `/feed.xml`

Changing any prefix requires an Amendment and a redirect via `app.Redirect()`.

---

## D3 — Forge features per module

**Post (`/devlog`):** `SitemapConfig{}`, `Social(OpenGraph, TwitterCard)`,
`Feed(FeedConfig{Title: "Forge Devlog"})`, `AIIndex(LLMsTxt, LLMsTxtFull, AIDoc)`,
`Cache(5 * time.Minute)`, `Templates("templates/devlog")`

**DocPage (`/docs`):** `SitemapConfig{}`, `AIIndex(LLMsTxt, LLMsTxtFull, AIDoc)`,
`Cache(10 * time.Minute)`, `Templates("templates/docs")` — no RSS feed.

**App:** `app.SEO(&forge.RobotsConfig{Sitemaps: true})`, `app.Health()`

Enabling/disabling a feature requires an Amendment. Cache TTLs do not.

---

## D4 — Storage: SQLite

`forge.SQLRepo[T]` backed by `modernc.org/sqlite` (pure-Go, no CGo).
Database path from `DATABASE_PATH` env var (default `./data/forge.db`).
WAL mode enabled at startup. File is gitignored.

`modernc.org/sqlite` is the only permitted third-party Go dependency.
Migrating to PostgreSQL is a one-line change (`forge-pgx`) — no application code changes.

---

## D5 — Deployment stack

```
internet → Caddy (:443, TLS auto) → Forge (:8080) → SQLite volume
```

Caddyfile:
```
forge-cms.dev {
    reverse_proxy localhost:8080
}
```

Hetzner CAX11 (~€4/month). Two-stage Docker build: `golang:1.22-alpine` → `alpine:latest` (~15MB).
Static files and templates embedded via `go:embed` — no volume mount needed.

---

## D6 — No frontend dependencies

No npm, no CSS frameworks, no JS libraries. Plain CSS (three files) and minimal
vanilla JS (copy-to-clipboard and nav toggle only).

External resources: Google Fonts (`<link>`) and analytics script (D8). Nothing else.
Any addition requires superseding this decision.

---

## D7 — Home page as static route

`/` is a `forge.Handle` route — not a `Module[T]`. Renders `templates/home/home.html`
with a manually constructed `forge.Head{}`. Not in sitemap (correct — it is the root).

---

## D8 — Analytics: third-party interim

Privacy-first, cookieless, EU-hosted analytics script until `forge.Analytics`
ships (Phase 2). Provider must be cookieless, EU-hosted, GDPR-compliant without
a consent banner. Candidates: Plausible, Umami, Fathom.

One `<script defer src="...">` in `templates/base.html`. No config files, no data
leaving the EU, no consent banner. Specific provider recorded as Amendment S1 when chosen.
Superseded when `forge.Analytics` ships.

---

## D9 — Health endpoint

`app.Health()` from Forge ≥ v1.0.6 (Amendment A42). Mounts `GET /_health`.

```go
app.Health()
```

`Config.Version` set from build-time ldflags. Response: `{"status":"ok","version":"X.Y.Z"}`.
Caddy health check points to `/_health`.

---

### Amendment S2 — Dockerfile build decisions (amends D5)

**Decision:** The following choices are locked for the Dockerfile:

| Choice | Value | Rationale |
|--------|-------|-----------|
| Build stage image | `golang:1.26-alpine` | Matches installed toolchain; `go.mod` declares `go 1.26`. D5 referenced the stale `1.22`. |
| `CGO_ENABLED` | `0` | `modernc.org/sqlite` is pure-Go — no C toolchain needed. Produces a fully static binary. |
| Build flags | `-trimpath -ldflags "-s -w"` | Strips local file paths and debug symbols; keeps the runtime image small. |
| Runtime user | `app` (uid 1000) | Non-root for container security hardening. |
| `VERSION` injection | `ARG VERSION=dev` → `-X main.Version=${VERSION}` | Build-time version flows into `Config.Version` and the `/_health` response. Defaults to `dev` when not supplied. |
| Data volume | `VOLUME ["/app/data"]` | Signals the SQLite mount point to Docker; wired to a named volume in `docker-compose.yml`. |

**Consequences:** `Dockerfile` only. No application code changes beyond what was already decided in D4 and D9.

---

### Amendment S3 — docker-compose configuration decisions (amends D5, S2, D9)

**Decision:** The following choices are locked for `docker-compose.yml`:

| Choice | Value | Rationale |
|--------|-------|-----------|
| Caddy topology | Caddy runs on the **host machine**, not in a Docker container. The app container is the only containerised process. | Keeps the setup simple: one container, one binary, no inter-container networking. Caddy manages TLS and certificates directly on the host. |
| Port binding | `127.0.0.1:8080:8080` | Loopback-only binding is the direct consequence of the Caddy-on-host topology. Only the host Caddy can reach port 8080; the port is not reachable from the public internet or other containers. |
| Container healthcheck | `wget -qO- http://localhost:8080/_health` every 30 s | Docker daemon container check, separate from Caddy's health check (D9). Uses `wget` — present in the `alpine` base image without extra packages. |
| SQLite volume name | `forge_data` | Named volume backing `/app/data` in the app container. Consistent with S2's `VOLUME ["/app/data"]`. |
| `SECRET` env var | No default; compose errors if unset | `requireEnv("SECRET")` in `main.go` crashes the process if `SECRET` is empty. No default in compose enforces this at the orchestration layer. Operators must supply `SECRET` via `.env` or the host environment. |

**Consequences:** `docker-compose.yml` and `Caddyfile` (Caddy-on-host architecture).

---

### Amendment S4 — Tags storage and content type schema (amends D1, D4)

**Decision:** The following storage choices are locked for content types:

| Choice | Value | Rationale |
|--------|-------|-----------|
| `Post.Tags` storage | JSON TEXT column (`["forge","go"]`) | SQLite has no native array type. Tags are serialised/deserialised via a `JSONStringSlice` custom type implementing `driver.Valuer` + `sql.Scanner`. Field tagged `db:"tags"`. |
| Table names | `posts`, `doc_pages` | Auto-derived by `SQLRepo` (`Post` → `posts`, `DocPage` → `doc_pages`). No `forge.Table()` override needed. |
| Schema creation | Manual `CREATE TABLE IF NOT EXISTS` in `schema.go`, called from `main()` before modules are wired | Forge does not auto-create schema. Upsert (`ON CONFLICT (id) DO UPDATE`) requires tables to exist at startup. |
| Type parameter | Pointer type throughout: `forge.NewSQLRepo[*Post](db)`, proto `(*Post)(nil)`, `Module[*Post]` | `NewModule` infers T from the proto value. `Repo[T]` must match — `Repository[Post]` does not satisfy `Repository[*Post]`. All Forge examples and tests use pointer types consistently. |

**Consequences:** `post.go`, `docpage.go`, `stringslice.go`, `schema.go`, `main.go`.

---

### Amendment S6 — Forge v1.0.6 template workarounds (amends D3, D7)

**Decision:** Three known Forge v1.0.6 limitations require explicit workarounds
in forge-site. Each will be removed when the corresponding Forge feature ships.

| Limitation | Workaround | Removes when |
|------------|------------|--------------|
| No shared template partials | Nav and footer duplicated in all four module templates | Forge shared partials (Phase 2) |
| `forgeHeadTmpl` is package-private | `base.html` uses manual `<head>` meta tags | `forge.HeadPartial()` or equivalent (Phase 2) |
| `forge.New` accepts invalid config silently | `main.go` uses `forge.MustConfig` explicitly | `forge.New` enforces validation internally (Phase 2) |

**Consequences:** Any change to nav, footer, or `<head>` meta must be applied
in `base.html` + all four module templates until shared partials ship.

---

### Amendment S5 — Template delivery: disk for modules, embed for home (amends D3, D5, D7)

**Decision:** Forge v1.0.6 loads module templates from the OS filesystem
via `os.Stat` + `template.ParseFiles`. There is no `embed.FS` support in
`Templates()` or `TemplatesOptional()`. This has two consequences:

| Concern | Decision |
|---------|----------|
| Module templates (`templates/devlog/`, `templates/docs/`) | Copied into the Docker runtime image via `COPY templates/ /app/templates/` in `Dockerfile`. Forge reads them from disk at startup. |
| Home page (`templates/base.html`, `templates/home/home.html`) | Parsed at startup from the embedded FS via `template.ParseFS(templates, ...)`. The `//go:embed templates` directive includes all template files in the binary. |
| Home page `<head>` meta | Written as manual HTML tags in `base.html`. `forgeHeadTmpl` (the partial Forge injects into module templates) is package-private and not accessible from outside the `forge` package. This is intentional — not a bug. |

**Consequences:** `Dockerfile` gains `COPY templates/ /app/templates/`. The
`//go:embed` directive on `main.go` changes from `all:templates` to `templates`
(no hidden files remain). Any nav/footer/head change must be applied in five
places: `base.html` + `devlog/list.html` + `devlog/show.html` +
`docs/list.html` + `docs/show.html` (per S6).
