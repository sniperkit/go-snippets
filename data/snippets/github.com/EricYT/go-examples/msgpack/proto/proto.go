package proto

//go:generine msgp

// test map type
type Foo struct {
	//	Actions map[string]Runer `msg:"action"`
}

type Oper1 struct {
	Key   string `msg:"key"`
	Value string `msg:"value"`
}

type Oper2 struct {
	Key   string `msg:"key"`
	Value string `msg:"value"`
}

// test struct pointer
type FooWithPointer struct {
	Key   string `msg:"key"`
	*Bar1 `msg:"bar1"`
	*Bar2 `msg:"bar2"`
}

type Bar1 struct {
	Value string `msg:"value"`
}

type Bar2 struct {
	Value string `msg:"value"`
}
