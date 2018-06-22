package main

import (
        "container/ring"
        "fmt"
        "time"
)
// golang container模块除了ring（环形链表）
// 还有list 双向链表  与 heap 堆
func main() {
        coffee := []string{"kenya", "guatemala", "ethiopia"}

        // create a ring and populate it with some values
        r := ring.New(len(coffee))
        for i := 0; i < r.Len(); i++ {
                r.Value = coffee[i]
                r = r.Next()
        }

        // print all values of the ring, easy done with ring.Do()
        r.Do(func(x interface{}) {
                fmt.Println(x)
        })

        // .. or each one by one.
        for _ = range time.Tick(time.Second * 1) {
                r = r.Next()
                fmt.Println(r.Value)
        }
}
