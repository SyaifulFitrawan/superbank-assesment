#!/bin/sh
set -e

echo "Running database migrations..."
/app/server migrate

echo "Running database seeder..."
/app/server seed

echo "Starting the server..."
exec /app/server