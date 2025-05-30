//go:build lambda
// +build lambda

package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/binxly/gofolio/blog"
	"github.com/gorilla/mux"
)

func main() {
	blogService = blog.NewBlogService("blog")
	if err := blogService.LoadPosts(); err != nil {
		// TODO: same as main.go
		log.Printf("Error loading blog posts: %v", err)
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/blog", blogHandler).Methods("GET")
	r.HandleFunc("/blog/{slug}", blogPostHandler).Methods("GET")
	r.HandleFunc("/about", aboutHandler).Methods("GET")
	r.HandleFunc("/rss", rssHandler).Methods("GET")

	startLambda(r)
}

func startLambda(r *mux.Router) {
	algnhsa.ListenAndServe(r, nil)
}
