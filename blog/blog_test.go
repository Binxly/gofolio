package blog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBlogService(t *testing.T) {
	bs := NewBlogService("test-content")
	if bs.contentDir != "test-content" {
		t.Errorf("Expected contentDir to be 'test-content', got %s", bs.contentDir)
	}
	if bs.posts == nil {
		t.Error("Expected posts slice to be initialized")
	}
}

func TestParsePost(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-post.md")

	testContent := `---
title: "Test Post"
date: "2024-01-15"
category: "testing"
excerpt: "This is a test post"
slug: "test-post"
---

# Test Content

This is the body of the test post.`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	bs := NewBlogService(tempDir)
	post, err := bs.parsePost(testFile)
	if err != nil {
		t.Fatalf("Failed to parse post: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %s", post.Title)
	}

	expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
	if !post.Date.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, post.Date)
	}

	if post.Category != "testing" {
		t.Errorf("Expected category 'testing', got %s", post.Category)
	}

	if post.Excerpt != "This is a test post" {
		t.Errorf("Expected excerpt 'This is a test post', got %s", post.Excerpt)
	}

	if post.Slug != "test-post" {
		t.Errorf("Expected slug 'test-post', got %s", post.Slug)
	}

	if len(post.Content) == 0 {
		t.Error("Expected content to be rendered")
	}
}

func TestParsePostWithMissingMetadata(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "minimal-post.md")

	testContent := `---
title: "Minimal Post"
---

# Minimal Content`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	bs := NewBlogService(tempDir)
	post, err := bs.parsePost(testFile)
	if err != nil {
		t.Fatalf("Failed to parse post: %v", err)
	}

	if post.Title != "Minimal Post" {
		t.Errorf("Expected title 'Minimal Post', got %s", post.Title)
	}

	if post.Slug != "minimal-post" {
		t.Errorf("Expected slug 'minimal-post', got %s", post.Slug)
	}
}

func TestParsePostInvalidFile(t *testing.T) {
	bs := NewBlogService("test-content")
	_, err := bs.parsePost("nonexistent-file.md")
	if err == nil {
		t.Error("Expected error when parsing nonexistent file")
	}
}

func TestLoadPosts(t *testing.T) {
	// Create temporary test directory structure
	tempDir := t.TempDir()
	blogDir := filepath.Join(tempDir, "blog")
	err := os.MkdirAll(blogDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create blog directory: %v", err)
	}

	// Create test posts
	posts := []struct {
		filename string
		content  string
	}{
		{
			"post1.md",
			`---
title: "First Post"
date: "2024-01-15"
category: "tech"
---
# First Post Content`,
		},
		{
			"post2.md",
			`---
title: "Second Post"
date: "2024-01-10"
category: "personal"
---
# Second Post Content`,
		},
		{
			"not-markdown.txt",
			"This should be ignored",
		},
	}

	for _, post := range posts {
		filePath := filepath.Join(blogDir, post.filename)
		err := os.WriteFile(filePath, []byte(post.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", post.filename, err)
		}
	}

	bs := NewBlogService(tempDir)
	err = bs.LoadPosts()
	if err != nil {
		t.Fatalf("Failed to load posts: %v", err)
	}

	if len(bs.posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(bs.posts))
	}

	if bs.posts[0].Title != "First Post" {
		t.Errorf("Expected first post to be 'First Post', got %s", bs.posts[0].Title)
	}
}

func TestGetAllPosts(t *testing.T) {
	bs := NewBlogService("test-content")
	bs.posts = []Post{
		{Title: "Post 1"},
		{Title: "Post 2"},
	}

	posts := bs.GetAllPosts()
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
}

func TestGetPostBySlug(t *testing.T) {
	bs := NewBlogService("test-content")
	bs.posts = []Post{
		{Title: "Post 1", Slug: "post-1"},
		{Title: "Post 2", Slug: "post-2"},
	}

	post, err := bs.GetPostBySlug("post-1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if post.Title != "Post 1" {
		t.Errorf("Expected title 'Post 1', got %s", post.Title)
	}

	_, err = bs.GetPostBySlug("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent post")
	}
}

func TestGetRecentPosts(t *testing.T) {
	bs := NewBlogService("test-content")
	bs.posts = []Post{
		{Title: "Post 1"},
		{Title: "Post 2"},
		{Title: "Post 3"},
		{Title: "Post 4"},
	}

	recent := bs.GetRecentPosts(2)
	if len(recent) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(recent))
	}

	recent = bs.GetRecentPosts(10)
	if len(recent) != 4 {
		t.Errorf("Expected 4 posts, got %d", len(recent))
	}

	recent = bs.GetRecentPosts(0)
	if len(recent) != 0 {
		t.Errorf("Expected 0 posts, got %d", len(recent))
	}
}

func TestLoadPostsEmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	blogDir := filepath.Join(tempDir, "blog")
	err := os.MkdirAll(blogDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create blog directory: %v", err)
	}

	bs := NewBlogService(tempDir)
	err = bs.LoadPosts()
	if err != nil {
		t.Fatalf("Failed to load posts from empty directory: %v", err)
	}

	if len(bs.posts) != 0 {
		t.Errorf("Expected 0 posts from empty directory, got %d", len(bs.posts))
	}
}

func TestLoadPostsNonexistentDirectory(t *testing.T) {
	bs := NewBlogService("nonexistent-directory")
	err := bs.LoadPosts()
	if err == nil {
		t.Error("Expected error when loading from nonexistent directory")
	}
}
