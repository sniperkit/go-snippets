// Copyright 2016 Gavin "Groovy" Grover. All rights reserved.
// Use of this source code is governed by the same BSD-style
// license as Go that can be found in the LICENSE file.

package parser_test

import (
	"github.com/gavingroovygrover/qu/parser"
	"go/token"
	"go/format"
	"go/types"
	"go/ast"
	"go/importer"
	"testing"
	"fmt"
)

//================================================================================================================================================================
type StringWriter struct{ data string }
func (sw *StringWriter) Write(b []byte) (int, error) {
	sw.data += string(b)
	return len(b), nil
}

func TestKanji(t *testing.T) {
	for i, tst:= range kanjiTests {
		src:= tst.key
		dst:= tst.val
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "", src, 0)
		if err != nil {
			t.Errorf("parse error in %d: %q", i, err)
		} else {
			var conf = types.Config{
				Importer: importer.Default(),
			}
			info := types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Defs:  make(map[*ast.Ident]types.Object),
				Uses:  make(map[*ast.Ident]types.Object),
			}
			_, err = conf.Check("testing", fset, []*ast.File{f}, &info)
			if err != nil {
				t.Errorf("type check error in %d: %q", i, err)
			}
			sw:= StringWriter{""}
			_= format.Node(&sw, fset, f)
			if sw.data != dst {
				t.Errorf("unexpected Go source for %d: received source: %q; expected source: %q", i, sw.data, dst)
				//t.Errorf("unexpected Go source for %d: received source:\n%s\nexpected source: %q", i, sw.data, dst)
			}
		}
	}
}

