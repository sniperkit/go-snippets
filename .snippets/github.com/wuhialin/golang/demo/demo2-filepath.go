/*
* @Author: wuhailin
* @Date:   2017-10-06 14:25:15
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-06 14:47:09
 */

package main

import (
	"path/filepath"
)

func main() {
	if matches, err := filepath.Glob("/tmp/*"); err == nil {
		for _, m := range matches {
			println(m)
		}
	}
}
