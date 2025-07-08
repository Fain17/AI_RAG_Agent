# Linting Tools

This directory contains all the linting and code quality tools for the project.

## Files

- **`lint-requirements.txt`** - Python dependencies for linting tools
- **`setup-lint.sh`** - One-time setup script to install dependencies
- **`lint.sh`** - Main linting script that runs all code quality checks
- **`README.md`** - This file

## Setup (One-time)

```bash
cd linting
./setup-lint.sh
```

This will:
- Install all required linting dependencies
- Make the lint script executable
- Create a convenient symlink in the root directory

## Usage

From the root directory (recommended):
```bash
./lint          # Format and check code
./lint --check  # Check-only mode (no modifications)
```

Or directly from the linting directory:
```bash
cd linting
./lint.sh       # Format and check code
./lint.sh --check  # Check-only mode
```

## Tools Included

1. **isort** - Import sorting
2. **black** - Code formatting
3. **mypy** - Type checking
4. **flake8** - Code linting
5. **bandit** - Security scanning
6. **pytest** - Test running

## Configuration Files

Configuration files are kept in the root directory where tools expect them:
- `.flake8` - flake8 configuration
- `mypy.ini` - mypy configuration
- `.gitignore` / `.dockerignore` - ignore files

## Before Pushing Code

Always run the check-only mode before pushing:
```bash
./lint --check
```

This ensures your code meets all quality standards without modifying any files. 