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
