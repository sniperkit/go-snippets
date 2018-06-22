package main
//
// import (
//   "github.com/eaigner/shield"
// )
//
// type classify struct {
//   sh *shield.Shield
// }
//
// func newClassify() *classify {
//   var cf classify
//   cf.sh := &shield.New(
//     shield.NewEnglishTokenizer(),
//     shield.NewRedisStore("127.0.0.1:6379", "", nil, ""),
//   )
//   return &classify
// }
//
// func classify() {
//
//   sh.Learn("good", "sunshine drugs love sex lobster sloth")
//   sh.Learn("bad", "fear death horror government zombie god")
//
//   c, _ := sh.Classify("sloths are so cute i love them")
//   if c != "good" {
//     panic(c)
//   }
//
//   c, _ = sh.Classify("i fear god and love the government")
//   if c != "bad" {
//     panic(c)
//   }
// }
