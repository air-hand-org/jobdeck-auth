#!/usr/bin/env bash
set -euo pipefail

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

go test -coverpkg=./... -coverprofile="$work/raw.out" "$@" ./...

# go tool cover embeds source. Validate every published file against this
# checkout's committed bytes, and reject paths outside the production trees.
module=$(go list -m)
awk -v prefix="$module/" '
  NR == 1 { print; next }
  index($0, prefix) != 1 { exit 1 }
  { print "./" substr($0, length(prefix) + 1) }
' "$work/raw.out" > "$work/local.out"
awk -F: 'NR > 1 { print $1 }' "$work/local.out" | sort -u > "$work/files"
while IFS= read -r file; do
  [[ "$file" =~ ^\./(cmd|internal)/([a-zA-Z0-9_-]+/)*[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*\.go$ &&
     "$file" != *_test.go && "$file" != */testdata/* ]] || {
    printf 'Unexpected coverage source: %s\n' "$file" >&2
    exit 1
  }
  [[ "$file" != ./internal/appmock/* ]] || continue
  [[ "$(realpath -- "$file")" == "$(pwd -P)/${file#./}" && -f "$file" ]]
  committed=$(git rev-parse --verify "HEAD:${file#./}")
  [[ "$(git hash-object --no-filters -- "$file")" == "$committed" ]]
done < "$work/files"
awk 'NR == 1 || $0 !~ /^\.\/internal\/appmock\//' "$work/local.out" > "$work/coverage.out"

go tool cover -func="$work/coverage.out" > "$work/functions"
total=$(awk '$1 == "total:" { print $NF }' "$work/functions")
[[ "$total" =~ ^[0-9]+\.[0-9]%$ ]]
printf 'Total statement coverage: %s\n' "$total"
mkdir -p pages/coverage
go tool cover -html="$work/coverage.out" -o pages/coverage/index.html
