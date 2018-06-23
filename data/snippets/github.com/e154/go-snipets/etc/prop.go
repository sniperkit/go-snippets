package main

import (
	"git.suretly.io/backend/backend/common/currency"
	"fmt"
)

func main() {
	var trs = []float64{100,200,500,300,200,400,900}
	var sumToReturn float64 = 2000
	totalBlock := &currency.Currency{}
	totalPerc := &currency.Currency{}
	sumToReturnCheck := &currency.Currency{}

	for _, tr := range trs {
		totalBlock.Add(currency.NewFromFloat(tr))
	}

	fmt.Println("totalBlock",totalBlock)
	onePerc := totalBlock.Div(currency.NewFromFloat(100))
	fmt.Println("onePerc", onePerc)

	for _, tr := range trs {
		sumOnePerc := currency.NewFromFloat(sumToReturn).Div(currency.NewFromFloat(100))
		qwe := currency.NewFromFloat(tr).Div(onePerc)
		strc := sumOnePerc.Mul(qwe).RoundCash(5)
		fmt.Println(qwe, "-", strc)
		totalPerc.Add(qwe)
		sumToReturnCheck.Add(strc)
	}

	fmt.Println(totalPerc.RoundCash(5), "%")
	fmt.Println(sumToReturnCheck.RoundCash(50))

	var promo int64 = 700
	var sum int64 = 500
	var i, ii int64

	if promo - sum > 0 {
		i = sum
	} else {
		i = promo
		ii = sum - promo
	}
	fmt.Println("s", i, ii)


}

