#!/bin/bash
set -e

BINDIR="${BINDIR:-/usr/local/bin}"
BINNAME="${BINNAME:-clipper}"

echo "Installing $BINNAME..."

go install ./...

GOBIN=$(go env GOBIN)
[ -z "$GOBIN" ] && GOBIN="$(go env GOPATH)/bin"

sudo mv "$GOBIN/$BINNAME" "$BINDIR/$BINNAME"

echo "Done."
