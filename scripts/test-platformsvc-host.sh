#!/bin/sh
# Exercise Android's pure-Go bridge on a Linux host, including race checks.
# Only the disposable copy loses its android build constraints.
set -eu
repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT HUP INT TERM
mkdir -p "$test_dir/pkg" "$test_dir/android/app/src/main/kotlin/com/nw5w/graywolf"
cp "$repo_dir/go.mod" "$repo_dir/go.sum" "$test_dir/"
cp -R "$repo_dir/pkg/platformsvc" "$repo_dir/pkg/platformproto" "$test_dir/pkg/"
cp "$repo_dir/android/app/src/main/kotlin/com/nw5w/graywolf/GraywolfService.kt" \
  "$test_dir/android/app/src/main/kotlin/com/nw5w/graywolf/"
for source in "$test_dir/pkg/platformsvc/"*.go; do
  sed '/^\/\/go:build android$/d' "$source" > "$source.host"
  mv "$source.host" "$source"
done
cd "$test_dir"
go test -race ./pkg/platformsvc
