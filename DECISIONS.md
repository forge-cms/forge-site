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

---

### Amendment S8 — OG image fallback workaround (amends D3)

**Decision:** Forge v1.0.6 does not support app-level OG image fallback.
The static fallback `og:image` is hardcoded in templates as a workaround.

| Template | Approach |
|----------|----------|
| `base.html` | Static `<meta property="og:image">` hardcoded directly (no dynamic Head) |
| 4 module templates | Conditional fallback after `{{template "forge:head" .Head}}`: `{{if not .Head.Image.URL}}` |

**Consequences:** All five templates carry the `og:image` fallback. Removed when
`forge.OGDefaults{}` ships (Phase 2 — tracked in forge BACKLOG.md).

---

### Amendment S9 — HeadFunc for list pages and Organization JSON-LD (amends D3, D7)

**Decision:** `HeadFunc` added to Post and DocPage modules so list pages have
titles and descriptions. Organization JSON-LD hardcoded in `base.html` —
`forge.Handle` routes cannot use `SchemaFor` (Forge core limitation). Removed
when `forge.AppSchema{}` ships (Phase 2).

**Consequences:** `main.go` (HeadFunc on both modules), `base.html` (JSON-LD script).

---

### Amendment S1 — Analytics provider: Plausible Cloud (amends D8)

**Decision:** Plausible Cloud selected as analytics provider. Script injected in
`templates/base.html` inside `<head>`. Uses `async` (equivalent to D8's `defer`
for end-of-head placement). Cookieless, EU-hosted, GDPR-compliant — no consent
banner required.

**Consequences:** `base.html` (Plausible script tag). Superseded when
`forge.Analytics` ships (Phase 2).

---

### Amendment S10 — Wire Authenticate middleware (amends D1, D3)

**Decision:** `app.Use(forge.Authenticate(forge.BearerHMAC(secret)))` added to
`main.go` after `forge.New()` and before `app.Health()`. Without this line all
requests were treated as `GuestUser` regardless of the `Authorization` header,
making the admin API effectively unauthenticated.

**Consequences:** `main.go` (one line added). No template or CSS changes.

**Amended by S10 (v1.0.8):** `app.Use(forge.Authenticate(...))` removed —
forge v1.0.8 wires BearerHMAC authentication automatically inside `forge.New()`
when `Config.Secret` is set. `go.mod` updated to v1.0.8.

---

### Amendment S11 — Rename DocPage.Order column to sort_order (amends D1, S4)

**Problem:** `"order"` is a reserved SQL keyword. `SQLRepo` generates
`INSERT`/`UPDATE` SQL without quoting column names, causing a syntax error when
saving any `DocPage`.

**Decision:** Rename the column to `sort_order` in both the `db` struct tag
(`docpage.go`) and the `CREATE TABLE` statement (`schema.go`). No quoting
required.

**Consequences:** `docpage.go` (`db:"sort_order"`), `schema.go` (`sort_order`
column). Existing database volumes must be reset: `docker-compose down -v &&
docker-compose up -d --build`.

---

### Amendment S12 — Markdown rendering via forge_markdown in show templates (amends D3)

**Decision:** forge v1.0.9 activates `forge_markdown` in `forge.TemplateFuncMap()`.
Both `templates/devlog/show.html` and `templates/docs/show.html` already call
`{{forge_markdown .Content.Body}}` inside `<div class="prose">` — no template
changes required. Forge injects `TemplateFuncMap()` internally when parsing
module templates via `forge.Templates()`.

**Consequences:** `go.mod` updated to v1.0.9. Body fields now render as HTML.

---

### Amendment S13 — Docker container user 1000:1000 (amends D5)

**Problem:** The `forge_data` volume was owned by root after first creation,
causing the app container to fail to write `forge.db` at startup.

**Decision:** Add `user: "1000:1000"` to the `app` service in
`docker-compose.yml`. The container runs as uid 1000 from the start and can
write to the mounted volume without a manual `chown` step.

**Consequences:** `docker-compose.yml` (`user: "1000:1000"` on app service).

---

### Amendment S14 — Upgrade forge to v1.0.10 (amends D3)

**Decision:** Upgrade `github.com/forge-cms/forge` from v1.0.9 to v1.0.10.
No template or code changes required.

**Consequences:** `go.mod` / `go.sum` updated.

---

### Amendment S15 — Upgrade forge to v1.0.11 (amends D3)

**Decision:** Upgrade `github.com/forge-cms/forge` from v1.0.10 to v1.0.11.
No template or code changes required.

**Consequences:** `go.mod` / `go.sum` updated.

---

### Amendment S16 — OG and Twitter Card tags in module show templates (amends S8)

**Problem:** `forge:head` emits only `<title>`, `<meta name="description">`,
and `<link rel="canonical">`. Module show templates are standalone documents
that do not extend `base.html` (Amendment S6), so `base.html`'s OG block
never renders for `/devlog/{slug}` or `/docs/{slug}` pages.

**Decision:** Add explicit `og:*` and `twitter:*` meta tags to
`templates/devlog/show.html` and `templates/docs/show.html`, immediately after
`{{template "forge:head" .Head}}`. Data pulled from `.Head.Title`,
`.Head.Description`, and `.Head.Canonical` — all correctly populated by
`Post.Head()` / `DocPage.Head()`. OG image uses site-level fallback (same as
S8). `og:type` is `article` for both types.

**Consequences:** `templates/devlog/show.html` and `templates/docs/show.html`
(10 meta lines added each). Superseded by shared partials in Phase 2.

---

### Amendment S17 — Add twitter:card override in show templates (amends S16)

**Problem:** `forge:head` emits `twitter:card` as `summary`. The S16 block adds
`twitter:image` but without overriding `twitter:card`, X/Twitter ignores the
image and renders the small card format.

**Decision:** Add `<meta name="twitter:card" content="summary_large_image">`
to the S16 block in both `templates/devlog/show.html` and
`templates/docs/show.html`, immediately before `twitter:image`. The later tag
overrides the `summary` emitted by `forge:head`.

**Consequences:** `templates/devlog/show.html` and `templates/docs/show.html`
(one line added each).

---

### Amendment S18 — Upgrade forge to v1.1.1, remove S16/S17 OG workarounds (amends S16, S17)

**Decision:** forge v1.1.1 fixes `forge:head` to emit absolute `og:url`,
`og:image`, and `twitter:card: summary_large_image` natively when `Type:
forge.Article` is set. The S16 and S17 workaround blocks (og:url, og:image,
twitter:card, twitter:image overrides) are removed from both show templates.

**Consequences:** `go.mod` / `go.sum` updated to v1.1.1. S16/S17 override
blocks removed from `templates/devlog/show.html` and
`templates/docs/show.html`.

---

### Amendment S19 — Absolute og:url + og:image via siteBaseURL package var (amends S18)

**Context:** Despite S18, forge:head still emits the `Canonical` value
verbatim. Since `forge.URL()` returns a root-relative path, `og:url` was
still relative. No `og:image` was emitted because `Head.Image` was zero.

**Decision:** Introduce a package-level `var siteBaseURL string` in `main.go`,
set to `BASE_URL` at startup. `post.go` and `docpage.go` `Head()` methods
prefix `Canonical` with `siteBaseURL` and populate `Head.Image` with the
existing `static/Forge-logo-OG1200.png` asset (1200×630). List page
`HeadFunc` canonicals updated to use `baseURL + forge.URL(...)` via closure.

**Consequences:** `og:url` now emits an absolute URL; `og:image`, `og:image:width`,
and `og:image:height` tags now appear on all content pages.
