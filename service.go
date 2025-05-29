package main

import "github.com/binxly/gofolio/blog"

// blogService is a shared instance used across both lambda and non-lambda builds
var blogService *blog.BlogService
