# goXA - Transaction Handler

The go package that easily handle the transactions across multiple chained services.


## How to use

### Receiver side:

```go
import (
    "goxa"
)

func main() {
    xa, err := goxa.Listen("tcp", ":8080")
    if err != nil {
        // error handling
    }
    buffer := make([]byte, 1024)
    id, buffer, count, err := goxa.Receive(xa)
    if err != nil {
        // error handling
    }

    well_handled := true
    // do something here

    if well_handled {
        goxa.Commit(xa, id)
    } else {
        goxa.Rollback(xa, id)
    }
    goxa.Close(xa)
}
```

### Sender side:

```go
import (
    "goxa"
)

func main() {
    xa, err := goxa.Connect("tcp", "127.0.0.1:8080")
    if err != nil {
        // error handling
    }
    count, err := goxa.Send(xa, xa, []byte("hello"), true)
    if err != nil {
        // error handling
    }
    goxa.Disconnect(xa)
}
```

### Run the test

```bash
$cd LOCATION_OF_GOXA_PACKAGE
$go test
```
