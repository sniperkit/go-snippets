#!/bin/sh
BIN_FILE="/opt/bin/data-profiler-server"
CONF_FILE="/etc/data-profiler"
SQL_FILE="/opt/src/github.com/giagiannis/data-profiler/data-profiler-server/database.sql"
DATABASE="$(grep database $CONF_FILE | awk '{print $2}')"

# if the database does not exist, create it!
[ ! -f "$DATABASE" ] && sqlite3 $DATABASE < $SQL_FILE

# run the daemon
exec "$BIN_FILE" "$CONF_FILE"
