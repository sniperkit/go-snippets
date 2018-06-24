#!/usr/bin/env python

import os
import subprocess
import sys

# The top-level entity-extractor application directory
topdir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Path (if needed) to psql command.
psql_cmd = "psql"


def in_development_env():
    """Are we running in a development environment?

    """
    return os.environ.get("ENV", "") == "development"

def psql_subprocess_args():
    """Arguments required to call psql and connect to database

    Returns a list of arguments suitable for passing to subprocess.call().

    """
    if in_development_env():
        return [psql_cmd, "entity-extractor_development"]
    else:
        dbname = "entity-extractor_production"
        dbuser = "entity-extractor"
        dbhost = "postgresql-master-1"
        return [psql_cmd, dbname, dbuser, "-h", dbhost]

def create_tables():
    """Create the database tables by reading the schema.sql file.

    Will report an error to stderr if the table already exists.

    """
    schemafile = os.path.join(topdir, "db", "schema.sql")
    subprocess.call(psql_subprocess_args() + ["-f", schemafile])

def replace_entities_table(csv_filename):
    """Replace the contents of the entities table.

    :param csv_filename: Path to a file containing all the entities in CSV
    format.

    """
    commands = """
	begin;
        delete from entities;
        copy entities from STDIN with (format csv);
        commit;
    """
    with open(csv_filename) as csv_file_object:
        subprocess.call(
            psql_subprocess_args() + ["-c", commands],
            stdin=csv_file_object,
        )


if __name__ == "__main__":
    set_schema()
    replace_entities_table(sys.argv[1])
