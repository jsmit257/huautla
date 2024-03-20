#!/bin/bash

. ./bin/vars.sh

"${psql_cmd[@]}" --echo-all -v ON_ERROR_STOP=1 <./sql/seed-system-test.pgsql
