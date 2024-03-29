#!/bin/sh

set -e
echo "Running Db Migrations"
source /app/app.env
/app/migrate --path /app/migrations --database "$DB_SOURCE" --verbose up

echo "Starting The App"
exec "$@"
