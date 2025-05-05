#!/bin/bash

exec 1>&2

if test "`ls -1 /var/lib/postgresql/data | wc -l`" -ne "0"; then
  cat <<-EOF
  refusing to migrate over an existing database; remove 'data/' 
  directory and try again
EOF
  exit 1
fi

/usr/local/bin/docker-entrypoint.sh postgres & 

sleep 2s

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

pg_dump \
  -h${HUAUTLA_MIG_SOURCE_HOST} \
  -p${HUAUTLA_MIG_SOURCE_PORT:-5432} \
  -U${HUAUTLA_MIG_SOURCE_USER:-postgres} \
  huautla \
| "${psql_cmd[@]}" -dhuautla
