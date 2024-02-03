#!/bin/sh

# make these configurable
pghost="${POSTGRES_HOST:-huautla}"
pgport="${POSTGRES_PORT:-5432}"
pguser="${POSTGRES_USER:-postgres}"
pgpass="${POSTGRES_PASSWORD:-root}"

psql_cmd=("psql" 'postgresql://"${pguser}:${pgpass}@${pghost}:${pgport}"')

echo "updating apt-get"
apt-get update
echo "installing postgres client"
apt-get -y upgrade postgresql-client

cat ./sql/create.pgsql ./sql/init.pgsql ./sql/seed.pgsql \
| "${psql_cmd[@]}" --echo-all -v ON_ERROR_STOP=1 \
&& "${psql_cmd[@]}"
