#!/bin/sh
# wait-for-postgres.sh

set -e
echo "Postgres handler"
host="$1"
shift

until PGPASSWORD=$DATABASE_PASSWORD psql -h "$DATABASE_HOST" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec "$@"