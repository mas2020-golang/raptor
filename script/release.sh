#!/usr/bin/env bash
set -euo pipefail

# Configuration
VERSION_FILE="packages/utils/utils.go"   # adjust if needed
VERSION_ASSIGN_REGEX='[[:space:]]*Version[[:space:]]*=[[:space:]]*"'

die() { echo "ERROR: $*" >&2; exit 1; }
require_cmd() { command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"; }

require_cmd git
require_cmd awk
require_cmd sed

# Check on main branch
current_branch="$(git rev-parse --abbrev-ref HEAD)"
[ "$current_branch" = "main" ] || die "You are on '$current_branch'. Switch to 'main' to release."

# Check clean tree
if [ -n "$(git status --porcelain)" ]; then
  die "Working tree not clean. Commit or stash changes before releasing."
fi

[ -f "$VERSION_FILE" ] || die "Version file not found: $VERSION_FILE"

# Extract current version from a line like:  Version = "0.3.0-SNAPSHOT"
current_version="$(awk -v rx="$VERSION_ASSIGN_REGEX" '
  match($0, rx"([0-9]+\\.[0-9]+\\.[0-9]+)(-SNAPSHOT)?(\")", m) { print m[1] (m[2]?m[2]:""); exit }
' "$VERSION_FILE")"

[ -n "$current_version" ] || die "Could not parse current version from $VERSION_FILE"
case "$current_version" in
  *-SNAPSHOT) ;;
  *) die "Current version '$current_version' does not end with -SNAPSHOT. Aborting." ;;
esac

release_version="${current_version%-SNAPSHOT}"
echo "Releasing version: $release_version"

# sed-in-place portable helper
inplace_sed() {
  local script="$1" file="$2" tmp
  tmp="$(mktemp "${file}.XXXX")"
  sed "$script" "$file" > "$tmp"
  mv "$tmp" "$file"
}

# 1) Write release version (drop -SNAPSHOT) in VERSION_FILE
# Replace only the first matching Version = "..."
# BSD/GNU sed portable: do a pattern-based substitution
sed_release="/$VERSION_ASSIGN_REGEX[0-9]\\+\\.[0-9]\\+\\.[0-9]\\+\\(-SNAPSHOT\\)\\?\\\"/{
  s//Version = \\\"$release_version\\\"/
}"
inplace_sed "$sed_release" "$VERSION_FILE"

# Verify change
after_release_version="$(awk -v rx="$VERSION_ASSIGN_REGEX" '
  match($0, rx"([0-9]+\\.[0-9]+\\.[0-9]+)(-SNAPSHOT)?(\")", m) { print m[1] (m[2]?m[2]:""); exit }
' "$VERSION_FILE")"
[ "$after_release_version" = "$release_version" ] || die "Failed to set release version in $VERSION_FILE"

# 2) Commit release on main
git add "$VERSION_FILE"
git commit -m "release: v$release_version"

# 3) Tag the release
git tag -a "v$release_version" -m "Release v$release_version"

# 4) Compute next dev version: bump MINOR, set PATCH=0, add -SNAPSHOT
IFS='.' read -r major minor patch <<<"$release_version"
minor=$((minor + 1))
next_dev_version="${major}.${minor}.0-SNAPSHOT"

echo "Next development version: $next_dev_version"

# 5) Write next dev version
sed_next="/$VERSION_ASSIGN_REGEX[0-9]\\+\\.[0-9]\\+\\.[0-9]\\+\\(-SNAPSHOT\\)\\?\\\"/{
  s//Version = \\\"$next_dev_version\\\"/
}"
inplace_sed "$sed_next" "$VERSION_FILE"

# Verify next dev
after_next_version="$(awk -v rx="$VERSION_ASSIGN_REGEX" '
  match($0, rx"([0-9]+\\.[0-9]+\\.[0-9]+)(-SNAPSHOT)?(\")", m) { print m[1] (m[2]?m[2]:""); exit }
' "$VERSION_FILE")"
[ "$after_next_version" = "$next_dev_version" ] || die "Failed to set next development version in $VERSION_FILE"

# 6) Commit next dev
git add "$VERSION_FILE"
git commit -m "chore: start next development cycle v$next_dev_version"

cat <<EOF

Success!
- Release commit:   v$release_version
- Release tag:      v$release_version
- Next dev version: $next_dev_version

This script does NOT push to remote.
To publish:
  git push origin main
  git push origin v$release_version
EOF
