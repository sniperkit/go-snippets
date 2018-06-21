//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package alphaNumSort

import (
	"errors"
	"math"
	"sort"
	"strings"
)

const (
	MAX_CHAR_ASCII_VAL = 123 //ASCII value of 'z' + 1
)

/* Simple sort routine that sorts alpha numeric strings
   containing the following runes (0-9, a-z, A-Z, _).
*/
func Sort(strList []string) []string {
	if (strList == nil) || (len(strList) == 1) {
		return strList
	}

	var outList []string = make([]string, len(strList))
	var strMap map[float64]string = make(map[float64]string)
	for _, str := range strList {
		wt := computeWeight(str)
		strMap[wt] = str
	}
	keySlice := make([]float64, len(strMap))
	idx := 0
	for key, _ := range strMap {
		keySlice[idx] = key
		idx++
	}
	sort.Float64s(keySlice)
	for idx, key := range keySlice {
		outList[idx] = strMap[key]
	}
	return outList
}

// IsLess reports whether s is less than t.
func IsLess(s1, s2 string) bool {
	return lessRunes([]rune(s1), []rune(s2))
}

// lessRunes reports whether s1 is less than s2.
func lessRunes(s1, s2 []rune) bool {
	nprefix := commonPrefix(s1, s2)
	if len(s1) == nprefix && len(s2) == nprefix {
		// equal
		return false
	}
	s1End := leadDigits(s1[nprefix:]) + nprefix
	s2End := leadDigits(s2[nprefix:]) + nprefix
	if s1End > nprefix || s2End > nprefix {
		start := trailDigits(s1[:nprefix])
		if s1End-start > 0 && s2End-start > 0 {
			s1n := atoi(s1[start:s1End])
			s2n := atoi(s2[start:s2End])
			if s1n != s2n {
				return s1n < s2n
			}
		}
	}
	switch {
	case len(s1) == nprefix:
		return true
	case len(s2) == nprefix:
		return false
	default:
		return s1[nprefix] < s2[nprefix]
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func atoi(r []rune) uint64 {
	if len(r) < 1 {
		panic(errors.New("atoi got an empty slice"))
	}
	const cutoff = uint64((1<<64-1)/10 + 1)
	const maxVal = 1<<64 - 1

	var n uint64
	for _, d := range r {
		v := uint64(d - '0')
		if n >= cutoff {
			return 1<<64 - 1
		}
		n *= 10
		n1 := n + v
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return 1<<64 - 1
		}
		n = n1
	}
	return n
}

func commonPrefix(s, t []rune) int {
	for i := range s {
		if i >= len(t) {
			return len(t)
		}
		if s[i] != t[i] {
			return i
		}
	}
	return len(s)
}

func trailDigits(r []rune) int {
	for i := len(r) - 1; i >= 0; i-- {
		if !isDigit(r[i]) {
			return i + 1
		}
	}
	return 0
}

func leadDigits(r []rune) int {
	for i := range r {
		if !isDigit(r[i]) {
			return i
		}
	}
	return len(r)
}

/* Remove special chars in the string.
 * "#", ":", "."
 * Add other special chars as they become part of key string in object.
 */
func Normalize(str string) string {
	strs := strings.Split(strings.TrimSpace(str), "#")
	pass1 := ""
	for i := range strs {
		if i > 0 {
			pass1 = pass1 + strs[i]
		}
	}
	strs = strings.Split(pass1, ":")
	pass2 := ""
	for i := range strs {
		pass2 = pass2 + strs[i]
	
	}
	strs = strings.Split(pass1, ".")
	pass3 := ""
	for i := range strs {
		pass3 = pass3 + strs[i]
	
	}
	return pass3
}

/* Computes weight of give string. Max char ascii val ('z') = 122 */
func computeWeight(str string) float64 {
	var wt float64
	l := len(str)
	for idx, val := range str {
		wt += math.Pow(float64(MAX_CHAR_ASCII_VAL), float64(l - idx)) * float64(val)
	}
	return wt
}
