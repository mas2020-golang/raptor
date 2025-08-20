#!/bin/bash

# Exit on any error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${GREEN}[STEP]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if utils.go exists
if [ ! -f "utils.go" ]; then
    print_error "utils.go file not found in current directory"
    exit 1
fi

# Get current branch name
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" = "main" ]; then
    print_error "Cannot run this script from main branch"
    exit 1
fi

print_step "Current branch: $CURRENT_BRANCH"

# Extract current version from utils.go
CURRENT_VERSION=$(grep -o 'Version = "[^"]*"' utils.go | sed 's/Version = "\([^"]*\)"/\1/')
print_step "Current version: $CURRENT_VERSION"

# Remove -SNAPSHOT suffix for release version
RELEASE_VERSION=$(echo "$CURRENT_VERSION" | sed 's/-SNAPSHOT//')
print_step "Release version: $RELEASE_VERSION"

# Check if we have uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    print_warning "You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

# Step 1: Remove -SNAPSHOT from utils.go
print_step "Removing -SNAPSHOT from utils.go"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/Version = \"$CURRENT_VERSION\"/Version = \"$RELEASE_VERSION\"/" utils.go
else
    # Linux
    sed -i "s/Version = \"$CURRENT_VERSION\"/Version = \"$RELEASE_VERSION\"/" utils.go
fi

# Commit the version change
git add utils.go
git commit -m "Release version $RELEASE_VERSION"

# Step 2: Create git tag
print_step "Creating git tag: v$RELEASE_VERSION"
git tag "v$RELEASE_VERSION"

# Step 3: Switch to main and merge current branch
print_step "Switching to main branch"
git checkout main
git pull origin main

print_step "Merging $CURRENT_BRANCH into main"
git merge "$CURRENT_BRANCH" --no-ff -m "Merge branch '$CURRENT_BRANCH' for release $RELEASE_VERSION"

# Step 4: Delete the current branch locally and remotely
print_step "Deleting branch $CURRENT_BRANCH locally and remotely"
git branch -d "$CURRENT_BRANCH"

# Check if remote branch exists before trying to delete it
if git ls-remote --heads origin "$CURRENT_BRANCH" | grep -q "$CURRENT_BRANCH"; then
    git push origin --delete "$CURRENT_BRANCH"
else
    print_warning "Remote branch $CURRENT_BRANCH does not exist, skipping remote deletion"
fi

# Step 5: Increment minor version and add -SNAPSHOT
print_step "Incrementing version for next development cycle"

# Parse version components (assuming semantic versioning: major.minor.patch)
IFS='.' read -ra VERSION_PARTS <<< "$RELEASE_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Increment minor version
NEW_MINOR=$((MINOR + 1))
NEXT_DEV_VERSION="$MAJOR.$NEW_MINOR.0-SNAPSHOT"

print_step "Next development version: $NEXT_DEV_VERSION"

# Update utils.go with new development version
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/Version = \"$RELEASE_VERSION\"/Version = \"$NEXT_DEV_VERSION\"/" utils.go
else
    # Linux
    sed -i "s/Version = \"$RELEASE_VERSION\"/Version = \"$NEXT_DEV_VERSION\"/" utils.go
fi

# Step 6: Create commit for next development version
print_step "Creating commit for next development version"
git add utils.go
git commit -m "Starting the new development version: $NEXT_DEV_VERSION"

# Step 7: Push changes and tags remotely
print_step "Pushing changes and tags to remote"
git push origin main
git push origin "v$RELEASE_VERSION"

print_step "Release process completed successfully!"
print_step "Released version: $RELEASE_VERSION"
print_step "Next development version: $NEXT_DEV_VERSION"
print_step "Current branch: main"

# Cleanup: The script note mentions "delete the previous sprint" but that's already done
# when we deleted the feature branch. If you meant something else, please clarify.

echo -e "${GREEN}✅ Release script completed successfully!${NC}"