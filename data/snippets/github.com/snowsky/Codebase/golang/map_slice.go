package main

import "fmt"

func main() {
	var s []int
	s = append(s, 1)

	var m map[string]int
	m["one"] = 1 //error

	// 可以指定map大小，但运行时不能cap
	m := make(map[string]int, 99)
	//cap(m) //error
	fmt.Println(len(m)) //0

	// cannot use nil as type string in assignment
	//invalid operation: x == nil (mismatched types string and nil)
	/* should be
		  var x string //defaults to "" (zero value)

	  	if x == "" {
	  	    x = "default"
	  	}
	*/
	var x string = nil //error
	if x == nil {      //error
		x = "default"
	}
}
