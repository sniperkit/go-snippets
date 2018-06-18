package main

import (
	"common/functions"
	"errorpkg/handling"
	"expChannel/channel"
	"fmt"
	"mathcommon/math"
	"net/http"
	"www/wwwweb"
)

func main() {
	//Ornek 1
	fmt.Println("Hello World!")

	//Ornek 2
	var sample string
	sample = "Yaşa Mustafa Kemal Paşa, Yaşa!!!"
	fmt.Println(sample)

	//Ornek 3
	//kendi yazdıgımız dısarıdan bir paket ile calısmak
	fmt.Println(functions.GetValueFromAnotherPackageFuntion())

	//Ornek 4
	xs := []float64{1, 2, 3, 4}
	avg := math.Average(xs)
	fmt.Println(avg)

	//Ornek 5 Karakter Karsılıkları
	for i := 0; i < len(sample); i++ {
		fmt.Printf("%x ", sample[i])
		fmt.Print("")
	}
	fmt.Println("")

	//Ornek 6 channel lar
	channel.ExpChannel()

	//Ornek 7 birden fazla value return'u sadece go da bulunur
	fmt.Println(functions.Swap(3, 5))

	//Ornek 8 exception handling
	//golang de try catch ile exception handling yapılmıyor
	//bunun nedeni async calısma operasyonlarını yavaşlatmasından kaynaklı
	//bundan dolayı error mekanizması biraz farklı

	//bunun anlamı error var ise burayı yap değilse devam et
	if x1, e := handling.ErrorExampleMethod(5); e != nil {
		fmt.Println(x1)
	}

	//Ornek 9 Struct yapıları yani class
	f := Foo{a: 555}
	fmt.Printf("Struct valuesunu gösterelim : %d", f)
	fmt.Println("")
	fmt.Println(f.GetName())

	//Ornek 10 Interface Kullanım örneği
	f.A()

	//Ornek 11 simple web server
	//go lang ile web site oluşturma ornekleri
	http.HandleFunc("/", wwwweb.Handler)
	http.ListenAndServe(":1907", nil)
}

type Foo struct {
	a int
}

func (f Foo) GetName() string {
	return "Foo"
}

type x interface {
	A()
	B() int
}

func (f Foo) A() {
	fmt.Printf("Interface methodu implement edip Foo ya dahil ettik.")
}

// type x2 interface {
// 	A()
// 	B() int
// }

// func (f Foo) A() {
// 	fmt.Printf("Interface methodu implement edip Foo ya dahil ettik.")
// }
