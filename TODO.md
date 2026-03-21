# forge-site — TODO

Work queue for forge-cms.dev. See GOVERNANCE.md for how this file is used.

Items in **Backlog** are high-level categories — not yet broken into tasks.
When an item moves to **Up next**, it is broken into atomic tasks before work begins.
Maximum 2 items **In progress** at once.

Last updated: 2026-03-21 (S39)

---

## In progress

---

## Up next

---

## Done

- [x] Re-enable Caddy health_uri — forge v1.1.7 A59 exempts /_health from HTTPS redirect (S39) — 2026-03-21
- [x] Remove ldflags versioning, fix HeadFunc Canonicals, clean docker-compose (S38) — 2026-03-20
- [x] Upgrade forge v1.1.7 (S37) — 2026-03-20
- [x] Route home handler 404 through forge.WriteError (S36) — 2026-03-20
- [x] Add styled 404 error template for module routes (S35) — 2026-03-20
- [x] Upgrade forge v1.1.4, replace siteBaseURL string concat with forge.AbsURL() (S34) — 2026-03-20
- [x] CSS fixes: .prose margin-bottom, .post__date fixed-width column, date span always rendered — 2026-03-20
- [x] Upgrade forge-mcp v1.0.5, delete_post tool (S33) — 2026-03-19
- [x] Upgrade forge-mcp v1.0.4, MCP content wrapper (S32) — 2026-03-19
- [x] Upgrade forge-mcp v1.0.3, list response format fix (S31) — 2026-03-19
- [x] Build-time version via ldflags (S30) — 2026-03-19
- [x] Upgrade forge-mcp v1.0.2, admin read tools automatically exposed (S29) — 2026-03-19
- [x] Fix horizontal scroll on narrow mobile viewports (S28) — 2026-03-18
- [x] Fix absolute og:url and add og:image (S19) — 2026-03-17
- [x] Wire forge-mcp — MCP read+write endpoints (S20) — 2026-03-17
- [x] cmd/mcp proxy for Claude Desktop (S21) — 2026-03-17
- [x] Upgrade forge v1.1.2 + forge-mcp v1.0.1, array-aware tags (S22) — 2026-03-18
- [x] Remove internal workaround comments from templates + main.go (S23) — 2026-03-18
- [x] Add BlogPosting + TechArticle JSON-LD to show templates (S24) — 2026-03-18
- [x] Add required image field to BlogPosting JSON-LD (S25) — 2026-03-18
- [x] Upgrade forge v1.1.3, A53 content negotiation fix (S26) — 2026-03-18
- [x] Extend README + clarify copilot-instructions for public repo (S27) — 2026-03-18
- [x] Fix horizontal scroll on narrow mobile viewports (S28) — 2026-03-18

## Backlog

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
