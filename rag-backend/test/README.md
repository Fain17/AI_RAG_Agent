# Testing

This directory contains the test files and scripts for the RAG Backend Go application.

## Files

- **handlers_test.go**: Unit tests for all API handlers with mock database queries
- **main_test.go**: Basic tests for the main application functionality
- **run_tests.sh**: Shell script to run the testing process with proper setup and coverage

## Usage

### Running Tests Locally

```bash
# Make sure you're in the rag-backend directory
cd rag-backend

# Make the script executable (if not already)
chmod +x test/run_tests.sh

# Run the testing process
./test/run_tests.sh
```

### Running Specific Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test file
go test -v ./test/handlers_test.go

# Run specific test function
go test -v -run TestGetHandler_Success ./test/
```

## Test Structure

### Handler Tests (`handlers_test.go`)

The handler tests use a mock-based approach:

- **MockQueries**: Mock implementation of the database queries interface
- **Test Setup**: Each test creates a mock router and queries
- **Test Coverage**: Tests cover:
  - Success scenarios
  - Error handling
  - Invalid input validation
  - HTTP status codes
  - Response formatting

### Test Categories

1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test API endpoints with mocked dependencies
3. **Error Handling**: Test various error scenarios
4. **Validation Tests**: Test input validation and sanitization

### Mock Implementation

The tests use `testify/mock` for mocking:

- **Database Queries**: All database operations are mocked
- **HTTP Requests**: Use `httptest` for HTTP request simulation
- **Response Validation**: Assert HTTP status codes and response bodies

## Coverage

The test suite includes coverage reporting:

- **Coverage Report**: Generated as `coverage.out`
- **HTML Report**: Generated as `coverage.html`
- **Coverage Upload**: Automatically uploaded to Codecov in CI/CD

### Coverage Goals

- **Minimum Coverage**: Aim for >80% code coverage
- **Handler Coverage**: All API handlers should be tested
- **Error Path Coverage**: Error handling paths should be covered
- **Edge Cases**: Include tests for edge cases and boundary conditions

## CI/CD Integration

The testing process is automatically run in the GitHub Actions workflow:

1. **Trigger**: Runs on every push and pull request to the `develop` branch
2. **Dependencies**: Automatically downloads and verifies Go modules
3. **Environment**: Sets up test environment variables
4. **Race Detection**: Runs tests with race detection enabled
5. **Coverage**: Generates and uploads coverage reports
6. **Artifacts**: Uploads test results and coverage reports

## Test Environment

### Environment Variables

The test script sets up the following environment variables:

- `GIN_MODE=test`: Sets Gin to test mode
- `DATABASE_URL`: Test database connection string

### Test Database

For integration tests that require a database:

1. Use Docker to spin up a test database
2. Run migrations against the test database
3. Clean up after tests complete

## Adding New Tests

### For New Handlers

1. Add mock methods to `MockQueries` struct
2. Create test functions following the naming convention `TestHandlerName_Scenario`
3. Include both success and failure scenarios
4. Test input validation and error handling

### For New Features

1. Create separate test files for new packages
2. Follow the same mock-based approach
3. Ensure good test coverage
4. Add integration tests if needed

### Test Naming Convention

- **Test Functions**: `TestFunctionName_Scenario`
- **Test Files**: `*_test.go`
- **Mock Objects**: `Mock*` prefix

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all dependencies are in `go.mod`
2. **Mock Issues**: Verify mock setup matches actual interface
3. **Coverage Issues**: Check that all code paths are tested
4. **Race Conditions**: Use proper synchronization in tests

### Running Tests with Debug

```bash
# Run tests with debug output
go test -v -debug ./...

# Run tests with CPU profiling
go test -cpuprofile cpu.prof ./...

# Run tests with memory profiling
go test -memprofile mem.prof ./...
```

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Clear Assertions**: Use descriptive assertion messages
3. **Mock Cleanup**: Always assert mock expectations
4. **Test Data**: Use realistic test data
5. **Error Testing**: Test both success and error paths
6. **Performance**: Keep tests fast and efficient

## Future Improvements

- Add integration tests with real database
- Add performance benchmarks
- Add end-to-end tests
- Add property-based testing
- Add mutation testing 