var kanjiTests = map[int]struct{key string; val string} {

// ========== ========== ========== ==========
//test keywords: 功
//test keyword scoping: 入
1:{`
package main;入"fmt"
功main(){
  fmt.Printf("Hi!\n") // comment here
}`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"

func _main() {
	_fmt.Printf("Hi!\n")
}
`},

// ========== ========== ========== ==========
//test keyword: 回
//test keyword scoping: 功
//test specids: 度,串,整,整64,漂32,漂64,复,复64,复128
2:{`
package main;入"fmt";
功main(){
  fmt.Printf("Len: %d\n", 度(fs("abcdefg")))
}
功fs(a串)串{回a+"xyz"}
功ff(a漂32)漂64{回漂64(a)}
功fc(a复64)复128{回复128(a)+复(1,1)}
功fi(a整)整64{回整64(a)}
`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"

func _main() {
	_fmt.Printf("Len: %d\n", len(_fs("abcdefg")))
}
func _fs(_a string) string        { return _a + "xyz" }
func _ff(_a float32) float64      { return float64(_a) }
func _fc(_a complex64) complex128 { return complex128(_a) + complex(1, 1) }
func _fi(_a int) int64            { return int64(_a) }
`},

// ========== ========== ========== ==========
//test keyword scoping: 变,如,否
//test specids: 真,假
3:{`
package main;入"fmt"
import 吧"fmt"
import 哪_fg"fmt"
入㕤hij"fmt"
入卟"unicode/utf8"
入叨叩kl"unicode/utf8"
var n = 50
变p=70
变string=170
功main(){
  如真{
    fmt.Printf("Len: %d\n", 度("abcdefg") + p)
  }
}
func deputy(){
  if真{
    fmt.Printf("Len: %d\n", 度("abcdefg") + n)
  }
  如假{
    fg.Printf("Len: %d\n", 度("hijk") + p)
  }否{
    hij.Printf("Len: %d\n", 度("hi") + p)
  }
  fr,_:= _utf8.DecodeRune([]byte("lmnop"))
  fmt.Printf("1st rune: %s; Len: %d\n", fr, len("lmnop") + n)
  让_,_=kl.DecodeRune([]节("lmnop"))
  㕤Printf("Fifty: %d\n", n)
  哪Printf("Fifty: %d\n", n)
  吧Printf("Fifty: %d\n", n)
}
`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"
import "fmt"
import _fg "fmt"
import _hij "fmt"
import _utf8 "unicode/utf8"
import _kl "unicode/utf8"

var n = 50
var _p = 70
var _string = 170

func _main() {
	if true {
		_fmt.Printf("Len: %d\n", len("abcdefg")+_p)
	}
}
func deputy() {
	if true {
		fmt.Printf("Len: %d\n", len("abcdefg")+n)
	}
	if false {
		_fg.Printf("Len: %d\n", len("hijk")+_p)
	} else {
		_hij.Printf("Len: %d\n", len("hi")+_p)
	}
	fr, _ := _utf8.DecodeRune([]byte("lmnop"))
	fmt.Printf("1st rune: %s; Len: %d\n", fr, len("lmnop")+n)
	_, _ = _kl.DecodeRune([]byte("lmnop"))
	_hij.Printf("Fifty: %d\n", n)
	_fg.Printf("Fifty: %d\n", n)
	fmt.Printf("Fifty: %d\n", n)
}
`},

// ========== ========== ========== ==========
//test keyword: 构
//test keyword scoping: 种,久
//test specids: 整8,整16,整32
4:{`
package main
type _string string
种A struct{a string; b 整8}
种B struct{a string; b 整16}
种C构{a string; b 整32}
type D struct{a string; b 整32}
种E构{a串;b整32}
久a=3.1416
const b=2.72
`,

// ---------- ---------- ---------- ----------
`package main

type _string string
type A struct {
	_a _string
	_b int8
}
type B struct {
	_a _string
	_b int16
}
type C struct {
	_a _string
	_b int32
}
type D struct {
	a string
	b int32
}
type E struct {
	_a string
	_b int32
}

const _a = 3.1416
const b = 2.72
`},

// ========== ========== ========== ==========
//test keywords: 围,为,继,破
//test keyword scoping: 入,图
//test specids: 节,字
//TODO: fix scoping for 图 and slices/arrays
5:{`
package main;入"fmt"
import "fmt"
type _byte byte
type _rune rune
var (
	_a= 图[字]节{'a': 127, 'b': 0, '7':7}
	_b= byte(7)
	c= 图[字]节{'a': b}
	cc= map[rune]byte{'a': _b}
	d= []节{127, b, 0}
	dd= [3]节{127, b, 0}
	e= 图[byte]rune{byte(b):'a'}
	_f= 7
	g= 构{a串;b整}{"abc",f} //TODO: why does this work, i.e. f become _f ?
	h= 功(a串)串{a="def"; return a}
	i 面{doIt(a串)串}
)
func main(){
	_zx:为i:=0;i<19;i++{
		if i==3 {继}; if i==6{破}
		如 i== 16{破zx}
		如 i== 17{继zx}
		fmt.Print(i," ")
	}
	fmt.Println("abc")
	为i:=围a{fmt.Print(i," ")}
	for i:= 0; i<28; i++ {
		if i==3 { continue }
		if i==6 { break }
		fmt.Print(i, " ")
	}
}
`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"
import "fmt"

type _byte byte
type _rune rune

var (
	_a = map[rune]byte{'a': 127, 'b': 0, '7': 7}
	_b = byte(7)
	c  = map[rune]byte{'a': _b}
	cc = map[rune]byte{'a': _b}
	d  = []byte{127, _b, 0}
	dd = [3]byte{127, _b, 0}
	e  = map[_byte]_rune{_byte(_b): 'a'}
	_f = 7
	g  = struct {
		_a string
		_b int
	}{"abc", _f}
	h = func(_a string) string { _a = "def"; return _a }
	i interface {
		_doIt(_a string) string
	}
)

func main() {
_zx:
	for _i := 0; _i < 19; _i++ {
		if _i == 3 {
			continue
		}
		if _i == 6 {
			break
		}
		if _i == 16 {
			break _zx
		}
		if _i == 17 {
			continue _zx
		}
		_fmt.Print(_i, " ")
	}
	fmt.Println("abc")
	for _i := range _a {
		_fmt.Print(_i, " ")
	}
	for i := 0; i < 28; i++ {
		if i == 3 {
			continue
		}
		if i == 6 {
			break
		}
		fmt.Print(i, " ")
	}
}
`},

// ========== ========== ========== ==========
//test keywords: 掉
//test keyword scoping: 择,事,别,面
//test specids: 双,空,绝,绝8,绝16,绝32,绝64
6:{`package main;入"fmt"
type _uint16 uint16
type A interface {
  aMeth()绝
}
种B面{
  bMeth()绝8
}
种C interface{
  cMeth(theC uint16)绝16
}
type D面{
  dMeth(theD绝32)绝64
}
func abc()*双{回空}
func main(){
	_a:=2
	择a{
	事1:
		fmt.Print('a');
	事2:
		fmt.Print('b')
		掉
	事3:
		fmt.Print('c')
	别:
		fmt.Print('d')
	}
}
`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"

type _uint16 uint16
type A interface {
	aMeth() uint
}
type B interface {
	_bMeth() uint8
}
type C interface {
	_cMeth(_theC _uint16) uint16
}
type D interface {
	_dMeth(_theD uint32) uint64
}

func abc() *bool { return nil }
func main() {
	_a := 2
	switch _a {
	case 1:
		_fmt.Print('a')
	case 2:
		_fmt.Print('b')
		fallthrough
	case 3:
		_fmt.Print('c')
	default:
		_fmt.Print('d')
	}
}
`},

// ========== ========== ========== ==========
//test keyword scoping: 选,去,通
//test special identifier: 正
7:{`包正;入("math/rand";"sync/atomic")
种readOp构{key整;resp通整}
种writeOp构{key整;val整;resp通双}
功正(){
    变ops整64=0
    reads:=造(通*readOp)
    writes:=造(通*writeOp)
    去功(){
        变state=造(图[整]整)
        为{选{事read:=<-reads:read.resp<-state[read.key]
              事write:=<-writes:state[write.key]=write.val;write.resp<-真
             }}}()
    为r:=0;r<100;r++{
        去功(){
            为{read:=&readOp{key:rand.Intn(5),resp:造(通整)}
               reads<-read
               <-read.resp
               atomic.AddInt64(&ops,1)
              }}()}
    为w:=0;w<10;w++{
        去功(){
            为{write:=&writeOp{key:rand.Intn(5),val:rand.Intn(100),resp:造(通双)}
               writes<-write
               <-write.resp
               atomic.AddInt64(&ops,1)
              }}()}
    时Sleep(时Second)
    opsFinal:=atomic.LoadInt64(&ops)
    形Println("ops:",opsFinal)

    让range:="abc" //when used with 让, Go keywords like "range" can be used as identifiers
    让range="abcdefg"
	形Printf("range: %v\n",range)
}
`,

// ---------- ---------- ---------- ----------
`package main

import (
	fmt "fmt"
	_rand "math/rand"
	_atomic "sync/atomic"
	time "time"
)

type _readOp struct {
	_key  int
	_resp chan int
}
type _writeOp struct {
	_key  int
	_val  int
	_resp chan bool
}

func main() {
	var _ops int64 = 0
	_reads := make(chan *_readOp)
	_writes := make(chan *_writeOp)
	go func() {
		var _state = make(map[int]int)
		for {
			select {
			case _read := <-_reads:
				_read._resp <- _state[_read._key]
			case _write := <-_writes:
				_state[_write._key] = _write._val
				_write._resp <- true
			}
		}
	}()
	for _r := 0; _r < 100; _r++ {
		go func() {
			for {
				_read := &_readOp{_key: _rand.Intn(5), _resp: make(chan int)}
				_reads <- _read
				<-_read._resp
				_atomic.AddInt64(&_ops, 1)
			}
		}()
	}
	for _w := 0; _w < 10; _w++ {
		go func() {
			for {
				_write := &_writeOp{_key: _rand.Intn(5), _val: _rand.Intn(100), _resp: make(chan bool)}
				_writes <- _write
				<-_write._resp
				_atomic.AddInt64(&_ops, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	_opsFinal := _atomic.LoadInt64(&_ops)
	fmt.Println("ops:", _opsFinal)

	_range := "abc"
	_range = "abcdefg"
	fmt.Printf("range: %v\n", _range)
}
`},

// ========== ========== ========== ==========
//test keyword scoping: 围
//test special keyword: 做
8:{`package main

import (
	"fmt"
	"github.com/gavingroovygrover/qu/parser"
	"go/token"
	"go/format"
	//"go/ast"
	//"os"
)

func main() {
	形Printf("Hi!\n")

	fset := token.NewFileSet() // positions are relative to fset

	//_f, err := parser.ParseFile(fset, "src/github.com/gavingroovygrover/qu/first.go", nil, parser.ImportsOnly)
	//_f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	_f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := 围 f.Imports {
		fmt.Println(s.Path.Value)
	}

	//ast.Print(fset, _f)
	//_= format.Node(os.Stdout, fset, _f)

	{	sw:= StringWriter{""}
		_= format.Node(&sw, fset, _f)
		fmt.Println(sw.data)
		fmt.Println(sw.data == dst)
	}
	做{ abc:= "abc"
		_ = abc
	}
}

type StringWriter struct{ data string }
func (sw *StringWriter) Write(b []byte) (int, error) {
	sw.data += string(b)
	return len(b), nil
}

var src string = "package main"
var dst string = "package main"`,

// ---------- ---------- ---------- ----------
`package main

import (
	"fmt"
	"github.com/gavingroovygrover/qu/parser"
	"go/format"
	"go/token"
)

func main() {
	fmt.Printf("Hi!\n")

	fset := token.NewFileSet()

	_f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range _f.Imports {
		fmt.Println(s.Path.Value)
	}

	{
		sw := StringWriter{""}
		_ = format.Node(&sw, fset, _f)
		fmt.Println(sw.data)
		fmt.Println(sw.data == dst)
	}
	{
		_abc := "abc"
		_ = _abc
	}
}

type StringWriter struct{ data string }

func (sw *StringWriter) Write(b []byte) (int, error) {
	sw.data += string(b)
	return len(b), nil
}

var src string = "package main"
var dst string = "package main"
`},

// ========== ========== ========== ==========
//test keyword scoping: 包
//test special keywords: 让,任
//test using keyword as id in kanji-context (both LHS and RHS)
9:{`
package main;入"fmt"
import "fmt"
const a=6
功main(){
  让b:=7
  让func:=8
  fmt.Printf("Hi, nos.%s and %s!\n", b, func)
  英{
    fmt.Printf("Hi, no. %s.\n", a)
  }
}
func baba(){
  变b任= 17
  fmt.Printf("Hi, no.%s!\n", _b)
}`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"
import "fmt"

const a = 6

func _main() {
	_b := 7
	_func := 8
	_fmt.Printf("Hi, nos.%s and %s!\n", _b, _func)
	{
		fmt.Printf("Hi, no. %s.\n", a)
	}
}
func baba() {
	var _b interface{} = 17
	fmt.Printf("Hi, no.%s!\n", _b)
}
`},

// ========== ========== ========== ==========
10:{`
包main;入"fmt"
功main(){
  让b:=7
  让func:=8
  fmt.Printf("Hi, nos.%s and %s!\n", b, func) // different comment here
}`,

// ---------- ---------- ---------- ----------
`package main

import _fmt "fmt"

func _main() {
	_b := 7
	_func := 8
	_fmt.Printf("Hi, nos.%s and %s!\n", _b, _func)
}
`},

// ========== ========== ========== ==========
//test prohibited kanji use on LHS
11:{`package main
func main() {
	a:= true
	b:= 真
	nil:= true
	iota:= 真
	//假:= true //parse error
	形Printf("a: %v, b: %v, nil: %v, iota: %v\n", a, b, nil, iota)

	变abc图[串]整;
	abc[串("def")]=789

	var _z= "abc"
	形Println(串(z))
}
`,

// ---------- ---------- ---------- ----------
`package main

import fmt "fmt"

func main() {
	a := true
	b := true
	nil := true
	iota := true

	fmt.Printf("a: %v, b: %v, nil: %v, iota: %v\n", a, b, nil, iota)

	var abc map[string]int
	abc[string("def")] = 789

	var _z = "abc"
	fmt.Println(string(_z))
}
`},

// ========== ========== ========== ==========
12:{`
package main
type A int
func main(){
  _a:=123
  形Println(整64(a))
  var b A
  var c这A
  形Println(b, c)
}
`,

// ---------- ---------- ---------- ----------
`package main

import fmt "fmt"

type A int

func main() {
	_a := 123
	fmt.Println(int64(_a))
	var b A
	var c A
	fmt.Println(b, c)
}
`},

// ========== ========== ========== ==========
13:{`包正;种A整;功正(){a:=123;形Println(整64(a));变b这A;变c这A;形Println(b,c)}`, //example of totally spaceless code

// ---------- ---------- ---------- ----------
`package main

import fmt "fmt"

type A int

func main() { _a := 123; fmt.Println(int64(_a)); var _b A; var _c A; fmt.Println(_b, _c) }
`},

// ========== ========== ========== ==========
999:{`
package main
`,

// ---------- ---------- ---------- ----------
`package main
`},

// ========== ========== ========== ==========
}

//================================================================================================================================================================
func TestKanjiParseError(t *testing.T) {
	for i, tst:= range kanjiParseErrorTests {
		src:= tst.key
		dse:= tst.val
		fset := token.NewFileSet()
		/*defer func(){
			if x:= recover(); fmt.Sprintf("%v", x) != msg {
				tt.Errorf("assert %d failed.\n" +
					"....found recover:%v\n" +
					"...expected panic:%v\n", i, fmt.Sprintf("%v", x), msg)
			}
		}*/

		_, err := parser.ParseFile(fset, "", src, 0)
		if fmt.Sprintf("%v", err) != dse {
			t.Errorf("unexpected parse error in %d: received error: %q; expected error: %q", i, err, dse)
		}
	}
}

var kanjiParseErrorTests = map[int]struct{key string; val string} {

// ========== ========== ========== ==========
1001:{`
package main
func main() {
	//a:= true
	//b:= 真
	//nil:= true
	//iota:= 真
	假:= true //this generates parse error "non-kanji special identifier on left hand side"
}
`,

// ---------- ---------- ---------- ----------
`8:5: non-kanji special identifier 假 on left hand side (and 1 more errors)`},

// ========== ========== ========== ==========
1002:{`
package main
func main() {
	串["def"]=789
}
`,

// ---------- ---------- ---------- ----------
`4:12: non-kanji special identifier 串 on left hand side (and 1 more errors)`},

// ========== ========== ========== ==========
1999:{`
package
`,

// ---------- ---------- ---------- ----------
`2:9: expected 'IDENT', found 'EOF'`},

// ========== ========== ========== ==========
}

//================================================================================================================================================================
/*
restrict kanji in decl/import idents
fix error where pre-existing imports aren't re-added with different name
Why is keyword scoping of struct values working?

test keyword scoping: 为,终,回,破,继,跳,构
test specids: 能,实,虚,造,新,关,加,副,删,丢,抓,写,线,毫,镇,错

more tests:
  让做任英这
  non-kanji on lhs
  keyword scoping: 包
  specid scoping: all except 真,假,空,毫
  default packages: 数,大,网,序
  keywords as labels in kanji-context
  blank: _
  lhs 口 radical kanji on imports
*/

