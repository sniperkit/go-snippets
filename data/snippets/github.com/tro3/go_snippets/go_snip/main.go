package main

/*
Example usages in go file:

//go:generate go_snip testhelper pkgname
//go:generate go_snip typedlist pkgname ObjectName
//go:generate go_snip typedlistmap pkgname ObjectName TargetName
*/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var root string = "src/github.com/tro3/go_snippets"

var router map[string]func() = map[string]func(){
	"testhelper":   handleTestHelper,
	"typedlist":    handleTypedList,
	"typedlistmap": handleTypedListMap,
}

func main() {
	fn, ok := router[os.Args[1]]
	if !ok {
		panic("Unknown snippet: " + os.Args[1])
	}
	fn()
}

func capVariants(in string) (orig, cap, low string) {
	orig = in
	low = strings.ToLower(orig)
	cap = strings.Title(low)
	return
}

func generateFile(in, out string, replace map[string]string) {
	c := getAndReplace(in, replace)
	_ = ioutil.WriteFile(out, c, 0666)
	fmt.Printf("Generated %s\n", out)
}

func getAndReplace(fname string, replace map[string]string) []byte {
	c, _ := ioutil.ReadFile(filepath.Join(os.Getenv("GOPATH"), root, fname))
	for key, val := range replace {
		c = bytes.Replace(c, []byte(key), []byte(val), -1)
	}
	return c
}

func handleTestHelper() {
	proj := os.Args[2] // Target project name

	generateFile(
		"testhelper.go",
		"testhelper.go",
		map[string]string{
			"snippets": proj,
		},
	)
}

func handleTypedList() {
	proj := os.Args[2]                            // Target project name
	kind, ckind, lkind := capVariants(os.Args[3]) // Target object name

	generateFile(
		"typedlist.go",
		fmt.Sprintf("%s_list.go", lkind),
		map[string]string{
			"snippets": proj,
			"TYPEA":    kind,
			"TypeA":    ckind,
		},
	)
}

func handleTypedListMap() {
	proj := os.Args[2]                            // Target project name
	kind, ckind, lkind := capVariants(os.Args[3]) // Target object name
	targ, ctarg, ltarg := capVariants(os.Args[4]) // Target resulting object name

	generateFile(
		"typedlistmap.go",
		fmt.Sprintf("%s_list_map_%s.go", lkind, ltarg),
		map[string]string{
			"snippets": proj,
			"TYPEA":    kind,
			"TypeA":    ckind,
			"TYPEB":    targ,
			"TypeB":    ctarg,
		},
	)
}
