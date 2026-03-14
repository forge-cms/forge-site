package main

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/forge-cms/forge"
	_ "modernc.org/sqlite"
)

//go:embed all:templates
var templates embed.FS

//go:embed static
var static embed.FS

// Version is set at build time via -ldflags "-X main.Version=x.y.z".
var Version string

func main() {
	secret := requireEnv("SECRET")
	baseURL := envOr("BASE_URL", "http://localhost:8080")
	dbPath := envOr("DATABASE_PATH", "./data/forge.db")
	port := envOr("PORT", "8080")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("forge-site: create data dir: %v", err)
	}
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("forge-site: open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("forge-site: ping db: %v", err)
	}
	defer db.Close()

	app := forge.New(forge.Config{
		BaseURL: baseURL,
		Secret:  []byte(secret),
		Version: Version,
		DB:      db,
		HTTPS:   strings.HasPrefix(baseURL, "https"),
	})

	app.Health()
	app.SEO(&forge.RobotsConfig{Sitemaps: true})

	// TODO(content-types): register Post module at /devlog
	// TODO(content-types): register DocPage module at /docs
	// TODO(templates): register home page handler at /

	maybeLogAdminToken(secret)

	addr := ":" + port
	log.Printf("forge-site: listening on %s", addr)
	if err := app.Run(addr); err != nil {
		log.Fatalf("forge-site: %v", err)
	}
}

// maybeLogAdminToken generates and logs a non-expiring admin bearer token when
// ADMIN_TOKEN is unset. This lets the operator seed content via post.http on a
// fresh deployment without a separate token-generation step.
func maybeLogAdminToken(secret string) {
	if os.Getenv("ADMIN_TOKEN") != "" {
		return
	}
	admin := forge.User{
		ID:    "admin",
		Name:  "Admin",
		Roles: []forge.Role{forge.Admin},
	}
	token, err := forge.SignToken(admin, secret, 0)
	if err != nil {
		log.Fatalf("forge-site: generate admin token: %v", err)
	}
	log.Printf("ADMIN_TOKEN not set — generated token (no expiry):\n  %s", token)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("forge-site: required env var %s is not set", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
