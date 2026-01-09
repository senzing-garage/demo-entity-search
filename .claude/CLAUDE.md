# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`demo-entity-search` is a Go CLI application demonstrating entity search capabilities using the Senzing SDK. It is part of the [senzing-tools](https://github.com/senzing-garage/senzing-tools) suite of tools.

## Prerequisites

The Senzing C library must be installed before development:
- `/opt/senzing/er/lib` - Shared objects
- `/opt/senzing/er/sdk/c` - SDK header files
- `/etc/opt/senzing` - Configuration

See [How to Install Senzing for Go Development](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md).

## Common Commands

```bash
# Install development dependencies (one-time)
make dependencies-for-development

# Update Go dependencies
make dependencies

# Build for current platform
make clean build

# Build for all platforms
make build-all

# Run linting (golangci-lint + govulncheck + cspell)
make lint

# Run tests (requires setup first)
make clean setup test

# Run a single test
go test -v -run TestBasicHTTPServer_Serve ./httpserver/...

# Run tests with coverage
make clean setup coverage

# Check coverage meets thresholds
make check-coverage

# Apply lint auto-fixes
make fix

# Run the application
make run
# Or directly: go run main.go
```

## Architecture

The application follows a layered architecture:

- **cmd/** - Cobra CLI commands. `root.go` defines the main command that starts the HTTP server. OS-specific context variables are in `context_*.go` files.

- **httpserver/** - HTTP server implementation (`BasicHTTPServer`). Sets up routes:
  - `/api/` - Senzing REST API
  - `/entity-search/` - Entity search UI
  - `/entity-search/api/` - Reverse proxy to Senzing API
  - `/` - Static HTML files

- **entitysearchservice/** - Entity search service (`BasicHTTPService`). Serves the Angular-based entity search UI from embedded static files.

- **internal/log/** - Internal logging constants and utilities.

- **makefiles/** - OS-specific Makefile includes (`linux.mk`, `darwin.mk`, `windows.mk`).

## Code Style

- Maximum line length: 120 characters
- Uses `gofumpt` for formatting (stricter than `gofmt`)
- Exhaustive linting via golangci-lint with 100+ linters enabled
- Test files use `_test` package suffix (e.g., `httpserver_test`)
- Tests should use `test.Parallel()` when possible

## Environment Variables

- `SENZING_TOOLS_DATABASE_URL` - Database connection (default: `sqlite3://na:na@nowhere/tmp/sqlite/G2C.db`)
- `SENZING_TOOLS_AVOID_SERVING` - Skip starting HTTP server (for testing)
- `SENZING_LOG_LEVEL` - Log level (e.g., `TRACE`)
- `LD_LIBRARY_PATH` - Path to Senzing libraries (default: `/opt/senzing/er/lib`)

## Testing Notes

Tests require setup to create a SQLite database:
```bash
make setup  # Creates /tmp/sqlite/G2C.db from testdata
```

Binary outputs go to `target/<os>-<arch>/demo-entity-search`.
