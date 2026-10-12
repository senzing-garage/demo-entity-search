# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`demo-entity-search` is a Go CLI application demonstrating entity search capabilities using the Senzing SDK. It is part of the [senzing-tools](https://github.com/senzing-garage/senzing-tools) suite of tools.

## Prerequisites

Go 1.27.1 or later. The Senzing SDK is not required: this repository does not link the Senzing C library.

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

# Run tests
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
  - `/entity-search/` - Entity search UI
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

- `SENZING_TOOLS_AVOID_SERVING` - Skip starting HTTP server (for testing)

## Testing Notes

Tests need no database or other setup; `make setup` is a no-op kept for consistency with other repositories.

Binary outputs go to `target/<os>-<arch>/demo-entity-search`.
