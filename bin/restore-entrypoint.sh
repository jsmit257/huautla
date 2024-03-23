#!/bin/sh

exec 1>&2 # all the noise

cd "${BACKUP_DIR:-/pgbackup}" || exit 1

if ! test -e "$RESTORE_POINT"; then 
  echo "the requested restore point: '$RESTORE_POINT' does not exist"
  exit 2
fi

echo "$DEST_HOST:${DEST_PORT:-5432}:huautla:${DEST_USER:-postgres}:${POSTGRES_PASSWORD}" >~/.pgpass
chmod 600 ~/.pgpass

pg_restore \
  -h"${DEST_HOST}" \
  -p"${DEST_PORT:-5432}" \
  -U"${DEST_USER:-postgres}" \
  -Fc \
  --clean \
  -dhuautla \
  -ev \
  "$RESTORE_POINT"

echo "$RESTORE_POINT: finished processing archive with result '$?'"
