package main

import (
	"github.com/forge-cms/forge"
)

// Post is the devlog content type. Mounted at /devlog.
// Fields match D1. Storage decisions in S4.
type Post struct {
	forge.Node
	Title string          `forge:"required,min=3" db:"title"`
	Body  string          `forge:"required"       db:"body"`
	Tags  JSONStringSlice `db:"tags"`
}

func (p *Post) Head() forge.Head {
	return forge.Head{
		Title:       p.Title,
		Description: forge.Excerpt(p.Body, 160),
		Tags:        []string(p.Tags),
		Type:        forge.Article,
		Canonical:   siteBaseURL + forge.URL("/devlog/", p.Slug),
		Image:       forge.Image{URL: siteBaseURL + "/static/Forge-logo-OG1200.png", Alt: "Forge", Width: 1200, Height: 630},
		Breadcrumbs: forge.Crumbs(
			forge.Crumb("Home", "/"),
			forge.Crumb("Devlog", "/devlog"),
			forge.Crumb(p.Title, "/devlog/"+p.Slug),
		),
	}
}

func (p *Post) Markdown() string { return p.Body }

func (p *Post) AISummary() string { return forge.Excerpt(p.Body, 120) }
