package main

import (
	"net/http"

	"github.com/binxly/gofolio/templates/layouts"
	"github.com/binxly/gofolio/templates/pages"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	recentPosts := blogService.GetRecentPosts(3) // TODO: move this value to a config file or something idk
	component := layouts.BaseLayout("Home", pages.HomePage(recentPosts))
	component.Render(r.Context(), w)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	posts := blogService.GetAllPosts()
	component := pages.BlogPage(posts)
	component.Render(r.Context(), w)
}

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	post, err := blogService.GetPostBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	component := pages.BlogPostPage(*post)
	component.Render(r.Context(), w)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.AboutPage()
	component.Render(r.Context(), w)
}
