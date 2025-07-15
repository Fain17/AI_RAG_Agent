#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

# Set script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

print_status "Starting test process..."
print_status "Project directory: $PROJECT_DIR"

# Change to project directory
cd "$PROJECT_DIR" || exit 1

# Check Go version
GO_VERSION=$(go version)
print_status "Go version: $GO_VERSION"

# Download dependencies
print_status "Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    print_error "Failed to download dependencies"
    exit 1
fi

# Verify dependencies
print_status "Verifying dependencies..."
go mod verify
if [ $? -ne 0 ]; then
    print_error "Failed to verify dependencies"
    exit 1
fi

# Clean test cache
print_status "Cleaning test cache..."
go clean -testcache

# Set test environment variables
export GIN_MODE=test
export DATABASE_URL="postgres://test:test@localhost/test"

print_status "Environment variables set:"
print_status "  GIN_MODE: $GIN_MODE"
print_status "  DATABASE_URL: $DATABASE_URL"

# Run tests with race detection and coverage
print_test "Running tests with race detection and coverage..."
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_status "Tests completed successfully! ✅"
    
    # Generate coverage report
    print_status "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    
    # Show coverage summary
    print_status "Coverage summary:"
    go tool cover -func=coverage.out | tail -1
    
    print_status "Coverage report generated: coverage.html"
    
    # Run benchmarks
    print_test "Running benchmark tests..."
    go test -bench=. -benchmem ./test > benchmark_results.txt 2>&1
    BENCH_EXIT_CODE=$?
    
    if [ $BENCH_EXIT_CODE -eq 0 ]; then
        print_status "Benchmarks completed successfully! ✅"
        print_status "Benchmark results saved to: benchmark_results.txt"
        
        # Show benchmark summary
        if [ -f benchmark_results.txt ]; then
            print_status "📊 Benchmark Summary:"
            echo "----------------------------------------"
            grep "^Benchmark" benchmark_results.txt | head -10
            echo "----------------------------------------"
            print_status "Full benchmark results in benchmark_results.txt"
        fi
    else
        print_warning "⚠️  Benchmarks had issues (exit code: $BENCH_EXIT_CODE)"
        print_warning "Check benchmark_results.txt for details"
    fi
    
    print_status "All tests and benchmarks completed! 🎉"
else
    print_error "Tests failed with exit code $TEST_EXIT_CODE ❌"
    exit $TEST_EXIT_CODE
fi 