package main
import (
		"fmt"
		"gopkg.in/vmihailenco/msgpack.v2"
                "sort"
       ) 

func main() {
in := map[string]interface{}{"foo": 1, "hello": "world"}
    b, err := msgpack.Marshal(in)
	    if err != nil {
		    fmt.Printf("The error is %v", err);
	    }
    fmt.Printf("The marshal is %v", b);
    var out map[string]interface{}
    err = msgpack.Unmarshal(b, &out)

	    var outKeys []string
	    for k := range out {
		    outKeys = append(outKeys, k)
	    }
    sort.Strings(outKeys)

	    fmt.Printf("err: %v\n", outKeys)
}
