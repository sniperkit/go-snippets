---
layout: post
title: "PostgreSQL : from a shell"
tagline: "PostgreSql... is awkward from the command line."
category : notes
published: false
tags : [PostgreSQL]
---

* [SQL language &ndash; postgresql.org/docs/current/static/sql.html](https://www.postgresql.org/docs/current/static/sql.html)
    - [Data types &ndash; postgresql.org/docs/current/static/datatype.html](https://www.postgresql.org/docs/current/static/datatype.html)
    - [Constraints &ndash; postgresql.org/docs/current/static/ddl-constraints.html](https://www.postgresql.org/docs/current/static/ddl-constraints.html)
* [postgresql.org/docs/current/static/sql-comment.html](https://www.postgresql.org/docs/current/static/sql-comment.html)
* [PostgreSQL @ ArchLinux wiki](https://wiki.archlinux.org/index.php/PostgreSQL)

## Install

```bash
# pacman -S postgresql
# su - postgres   (Or: `sudo -u postgres -i`)
```

Initialize database. Here __$LANG = en_US.UTF-8__.

```bash
[postgres@wall ~]$ initdb --locale $LANG -E UTF8 -D '/var/lib/postgres/data'
```

Back as __root__, start and enable the systemd unit:

```bash
[root@wall ~]# systemctl start postgresql
[root@wall ~]# systemctl status postgresql
[root@wall ~]# systemctl enable postgresql
```

## Create user

* [createuser](https://www.postgresql.org/docs/current/static/app-createuser.html)
* [Database roles §21.1](https://www.postgresql.org/docs/10/static/database-roles.html)

```bash
[postgres@wall ~]$ createuser --pwprompt --echo \
  [--createdb] [--createrole] [--superuser]
    fabi

Enter password for new role: XXXX
Enter it again: XXXX

SELECT pg_catalog.set_config('search_path', '', false)
CREATE ROLE fabi PASSWORD 'md5f56...' NOSUPERUSER CREATEDB NOCREATEROLE INHERIT LOGIN;
```

Test connectivity by listing databases as user `fabi`:

```bash
$ psql -l
```

Create the "default" database for role/user `fabi` (_optional_) :

```bash
[fabi@wall] ~ $ createdb -e fabi
```

## List all databases

```bash
$ psql -l
```

## Database administration

* todo: users and roles, confusion.

```bash
$ sudo -u postgres -i
```

Create a new database named `clong` owned by ~~user~~ role `fabi` :

```bash
$ createdb --owner=fabi clong
```

## Pointers

* <http://www.sqlines.com/postgresql/datatypes/serial>
* [Postgres-cheatsheet by Kartones (Gist)](https://gist.github.com/Kartones/dd3ff5ec5ea238d4c546)

__EOF__
