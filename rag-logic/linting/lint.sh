#!/bin/bash

# Comprehensive linting script for local development
# Run this before pushing code to ensure quality standards
#
# SETUP: Install linting dependencies first:
#   cd linting && ./setup-lint.sh
#
# USAGE (from root directory):
#   ./lint          # Format and check code
#   ./lint --check  # Check-only mode (no modifications)
#
# USAGE (from linting directory):
#   ./lint.sh       # Format and check code
#   ./lint.sh --check  # Check-only mode (no modifications)

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored status messages
print_status() {
    echo -e "${BLUE}🔧 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if we're in check-only mode
CHECK_ONLY=false
if [[ "$1" == "--check" ]]; then
    CHECK_ONLY=true
    print_warning "Running in CHECK-ONLY mode (no files will be modified)"
fi

echo "==========================================="
echo "🚀 Starting Python Code Quality Checks"
echo "==========================================="

# 1. Import sorting with isort
print_status "Running isort (import sorting)..."
if [[ "$CHECK_ONLY" == true ]]; then
    if ! python -m isort --check-only --diff app/ tests/; then
        print_error "isort found incorrectly sorted imports!"
        exit 1
    fi
else
    python -m isort --profile black app/ tests/
fi
print_success "isort completed successfully!"

# 2. Code formatting with black
print_status "Running black (code formatting)..."
if [[ "$CHECK_ONLY" == true ]]; then
    if ! python -m black --check app/ tests/; then
        print_error "black found formatting issues!"
        exit 1
    fi
else
    python -m black --line-length 88 app/ tests/
fi
print_success "black completed successfully!"

# 3. Type checking with mypy
print_status "Running mypy (type checking)..."
if ! python -m mypy app/; then
    print_error "mypy found type errors!"
    exit 1
fi
print_success "mypy completed successfully!"

# 4. Linting with flake8 (optional, can be disabled)
print_status "Running flake8 (code linting)..."
if command -v flake8 &> /dev/null; then
    # Configuration is loaded from .flake8 file
    if ! python -m flake8 app/ tests/; then
        print_error "flake8 found linting issues!"
        exit 1
    fi
    print_success "flake8 completed successfully!"
else
    print_warning "flake8 not installed, skipping..."
fi

# 5. Security check with bandit (optional)
print_status "Running bandit (security checking)..."
if command -v bandit &> /dev/null; then
    if ! python -m bandit -r app/ -f json -o bandit-report.json; then
        print_warning "bandit found potential security issues (check bandit-report.json)"
    else
        print_success "bandit completed successfully!"
    fi
else
    print_warning "bandit not installed, skipping..."
fi

# 6. Run tests (optional)
if [[ -d "tests/" ]] && command -v pytest &> /dev/null; then
    print_status "Running tests..."
    if ! python -m pytest tests/ -v; then
        print_error "tests failed!"
        exit 1
    fi
    print_success "tests completed successfully!"
fi

echo "==========================================="
print_success "🎉 All code quality checks passed!"
echo "==========================================="

if [[ "$CHECK_ONLY" == true ]]; then
    echo "Your code is ready to be pushed! 🚀"
else
    echo "Code has been formatted and is ready to be pushed! 🚀"
fi 