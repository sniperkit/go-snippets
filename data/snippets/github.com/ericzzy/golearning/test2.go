package main

import (
    "bytes"
    "fmt"
    "github.com/ugorji/go/codec"
    "encoding/json"
)

func main() {
    mh := &codec.MsgpackHandle{RawToString: true}
    data := map[string]interface{}{"abc": "hello", "abc2": 12345, "abc3": 1.2345}
    buf := &bytes.Buffer{}
    enc := codec.NewEncoder(buf, mh)
    enc.Encode(data)
    fmt.Printf("%x", buf)
    var v map[string]interface{}
    dec := codec.NewDecoder(buf, mh)
    dec.Decode(&v);
    fmt.Printf("%v", v)
    //fmt.Printf("%s", v["abc"])
    b, _ := json.Marshal(v)
    fmt.Printf("%s", b)
}
