# forge-cms.dev — CSS Architecture

This document defines the CSS conventions for forge-cms.dev.
It is the authoritative reference for all template and stylesheet work.
Copilot must read this before writing any HTML template or CSS.

**Format:** decisions are immutable once agreed. Revisions require a new entry
that supersedes the original.

---

## Decision index

| # | Title | Status | Date |
|---|-------|--------|------|
| C1 | No CSS framework or preprocessor | Locked | 2026-03-12 |
| C2 | Three-layer file structure | Locked | 2026-03-12 |
| C3 | BEM naming convention | Locked | 2026-03-12 |
| C4 | `.prose` wrapper for forge_markdown output | Locked | 2026-03-12 |
| C5 | Design tokens | Locked | 2026-03-12 |
| C6 | Dark mode as default | Locked | 2026-03-12 |
| C7 | Typography | Locked | 2026-03-12 |
| C8 | No JavaScript for layout or styling | Locked | 2026-03-12 |

---

## Decision C1 — No CSS framework or preprocessor

**Status:** Locked
**Date:** 2026-03-12

**Decision:** forge-cms.dev uses plain CSS only. No Tailwind, no Bootstrap,
no Sass, no PostCSS, no build step of any kind.

**Rationale:**
forge-cms.dev is itself a Forge application. The site must demonstrate that
a Forge application requires no frontend toolchain. Introducing Tailwind or a
build step would contradict the zero-dependency philosophy of the framework and
send a mixed message to developers evaluating Forge.

A Go developer reading the source of this site should be able to understand,
run, and modify every file without installing anything beyond the Go toolchain.

**Consequences:**
- All styles are written in standard CSS with custom properties
- No `package.json`, no `node_modules`, no build commands
- CSS is served as static files — no compilation step
- Browser support target: the two most recent versions of Chrome, Firefox,
  and Safari. No IE, no legacy Edge.

---

## Decision C2 — Three-layer file structure

**Status:** Locked
**Date:** 2026-03-12

**Decision:** CSS is split into exactly three files, loaded in this order:

```
static/css/tokens.css      — design tokens (custom properties only)
static/css/base.css        — element resets and global typography
static/css/components.css  — BEM components
```

No fourth file. No per-page stylesheets. No inline `<style>` blocks in
templates, except for the landing page hero (see note below).

**File responsibilities:**

`tokens.css` defines all CSS custom properties on `:root`. Nothing else.
No selectors other than `:root`. No rules. Tokens only.

`base.css` styles bare HTML elements: `html`, `body`, `h1`–`h6`, `p`, `a`,
`code`, `pre`, `ul`, `ol`, `li`, `blockquote`, `hr`, `img`. Also defines the
`.prose` component (Decision C4). No class selectors except `.prose`.

`components.css` contains all BEM components (Decision C3). Every reusable
UI pattern lives here. Order within the file: layout components first
(`.nav`, `.footer`, `.layout`), then content components (`.post`, `.doc`),
then interactive components (`.btn`, `.copy-btn`), then utility overrides last.

**Note on the landing page:**
The landing page (`/`) may use an additional `<style>` block in its template
for hero-specific styles that are not reused elsewhere. This is the only
permitted exception to the no-inline-styles rule.

**Rejected alternatives:**
- *Per-component files:* Adds filesystem complexity with no meaningful benefit
  at this scale. One file per layer is easier to navigate and sufficient.
- *Single file:* Conflates tokens, resets, and components. Makes it harder
  to reason about cascade order and specificity.

---

## Decision C3 — BEM naming convention

**Status:** Locked
**Date:** 2026-03-12

**Decision:** All class names follow BEM: Block, Element, Modifier.

```
.block
.block__element
.block--modifier
.block__element--modifier
```

**Rules:**

- Block names are short, lowercase, hyphenated nouns: `.nav`, `.post`,
  `.doc-page`, `.code-block`
- Element names describe their role within the block: `.post__title`,
  `.nav__link`, `.code-block__bar`
- Modifier names describe state or variant: `.nav__link--active`,
  `.btn--primary`, `.lc-step--active`
- No nesting beyond one level of element. `.post__title` is correct.
  `.post__body__paragraph` is not.
- No ID selectors in components.css. IDs are reserved for anchor targets only.
- No element selectors in components.css except inside `.prose` (Decision C4).

**Examples from the site:**

