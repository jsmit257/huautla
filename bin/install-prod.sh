#!/bin/bash

. ./bin/vars.sh

cat ./sql/create.sql ./sql/init.sql ./sql/seed.sql \
| "${psql_cmd[@]}" --echo-all -v ON_ERROR_STOP=1
