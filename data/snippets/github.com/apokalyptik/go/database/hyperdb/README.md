HyperDB
=======

## Documentation

<a href="http://godoc.org/github.com/apokalyptik/go/database/hyperdb">godoc.org/github.com/apokalyptik/go/database/hyperdb</a>

## Example Usage

```go
var hdb *hyperdb.DB
var dbHandle *sql.DB
hdb = hyperdb.New("dc1")
hdb.AddServer(hyperdb.Server{
  Dataset: "users",
  Datacenter: "dc1",
  // etc...
})
dbh = hyperdb.Read("users")
```
