# CLAUDE.md for eirka-index

## Build & Test Commands
```
# Run all tests
go test ./...

# Run specific test
go test -v -run=TestDetailsSQL ./middleware

# Run tests with coverage report
go test -cover ./...

# Build the project
go build -o eirka-index

# Run static code analysis
go vet ./...

# Format code
gofmt -s -w .
```

## Code Style Guidelines

1. **Imports**: Group imports with standard library first, third-party next, then local packages. Separate groups with blank lines.

2. **Formatting**: Use tabs for indentation. Run `gofmt -s -w .` before committing.

3. **Naming**:
   - CamelCase for exported functions, types, and variables
   - Follow Controller naming pattern for handlers (e.g., `IndexController`)

4. **Error Handling**:
   - Use `c.Error(err).SetMeta()` in middleware with appropriate status codes
   - Panic only for critical startup errors

5. **Testing**:
   - Use testify for assertions and go-sqlmock for database mocks
   - Follow `TestFunctionName` naming convention
   - Create helper functions for test setup and teardown
   - Clear global state between tests (caches, mocks)
   - Test both success and error paths
   - For middleware tests, use simplified handlers that return JSON instead of HTML

6. **Architecture**:
   - MVC-like pattern (controllers, middleware)
   - Configuration in separate package