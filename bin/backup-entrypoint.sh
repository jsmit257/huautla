#!/bin/sh

exec 1>&2

cd "${BACKUP_DIR:-/pgbackup}" || exit 1

echo "$SOURCE_HOST:${SOURCE_PORT:-5432}:huautla:${SOURCE_USER:-postgres}:$POSTGRES_PASSWORD" >~/.pgpass
chmod 600 ~/.pgpass

timestamp="`date +'%Y%m%dT%H%M%SZ'`"

pg_dump \
  -h"${SOURCE_HOST}" \
  -p"${SOURCE_PORT:-5432}" \
  -U"${SOURCE_USER:-postgres}" \
  -Fc \
  -f"$timestamp" \
  --compress=9 \
  --clean \
  -v \
  huautla \

# guessing this works right across OSes; either way, the timestamped ones work
ln -svf "$timestamp" "latest"

du -sh "$timestamp"
