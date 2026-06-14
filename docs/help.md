---
title: Search
description: Web search aggregator
---

# Search

Search the web using multiple search engines simultaneously and view combined results.

Full API spec: swagger.json

## API Endpoints

Method | Path | Description
GET | /api/search?q={query} | Search across all engines

## Web Search

Submit a query and receive results from DuckDuckGo, Brave Search, and Bing in parallel. The app scrapes public search result pages and returns structured results.

    GET /api/search?q=localitas+distributed+cluster

## Search Engines

The app queries three search engines concurrently:

- **DuckDuckGo** - Privacy-focused search with no tracking
- **Brave Search** - Independent search index with privacy features
- **Bing** - Microsoft's search engine with broad web coverage

Results from each engine are returned separately so you can compare rankings across providers.

## Response Format

Each search result includes:
- Title of the page
- URL link
- Snippet or description text
- Source engine name

Errors from individual engines are reported in the response without blocking results from other engines.

## Web Interface

The browser UI provides a search box and displays results from all engines in a unified view. Results are grouped by engine for easy comparison.

## Rate Limiting

The app uses standard HTTP clients with reasonable timeouts. Avoid sending excessive requests to prevent temporary blocks from search engines.

## Build & Deploy

### Version

```bash
./search-server --version
```

### Build from source

```bash
# Development (native)
cd apps/search && go build -o bin/search-server ./cmd/search-server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o bin/search-server-linux-amd64 ./cmd/search-server
```

### Docker

Build a Docker image directly from the binary:

```bash
# Default base image (debian:12-slim)
./search-server docker-build

# Custom base image
./search-server docker-build --base ubuntu:24.04

# Custom Dockerfile
./search-server docker-build --dockerfile ./my.Dockerfile

# Tag and push to registry
./search-server docker-build --tag ghcr.io/localitas/search:latest --push
```

The `docker-build` command requires a Linux amd64 binary in the same directory. Run `make deploy-build` from the project root first.

### Download

Pre-built binaries are available on the [GitHub releases page](https://github.com/localitas/localitas/releases).

Each release includes three builds per app:
- `search-server-darwin-arm64` (macOS Apple Silicon)
- `search-server-linux-amd64` (Linux x86_64)
- `search-server-linux-arm64` (Linux ARM64)

Download with the GitHub CLI:

    gh release download --repo localitas/localitas --pattern 'search-server-*'

### Release

All app binaries are published to GitHub releases as part of `make deploy-upload-image`.
