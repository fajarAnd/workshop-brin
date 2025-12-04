#!/bin/sh
set -e

echo "Starting container..."

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h ${DB_HOST:-postgres-brin} -p ${DB_PORT:-5432} -U ${POSTGRES_WA_SERVICE_USER:-wa_service} > /dev/null 2>&1; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations if DATABASE_URL is set
if [ -n "$DATABASE_URL" ]; then
  echo "Running database migrations..."
  if migrate -path /app/migrations -database "$DATABASE_URL" up; then
    echo "Migrations completed successfully"
  else
    echo "Warning: Migrations failed or no migrations to run"
    # Don't exit, let the app start anyway
  fi
else
  echo "DATABASE_URL not set, skipping migrations"
fi

echo "Starting application..."
exec "$@"