#!/bin/bash
set -e

# Usage: GHCR_PAT=your_token ./ghcr-cleanup.sh
# Requires: gh, jq

OWNER="Fain17"

if [ -z "$GHCR_PAT" ]; then
  echo "❌ Please set the GHCR_PAT environment variable with a Personal Access Token (PAT) that has read:packages and delete:packages scopes."
  exit 1
fi

echo "$GHCR_PAT" | gh auth login --with-token

echo "📦 Fetching container packages for user: $OWNER"

packages=$(gh api -H "Accept: application/vnd.github+json" \
  /users/$OWNER/packages?package_type=container | jq -r '.[].name')

if [ -z "$packages" ]; then
  echo "❌ No packages found or failed to fetch packages. Check your GitHub token scopes (needs read:packages & delete:packages)"
  exit 1
fi

for package in $packages; do
  echo "🔄 Processing package: $package"

  versions=$(gh api -H "Accept: application/vnd.github+json" \
    /users/$OWNER/packages/container/$package/versions | jq -r '.[].id')

  if [ -z "$versions" ]; then
    echo "⚠️  No versions found for $package — skipping."
    continue
  fi

  count=0
  for version in $versions; do
    count=$((count + 1))
    if [ $count -le 2 ]; then
      echo "✅ Keeping version $version of $package"
    else
      echo "🗑️ Deleting version $version of $package"
      gh api --method DELETE -H "Accept: application/vnd.github+json" \
        /users/$OWNER/packages/container/$package/versions/$version
    fi
  done
done 