#!/bin/bash

# make these configurable
pghost="${POSTGRES_HOST:-huautla}"
pgport="${POSTGRES_PORT:-5432}"
pguser="${POSTGRES_USER:-postgres}"
pgpass="${POSTGRES_PASSWORD:-root}"

psql_cmd=( "psql" "postgresql://${pguser}:${pgpass}@${pghost}:${pgport}" )
