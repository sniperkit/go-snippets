package utils

import (
	"log"
	"math"
)

//CheckErrorMsg checks error value and logs a custom message with stact trace and exits program
func CheckErrorMsg(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

//CheckError checks error value and logs a the stack trace and exits program
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Round rounds given number
func Round(val float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(0.5, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
