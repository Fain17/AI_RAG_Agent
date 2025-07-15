# Linting

This directory contains the linting configuration and scripts for the RAG Backend Go application.

## Files

- **golangci-lint.yml**: Configuration file for golangci-lint with comprehensive linting rules
- **lint.sh**: Shell script to run the linting process with proper setup and error handling

## Usage

### Running Linting Locally

```bash
# Make sure you're in the rag-backend directory
cd rag-backend

# Make the script executable (if not already)
chmod +x linting/lint.sh

# Run the linting process
./linting/lint.sh
```

### Linting Configuration

The `golangci-lint.yml` file includes:

- **Enabled Linters**: A comprehensive set of linters including:
  - `errcheck`: Checks for unchecked errors
  - `govet`: Examines Go source code and reports suspicious constructs
  - `staticcheck`: Advanced static analysis
  - `gosec`: Security-focused linter
  - `gofmt`: Checks code formatting
  - `goimports`: Checks import organization
  - And many more...

- **Custom Settings**: Configured for:
  - Reasonable complexity limits
  - Project-specific import paths
  - Test file exceptions
  - Performance optimizations

### CI/CD Integration

The linting process is automatically run in the GitHub Actions workflow:

1. **Trigger**: Runs on every push and pull request to the `develop` branch
2. **Dependencies**: Automatically installs golangci-lint if not present
3. **Caching**: Uses Go module caching for faster builds
4. **Artifacts**: Uploads linting results for review

### Linting Rules

The configuration includes rules for:

- **Code Quality**: Detecting unused variables, dead code, and inefficient constructs
- **Security**: Identifying potential security issues
- **Style**: Ensuring consistent code formatting and naming conventions
- **Performance**: Detecting performance anti-patterns
- **Maintainability**: Enforcing good practices for readable code

### Troubleshooting

If linting fails:

1. Check the output for specific errors
2. Use `golangci-lint run --fix` to automatically fix some issues
3. Review the configuration in `golangci-lint.yml` if needed
4. Check the GitHub Actions logs for CI/CD failures

### Adding New Rules

To add new linting rules:

1. Edit `golangci-lint.yml`
2. Add the linter to the `linters.enable` section
3. Configure any specific settings in `linters-settings`
4. Test locally before pushing

### Local Development

For development, you can also run specific linters:

```bash
# Run only specific linters
golangci-lint run --enable=errcheck,govet

# Run with auto-fix
golangci-lint run --fix

# Run with detailed output
golangci-lint run --verbose
``` 