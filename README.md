# godb-tutorial
database tutorial with Golang, especially for key-value storage.

When running redis client, `redis` must be started before.

```bash
brew update && brew install redis && brew services run redis
```
## dependency 
1. dep, <https://github.com/golang/dep>
2. go-redis/redis, <https://github.com/go-redis/redis> (godoc -- <https://godoc.org/github.com/go-redis/redis>)
3. boltdb/bolt, <https://github.com/boltdb/bolt> (godoc -- <https://godoc.org/github.com/boltdb/bolt>)


## resources

`bolt` related:
> 1. Using Boltdb as a fast, persistent key value store, <https://awmanoj.github.io/2016/08/03/using-boltdb-as-a-fast-persistent-kv-store/> 
> 2. Intro to BoltDB: Painless Performant Persistence, <https://npf.io/2014/07/intro-to-boltdb-painless-performant-persistence/>
> 3. WTF Dial: Data storage with BoltDB, <https://medium.com/wtf-dial/wtf-dial-boltdb-a62af02b8955>


`redis` related:
> 1. Commands, <https://redis.io/commands>
> 2. Transactions, <https://redis.io/topics/transactions>
> 3. Working with Redis in GO, <http://www.alexedwards.net/blog/working-with-redis>