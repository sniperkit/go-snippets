package controllers

import "strconv"

func GenerateKey(x, y int) string {
	return strconv.Itoa(x) + "-" + strconv.Itoa(y)
}
