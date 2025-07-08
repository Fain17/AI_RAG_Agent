#!/bin/bash

# Setup script for installing linting dependencies
# Run this once to set up your local development environment
# Usage: cd linting && ./setup-lint.sh

echo "🔧 Setting up linting tools..."

# Install linting dependencies (from current directory)
pip install -r lint-requirements.txt

# Make lint script executable
chmod +x lint.sh

# Create a convenient symlink in the root directory
cd ..
ln -sf linting/lint.sh lint
echo "📎 Created symlink: ./lint -> linting/lint.sh"

echo "✅ Setup complete!"
echo ""
echo "Now you can run from the root directory:"
echo "  ./lint          # Format and check code"
echo "  ./lint --check  # Check-only mode (no modifications)"
echo ""
echo "Or run directly from linting/:"
echo "  cd linting && ./lint.sh"
echo ""
echo "💡 Tip: Run './lint --check' before pushing to ensure code quality!" 