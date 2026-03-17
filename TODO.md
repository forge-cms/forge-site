# forge-site — TODO

Work queue for forge-cms.dev. See GOVERNANCE.md for how this file is used.

Items in **Backlog** are high-level categories — not yet broken into tasks.
When an item moves to **Up next**, it is broken into atomic tasks before work begins.
Maximum 2 items **In progress** at once.

Last updated: 2026-03-17 (S20)

---

## In progress

---

## Up next

---

## Done

- [x] Fix absolute og:url and add og:image (S19) — 2026-03-17
- [x] Wire forge-mcp — MCP read+write endpoints (S20) — 2026-03-17
- [x] Wire forge-mcp — MCP read+write endpoints (S20) — 2026-03-17

## Backlog

- [ ] Markdown rendering — Body field renders raw markdown; needs markdown→HTML
  in templates via `forge.TemplateFuncMap` (Amendment A46 in core, then
  wire `{{ .Item.Body | markdown }}` in devlog/docs show templates)
- [ ] Health check — Caddy `health_uri` removed as workaround; Forge HTTPS
  redirect breaks internal `/_health` calls (Amendment S10 context)
- [ ] ADMIN_TOKEN — set a persistent token in `.env` so it survives restarts
  without re-fetching from logs

---

## Done

- [x] OG + Twitter Card on show pages — missing meta added to devlog/docs show templates, Amendment S16 — 2026-03-17
- [x] Launch — 6 content items published (3 devlog + 3 docs), `llms.txt` and sitemap verified — 2026-03-15
- [x] Deployment — Hetzner CX23 provisioned, Docker volume, TLS live via Caddy, `forge-cms.dev` DNS pointed — 2026-03-15
- [x] Analytics — Plausible Cloud added to `base.html`, Amendment S1 — 2026-03-15
- [x] Static assets — `io/fs.Sub` + `http.FileServerFS` wired; `tokens.css`, `base.css`, `components.css` all return 200 — 2026-03-14

---

## [Older Done entries below]

- [x] Templates — `base.html`, `home.html`, devlog + docs list/show; home handler wired; `forge.Templates` promoted; amendments S5+S6 — 2026-03-14
- [x] Content types — `post.go`, `docpage.go`, `stringslice.go`, `schema.go`, `seed.go`, main.go wired; Amendment S4 — 2026-03-14
- [x] Scaffold — `main.go`, `go.mod/go.sum`, `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `post.http`, `README.md`; admin token, amendments S2+S3 — 2026-03-14
- [x] Design system defined — `DESIGN.md`, `tokens.css`, `base.css`, `components.css` — 2026-03-12
- [x] Governance model defined — `GOVERNANCE.md`, `DECISIONS.md`, `TODO.md` — 2026-03-12
- [x] Copilot instructions written — `copilot-instructions.md` + template — 2026-03-12
- [x] Landing page mockup approved — hero copy, colour palette, typography — 2026-03-12
- [x] Health endpoint — `app.Health()` available in Forge v1.0.6 (Amendment A42, D9) — 2026-03-12 — `post.go`, `docpage.go`, `stringslice.go`, `schema.go`, `seed.go`, main.go wired; Amendment S4 — 2026-03-14
- [x] Scaffold — `main.go`, `go.mod/go.sum`, `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `post.http`, `README.md`; admin token, amendments S2+S3 — 2026-03-14
- [x] Design system defined — `DESIGN.md`, `tokens.css`, `base.css`, `components.css` — 2026-03-12
- [x] Governance model defined — `GOVERNANCE.md`, `DECISIONS.md`, `TODO.md` — 2026-03-12
- [x] Copilot instructions written — `copilot-instructions.md` + template — 2026-03-12
- [x] Landing page mockup approved — hero copy, colour palette, typography — 2026-03-12
- [x] Health endpoint — `app.Health()` available in Forge v1.0.6 (Amendment A42, D9) — 2026-03-12
