package main

import (
	"fmt"
	"git.suretly.io/backend/backend/common/currency"
)

var (
	grades = []float64{6,8.87,12.21,14.88,17.12,20.2,25.65}
)

func getProfit(sum float64, returnPercent float64) {

	profitPercent := 2.0
	profit := (((100 + profitPercent)/returnPercent) - 1)*100
	fmt.Println("profit:", currency.NewFromFloat(profit).RoundCash(5).Float64(), "%")
	fmt.Println("total:", returnPercent, "-", currency.NewFromFloat((sum/100.0)*profit).RoundCash(5).Float64())
}

func main() {

	k := &currency.Currency{}

	for i:=0;i<1000;i++ {
		k = k.Add(currency.NewFromFloat(0.1))
	}

	fmt.Println(k)

	fmt.Println("-------")

	sum := 500.0
	for _, grade := range grades {
		getProfit(sum, 100.0 - grade)
	}

	//mSum := currency.NewFromFloat(sum)
	//mProfitPercent := currency.NewFromFloat(profitPercent)
	//mReturnPercent := currency.NewFromFloat(returnPercent)
	//mIOO := currency.NewFromFloat(100)
	//
	//mProfit := &currency.Currency{}
	//mProfit.Add(mIOO).Add(mProfitPercent).Div(mReturnPercent).Sub(currency.NewFromFloat(1)).Mul(mIOO)
	//
	//fmt.Println(mSum.Div(mIOO).Mul(mProfit))
}
