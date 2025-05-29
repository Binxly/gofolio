package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func rssHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://binx.page"

	feed := &feeds.Feed{
		Title:       "Binx Bytes",
		Link:        &feeds.Link{Href: baseURL + "/blog"},
		Description: "Thoughts on technology, learning, and building things.",
		Author:      &feeds.Author{Name: "Zac Bagley"},
		Created:     time.Now(),
		Updated:     time.Now(),
		Image: &feeds.Image{
			Url:   baseURL + "/static/favicon.ico",
			Title: "Binx Bytes",
			Link:  baseURL,
		},
	}

	posts := blogService.GetAllPosts()

	for _, post := range posts {
		item := &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/blog/%s", baseURL, post.Slug)},
			Description: post.Excerpt,
			Author:      feed.Author,
			Created:     post.Date,
			Updated:     post.Date,
			Id:          fmt.Sprintf("%s/blog/%s", baseURL, post.Slug),
			Content:     string(post.Content),
		}

		if post.Category != "" {
			item.Description = fmt.Sprintf("Category: %s\n\n%s", post.Category, post.Excerpt)
		}

		feed.Items = append(feed.Items, item)
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")

	rss, err := feed.ToRss()
	if err != nil {
		http.Error(w, "Failed to generate RSS feed", http.StatusInternalServerError)
		log.Printf("Error generating RSS feed: %v", err)
		return
	}

	_, err = w.Write([]byte(rss))
	if err != nil {
		log.Printf("Error writing RSS feed to response: %v", err)
	}
}