```css
/* Navigation */
.nav { }
.nav__logo { }
.nav__links { }
.nav__link { }
.nav__link--active { }
.nav__cta { }

/* Post card (devlog list) */
.post { }
.post__date { }
.post__title { }
.post__excerpt { }
.post__arrow { }

/* Lifecycle diagram */
.lc { }
.lc__step { }
.lc__step--active { }
.lc__name { }
.lc__desc { }

/* Button */
.btn { }
.btn--primary { }
.btn--ghost { }

/* Code block */
.code-block { }
.code-block__bar { }
.code-block__dots { }
.code-block__filename { }
.code-block__body { }
```

**Rejected alternatives:**
- *Utility classes (Tailwind-style):* See Decision C1.
- *Flat namespaced classes (`.forge-nav`, `.forge-post`):* Verbose and
  unnecessary for a single-site codebase.
- *Scoped CSS via data attributes:* Adds complexity without benefit.

---

## Decision C4 — `.prose` wrapper for forge_markdown output

**Status:** Locked
**Date:** 2026-03-12

**Decision:** All output from `{{forge_markdown}}` must be wrapped in a
`<div class="prose">` element. Element selectors inside `.prose` are the
only permitted element selectors in `base.css`.

**Rationale:**
`forge_markdown` generates raw HTML elements without class attributes:
`<p>`, `<h1>`–`<h6>`, `<strong>`, `<em>`, `<code>`, `<a>`, `<ul>`, `<li>`,
`<blockquote>`. These elements cannot be styled via BEM class selectors.

The `.prose` wrapper scopes element-level typography rules to markdown content
only, preventing them from leaking into the rest of the page.

**Template usage:**

```html
<div class="prose">
  {{forge_markdown .Content.Body}}
</div>
```

**CSS structure:**

```css
/* base.css */
.prose p        { margin: 0 0 1rem; }
.prose h2       { font-size: 1.25rem; margin: 2rem 0 0.5rem; }
.prose h3       { font-size: 1.1rem; margin: 1.5rem 0 0.5rem; }
.prose a        { color: var(--accent); }
.prose a:hover  { text-decoration: underline; }
.prose code     { font-family: var(--mono); font-size: 0.875em;
                  background: var(--bg-3); padding: 0.1em 0.3em;
                  border-radius: 2px; }
.prose ul       { padding-left: 1.5rem; margin-bottom: 1rem; }
.prose li       { margin-bottom: 0.25rem; }
.prose blockquote {
  border-left: 3px solid var(--border-2);
  padding-left: 1rem;
  color: var(--text-muted);
  margin: 1.5rem 0;
}
```

**Consequences:**
- Every template that renders markdown must wrap output in `.prose`
- `.prose` must never be used for non-markdown content
- Element selectors outside `.prose` are not permitted in any CSS file

---

## Decision C5 — Design tokens

**Status:** Locked
**Date:** 2026-03-12

**Decision:** All colours, font families, and spacing constants are defined
as CSS custom properties in `tokens.css`. No hardcoded values in
`base.css` or `components.css`.

**Full token list:**

```css
:root {
  /* Backgrounds */
  --bg:           #0a0b0d;   /* page background */
  --bg-2:         #111316;   /* card / code block background */
  --bg-3:         #181b1f;   /* elevated surface (code bar, nav) */

  /* Borders */
  --border:       #1f2329;   /* default border */
  --border-2:     #2a3038;   /* slightly more visible border */

  /* Text */
  --text:         #dde1e7;   /* primary text */
  --text-muted:   #6b7280;   /* secondary text, labels */
  --text-dim:     #3d4450;   /* decorative text, arrows, separators */

  /* Accent */
  --accent:       #e8702f;   /* primary accent — links, active states */
  --accent-dim:   #e8702f1a; /* accent background tint */
  --accent-border:#e8702f33; /* accent border tint */

  /* Semantic */
  --green:        #3fb950;   /* success, live status */

  /* Typography */
  --mono: 'JetBrains Mono', 'Fira Code', monospace;
  --sans: 'Inter', system-ui, sans-serif;

  /* Spacing scale */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;
  --space-16: 4rem;
  --space-20: 5rem;

  /* Layout */
  --max-width: 900px;
  --nav-height: 52px;
}
```

**Rules:**
- Token names are always referenced via `var(--name)` — never hardcoded hex
- Adding a new token requires updating this document
- Tokens are never defined outside `tokens.css`

---

## Decision C6 — Dark mode as default

