/*
* @Author: wuhailin
* @Date:   2017-09-29 17:56:52
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-05 17:14:39
 */

package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now().Unix()

	fmt.Println(time.Since(start).Seconds())
}
