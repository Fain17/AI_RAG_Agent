#!/bin/bash

# Unified linting script that runs all tools in the correct order
# This prevents conflicts between isort, black, and flake8

echo "🔧 Running isort..."
python -m isort --profile black app/ tests/

echo "🎨 Running black..."
python -m black --line-length 79 app/ tests/

echo "🔍 Running flake8..."
python -m flake8 app/ tests/

echo "✅ All linting tools completed successfully!" 