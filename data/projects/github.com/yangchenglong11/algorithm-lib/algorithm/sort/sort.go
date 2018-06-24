/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/11/02        Yang Chenglong
 */

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	var (
		num   int
		arr01 []int
		sys   []int
	)
	fmt.Println("Please input the sum of the number:")
	fmt.Scanf("%d", &num)

	arr01 = make([]int, num)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		arr01[i] = r.Intn(num)
	}
	sys = arr01

	now1 := time.Now()
	mergeSort1(arr01, 0, len(arr01)-1)
	finish := time.Since(now1)
	fmt.Print("mergeSort1:")
	fmt.Println(finish)

	now2 := time.Now()
	mergeSort2(arr01, 0, len(arr01)-1)
	finish2 := time.Since(now2)
	fmt.Print("mergeSort2:")
	fmt.Println(finish2)

	s_now := time.Now()
	sort.Ints(sys)
	s_finish := time.Since(s_now)
	fmt.Print("System sort:")
	fmt.Println(s_finish)
}

func mergeSort1(arr []int, low, high int) {
	if low < high {
		mid := (low + high) / 2
		mergeSort1(arr, low, mid)
		mergeSort1(arr, mid+1, high)
		merge1(arr, low, mid, high)
	}
}

func mergeSort2(arr []int, low, high int) {
	if low < high {
		mid := (low + high) / 2
		mergeSort2(arr, low, mid)
		mergeSort2(arr, mid+1, high)
		merge2(arr, low, mid, high)
	}
}

func merge1(arr []int, low, mid, high int) {
	leftLen := mid - low + 1
	rightLen := high - mid

	arrLeft := make([]int, leftLen+1)
	for i := 0; i < leftLen; i++ {
		arrLeft[i] = arr[low+i]
	}
	arrLeft[leftLen] = 999999999 //哨兵牌

	arrRight := make([]int, rightLen+1)
	for j := 0; j < rightLen; j++ {
		arrRight[j] = arr[mid+j+1]
	}
	arrRight[rightLen] = 999999999 //哨兵牌

	i, j := 0, 0
	for k := low; k <= high; k++ {
		if arrLeft[i] <= arrRight[j] {
			arr[k] = arrLeft[i]
			i++
		} else {
			arr[k] = arrRight[j]
			j++
		}
	}
}

func merge2(arr []int, low, mid, high int) {
	leftLen := mid - low + 1
	rightLen := high - mid

	arrLeft := make([]int, leftLen)
	for i := 0; i < leftLen; i++ {
		arrLeft[i] = arr[low+i]
	}

	arrRight := make([]int, rightLen)
	for j := 0; j < rightLen; j++ {
		arrRight[j] = arr[mid+j+1]
	}

	i, j, k := 0, 0, low
	for ; k <= high && i < leftLen && j < rightLen; k++ {
		if arrLeft[i] <= arrRight[j] {
			arr[k] = arrLeft[i]
			i++
		} else {
			arr[k] = arrRight[j]
			j++
		}
	}

	for ; i < leftLen && k <= high; k++ {
		arr[k] = arrLeft[i]
		i++
	}

	for ; j < rightLen && k <= high; k++ {
		arr[k] = arrRight[j]
		j++
	}
}
