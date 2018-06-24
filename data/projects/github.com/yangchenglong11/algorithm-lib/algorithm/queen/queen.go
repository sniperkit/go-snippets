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
 *     Initial: 2017/11/15        Yang Chenglong
 */

package main

import (
	"fmt"
	"math"
)

const QUEEN = 5
const INITIAL = -10000

var a [QUEEN]int //一维数组表示棋盘

//对棋盘进行初始化
func initq() {
	for i := 0; i < QUEEN; i++ {
		a[i] = INITIAL
	}
}

//判断第row行第col列是否可以放置皇后
func valid(row, col int) bool {
	var i int
	//对棋盘进行扫描
	for i = 0; i < QUEEN; i++ {
		//判断列冲突与斜线上的冲突
		if a[i] == col || math.Abs(float64(i-row)) == math.Abs(float64(a[i]-col)) {
			return false
		}
	}

	return true
}

//打印输出N皇后的一组解
func printq() {
	var i, j int
	for i = 0; i < QUEEN; i++ {
		for j = 0; j < QUEEN; j++ {
			//a[i]为初始值
			if a[i] == j {
				fmt.Print(j)
			}
		}
	}
	fmt.Println()
}

func queen() {
	n, i, j := 0, 0, 0
	for i < QUEEN {
		//对i行的每一列进行探测，看是否可以放置皇后
		for j < QUEEN {
			//该位置可以放置皇后
			if valid(i, j) {
				//第i行放置皇后
				a[i] = j
				//第i行放置皇后以后，需要继续探测下一行的皇后位置，所以此处将j清零，从下一行的第0列开始逐列探测
				j = 0
				break
			} else {
				//继续探测下一列
				j++
			}
		}

		//第i行没有找到可以放置皇后的位置
		if a[i] == INITIAL {
			//回溯到第一行，仍然无法找到可以放置皇后的位置，则说明已经找到所有的解，程序终止
			if i == 0 {
				break
			} else {
				i--
				//把上一行皇后的位置往后移一列
				j = a[i] + 1
				//把上一行皇后的位置清除，重新探测
				a[i] = INITIAL
				continue
			}
		}
		//最后一行找到了一个皇后位置，说明找到一个结果，打印出来
		if i == QUEEN-1 {
			fmt.Printf("answer%d:", n+1)
			printq()
			//不能在此处结束程序，因为我们要找的是N皇后问题的所有解，此时应该清除该行的皇后，从当前放置皇后列数的下一列继续探测。
			j = a[i] + 1   //从最后一行放置皇后列数的下一列继续探测
			a[i] = INITIAL //清除最后一行的皇后位置
			n = n + 1
			continue
		}
		i++ //继续探测下一行的皇后位置
	}
}

func main() {
	initq()
	queen()
}
