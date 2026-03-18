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
- S13: Docker user 1000:1000 to fix volume permission errors
- S14: forge upgraded to v1.0.10
- S15: forge upgraded to v1.0.11
- S16: OG + Twitter Card tags added to devlog/docs show templates (S8 workaround extended)
- S17: twitter:card override attempted (first-tag-wins limitation discovered)
- S18: forge v1.1.1 — og:url/og:image/twitter:card now correct in forge:head; S16+S17 workarounds removed
- S19: siteBaseURL package var — Head() Canonical now absolute; Head.Image set to Forge-logo-OG1200.png (1200×630); og:url and og:image now present
- S20: forge-mcp v1.0.0 wired — MCP(MCPRead, MCPWrite) on Post+DocPage; GET /mcp + POST /mcp/message mounted; admin Bearer token logged at startup
- S21: cmd/mcp/main.go — stdlib stdio-to-SSE proxy for Claude Desktop; builds to forge-mcp-proxy.exe; token via MCP_TOKEN env var
- S22: forge v1.1.2 + forge-mcp v1.0.1 — array-aware MCP tags fixed (A52); no forge-site code changes
- S23: Remove internal workaround/amendment comments from all four module templates + main.go; no functional changes
- S24: BlogPosting JSON-LD in devlog/show.html + TechArticle JSON-LD in docs/show.html; forge_rfc3339 for datePublished
- S25: Add required image field to BlogPosting JSON-LD (Google rich results eligibility)
- S26: forge v1.1.3 — A53 content negotiation fix; crawlers now receive HTML not JSON
- S27: Extend README with API/MCP/AI/boilerplate sections; clarify SESSION_CONTEXT.md note in copilot-instructions
- S28: Fix horizontal scroll on ~360px mobile — overflow-x hidden on html/body, feat-grid minmax fix, code-block max-width 100%
- S29: forge-mcp v1.0.2 — admin read tools exposed automatically via existing forgemcp.New(app) wiring

## Open corepilot items (sent)
- A43: NewSQLRepo pointer type docs ✅ sent
- Phase 2 promoted: shared partials, forge:head helper, MustConfig enforcement ✅ sent
- OGDefaults + AppSchema → Phase 2 ✅ sent
- AppSchema + OGDefaults backlog entries ✅ sent

## Forge version in use
github.com/forge-cms/forge v1.1.3
forge-mcp v1.0.2

## Next up (backlog)
- Markdown rendering — `{{ .Item.Body | markdown }}` in show templates
- Health check — Caddy `health_uri` workaround (HTTPS redirect issue)
- ADMIN_TOKEN — persist in `.env`

## Design decisions (locked)
- Font: JetBrains Mono via Bunny Fonts (not Google Fonts)
- Logo: [forge] med glødende o i hullet — mørk og lys variant
- Accent: #e8702f, Baggrund: #0a0b0d (#000000 i logo-filer)
- Favicon: [f] pixel-art 16×16 og 32×32, fuldt logo 180×180