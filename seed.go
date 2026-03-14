package main

import (
	"context"
	"log"
	"time"

	"github.com/forge-cms/forge"
)

// seedDB inserts a single draft Post and a single draft DocPage if both
// tables are empty. This gives a minimal working dataset for local dev
// without polluting production deployments that already have content.
func seedDB(
	ctx context.Context,
	postRepo forge.Repository[*Post],
	docRepo forge.Repository[*DocPage],
) {
	posts, err := postRepo.FindAll(ctx, forge.ListOptions{})
	if err != nil {
		log.Printf("forge-site: seedDB: check posts: %v", err)
		return
	}
	docs, err := docRepo.FindAll(ctx, forge.ListOptions{})
	if err != nil {
		log.Printf("forge-site: seedDB: check doc_pages: %v", err)
		return
	}
	if len(posts) > 0 || len(docs) > 0 {
		return // content already exists — skip seeding
	}

	now := time.Now().UTC()

	post := &Post{
		Node: forge.Node{
			ID:        forge.NewID(),
			Slug:      "hello-forge",
			Status:    forge.Draft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title: "Hello, Forge",
		Body:  "This is the first devlog entry. Replace with real content before launch.",
		Tags:  JSONStringSlice{"forge", "go"},
	}
	if err := postRepo.Save(ctx, post); err != nil {
		log.Printf("forge-site: seedDB: save post: %v", err)
		return
	}

	doc := &DocPage{
		Node: forge.Node{
			ID:        forge.NewID(),
			Slug:      "getting-started",
			Status:    forge.Draft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:   "Getting Started",
		Body:    "Welcome to Forge. This guide covers installation and your first application.",
		Section: "introduction",
		Order:   1,
	}
	if err := docRepo.Save(ctx, doc); err != nil {
		log.Printf("forge-site: seedDB: save doc: %v", err)
		return
	}

	log.Println("forge-site: dev seed inserted (1 Post + 1 DocPage, both Draft)")
}
