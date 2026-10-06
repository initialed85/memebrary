#!/bin/bash

set -euo pipefail

pushd "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
function teardown() {
    popd >/dev/null 2>&1 || true
}
trap teardown EXIT

if ! command -v docker >/dev/null 2>&1; then
    echo "error: docker not found" >&2
    exit 1
fi

# Compile the API on the host so the image build stays small and reproducible.
mkdir -p ./tmp
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./tmp/api -trimpath ./cmd/api

BACKEND_IMAGE="${BACKEND_IMAGE:-initialed85/memebrary-backend:latest}"
FRONTEND_IMAGE="${FRONTEND_IMAGE:-initialed85/memebrary-frontend:latest}"

echo "building ${BACKEND_IMAGE}..."
docker build --platform=linux/amd64 -t "${BACKEND_IMAGE}" -f ./docker/api/Dockerfile .
echo "building ${FRONTEND_IMAGE}..."
docker build --platform=linux/amd64 -t "${FRONTEND_IMAGE}" -f ../frontend/Dockerfile ../frontend

docker image push "${BACKEND_IMAGE}"
docker image push "${FRONTEND_IMAGE}"
