#!/bin/sh

exec 1>&2 # all the noise

cd "${HUAUTLA_BACKUP_RESTORE_DIR:-/pgbackup}" || exit 1

if ! test -e "$RESTORE_POINT"; then 
  echo "the requested restore point: '$RESTORE_POINT' does not exist"
  exit 2
fi

: ${HUAUTLA_RESTORE_DEST_PORT:=5432}
: ${HUAUTLA_RESTORE_DEST_USER:=postgres}

echo "$HUAUTLA_RESTORE_DEST_HOST:${HUAUTLA_RESTORE_DEST_PORT}:huautla:${HUAUTLA_RESTORE_DEST_USER}:${HUAUTLA_RESTORE_DEST_PASS}" >~/.pgpass
chmod 600 ~/.pgpass

pg_restore \
  -h"${HUAUTLA_RESTORE_DEST_HOST}" \
  -p"${HUAUTLA_RESTORE_DEST_PORT}" \
  -U"${HUAUTLA_RESTORE_DEST_USER}" \
  -Fc \
  --clean \
  -dhuautla \
  -ev \
  "$RESTORE_POINT"

echo "$RESTORE_POINT: finished processing archive with result '$?'"
