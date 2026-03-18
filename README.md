# forge-site

The source for [forge-cms.dev](https://forge-cms.dev) — built as a Forge application.

---

## Local development

**Prerequisites:** Go 1.26+, no other tools required.

```powershell
$env:SECRET = "dev-secret-change-me"   # required — minimum 16 bytes
$env:BASE_URL = "http://localhost:8080" # optional, default shown
$env:DATABASE_PATH = "./data/forge.db"  # optional, default shown
go run .
```

The server starts on `:8080`. SQLite database is created automatically on first run.

---

## Admin token

On startup, if `ADMIN_TOKEN` is not set, the app generates and logs a non-expiring
admin bearer token:

```
forge-site: ADMIN_TOKEN not set — generated token (no expiry):
  eyJ...
```

Copy that token into `post.http`:

```
@token = eyJ...
```

Then use the VS Code REST Client extension to send requests from `post.http`.
The file walks through: health check → create Post → publish Post → create DocPage → verify public visibility.

---

## Docker deploy

1. **Edit `Caddyfile`** — replace `DOMAIN_PLACEHOLDER` with your domain (e.g. `forge-cms.dev`).

2. **Create a `.env` file** on the server (never commit this):

   ```
   SECRET=<random-32+-byte-string>
   BASE_URL=https://forge-cms.dev
   ```

3. **Build and start:**

   ```bash
   docker compose build --build-arg VERSION=1.0.0
   docker compose up -d
   ```

   The app binds to `127.0.0.1:8080`. Caddy runs on the host and proxies to it.
   TLS certificates are provisioned automatically via ACME.

4. **Check health:**

   ```bash
   curl https://forge-cms.dev/_health
   # {"status":"ok","version":"1.0.0"}
   ```

---

## Token rotation

Admin tokens have no expiry (`TTL=0`). If a token is compromised:

1. Change `SECRET` in `.env` to a new random value.
2. Restart the app (`docker compose up -d`).

All previously issued tokens are immediately invalidated — they are signed with
the old secret and will fail HMAC verification.

---

## Content types

Two content types power this site:

### Post — `/devlog`

Engineering notes and release announcements.
```go
type Post struct {
    forge.Node
    Title string   `forge:"required,min=3"`
    Body  string   `forge:"required"`
    Tags  []string
}
```

Implements `Headable`, `Markdownable`, `AIDocSummary`.

### DocPage — `/docs`

Framework documentation pages.
```go
type DocPage struct {
    forge.Node
    Title   string `forge:"required,min=3"`
    Body    string `forge:"required"`
    Section string
    Order   int    `db:"sort_order"`
}
```

Implements `Headable`, `Markdownable`.

---

## REST API

Forge auto-generates these endpoints for each module.
All write operations require a Bearer token with sufficient role.

### Post — /devlog

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /devlog | — | List published posts |
| GET | /devlog/{slug} | — | Get single post |
| POST | /devlog | Author+ | Create post |
| PUT | /devlog/{slug} | Author+ | Update post |
| DELETE | /devlog/{slug} | Editor+ | Delete post |

### DocPage — /docs

Same pattern at `/docs`.

**Auth:** `Authorization: Bearer <token>`
Generate a token: the app logs an admin token (no expiry) at startup.
See `post.http` for ready-to-run HTTP examples covering the full
create → publish → verify lifecycle.

---

## MCP

Two endpoints for AI assistant integration:

| Method | Path | Description |
|--------|------|-------------|
| GET | /mcp | SSE connection |
| POST | /mcp/message | JSON-RPC messages |

Both Post and DocPage support read and write via MCP.

**Claude Desktop:** build `cmd/mcp/main.go` as a stdio-to-SSE proxy:
```powershell
go build -o forge-mcp-proxy.exe ./cmd/mcp/
```

Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "forge-cms": {
      "command": "C:\\path\\to\\forge-mcp-proxy.exe",
      "env": { "MCP_TOKEN": "<bearer token from startup log>" }
    }
  }
}
```

**MCP Bearer token** is logged at startup:
```bash
docker logs <container> 2>&1 | grep "MCP Bearer token"
```

---

## AI endpoints

Forge auto-generates these for AI crawlers and assistants:

| Path | Description |
|------|-------------|
| /llms.txt | Compact content index |
| /llms-full.txt | Full markdown corpus |
| /devlog/{slug}/aidoc | Per-post AI summary |
| /docs/{slug}/aidoc | Per-doc AI summary |
| /sitemap.xml | Aggregate sitemap |
| /devlog/feed.xml | RSS feed |

---

## Using this as a boilerplate

forge-site is the canonical starting point for a Forge application.

1. Clone this repo
2. Rename the module in `go.mod`
3. Replace `Post` and `DocPage` with your own content types
4. Update `schema.go` with your table definitions
5. Update `templates/` for your content types
6. Remove or replace `seed.go`

The content lifecycle (Draft → Scheduled → Published → Archived),
sitemap, RSS feed, llms.txt, and MCP endpoints are all automatic once
your types embed `forge.Node` and implement the required interfaces.

See [github.com/forge-cms/forge](https://github.com/forge-cms/forge)
for the full framework documentation.
