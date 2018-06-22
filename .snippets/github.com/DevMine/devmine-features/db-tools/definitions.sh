#/bin/bash
set -o nounset
set -o errexit

DB=${DB-devmine}
DB_USER=${DB_USER-devmine}
SCORES_TABLE="SCORES"
USERS_TABLE="USERS"
FEATURES_TABLE="FEATURES"
