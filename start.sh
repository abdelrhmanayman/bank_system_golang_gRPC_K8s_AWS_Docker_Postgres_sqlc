#!/bin/sh

set -e

echo "Database migration"
source /app/.env
/app/migrate -path /app/migration -database "$DB_SOURCE_NAME" -verbose up


echo "Start App"
exec "$@"