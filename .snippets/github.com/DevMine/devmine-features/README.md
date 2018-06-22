devmine-features
================

Features computation


## Contribute

This project uses mainly `python` in version 3.
The easiest way to have every required libraries is to use a virtual
environment:

* Install `virtualenv` if necessary.
* Set up a virtual environment: `virtualenv -p python env` (replace `python`
  with `python3` if `python 3` is not your default `python` version).
* Activate it: `source env/bin/activate`.
* Install the required libraries using `pip`:
  `pip install -r requirements.txt`
* When contributing to the project, you also need to install development
  requirements:
  `pip install -r requirements_dev.txt`


When contributing, make sure that your changes are conform to PEP8 by running
`invoke pep8`. You may also want to do a static analysis of the code:
`invoke pyflakes`. To run a full check (both PEP8 and static analysis), run:
`invoke check`



## Setting up the dataset
First, follow the steps described under the `Contribute` section.

Now download and decompres the latest
[mysql dump](http://ghtorrent.org/downloads.html) from the
[GHTorrent](http://ghtorrent.org). This file should be placed inside the
dataset folder with the name `mysql`.

Run `invoke parse_mysql` to generate the tables.

## Retrieving information from tables

* To list the available tables:
  `invoke list_tables`

* To list the available fields in a given table:
  `invoke list_fields <table_name>`, for example `invoke list_fields users`

* To extract fields from a table:
  `invoke get_fields table_file output_file "field1 field2..."`,
  for example `invoke get_fields dataset/tables/users dataset/users.txt
  "id login"`
