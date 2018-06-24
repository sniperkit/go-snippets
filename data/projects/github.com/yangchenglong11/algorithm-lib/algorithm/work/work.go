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
 *     Initial: 2017/11/22        Yang Chenglong
 */

package main

import (
	"fmt"
	"sort"
)

var (
	s     []int //每台机器当前已分配的作业总耗时
	time  []int
)

func main() {
	var (
		NWork int //作业数
		NMac  int
		maxtime = 0
	)

	fmt.Println("请输入机器总数,作业总数")
	fmt.Scanf("%d%d", &NMac, &NWork)
	s = make([]int, NMac)
	time = make([]int,NWork)


	fmt.Println("请输入各作业用时")
	for i := 0; i < NWork; i++ {
		fmt.Scanf("%d", &time[i])
	}

	sort.Ints(time)

	if NMac >= NWork {
		maxtime = setwork1(time, NWork, NWork)
	} else {
		maxtime = setwork2(time, NWork, NMac)
	}

	fmt.Printf("最多耗费时间%d", maxtime)
	fmt.Println()
}

//机器数大于待分配作业数
func setwork1(t []int, n int, nwork int) int {
	var i int
	for i = 0; i < n; i++ {
		s[i] = t[i]
	}
	ma := max(s, nwork)
	return ma
}

// 机器数小于待分配作业数
func setwork2(t []int, n int, nmac int) int {
	var i int
	mi := 0
	for i = 0; i < n; i++ {
		mi = min(nmac)
		fmt.Printf("%d,时间和最小的机器号为%d.时间和为%d：\n", i+1, mi, s[mi])
		s[mi] = s[mi] + t[i]
	}
	ma := max(s, nmac)
	return ma
}

//求出目前处理作业的时间和最小的机器号
func min(m int) int {
	min := 0
	var i int
	for i = 1; i < m; i++ {
		if s[min] > s[i] {
			min = i
		}
	}

	return min
}

//求最终结果（最长处理时间）
func max(s []int, num int) int {
	max := s[0]
	var i int
	for i = 1; i < num; i++ {
		if max < s[i] {
			max = s[i]
		}
	}

	return max
}
