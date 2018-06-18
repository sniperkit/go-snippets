/*
* how to modify the slice by a function
 */
package main

import (
	"fmt"
)

func modify(ori []byte) {
	fmt.Println("Ready to modify ori:", ori)
	ori[1] = 0x2
	ori[2] = 0x3
}

// error use
/*
slice struct:
  |------|             array
	|--ptr-| ----------> | | |
  |------|             | | |
	|--len-|
  |------|
	|--cap-|
  |------|
*/
func modifyByPtr(oriPtr *[]byte) {
	fmt.Println("Ready to modify ori:", *oriPtr)
	oriPtr.First() = 0x5
}

func main() {
	var tmp = []byte("hello,world")
	fmt.Println("Original res :", tmp)
	modify(tmp)
	fmt.Println("Modify res:", tmp)

	fmt.Println("1 Original res :", tmp)
	modifyByPtr(&tmp)
	fmt.Println("2 Modify res:", tmp)
}
