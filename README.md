# Eirka-Index

A Go web service that generates index pages and handles CSRF protection for the Eirka imageboards platform.

## Overview

Eirka-Index serves as the frontend service for the Eirka imageboards platform. It provides:

- HTML templates for the SPA frontend
- Domain-specific routing and handling
- CSRF protection via secure cookies
- Cached imageboard settings and configuration

## Architecture

The application follows an MVC-like architecture:

- **Controllers**: Handle incoming requests and rendering templates
- **Middleware**: Process requests before they reach controllers
- **Config**: Manages application and imageboard settings

## Technology Stack

- [Go](https://golang.org) 1.22+
- [Gin](https://github.com/gin-gonic/gin) web framework
- [HTML templates](https://pkg.go.dev/html/template) for server-side rendering
- [Facebook Grace](https://github.com/facebookgo/grace) for zero-downtime restarts

## Development

### Prerequisites

- Go 1.22 or later
- MySQL/MariaDB database
- Configuration file at `/etc/pram/pram.conf` (optional, falls back to defaults)

### Building and Running

```bash
# Build the application
go build -o eirka-index

# Run the application
./eirka-index
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific tests
go test -v -run=TestDetailsSQL ./middleware

# Generate test coverage report
go test -cover ./...
```

## Configuration

Configuration is loaded from `/etc/pram/pram.conf` if available, otherwise defaults are used.
See `config/config.go` for configuration options and defaults.

## License

See [LICENSE](LICENSE) file for details.
