#!/bin/sh
set -eu

POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-NoColgate!11}"
POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_DB="${POSTGRES_DB:-memebrary}"
POSTGRES_SSLMODE="${POSTGRES_SSLMODE:-disable}"

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}"
export POSTGRES_USER POSTGRES_PASSWORD POSTGRES_HOST POSTGRES_PORT POSTGRES_DB POSTGRES_SSLMODE

# The API image owns its schema migrations. Retrying here makes a fresh
# Postgres pod safe to start before its service is accepting connections.
until /usr/local/bin/migrate -path /srv/migrations -database "$DATABASE_URL" up; do
  echo "waiting for PostgreSQL migrations..." >&2
  sleep 2
done

exec /srv/bin/api "${@:-serve}"
