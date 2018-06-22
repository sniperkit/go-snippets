

Package kv provides the common operations of the k/v database.
Usage:

```go
import (
	_ "github.com/fengyfei/nuts/kv/store/bbolt"
	"github.com/fengyfei/nuts/kv"
)
```

Use it like this:

```go
kv.DBStore.Put("user", 1234, "1234")
v, _ := kv.DBStore.Get("user", 1234)
kv.DBStore.Delete("user", test)
```

The default name of the db file is the name of database which you really imported and it is in the current directory. If you want to name the db file by youself , you can do this :

```go
kv.DBStore.DB("./user/user.db")

kv.DBStore.Put("test", 1234, "1234")
```

