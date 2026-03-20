package main

import (
	"github.com/forge-cms/forge"
)

// DocPage is the documentation content type. Mounted at /docs.
// Fields match D1. Storage decisions in S4.
type DocPage struct {
	forge.Node
	Title   string `forge:"required,min=3" db:"title"`
	Body    string `forge:"required"       db:"body"`
	Section string `db:"section"`
	Order   int    `db:"sort_order"`
}

func (d *DocPage) Head() forge.Head {
	return forge.Head{
		Title:       d.Title + " — Forge Docs",
		Description: forge.Excerpt(d.Body, 160),
		Type:        forge.Article,
		Canonical:   forge.AbsURL(siteBaseURL, forge.URL("/docs/", d.Slug)),
		Image:       forge.Image{URL: forge.AbsURL(siteBaseURL, "/static/Forge-logo-OG1200.png"), Alt: "Forge", Width: 1200, Height: 630},
		Breadcrumbs: forge.Crumbs(
			forge.Crumb("Home", "/"),
			forge.Crumb("Docs", "/docs"),
			forge.Crumb(d.Title, "/docs/"+d.Slug),
		),
	}
}

func (d *DocPage) Markdown() string { return d.Body }
