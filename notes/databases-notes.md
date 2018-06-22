---
layout: page
title: "Databases - notes"
tagline: "playground"
category : notes
tags : [draft, wide]
published: true
maths: true
---

* [SchemaSpy](git@github.com:fabic/SchemaSpy-ng.git) : how to & your repo remember ?

![NoSQL: Not Only SQL](/assets/nosql--not-only-sql.png "NoSQL: Not Only SQL.")

## Postgresql

```bash
$ sudo pacman -S postgresql{,-docs} php-pgsql phppgadmin
```

### Postgres :: dump database

This one command is what you actually want :

```bash
$ pg_dump my_database | bzip2 -c > tmp/my_database.db.dump.sql.bz2
```

There's also the `-F tar` argument to `pg_dump` which will produce a tarball:

```bash
$ pg_dump -F tar my_database > `date -Is`--my_database.db-dump.sql.tar
```

But note that it contains several files, probably one per table:

```bash
$ tar -tf my_database.db-dump.sql.tar
toc.dat
2660.dat
2662.dat
...
...
2733.dat
restore.sql
```

## MySQL / MariaDB

### Users, security, access control

#### List all users with remote access

```
MariaDB [(none)]> SELECT User, Host FROM mysql.user WHERE Host <> 'localhost';
+------+-----------+
| User | Host      |
+------+-----------+
| root | 127.0.0.1 |
| root | ::1       |
|      | wall      |
| root | wall      |
+------+-----------+
4 rows in set (0.00 sec)
```

#### List all users that do not have remote access

This would list users that _have_ localhost access _but no_ remote access:

```sql
SELECT u.User, u.Host, u.Password
FROM mysql.user u
WHERE NOT EXISTS (
    SELECT * FROM mysql.user uu
    WHERE uu.User = u.User
      AND uu.Host NOT IN ('localhost', '127.0.0.1', '::1')
  ) ;
```

```
+------+-----------+-------------------------------------------+
| User | Host      | Password                                  |
+------+-----------+-------------------------------------------+
| fabi | localhost | *12344E8CABCDEFA5812346BABCDFE90FEDCBA123 |
+------+-----------+-------------------------------------------+
1 row in set (0.00 sec)
```

Then one may add a corresponding user with remote access _using the same password_ :

```sql
CREATE USER fabi@'%'
  IDENTIFIED BY PASSWORD '*12344E8CABCDEFA5812346BABCDFE90FEDCBA123' ;
```

## Algorithms

### B-trees, B+ trees

__[Wikipedia : B+ tree](https://en.wikipedia.org/wiki/B%2B_tree)__

* [2014: Oracle B-Tree Index: From the concept to Internals](http://www.toadworld.com/platforms/oracle/w/wiki/11001.oracle-b-tree-index-from-the-concept-to-internals)
via [SO: How to lay out B-Tree data on disk?](https://stackoverflow.com/a/40740893)
* [2014: B+Trees – How SQL Server Indexes are Stored on Disk](http://sqlity.net/en/2445/b-plus-tree/)
