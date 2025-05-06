#!/bin/sh

exec 1>&2

: ${HUAUTLA_BACKUP_SOURCE_PORT:=5432}
: ${HUAUTLA_BACKUP_SOURCE_USER:=postgres}

cd "${HUAUTLA_BACKUP_RESTORE_DIR:-/pgbackup}" || exit 1

echo "$HUAUTLA_BACKUP_SOURCE_HOST:${HUAUTLA_BACKUP_SOURCE_PORT}:huautla:${HUAUTLA_BACKUP_SOURCE_USER}:$HUAUTLA_BACKUP_SOURCE_PASS" >~/.pgpass
chmod 600 ~/.pgpass

timestamp="`date +'%Y%m%dT%H%M%SZ'`"

if ! pg_dump \
  -h"${HUAUTLA_BACKUP_SOURCE_HOST}" \
  -p"${HUAUTLA_BACKUP_SOURCE_PORT}" \
  -U"${HUAUTLA_BACKUP_SOURCE_USER}" \
  -Fc \
  -f"$timestamp" \
  --compress=9 \
  --clean \
  -v \
  huautla; then
  echo "postgres backup command failed"
  exit 1
fi

# guessing this works right across OSes; either way, the timestamped ones work
ln -svf "$timestamp" "latest"

du -sh "$timestamp"
