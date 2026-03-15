package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/forge-cms/forge"
	_ "modernc.org/sqlite"
)

//go:embed templates
var templates embed.FS

//go:embed static
var static embed.FS

// Version is set at build time via -ldflags "-X main.Version=x.y.z".
var Version string

// homeData is the template data passed to templates/home/home.html.
// It mirrors forge.TemplateData but without generic type constraints so
// the home handler can be wired with a plain http.Handler.
type homeData struct {
	Head     forge.Head
	Request  *http.Request
	SiteName string
	Posts    []*Post
}

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

	if err := migrateDB(db); err != nil {
		log.Fatalf("forge-site: migrate db: %v", err)
	}

	postRepo := forge.NewSQLRepo[*Post](db)
	docRepo := forge.NewSQLRepo[*DocPage](db)
	seedDB(context.Background(), postRepo, docRepo)

	// Parse the home page template once at startup.
	// base.html + home.html are embedded in the binary via go:embed.
	// forge:head is NOT used here — forgeHeadTmpl is package-private in
	// the forge package. The home handler constructs <head> tags manually
	// from homeData.Head fields. This is intentional; see Amendment S5.
	homeTmpl, err := template.New("").Funcs(forge.TemplateFuncMap()).
		ParseFS(templates, "templates/base.html", "templates/home/home.html")
	if err != nil {
		log.Fatalf("forge-site: parse home template: %v", err)
	}

	app := forge.New(forge.Config{
		BaseURL: baseURL,
		Secret:  []byte(secret),
		Version: Version,
		DB:      db,
		HTTPS:   strings.HasPrefix(baseURL, "https"),
	})

	app.Use(forge.Authenticate(forge.BearerHMAC(secret)))

	app.Health()

	// Override MIME type for CSS files — Windows registry maps .css to
	// text/plain which causes browsers to reject stylesheets. This is a
	// no-op on Linux/Docker but required for local development on Windows.
	_ = mime.AddExtensionType(".css", "text/css")

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("forge-site: static sub: %v", err)
	}
	app.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	app.SEO(&forge.RobotsConfig{Sitemaps: true})

	app.Content(forge.NewModule((*Post)(nil),
		forge.Repo(postRepo),
		forge.At("/devlog"),
		forge.Templates("templates/devlog"),
		forge.SitemapConfig{},
		forge.Social(forge.OpenGraph, forge.TwitterCard),
		forge.Feed(forge.FeedConfig{Title: "Forge Devlog"}),
		forge.AIIndex(forge.LLMsTxt, forge.LLMsTxtFull, forge.AIDoc),
		forge.Cache(5*time.Minute),
		forge.HeadFunc(func(_ forge.Context, _ []*Post) forge.Head {
			return forge.Head{
				Title:       "Devlog — Forge",
				Description: "Engineering notes and release announcements from the Forge team.",
				Canonical:   forge.URL("/devlog"),
			}
		}),
	))

	app.Content(forge.NewModule((*DocPage)(nil),
		forge.Repo(docRepo),
		forge.At("/docs"),
		forge.Templates("templates/docs"),
		forge.SitemapConfig{},
		forge.AIIndex(forge.LLMsTxt, forge.LLMsTxtFull, forge.AIDoc),
		forge.Cache(10*time.Minute),
		forge.HeadFunc(func(_ forge.Context, _ []*DocPage) forge.Head {
			return forge.Head{
				Title:       "Docs — Forge",
				Description: "Documentation for the Forge Go web framework.",
				Canonical:   forge.URL("/docs"),
			}
		}),
	))

	hostname := strings.TrimPrefix(strings.TrimPrefix(baseURL, "https://"), "http://")
	app.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		all, err := postRepo.FindAll(r.Context(), forge.ListOptions{})
		var recent []*Post
		if err == nil {
			for _, p := range all {
				if p.Status == forge.Published {
					recent = append(recent, p)
					if len(recent) == 3 {
						break
					}
				}
			}
		}
		data := homeData{
			Head: forge.Head{
				Title:       "Forge — The Go web framework built for the age of AI",
				Description: "Forge is a Go web framework designed for developers, AI builders, human visitors, and AI agents consuming content.",
				Canonical:   baseURL + "/",
			},
			Request:  r,
			SiteName: hostname,
			Posts:    recent,
		}
		var buf bytes.Buffer
		if err := homeTmpl.ExecuteTemplate(&buf, "base", data); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Printf("forge-site: home template: %v", err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
	}))

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
