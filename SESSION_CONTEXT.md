# forge-site — Session Context

Updated by Copilot at every commit. Paste this into the Claude architect
session alongside PROJECT_BRIEF.md to restore full context.

## Current status
- Scaffold: ✅ done
- Content types: ✅ done
- Templates — base + home: ✅ done
- Templates — module (4 templates): ✅ done
- Smoke-test (all routes 200/404 as expected): ✅ done
- Static assets smoke-test: 🔲 up next

## Latest amendments
- S1: analytics provider (pending — not yet chosen)
- S2: Dockerfile Go 1.26 + build flags
- S3: docker-compose topology (Caddy on host, loopback bind)
- S4: Tags JSON storage, pointer types throughout
- S5: Dockerfile `COPY templates/` — Forge v1.0.6 reads module templates from disk
- S6: Forge v1.0.6 template workarounds (no shared partials, package-private `forgeHeadTmpl`, `MustConfig` explicit)

## Open upstream items (BACKLOG.md)
- Phase 2: shared partials, `forge.HeadPartial()`, `forge.New` enforces validation internally

## Forge version in use
github.com/forge-cms/forge v1.0.6

## Next up
- Static assets smoke-test (verify CSS and fonts served correctly on all 5 routes)
- Choose analytics provider (Amendment S1)
- Deployment sprint
