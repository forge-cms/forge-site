# forge-site — Session Context

Updated by Copilot at every commit. Paste into Claude alongside
PROJECT_BRIEF.md to restore full context.

## Current status
- Scaffold: ✅ done
- Content types: ✅ done
- Templates — base + home + module (4): ✅ done
- Favicon + OG + JSON-LD + Twitter meta: ✅ done
- Logo assets (Affinity): ✅ done
- Social profiles updated: ✅ done
- Deployment sprint: 🔲 not started

## Latest amendments
- S1: Plausible Cloud analytics added to base.html ✅
- S2: Dockerfile Go 1.26 + build flags
- S3: docker-compose topology (Caddy on host, loopback bind)
- S4: Tags JSON storage, pointer types throughout
- S5: Dockerfile COPY templates/
- S6: Forge v1.0.6 template workarounds (nav/footer duplication, forge:head)
- S7: Author and Published fields in Head()
- S8: OG image fallback in base.html (workaround)
- S9: HeadFunc for list pages + Organization JSON-LD workaround

## Open corepilot items (sent)
- A43: NewSQLRepo pointer type docs ✅ sent
- Phase 2 promoted: shared partials, forge:head helper, MustConfig enforcement ✅ sent
- OGDefaults + AppSchema → Phase 2 ✅ sent
- AppSchema + OGDefaults backlog entries ✅ sent

## Forge version in use
github.com/forge-cms/forge v1.0.6

## Pending before deployment
- Choose analytics provider → Amendment S1
- Write and publish content (6 items: 3 devlog + 3 docs) via post.http against prod
- Deployment sprint: Hetzner CX23, Docker, Caddy, domain DNS

## Content ready to publish (post to prod after go-live)
- Post 1: Why I built Forge (why-i-built-forge)
- Post 2: Forge v1.0.0 — what shipped (forge-v1-0-0)
- Post 3: forge-cms.dev is built with Forge (forge-cms-dev-is-built-with-forge)
- Doc 1: Getting Started (getting-started) — Section: Guides, Order: 1
- Doc 2: Content Types (content-types) — Section: Guides, Order: 2
- Doc 3: Content Lifecycle (content-lifecycle) — Section: Guides, Order: 3

## Design decisions (locked)
- Font: JetBrains Mono via Bunny Fonts (not Google Fonts)
- Logo: [forge] med glødende o i hullet — mørk og lys variant
- Accent: #e8702f, Baggrund: #0a0b0d (#000000 i logo-filer)
- Favicon: [f] pixel-art 16×16 og 32×32, fuldt logo 180×180