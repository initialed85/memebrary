#!/bin/bash

set -e

# this block ensures we can invoke this script from anywhere and have it automatically change to this folder first
pushd "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
function teardown() {
    popd >/dev/null 2>&1 || true
}
trap teardown exit

# we need docker to build the images
if ! command -v docker >/dev/null 2>&1; then
    echo "error: docker not found"
    exit 1
fi

# TODO: hack workaround for intermittent builds on darwin aarch64
mkdir -p ./tmp
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./tmp/api -trimpath ./cmd/api

docker build --platform=linux/amd64 -t kube-registry:5000/memebrary-api:latest -f ./docker/api/Dockerfile .

docker image push kube-registry:5000/memebrary-api:latest
