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

- [ ] Content types — `Post` and `DocPage` with all fields, interfaces, and seed data

---

## Backlog

- [ ] Templates — base layout, home page, devlog list/show, docs list/show
  > Note: remove `all:` prefix from `//go:embed all:templates` and delete `.gitkeep` files when real templates are in place.
- [ ] Static assets — `tokens.css`, `base.css`, `components.css` wired and served
- [ ] Analytics — choose provider (Plausible / Umami / Fathom), add script to base layout, record as Amendment S1 (D8)
- [ ] Deployment — Hetzner server provisioned, Docker volume, TLS live, domain pointed
  > Note: replace `DOMAIN_PLACEHOLDER` in `Caddyfile` with `forge-cms.dev`
- [ ] Launch — seed data complete, `llms.txt` verified, sitemap verified, go live

---

## Done

- [x] Scaffold — `main.go`, `go.mod/go.sum`, `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `post.http`, `README.md`; admin token, amendments S2+S3 — 2026-03-14
- [x] Design system defined — `DESIGN.md`, `tokens.css`, `base.css`, `components.css` — 2026-03-12
- [x] Governance model defined — `GOVERNANCE.md`, `DECISIONS.md`, `TODO.md` — 2026-03-12
- [x] Copilot instructions written — `copilot-instructions.md` + template — 2026-03-12
- [x] Landing page mockup approved — hero copy, colour palette, typography — 2026-03-12
- [x] Health endpoint — `app.Health()` available in Forge v1.0.6 (Amendment A42, D9) — 2026-03-12
