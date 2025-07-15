#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Set script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

print_status "Starting linting process..."
print_status "Project directory: $PROJECT_DIR"

# Change to project directory
cd "$PROJECT_DIR" || exit 1

# Set up golangci-lint path
GOPATH=$(go env GOPATH)
export PATH="$GOPATH/bin:$PATH"

# Check if golangci-lint is installed
GOLANGCI_LINT_PATH="$GOPATH/bin/golangci-lint"
if ! command -v golangci-lint &> /dev/null && [ ! -f "$GOLANGCI_LINT_PATH" ]; then
    print_error "golangci-lint is not installed"
    print_status "Installing golangci-lint..."
    
    # Install golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOPATH/bin v1.54.2
    
    # Check if installation was successful
    if [ ! -f "$GOLANGCI_LINT_PATH" ]; then
        print_error "Failed to install golangci-lint"
        exit 1
    fi
    
    print_status "golangci-lint installed successfully"
else
    print_status "golangci-lint is already installed"
fi

# Determine which golangci-lint to use
if command -v golangci-lint &> /dev/null; then
    GOLANGCI_LINT="golangci-lint"
else
    GOLANGCI_LINT="$GOLANGCI_LINT_PATH"
fi

print_status "Using golangci-lint: $GOLANGCI_LINT"

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

# Tidy up modules
print_status "Tidying up modules..."
go mod tidy
if [ $? -ne 0 ]; then
    print_error "Failed to tidy modules"
    exit 1
fi

# Format code
print_status "Formatting code..."
go fmt ./api/handlers/*.go
if [ $? -ne 0 ]; then
    print_error "Failed to format code"
    exit 1
fi

go fmt ./api/models/*.go
if [ $? -ne 0 ]; then
    print_error "Failed to format code"
    exit 1
fi

go fmt ./api/routes/*.go
if [ $? -ne 0 ]; then
    print_error "Failed to format code"
    exit 1
fi

# Skipping connection.go formatting (ignored file)
# go fmt ./db/connection.go

# Run basic linting tools
print_status "Running basic linting tools..."

# Check syntax with go vet
print_status "Running go vet..."
go vet ./api/handlers/*.go
if [ $? -ne 0 ]; then
    print_error "go vet failed for api/handlers"
    exit 1
fi

go vet ./api/models/*.go
if [ $? -ne 0 ]; then
    print_error "go vet failed for api/models"
    exit 1
fi

go vet ./api/routes/*.go
if [ $? -ne 0 ]; then
    print_error "go vet failed for api/routes"
    exit 1
fi

# Skipping connection.go vet check (ignored file)
# go vet ./db/connection.go

# Check imports with goimports
print_status "Checking imports with goimports..."
if command -v goimports &> /dev/null; then
    goimports -l ./api/handlers/*.go ./api/models/*.go ./api/routes/*.go
    # Skipping connection.go from goimports check (ignored file)
    if [ $? -ne 0 ]; then
        print_error "goimports check failed"
        exit 1
    fi
else
    print_warning "goimports not found, skipping import check"
fi

# Run golangci-lint with configuration
print_status "Running golangci-lint with configuration..."
CONFIG_FILE="$SCRIPT_DIR/golangci-lint.yml"
if [ -f "$CONFIG_FILE" ]; then
    $GOLANGCI_LINT run --config="$CONFIG_FILE" --timeout=5m
    LINT_EXIT_CODE=$?
    
    if [ $LINT_EXIT_CODE -ne 0 ]; then
        print_error "golangci-lint failed with exit code $LINT_EXIT_CODE"
        exit $LINT_EXIT_CODE
    fi
    
    print_status "golangci-lint passed successfully"
else
    print_warning "golangci-lint.yml not found at $CONFIG_FILE"
    print_status "Running golangci-lint with default configuration..."
    $GOLANGCI_LINT run --timeout=5m
    LINT_EXIT_CODE=$?
    
    if [ $LINT_EXIT_CODE -ne 0 ]; then
        print_error "golangci-lint failed with exit code $LINT_EXIT_CODE"
        exit $LINT_EXIT_CODE
    fi
fi

print_status "All linting checks passed! 🎉" 