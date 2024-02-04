#!/bin/bash

. ./bin/vars.sh

cat ./sql/create.pgsql ./sql/init.pgsql ./sql/seed.pgsql \
| "${psql_cmd[@]}" --echo-all -v ON_ERROR_STOP=1
