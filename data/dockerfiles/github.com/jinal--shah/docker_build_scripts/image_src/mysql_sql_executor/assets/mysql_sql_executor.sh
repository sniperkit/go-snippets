#!/bin/bash
# vim: ts=4 sw=4 et sr smartindent:
#
# To run, set the following vars in your environment:
#
# $DBHOST: location of remote DB
# $DBPORT: listening mysql port
# $DBUSER: the user to connect as
# $DBPASS: the user's pword.
# $SQL_S3_PATH: the full s3:// uri to the sql file
#
SQL_LOCAL=/var/tmp/my.sql
SQL_ZIPPED=$SQL_LOCAL.gz
MYSQL_OPTS="--host=$DBHOST --port=$DBPORT --user=$DBUSER --password=$DBPASS $DBNAME"
REQUIRED_VARS="
    AWS_REGION
    DBHOST
    DBPASS
    DBPORT
    DBUSER
    SQL_S3_PATH
"

function check_var_defined() {
    var_name="$1"
    var_val="${!var_name}"
    if [[ -z $var_val ]]; then
        echo "$0 ERROR: You must pass \$$var_name to this script" >&2
        FAILED_VALIDATION="you bet'cha"
        return 1
    fi
}

# ... validate required vars
for this_var in $REQUIRED_VARS; do
    check_var_defined $this_var
done

if [[ ! -z $FAILED_VALIDATION ]]; then
    echo "$0 ERROR: FAILURE. One of more required vars not passed to this script." >&2
    echo "$0 ERROR: required vars:" $REQUIRED_VARS >&2
    exit 1
fi

# ... test can connect to specific db
if ! mysql $MYSQL_OPTS -e 'select "1";'
then
    echo "$0 ERROR: couldn't connect to mysql host:port<$DBHOST:$DBPORT> db $DBNAME"
    exit 1
fi

# ... fetch sql
if ! aws --region $AWS_REGION s3 cp "$SQL_S3_PATH" $SQL_ZIPPED
then
    echo "$0 ERROR: couldn't download $SQL_S3_PATH." >&2
    exit 1
fi

# ... decompress
if ! gunzip $SQL_ZIPPED
then
    echo "$0 ERROR: couldn't gunzip local file $SQL_ZIPPED from $SQL_S3_PATH." >&2
    exit 1
fi

# ... run
mysql $MYSQL_OPTS < $SQL_LOCAL

