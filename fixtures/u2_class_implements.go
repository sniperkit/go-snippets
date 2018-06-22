package fixtures

type Testiface1 interface{
	Testfnc1()
}
type Testcls1 struct{}
func (t Testcls1) Testfnc1() {}
var  _ Testiface1 = Testcls1{}
