# Database schema creation script

The database in use is PostgresSQL 9.3+.
This script creates the commits table and the commit diff deltas table.
`repotool` also need access to the users and repositories table as created by
[crawld](http://devmine.ch/doc/crawld/).

To create the tables, use the following command (providing that you
already have the database as defined by `crawld`):

    psql -U user dbname < create_schema.sql
