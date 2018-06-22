package basic

import (
	"fmt"
	"testing"
)

type Data struct {
	value int
}

func (d Data) add(delta int) {
	d.value += delta
}

func (d *Data) minus(delta int) {
	d.value -= delta
}

func TestStructPointer(t *testing.T) {
	d := Data{10}
	dd := &Data{10}
	d.add(10)
	dd.add(10)
	d.minus(10)
	dd.minus(5)
	fmt.Println(d)
	fmt.Println(*dd)

}
