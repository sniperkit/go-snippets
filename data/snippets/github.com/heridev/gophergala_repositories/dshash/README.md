# dshash

A small research to learn about hashes and how to implement one on top of 
appengine datastore.

I want to build a library that will allow you to store a hash bigger than 10MB
which will be serialized and stored across multiple appengine datastore entities.

If I search for a key it should de-serialize the data only in the bucket containing
the key.

No idea if this is something useful this is just for fun.
