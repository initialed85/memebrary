#!/bin/bash

set -e

REDIS_URL=redis://localhost:6379
DJANGOLANG_API_ROOT=/api/
POSTGRES_DB=memebrary
POSTGRES_PASSWORD="NoColgate!11"

export REDIS_URL
export DJANGOLANG_API_ROOT
export POSTGRES_DB
export POSTGRES_PASSWORD

PAGER=cat PGPASSWORD="${POSTGRES_PASSWORD}" psql -h localhost -U postgres "${POSTGRES_DB}" -c 'TRUNCATE TABLE repository CASCADE;'

rm -fr ./pkg/job_executor/tmp >/dev/null 2>&1 || true

if [[ "${1}" == "" ]]; then
    go test -v -count=1 -failfast ./...

fi

go test -v -count=1 -failfast "${@}"
