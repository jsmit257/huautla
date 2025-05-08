#!/bin/bash

exec 1>&2

if ls /var/lib/postgresql/data/* >/dev/null 2>&1; then
  cat <<-EOF
  refusing to migrate over an existing database; remove 'data/' 
  directory and try again
EOF
  exit 1
fi

/usr/local/bin/docker-entrypoint.sh postgres & 

sleep 5s # saw random but frequent failures when waiting 2s, next time try 5

psql_cmd=( 
  "psql" 
  -h"${HUAUTLA_MIG_DEST_HOST:-localhost}" 
  -p"${HUAUTLA_MIG_DEST_PORT:-5432}" 
  -U"${HUAUTLA_MIG_DEST_USER:-postgres}" 
)

"${psql_cmd[@]}" <<-EOF
  CREATE USER huautla;
  CREATE DATABASE huautla;
  GRANT ALL PRIVILEGES ON DATABASE huautla TO huautla;
EOF

sc=$?

if [ "$sc" -ne "0" ]; then
  echo "command failed: ${psql_cmd[@]}"
  exit 1
fi

pg_dump \
  -h${HUAUTLA_MIG_SOURCE_HOST} \
  -p${HUAUTLA_MIG_SOURCE_PORT:-5432} \
  -U${HUAUTLA_MIG_SOURCE_USER:-postgres} \
  huautla \
  | "${psql_cmd[@]}" -dhuautla

sc=$?

if [ "$sc" -ne "0" ]; then
  echo "migrating data from source to persistence failed"
  exit 1
fi

