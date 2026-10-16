#!/bin/sh
# The test of this Airfield (befehle.test in .tower/airfield.json): go vet and go build,
# in the golang image of the version go.mod names. Exits non-zero if either fails.
#
# There is no Go where the TOWER runs its agents, but there is Docker, and the repo is
# mounted at its own path, so this needs the same path on the Docker host.
# -buildvcs=false: in a worktree .git points outside the mount, and VCS stamping would
# fail the build. Nothing is shipped from here; the build only has to succeed.
set -eu
cd "$(dirname "$0")/.."
WURZEL=$(pwd)
GO=$(sed -n 's/^go \([0-9][0-9.]*\)$/\1/p' go.mod)
[ -n "$GO" ] || { echo "ops/test.sh: no go version in go.mod" >&2; exit 2; }
docker run --rm -v "$WURZEL:$WURZEL" -w "$WURZEL" "golang:$GO" \
  sh -c 'go vet ./... && go build -buildvcs=false -o /dev/null . && echo "go vet, go build: ok (go $(go env GOVERSION))"'