**Status:** Locked
**Date:** 2026-03-12

**Decision:** The site is dark mode only. There is no light mode variant.
No `@media (prefers-color-scheme)` query is implemented.

**Rationale:**
Dark mode is the natural default for a developer-facing technical site.
Implementing both modes doubles the CSS surface area and creates maintenance
overhead. The target audience — Go developers — overwhelmingly uses dark
terminals, dark editors, and dark browsers.

**Consequences:**
- All token values assume a dark background
- No `prefers-color-scheme` media queries
- The `color-scheme: dark` meta tag is set in the base layout template

---

## Decision C7 — Typography

**Status:** Locked
**Date:** 2026-03-12

**Decision:** The site uses a mono-first typographic system. The primary
font is JetBrains Mono (loaded via Google Fonts). Inter is used for
prose body text only (`.prose`, `.hero-sub`, `.feat__val`).

**Rationale:**
The entire UI — navigation, labels, buttons, headings, code — renders in
monospace. This is a deliberate identity decision: the site reads like source
code. It signals that Forge is a developer tool, not a marketing product.

Inter is permitted for longer reading text because monospace body text at
small sizes reduces readability. The boundary is clear: UI chrome is mono,
reading content is sans.

**Type scale:**

```
0.68rem  — badge, fine print
0.70rem  — section labels, meta, timestamps
0.75rem  — secondary labels, code filenames
0.80rem  — primary UI text, buttons, code body
0.85rem  — post titles in lists, feature values
0.90rem  — nav logo
1.00rem  — hero sub (Inter)
hero h1  — clamp(1.8rem, 4vw, 3.2rem)
```

**Font loading:**

```html
<link rel="preconnect" href="https://fonts.bunny.net">
<link rel="stylesheet" href="https://fonts.bunny.net/css?family=jetbrains-mono:400,500,600|inter:300,400&display=swap">
```

**Consequences:**
- `font-family: var(--mono)` is set on `body` — everything inherits mono by default
- `font-family: var(--sans)` is applied explicitly only where prose is needed
- No third font is introduced without amending this decision

---

## Decision C8 — No JavaScript for layout or styling

**Status:** Locked
**Date:** 2026-03-12

**Decision:** JavaScript is not used for layout, animation, or styling.
CSS handles all visual state. JavaScript is permitted only for discrete
interactive behaviours: copy-to-clipboard, mobile navigation toggle.

**Rationale:**
Consistent with Decision C1. A Forge site demonstrates that a useful,
well-designed web experience does not require a JavaScript framework.

**Permitted JavaScript:**
- Copy-to-clipboard for the install command (`navigator.clipboard.writeText`)
- Mobile navigation toggle (add/remove a CSS class on `<nav>`)
- Nothing else without amending this decision

**Consequences:**
- No Alpine.js, htmx, or any other JS library
- Animations are CSS-only (`@keyframes`, `transition`, `animation`)
- No JavaScript files are served by the site except an optional
  single inline `<script>` in the base layout for the nav toggle

---

## Template conventions

These are not decisions — they are implementation rules that follow from the
decisions above. Copilot must follow them without being asked.

### Base layout structure

Every page template inherits from a base layout that provides:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="color-scheme" content="dark">
  {{template "forge:head" .Head}}
  <link rel="preconnect" href="https://fonts.bunny.net">
  <link rel="stylesheet" href="https://fonts.bunny.net/css?family=jetbrains-mono:400,500,600|inter:300,400&display=swap">
  <link rel="stylesheet" href="/static/css/tokens.css">
  <link rel="stylesheet" href="/static/css/base.css">
  <link rel="stylesheet" href="/static/css/components.css">
</head>
<body>
  {{template "nav" .}}
  <main class="layout">
    {{block "content" .}}{{end}}
  </main>
  {{template "footer" .}}
</body>
</html>
```

### Markdown rendering

Always wrap `forge_markdown` output in `.prose`:

```html
<div class="prose">{{forge_markdown .Content.Body}}</div>
```

### Section labels

Section headers use the `// label` convention in mono:

```html
<p class="section__label">// features</p>
```

### Dates

Always use `forge_date` for human-readable dates:

```html
<time class="post__date">{{.Content.PublishedAt | forge_date}}</time>
```

### Active nav links

The active nav link receives the `.nav__link--active` modifier.
This is set server-side by comparing `r.URL.Path` with the link target —
not via JavaScript.
