# Contributing to the Bunny.net Go API Client

Thank you for your interest in contributing to the Bunny.net Go API client! This document provides guidelines and instructions to help you contribute effectively.

## Code of Conduct

By participating in this project, you agree to respect all project contributors and maintain a positive and inclusive atmosphere.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/bunnynet-go-client.git
   cd bunnynet-go-client
   ```
3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/venom90/bunnynet-go.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

### Creating a Branch

Create a branch for your contribution:

```bash
git checkout -b feature/your-feature-name
```

Use prefixes like `feature/`, `bugfix/`, `docs/`, or `test/` to indicate the type of change.

### Making Changes

1. Implement your changes following the existing code style and patterns
2. Add or update tests to cover your changes
3. Ensure your code passes all tests
4. Update documentation as needed

## Testing Requirements

We place a strong emphasis on comprehensive testing. All contributions must include adequate test coverage.

### Test Structure

Tests should be organized in the following manner:

- Unit tests for each package in `*_test.go` files next to the implementation
- Integration tests in the `test/` directory
- Resource-specific tests in `test/resources/`
- Mock HTTP server tests in `test/mocks.go`

### Testing Standards

1. **Coverage Requirements**:

   - New code should have at least 80% test coverage
   - Critical functionality should aim for 100% coverage
   - Run coverage checks: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

2. **Mock HTTP Testing**:

   - Use the provided `MockServer` from `test/mocks.go` for API testing
   - Test both happy paths and error scenarios
   - Verify request headers, paths, and query parameters
   - Test response handling and error conditions

3. **Test Utilities**:

   - Use provided assertion helpers in `test/mocks.go`
   - Create table-driven tests for comprehensive test cases
   - Use meaningful test names: `TestServiceName_MethodName_Condition`

4. **Example Test Pattern**:

   ```go
   func TestResourceService_Method_Success(t *testing.T) {
       // Setup mock server with expected response
       server := test.MockServer(t, http.StatusOK, `{"expected": "response"}`, func(r *http.Request) {
           test.AssertRequestMethod(t, r, http.MethodGet)
           test.AssertRequestPath(t, r, "/expected/path")
           test.AssertRequestHasHeader(t, r, "AccessKey", "test-api-key")
           // Verify query parameters if applicable
       })
       defer server.Close()

       // Create client that uses the mock server
       client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

       // Call the method being tested
       result, err := client.ResourceName.Method(context.Background(), params)

       // Assert results
       assert.NoError(t, err)
       assert.NotNil(t, result)
       assert.Equal(t, expectedValue, result.SomeField)
   }

   func TestResourceService_Method_Error(t *testing.T) {
       // Setup mock server with error response
       server := test.MockServer(t, http.StatusNotFound, `{
           "ErrorKey": "resource.not_found",
           "Field": "ResourceId",
           "Message": "The requested resource was not found"
       }`, nil)
       defer server.Close()

       // Create client that uses the mock server
       client := bunnynet.NewClient("test-api-key", bunnynet.WithBaseURL(server.URL))

       // Call the method being tested
       result, err := client.ResourceName.Method(context.Background(), params)

       // Assert error handling
       assert.Error(t, err)
       assert.Nil(t, result)
       assert.Contains(t, err.Error(), "resource.not_found")
   }
   ```

5. **Edge Case Testing**:

   - Test pagination with multiple pages
   - Test error handling with various HTTP status codes
   - Test with malformed server responses
   - Test with context cancellation
   - Test with rate limiting responses (429)
   - Test with server errors (500s)

6. **Performance Testing**:
   - Include benchmarks for performance-critical code:
   ```go
   func BenchmarkResourceService_Method(b *testing.B) {
       client := bunnynet.NewClient("test-api-key")
       for i := 0; i < b.N; i++ {
           // Call method to benchmark
       }
   }
   ```

### Running Tests

Run the entire test suite:

```bash
go test ./...
```

For more detailed test output:

```bash
go test -v ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

Run tests for a specific package:

```bash
go test -v ./resources
```

Run a specific test:

```bash
go test -v ./resources -run TestResourceName_MethodName
```

### Adding New Resources

When adding a new Bunny.net API resource:

1. Create a new file in the `resources` directory
2. Implement the resource service with appropriate methods
3. Add comprehensive tests for all methods and edge cases
4. Create a test file in `test/resources/`
5. Add the resource to the main `Client` struct in `client.go`
6. Add usage examples in the `examples` directory

## Pull Request Process

1. Update your fork with the latest upstream changes:

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Ensure all tests pass:

   ```bash
   go test ./...
   ```

3. Push your changes to your fork:

   ```bash
   git push origin feature/your-feature-name
   ```

4. Create a pull request from your branch to the main repository

5. In your pull request description, include:

   - A clear description of the changes
   - Any related issues that are addressed
   - Test coverage information
   - Any breaking changes or considerations for users

6. Participate in the code review process by responding to feedback

## Code Standards

### Go Conventions

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` and `go vet` before submitting your code
- Use meaningful variable names and add comments for complex logic

### Documentation

- Document all exported functions, types, and constants
- Include usage examples for new features
- Update the README.md if necessary

## Resource Structure

When implementing a new Bunny.net API resource, follow this structure:

```go
// Define resource types
type ResourceName struct {
    // Fields...
}

// ResourceNameService handles operations on the resource
type ResourceNameService struct {
    client    *http.Client
    baseURL   string
    apiKey    string
    userAgent string
}

// New ResourceNameService creates a new ResourceNameService
func NewResourceNameService(client *http.Client, baseURL, apiKey, userAgent string) *ResourceNameService {
    // Implementation...
}

// List, Get, Create, Update, Delete methods...
```

## Attribution

Your contributions will be recognized in the project's contributors list.

## Questions?

If you have any questions or need help, feel free to open an issue with the "question" label.

Thank you for contributing to the Bunny.net Go API client!
