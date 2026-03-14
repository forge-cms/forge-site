# forge — Backlog

Issues and improvement ideas discovered during forge-site development.
These are upstream concerns for the forge framework itself, not this site.

---

## v2+ Roadmap

- **`forge.New` requires `MustConfig`** — `forge.New(forge.Config{...})` without `MustConfig` silently accepts invalid config (empty BaseURL, short Secret); consider making `New` call `MustConfig` internally so validation is not opt-in; discovered during forge-site templates sprint
