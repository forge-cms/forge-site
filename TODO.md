# forge-site — TODO

Work queue for forge-cms.dev. See GOVERNANCE.md for how this file is used.

Items in **Backlog** are high-level categories — not yet broken into tasks.
When an item moves to **Up next**, it is broken into atomic tasks before work begins.
Maximum 2 items **In progress** at once.

Last updated: 2026-03-12

---

## In progress

---

## Up next

### Scaffold

1. [x] `go.mod` — initialize module `forge-site`, add `forge` dependency, generate `go.sum`
2. [x] `main.go` skeleton — `forge.New()`, health endpoint, `//go:embed` directives for `templates/` and `static/`, `.Run()`; content module stubs left as TODO
3. [x] Admin token — `forge.SignToken` call in `main.go`; log token to stdout on first run when `ADMIN_TOKEN` env var is unset
4. [x] `Dockerfile` — multi-stage build (`golang:1.22-alpine` → `alpine` runtime); copy binary with embedded assets
5. [x] `docker-compose.yml` — app service, volume for `data/forge.db`, `PORT` + `ADMIN_TOKEN` env vars
6. [x] `Caddyfile` — reverse-proxy to app container, ACME TLS, domain placeholder
7. [x] `post.http` — REST Client file (VS Code); create a `Post`, create a `DocPage`, publish a `Post`
8. [x] `README.md` — project overview, local dev (`go run .`), `post.http` usage, Docker deploy steps; include note on admin token rotation — `TTL=0` means no automatic expiry, token must be manually rotated (change `SECRET`) if compromised
9. [x] Verify — `go build ./...` + `go vet ./...` pass clean

---

## Backlog

- [ ] Content types — `Post` and `DocPage` with all fields, interfaces, and seed data
- [ ] Templates — base layout, home page, devlog list/show, docs list/show
  > Note: remove `all:` prefix from `//go:embed all:templates` and delete `.gitkeep` files when real templates are in place.
- [ ] Static assets — `tokens.css`, `base.css`, `components.css` wired and served
- [ ] Analytics — choose provider (Plausible / Umami / Fathom), add script to base layout, record as Amendment S1 (D8)
- [ ] Deployment — Hetzner server provisioned, Docker volume, TLS live, domain pointed
  > Note: replace `DOMAIN_PLACEHOLDER` in `Caddyfile` with `forge-cms.dev`
- [ ] Launch — seed data complete, `llms.txt` verified, sitemap verified, go live

---

## Done

- [x] Design system defined — `DESIGN.md`, `tokens.css`, `base.css`, `components.css` — 2026-03-12
- [x] Governance model defined — `GOVERNANCE.md`, `DECISIONS.md`, `TODO.md` — 2026-03-12
- [x] Copilot instructions written — `copilot-instructions.md` + template — 2026-03-12
- [x] Landing page mockup approved — hero copy, colour palette, typography — 2026-03-12
- [x] Health endpoint — `app.Health()` available in Forge v1.0.6 (Amendment A42, D9) — 2026-03-12
