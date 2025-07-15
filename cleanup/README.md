# GHCR Cleanup Script

This script deletes all but the latest 2 images for each GitHub Container Registry (GHCR) package owned by the user `Fain17`.

## Usage

1. **Create a Personal Access Token (PAT):**
   - Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
   - Click "Generate new token"
   - Select scopes: `read:packages` and `delete:packages`
   - Copy the token (you won't see it again!)

2. **Install requirements:**
   - You need the GitHub CLI (`gh`) and `jq` installed.

3. **Run the script:**
   ```sh
   GHCR_PAT=your_token ./ghcr-cleanup.sh
   ```

## What it does
- Authenticates with your PAT
- Lists all container packages for the user
- For each package, keeps the 2 most recent images and deletes the rest

**Warning:** This will permanently delete old images from your GHCR account. 