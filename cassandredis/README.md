Cassandredis
============

A Redis <-> Cassandra bridge. which interprets the protocol and maps some commands to a Cassandra datastructure.

Concept
-------

You take your favorite Redis client, connect to Cassandredis, send it *supported* Redis commands (see below for all supported commands) and it will automagically transform it into a appropriate Cassandra CQL query.

You don't have to set up anything schema side on Cassandra except for the keyspace, everything else will be taken care of by Cassandredis.

You do have to follow a certain set of rules for your key names (see below).

Supported commands
------------------

* LPUSH
* LRANGE (limited to sending all the data, no real range)

Running it
----------

    go get github.com/gophergala/cassandredis

    cassandredis cassandra1.example.com:9042 mykeyspace

By default it listens on *0.0.0.0:8765*.

You now can use your favorite Redis client
