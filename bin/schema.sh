#!/bin/sh

# echo fuckoff >&2
# make these configurable
pghost="${POSTGRES_HOST:-huautla}"
pgport="${POSTGRES_PORT:-5432}"
pguser="${POSTGRES_USER:-postgres}"
pgpass="${POSTGRES_PASSWORD:-root}"

echo "updating apt-get"
apt-get update
echo "installing postgres client"
apt-get -y upgrade postgresql-client

echo psql -h"${pghost}" -p"${pgport}" -o"${pguser}" -W"${pgpass}" -c "select 1 from dual;" >&2
cat ./sql/init.pgsql ./sql/create.pgsql ./sql/seed.pgsql \
| psql postgresql://"${pguser}:${pgpass}@${pghost}:${pgport}" --echo-all

#psql -h"${pghost}" -p"${pgport}" -o -W -c "select 1 from dual" >&2
#pgsql -h"${pghost}" -p"${pgport}" -o "${pguser}" -W"${pgpass}" -1aeb 
