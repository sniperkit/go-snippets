// globalUtils.go
package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ConvertAsdotToAsplain(asdot string) (int, error) {
	asnum := 0
	nums := strings.Split(asdot, ".")
	if len(nums) != 2 {
		return asnum, errors.New(fmt.Sprintln("Invalid asdot ", asdot))
	}
	n1, _ := strconv.Atoi(nums[0])
	n2, _ := strconv.Atoi(nums[1])
	asnum = (n1 * 65536) + n2
	return asnum, nil
}

func GetAsNum(asdot string) (int, error) {
	asnum := 0
	if asdot == "" {
		asnum = 0
	} else if strings.ContainsAny(asdot, ".") {
		asnum, err := ConvertAsdotToAsplain(asdot)
		return asnum, err
	} else {
		asnum, err := strconv.Atoi(asdot)
		return asnum, err
	}
	return asnum, nil
}

func GetAsDot(asnum int) (string, error) {
	asdot := ""
	if asnum < 65536 {
		asdot = strconv.Itoa(asnum)
	} else {
		n1 := asnum / 65536
		n2 := (asnum - (65536 * n1)) % 65536
		asdot = strconv.Itoa(n1) + "." + strconv.Itoa(n2)
	}
	return asdot, nil
}
