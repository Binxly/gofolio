# Gofolio

Go + Templ + TailwindCSS

## Run

To start the server, you can use the commands found in the Makefile:

```bash
make dev
```

Or manually:

```bash
make css    # Build CSS
make run    # Generate templates and run server
```

Once the server is running, you can navigate to <http://localhost:3000> and see the home page.

## Stack

- Go - server logic
- Templ - templates for HTML components
- Goldmark - MD to HTML conversion
- TailwindCSS - styling
- Gorilla Mux - routing

## Development

Edit `.templ` files. CSS is located in `static/assets/css/input.css`.

Available commands:

- `make dev` - Start development server with auto-reload and CSS watching
- `make run` - Generate templates and run server once
- `make quick-deploy` - Attempts to deploy to AWS Lambda function  
- `make css` - Build CSS once
- `make build` - Build production binary
- `make clean` - Clean generated files

## Blog System

The blog system uses Goldmark for markdown parsing, Goldmark-meta for frontmatter / metadata parsing. The blog posts are contained in `content/blog/`. Their associated custom prose styles are defined inside `static/assets/input.css` && `output.css`

## TODO

### Blog Feed

- [X] RSS feed!

### Styles

- [ ] Scrape out Tailwind CDN in templ base component

### QoL

- [ ] Self-host Tailwind

- [ ] Split components

- [ ] Lazy Loading for Posts!! 
    - Make RecentPosts scan for metadata instead of running LoadPosts()

- [ ] S3

- [ ] Fix GH Action to automate update of AWS Lambda function.
