# copilot-instructions.md

This file provides guidance to GitHub Copilot when working with code in this repository.

## Overview

This is a KrakenD-based LLM model router implementation that acts as an egress API gateway to route between different LLM providers (OpenAI, Gemini, etc.). The project uses Go plugins to extend KrakenD functionality.

## Architecture

### Core Structure
- **KrakenD Base**: Uses KrakenD v2.10.1 as the API gateway foundation
- **Plugin System**: Go plugins built as `.so` files that extend KrakenD handlers
- **Shared Libraries**: Common utilities in `go/lib/` for header constants and logging

### Key Components
- `go/plugin/echo/`: Example plugin demonstrating request/response handling
- `go/lib/header/`: HTTP header constants
- `go/lib/logging/`: Structured logging with session tracking
- `Dockerfile`: Multi-stage build that compiles KrakenD and plugins

### Plugin Architecture
- Plugins implement the `HandlerRegisterer` interface
- Each plugin registers handlers via `RegisterHandlers()` function
- Plugins are loaded alphabetically and built as shared objects
- Configuration passed through KrakenD's `extra_config` mechanism

## Development Commands

### Building and Testing
```bash
# Build all plugins
make plugins

# Run tests for all Go code
make test

# Build Docker image
make image

# Start with Docker
make run
```

### Go-specific Commands (from go/ directory)
```bash
# Run tests with coverage
go test -cover ./...

# Build individual plugin
go build -buildmode=plugin -o ../build/echo.so ./plugin/echo

# Get dependencies
go get -t ./...
```

### Docker Operations
```bash
# Build Docker image
docker build -t agentic-layer/model-router-krakend .

# Run with environment variables
docker run -p 8080:8080 -e OPENAI_API_KEY=your_key_here agentic-layer/model-router-krakend
```

## Testing

The project uses Go's standard testing framework with testify assertions. Tests are located alongside source files with `_test.go` suffix.

## Plugin Development

When creating new plugins:
1. Follow the pattern in `go/plugin/echo/`
2. Implement `HandlerRegisterer` interface
3. Use shared logging from `go/lib/logging/`
4. Add plugin name to `PLUGINS` in `go/Makefile`
5. Configuration should be parsed from `extra_config` map

## Dependencies

- Go 1.24 (must match KrakenD CE requirements)
- KrakenD CE v2.10.1
- Key Go modules: `github.com/google/uuid`, `github.com/stretchr/testify`
- OpenTelemetry dependencies with specific version pins