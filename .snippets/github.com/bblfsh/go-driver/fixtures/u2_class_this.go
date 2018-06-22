package fixtures

type Testcls1 struct {
	a int
}
func (this Testcls1) Testfnc1() {
	this.a = 1
}

func (otherthis Testcls1) Testfnc2() {
	otherthis.a = 2
}
