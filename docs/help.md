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

Run these in the app's repository, `~/localitas-app-search`.

### Version

```bash
bin/search-server --version
```

### Build from source

```bash
# Native binary (bin/search-server)
make build

# Linux amd64 binary the Docker image is built from (search-server-linux-amd64)
make build-linux
```

Both run `make lint` (gofmt and go vet) first; `make test` runs the lint step and then the tests.

### Run natively

```bash
make start   # builds bin/search-server, writes the launchd service com.localitas.app.search for bin/search-server and loads it on port 9209
make stop
make logs    # follows ~/.localitas/logs/search/stdout.log
```

The app announces itself to core when it starts, and core serves it at `/apps/ext/search/`.

### Docker

```bash
make build-docker   # builds the Linux binary and the search:latest image
make start-docker   # runs container "search" on port 9209, labelled localitas.app=true
make stop-docker
make logs-docker
```

Core registers containers labelled `localitas.app=true` once Docker reports them healthy.

### Publish

```bash
# Run the tests, build the image, push ghcr.io/localitas/localitas-app-search:latest and :<VERSION>
make docker-push

# Build dist/search-server-darwin-arm64 and dist/search-server-darwin-amd64
# and create GitHub release vX.Y.Z on this repository with them
make release VERSION=vX.Y.Z
```

`VERSION` defaults to `git describe --tags --always --dirty`.

### Download

Pre-built macOS binaries are on the [releases page](https://github.com/localitas/localitas-app-search/releases):
- `search-server-darwin-arm64` (Apple Silicon)
- `search-server-darwin-amd64` (Intel)

Every push to `main` publishes a release: `.github/workflows/release.yml` runs the tests and releases both binaries as `v0.1.<commit count>.<UTC timestamp>`. `make release` publishes one by hand.

Download with the GitHub CLI:

    gh release download --repo localitas/localitas-app-search --pattern 'search-server-*'
