#!/bin/sh

echo "run db migration"
migrate -path /app/migration -database "$DATABASE_URL" -verbose up
