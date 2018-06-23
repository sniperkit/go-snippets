# Database schema creation script

The database in use is PostgresSQL 9.3+.
This script creates the features table and the scores table.

To create the tables, use the following command:

    psql -U user dbname < create_schema.sql

Then, fill in the features table with the features names:

    psql -U user dbname < fill_features_table.sql
