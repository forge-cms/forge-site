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
- Analytics: ✅ Plausible Cloud live
- Deployment sprint: ✅ live at forge-cms.dev
- Launch: ✅ 6 content items published

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
- S10: Authenticate middleware wired (v1.0.7), removed in v1.0.8 (auto-wired)
- S11: DocPage.Order renamed to sort_order (reserved SQL keyword)
- S12: forge_markdown active in v1.0.9; show templates already correct

## Open corepilot items (sent)
- A43: NewSQLRepo pointer type docs ✅ sent
- Phase 2 promoted: shared partials, forge:head helper, MustConfig enforcement ✅ sent
- OGDefaults + AppSchema → Phase 2 ✅ sent
- AppSchema + OGDefaults backlog entries ✅ sent

## Forge version in use
github.com/forge-cms/forge v1.0.9

## Next up (backlog)
- Markdown rendering — `{{ .Item.Body | markdown }}` in show templates
- Health check — Caddy `health_uri` workaround (HTTPS redirect issue)
- ADMIN_TOKEN — persist in `.env`

## Design decisions (locked)
- Font: JetBrains Mono via Bunny Fonts (not Google Fonts)
- Logo: [forge] med glødende o i hullet — mørk og lys variant
- Accent: #e8702f, Baggrund: #0a0b0d (#000000 i logo-filer)
- Favicon: [f] pixel-art 16×16 og 32×32, fuldt logo 180×180