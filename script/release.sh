#!/usr/bin/env bash
set -euo pipefail

# === Config ===
VERSION_FILE="packages/utils/utils.go"     # path to the file containing: Version = "X.Y.Z[-SNAPSHOT]"
VERSION_RE='Version[[:space:]]*=[[:space:]]*"'  # left side of the assignment (no number here)

# === Helpers ===
die() { echo "ERROR: $*" >&2; exit 1; }
require_cmd() { command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"; }

require_cmd git
require_cmd sed
require_cmd awk

# === Preconditions ===
current_branch="$(git rev-parse --abbrev-ref HEAD)"
[ "$current_branch" = "main" ] || die "You are on '$current_branch'. Switch to 'main' to release."

if [ -n "$(git status --porcelain)" ]; then
  die "Working tree not clean. Commit or stash changes before releasing."
fi

[ -f "$VERSION_FILE" ] || die "Version file not found: $VERSION_FILE"

# === Read current version ===
# Extracts the first occurrence of Version = "X.Y.Z(-SNAPSHOT)?"
current_version="$(sed -nE "s/.*(${VERSION_RE})([0-9]+\.[0-9]+\.[0-9]+)(-SNAPSHOT)?(\").*/\2\3/p" "$VERSION_FILE" | head -n1)"
[ -n "$current_version" ] || die "Could not parse current version from $VERSION_FILE (expected: Version = \"X.Y.Z-SNAPSHOT\")"

case "$current_version" in
  *-SNAPSHOT) ;;
  *) die "Current version '$current_version' does not end with -SNAPSHOT. Aborting." ;;
esac

release_version="${current_version%-SNAPSHOT}"
echo "Releasing version: $release_version"

# Portable in-place write (BSD/GNU compatible) — replace only the FIRST match in file
write_version() {
  # $1: new version string (e.g., 1.2.3 or 1.3.0-SNAPSHOT)
  local new_version="$1"
  local tmp
  tmp="$(mktemp "${VERSION_FILE}.XXXX")" || die "mktemp failed"
  # Pattern that matches the full right-hand side with the number
  local assign_pat='Version[[:space:]]*=[[:space:]]*"[0-9]+\.[0-9]+\.[0-9]+(-SNAPSHOT)?"'
  # Use awk to replace the FIRST occurrence only (POSIX awk)
  awk -v pat="$assign_pat" -v rep="Version = \"${new_version}\"" '
    $0 ~ pat && !done { sub(pat, rep); done=1 }
    { print }
  ' "$VERSION_FILE" > "$tmp" || die "Failed to rewrite $VERSION_FILE"
  mv "$tmp" "$VERSION_FILE"
}

# 1) Write release version (drop -SNAPSHOT)
write_version "$release_version"

# Verify write
after_release_version="$(sed -nE "s/.*(${VERSION_RE})([0-9]+\.[0-9]+\.[0-9]+)(-SNAPSHOT)?(\").*/\2\3/p" "$VERSION_FILE" | head -n1)"
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

# 5) Write next dev version and verify
write_version "$next_dev_version"
after_next_version="$(sed -nE "s/.*(${VERSION_RE})([0-9]+\.[0-9]+\.[0-9]+)(-SNAPSHOT)?(\").*/\2\3/p" "$VERSION_FILE" | head -n1)"
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
