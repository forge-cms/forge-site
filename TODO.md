# forge-site — TODO

Work queue for forge-cms.dev. See GOVERNANCE.md for how this file is used.

Items in **Backlog** are high-level categories — not yet broken into tasks.
When an item moves to **Up next**, it is broken into atomic tasks before work begins.
Maximum 2 items **In progress** at once.

Last updated: 2026-03-14

---

## In progress

---

## Up next

### Static assets

- [ ] Verify `tokens.css`, `base.css`, `components.css` are served correctly via `//go:embed static`
- [ ] Smoke-test all five routes in a browser; confirm fonts load and layout renders

---

## Backlog

- [ ] Analytics — choose provider (Plausible / Umami / Fathom), add `<script defer>` to all five templates, record as Amendment S1 (D8)
- [ ] Deployment — Hetzner server provisioned, Docker volume, TLS live, domain pointed
  > Note: replace `DOMAIN_PLACEHOLDER` in `Caddyfile` with `forge-cms.dev`
- [ ] Launch — seed data complete, `llms.txt` verified, sitemap verified, go live

---

## Done

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
