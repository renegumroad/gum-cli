#!/usr/bin/env bash

set -exo pipefail

# Check if there are any tags in the repository
if git show-ref --tags | grep -q .; then
    # If there are tags, attempt to get the most recent tag
    TAG=$(git describe --tags 2> /dev/null)
    VERSION=$TAG
else
    # Fallback to the current commit hash if no tag is found
    VERSION=$(git rev-parse HEAD)
fi

go build -v -o build/darwin_arm64/gum \
    -ldflags "-s -w \
    -X github.com/renehernandez/gum-cli/internal/version.VERSION=${VERSION}" \
    main.go
