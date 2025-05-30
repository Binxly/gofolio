package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Post struct {
	Title    string
	Date     time.Time
	Category string
	Excerpt  string
	Slug     string
	Content  template.HTML
	FilePath string
}

type BlogService struct {
	contentDir string
	posts      []Post
}

func NewBlogService(contentDir string) *BlogService {
	return &BlogService{
		contentDir: contentDir,
		posts:      make([]Post, 0),
	}
}

func (bs *BlogService) LoadPosts() error {
	blogDir := filepath.Join(bs.contentDir, "posts")

	err := filepath.WalkDir(blogDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, err := bs.parsePost(path)
		if err != nil {
			fmt.Printf("Error parsing post %s: %v\n", path, err)
			return nil
		}

		bs.posts = append(bs.posts, *post)
		return nil
	})

	if err != nil {
		return err
	}

	sort.Slice(bs.posts, func(i, j int) bool {
		return bs.posts[i].Date.After(bs.posts[j].Date)
	})

	return nil
}

func (bs *BlogService) parsePost(filePath string) (*Post, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta), // don't like this - needs load content to display previews on feed!
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // newlines are <br> tags
			html.WithXHTML(),     //xhtml standards bc i forget to close tags
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()

	if err := md.Convert(content, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	metaData := meta.Get(context)

	post := &Post{
		Content:  template.HTML(buf.String()),
		FilePath: filePath,
	}

	if title, ok := metaData["title"]; ok {
		if titleStr, ok := title.(string); ok {
			post.Title = titleStr
		}
	}

	if date, ok := metaData["date"]; ok {
		if dateStr, ok := date.(string); ok {
			if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
				post.Date = parsedDate
			}
		}
	}

	if category, ok := metaData["category"]; ok {
		if categoryStr, ok := category.(string); ok {
			post.Category = categoryStr
		}
	}

	if excerpt, ok := metaData["excerpt"]; ok {
		if excerptStr, ok := excerpt.(string); ok {
			post.Excerpt = excerptStr
		}
	}

	if slug, ok := metaData["slug"]; ok {
		if slugStr, ok := slug.(string); ok {
			post.Slug = slugStr
		}
	}

	if post.Slug == "" {
		filename := filepath.Base(filePath)
		post.Slug = strings.TrimSuffix(filename, ".md")
	}

	return post, nil
}

func (bs *BlogService) GetAllPosts() []Post {
	return bs.posts
}

func (bs *BlogService) GetPostBySlug(slug string) (*Post, error) {
	for _, post := range bs.posts {
		if post.Slug == slug {
			return &post, nil
		}
	}
	return nil, fmt.Errorf("post with slug '%s' not found", slug)
}

func (bs *BlogService) GetRecentPosts(count int) []Post {
	if count > len(bs.posts) {
		count = len(bs.posts)
	}
	return bs.posts[:count]
}
