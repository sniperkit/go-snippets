package handling

import (
	"errors"
)

func ErrorExampleMethod(x int) (int, error) {
	//burada bir diğer nokta da if lerde parantez içerisinde checking yapmak zorunlu değildir
	if x == 0 {
		return x, errors.New("bir hata dönelim!!!")
	} else {
		//bir diğer noktada ;(noktalı virgul kullanımı zorunlu değildir)
		x++
		return x, nil
	}
}
