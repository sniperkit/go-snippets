# ght2dm

`ght2dm` is a CLI tool for importing GHTorrent data into DevMine database.

**Important notice:** This is a specific tool for a specific need. It is not
meant to be generic or used to serve any other purpose than filling DevMine
database with GHTorrent dumps.

## Installation

```
go get -u github.com/DevMine/ght2dm
```

## Usage

*IMPORTANT:* Make sure to run the script `db/create_tmp_table.sql` before
running `ght2dm`. The latter assumes that this table is created and empty. Make
also sure to run the script `db/insert_from_tmp_tables.sql` after the execution
of `ght2dm`. Note that these two scripts are only useful when importing
`repositories`.

`ght2dm` usage is pretty simple: it only requires to pass a configuration file
as argument:

```
ght2dm ght2dm.conf
```

Here is a configuration file example:

```
{
    "ghtorrent_folders": [
        "/path/to/bson/users",
        "/path/to/bson/org_members",
        "/path/to/bson/repos",
        "/path/to/bson/repo_collaborators"
    ],
    "devmine_database": {
        "host": "localhost",
        "port": 5432,
        "user": "postgres",
        "password": "postgres",
        "database": "devmine",
        "ssl_mode": "disable"
    }
}
```

The `bson` files must be organized as follow:

```
.
├── org_members
│   └── 2012-09-29.bson
├── repo_collaborators
│   └── 2012-09-29.bson
├── repos
│   └── 2012-09-29.bson
└── users
    └── 2012-09-29.bson
```

Where each folder must have the name of the corresponding GHTorrent entity.
The currently supported entities are:

* users
* org_members
* repos
* repo_collaborators

Each `bson` dump must be named according to its creation date and using the
format `yyyy-mm-dd`. Files that does not respect this convention are skipped.
