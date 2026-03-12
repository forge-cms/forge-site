# forge-site — Governance

Last updated: 2026-03-12

Same principles as Forge itself: decisions are explicit, changes are intentional,
AI is a participant — not an autonomous agent. Lighter process than the framework.

---

## Documents

| File | Purpose |
|------|---------|
| `DECISIONS.md` | Architecture: content types, URLs, Forge features, deployment |
| `DESIGN.md` | CSS architecture, tokens, BEM, typography |
| `TODO.md` | Work queue |
| `.github/copilot-instructions.md` | Copilot working instructions |
| `README.md` | How to run and deploy |

---

## Amendments

Required when a change:
- Crosses a file boundary (e.g. adding a content type field touches `main.go`, template, CSS)
- Contradicts an existing decision
- Changes URL structure

Not required for: copy edits, CSS within an existing component, new posts/docs, typos.

Format — append to `DECISIONS.md`:
```
### Amendment S{N} — {title} (amends {ID})
**Decision:** {what is changing}
**Rationale:** {why}
**Consequences:** {files/behaviours affected}
```

`S` prefix distinguishes site amendments from Forge framework amendments (`A`).

---

## Commits

Copilot proposes — never commits autonomously.

```
{type}({scope}): {short description}

{body}

Decisions: {IDs or "none"}
```

Types: `feat`, `fix`, `content`, `style`, `docs`, `chore`

Run autonomously (no approval needed): `go build ./...`, `go vet ./...`,
`gofmt -l .`, any read-only file or git command.

---

## TODO.md format

```
## In progress   (max 2)
## Up next
## Backlog       (categories only — break into tasks when moving to Up next)
## Done
```

---

## Forge core issues

Real usage will surface Forge bugs and missing features. When found:

**Bug** → open issue in `forge-cms/forge`, fix there, update `go.mod` when tagged.

**Missing feature** → small: propose as Amendment in forge repo. Large: add to forge `BACKLOG.md`.

**Blocked and needs workaround** → document in `DECISIONS.md`:
- What it does
- Link to forge issue
- What removes it (fix + version)

Copilot must stop and report before writing any workaround. Never silently
reimplement something Forge should handle. The Claude project session is the
coordination layer — bring issues there before acting.

---

## Deployment

- All deployments manual until CI/CD is explicitly set up
- Docker image built locally, pushed to Hetzner via SSH
- SQLite in a named Docker volume — never inside the image
- No deploy without green `go build ./...`
- Tags: `vMAJOR.MINOR.PATCH`, annotated only
