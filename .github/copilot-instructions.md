# forge-site — Copilot Instructions

This is **forge-cms.dev** — the official website for the Forge Go web framework,
built as a Forge application. Read this file before writing any code, template,
CSS, or configuration.

---

## Forge framework rules (non-negotiable)

- `forge:head` owns all `<head>` SEO/social meta — never duplicate it
- `forge_markdown` output must always be wrapped in `<div class="prose">`
- Lifecycle is enforced by Forge — never filter by Status in templates
- Templates are always `list.html` + `show.html` in a named directory
- Never add a Go dependency outside stdlib or `forge`

---

## Content types

### Post — `/devlog`

```go
type Post struct {
    forge.Node
    Title  string   `forge:"required,min=3"`
    Body   string   `forge:"required"`
    Tags   []string
}

func (p *Post) Head() forge.Head {
    return forge.Head{
        Title:       p.Title,
        Description: forge.Excerpt(p.Body, 160),
        Tags:        p.Tags,
        Type:        forge.Article,
        Canonical:   forge.URL("/devlog/", p.Slug),
        Breadcrumbs: forge.Crumbs(
            forge.Crumb("Home",   "/"),
            forge.Crumb("Devlog", "/devlog"),
            forge.Crumb(p.Title,  "/devlog/"+p.Slug),
        ),
    }
}
func (p *Post) Markdown() string  { return p.Body }
func (p *Post) AISummary() string { return forge.Excerpt(p.Body, 120) }
```

Module options: `At("/devlog")`, `Templates("templates/devlog")`,
`SitemapConfig{}`, `Social(OpenGraph, TwitterCard)`,
`Feed(FeedConfig{Title: "Forge Devlog"})`,
`AIIndex(LLMsTxt, LLMsTxtFull, AIDoc)`

### DocPage — `/docs`

```go
type DocPage struct {
    forge.Node
    Title   string `forge:"required,min=3"`
    Body    string `forge:"required"`
    Section string
    Order   int
}

func (d *DocPage) Head() forge.Head {
    return forge.Head{
        Title:       d.Title + " — Forge Docs",
        Description: forge.Excerpt(d.Body, 160),
        Type:        forge.Article,
        Canonical:   forge.URL("/docs/", d.Slug),
        Breadcrumbs: forge.Crumbs(
            forge.Crumb("Home",  "/"),
            forge.Crumb("Docs",  "/docs"),
            forge.Crumb(d.Title, "/docs/"+d.Slug),
        ),
    }
}
func (d *DocPage) Markdown() string { return d.Body }
```

Module options: `At("/docs")`, `Templates("templates/docs")`,
`SitemapConfig{}`, `AIIndex(LLMsTxt, LLMsTxtFull, AIDoc)`

---

## URL structure

| Route | Template | Note |
|-------|----------|------|
| `/` | `templates/home/home.html` | Static `forge.Handle` — no `TemplateData[T]` |
| `/devlog` | `templates/devlog/list.html` | |
| `/devlog/{slug}` | `templates/devlog/show.html` | |
| `/docs` | `templates/docs/list.html` | |
| `/docs/{slug}` | `templates/docs/show.html` | |
| `/devlog/feed.xml` | auto | Forge |
| `/llms.txt` | auto | Forge |
| `/sitemap.xml` | auto | Forge |

Never change URL prefixes without an Amendment.

---

## Project structure

```
forge-site/
├── main.go
├── go.mod / go.sum
├── DECISIONS.md
├── DESIGN.md
├── GOVERNANCE.md
├── TODO.md
├── README.md
├── Dockerfile
├── docker-compose.yml
├── Caddyfile
├── templates/
│   ├── base.html
│   ├── home/home.html
│   ├── devlog/list.html + show.html
│   └── docs/list.html + show.html
├── static/css/
│   ├── tokens.css
│   ├── base.css
│   └── components.css
└── data/forge.db          # runtime only, gitignored
```

---

## CSS rules

Three files, strict separation — read `DESIGN.md` for full spec.

- `tokens.css` — CSS variables only
- `base.css` — element resets + `.prose`
- `components.css` — all BEM components

BEM: `.block`, `.block__element`, `.block--modifier`

Known components (do not invent new names):
```
.nav .nav__logo .nav__logo-bracket .nav__links .nav__link .nav__link--active .nav__cta
.layout
.section .section__label
.footer .footer__logo .footer__logo-bracket .footer__links .footer__link
.hero .hero__heading .hero__heading-accent .hero__sub .hero__actions .hero__meta
.code-block .code-block__bar .code-block__dots .code-block__dot .code-block__filename .code-block__body
.feat-grid .feat .feat__key .feat__val
.lc .lc__step .lc__step--active .lc__name .lc__desc
.post .post__date .post__title .post__arrow
.post-list
.doc-page .doc-page__sidebar .doc-page__content .doc-page__title .doc-page__meta
.doc-nav .doc-nav__section .doc-nav__section-label .doc-nav__link .doc-nav__link--active
.install-block .install-block__prompt .install-block__cmd
.btn .btn--primary .btn--ghost
.copy-btn
.badge
```

Typography: `var(--mono)` on everything. `var(--sans)` only on `.prose`, `.hero__sub`, `.feat__val`.

Patterns:
```html
<div class="prose">{{forge_markdown .Content.Body}}</div>
<p class="section__label">// section name</p>
{{if eq .Request.URL.Path "/devlog"}}
<a href="/devlog" class="nav__link nav__link--active">devlog</a>
{{else}}
<a href="/devlog" class="nav__link">devlog</a>
{{end}}
```

No inline styles except hero. No JS for layout or styling.

---

## Governance

Amendment required when: adding a content type field, changing a URL prefix,
adding JS, changing CSS file structure. See `GOVERNANCE.md`.

**Copilot must self-identify amendments.** Before reporting task completion,
scan what was implemented and ask: does this contradict or extend an existing
decision in `DECISIONS.md`? If yes, propose the Amendment text before proposing
a commit. Do not wait for the user to notice.

Examples of things that trigger a self-identified amendment:
- A version, tool, or dependency that differs from what `DECISIONS.md` states
- A build flag, env var, or configuration choice not covered by an existing decision
- A structural choice (file layout, naming, port) that extends a locked decision

Commit format:
```
{type}({scope}): {short description}
Decisions: {IDs or "none"}
```
Types: `feat`, `fix`, `content`, `style`, `docs`, `chore`

**TODO.md must be updated as part of every task completion.**

When a backlog item is fully done and committed:
- Move it from its current section to `## Done` with today's date
- Move the next Backlog item to `## Up next` if In progress drops below 2

This update is included in the same commit as the work it tracks — never a separate commit.

**SESSION_CONTEXT.md tracks operational state. Update Current status,
Latest amendments, and Next up to reflect the state after each commit.**
This update is included in the same commit as the work it tracks.

**Never commit without explicit user approval.**
Run autonomously: `go build ./...`, `go vet ./...`, `gofmt -l .`, any read-only git command.

---

## Never

- Add npm, pip, or non-stdlib Go dependency
- Use any CSS/JS framework
- Implement sitemap, robots.txt, llms.txt, feed, or AIDoc manually
- Filter by Status in templates
- Duplicate `forge:head` output
- Use `forge_markdown` without `.prose` wrapper
- Add JS for layout/animation/styling
- Change URL prefixes or content type fields without an Amendment
- Commit without approval
