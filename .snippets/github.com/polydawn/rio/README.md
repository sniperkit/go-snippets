rio
===

**R**epeatable **I**/**O**.

Rio is part of the [Timeless Stack](https://github.com/polydawn/timeless) --
it handles packing, unpacking, syncing, copying, and hashing of Filesets.

In other words, Rio is tooling for repeatable, reproducible, filesystem I/O and transport.
It's both a command line tool, and usable as a series of libraries.

(You may also want to check out [Repeatr](https://github.com/polydawn/repeatr),
which is the container executor in the Timeless Stack -- it uses Rio under the hood
to provide snapshotting and decentralized sync of container filesystem images.)


What?
-----

Okay, by example:

```
rio pack tar /tmp/something/ --target=file://output.tgz
```

That just created a tar pack of the files you aimed it at.

It also gave you a hash on stdout.  That's a
[content-addressable](https://en.wikipedia.org/wiki/Content-addressable_storage)
ID of the pack you just created.

You can use that ID later to unpack things:

```
rio unpack tar:879UrF8j7E[...]udF57KpF8 \
   /tmp/unpackhere/ --source=file://output.tgz
```

It's long.  That's because it's a [cryptographic hash](https://en.wikipedia.org/wiki/Cryptographic_hash_function).

Why is this neat?

You can put that packed data on any server, anywhere, and fetch it again by hash.
The hash makes it immutable, and reproducible, even if you don't control the storage.

Turning it up to 11
-------------------

```
rio pack tar /tmp/something/ --target=ca+32://mybucket.s3.amazonaws.com/assets/
```

Same drill except... we just:

- uploaded to a cloud storage host
- without a specific name -- the hash is used to organize storage
- and still got the same hash.

Now, hand that hash to someone else:

```
rio unpack tar:879UrF8j7E[...]udF57KpF8 \
   /tmp/unpackhere/ --source=ca+https:///mybucket.s3.amazonaws.com/assets/
```

They get your files.  Boom.  Huge amounts of data.  Just one handle to copypaste: that hash.
