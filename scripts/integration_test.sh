#!/usr/bin/env bash
set -e

# Bring up API + DB
docker compose up -d

# Wait for Postgres to be ready
echo "Waiting for Postgres…"
until docker compose exec -T db pg_isready -U postgres; do
  sleep 1
done

# Wait for API to be ready
echo "Waiting for API…"
until curl -sf http://localhost:8080/healthz > /dev/null; do
  sleep 1
done

# Run integration tests
go test ./integration

# Tear down
docker compose down
