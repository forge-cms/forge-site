package main

import (
	"context"
	"sort"

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

// sortedPostRepo wraps SQLRepo[*Post] and overrides FindAll to return posts
// sorted newest-first: published posts by PublishedAt DESC, then unpublished
// (draft/scheduled/archived) by CreatedAt DESC.
type sortedPostRepo struct {
	*forge.SQLRepo[*Post]
}

func (r *sortedPostRepo) FindAll(ctx context.Context, opts forge.ListOptions) ([]*Post, error) {
	posts, err := r.SQLRepo.FindAll(ctx, opts)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(posts, func(i, j int) bool {
		pi, pj := posts[i], posts[j]
		iPublished := !pi.PublishedAt.IsZero()
		jPublished := !pj.PublishedAt.IsZero()
		if iPublished != jPublished {
			return iPublished // published entries before unpublished
		}
		if iPublished {
			return pi.PublishedAt.After(pj.PublishedAt) // newest published first
		}
		return pi.CreatedAt.After(pj.CreatedAt) // newest created first for unpublished
	})
	return posts, nil
}
