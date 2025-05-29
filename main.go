//go:build !lambda
// +build !lambda

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binxly/gofolio/blog"
	"github.com/gorilla/mux"
)

func main() {
	blogService = blog.NewBlogService("content")
	if err := blogService.LoadPosts(); err != nil {
		// TODO: replace with a method for loading metadata of 'first' 3 posts only.
		// very ugly and hacked together because I got lazy when I opted for
		// a goldmark extension to help with .md frontmatter parsing.
		// ref lambda.go
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
