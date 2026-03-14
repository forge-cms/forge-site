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
