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
 *     Initial: 2017/11/14        Yang Chenglong
 */

package main

import "fmt"

func LCSLength(lengthS1 int, lengthS2 int, s1 string, s2 string, result [][]int, trace [][]int) int {
	trace[0][0] = 0 //trace[i][0]与trace[0][j]不用，可以用作标志位

	for i := 1; i <= lengthS1; i++ {
		result[i][0] = 0
	}
	for j := 1; j <= lengthS2; j++ {
		result[0][j] = 0
	}

	for i := 1; i <= lengthS1; i++ {
		for j := 1; j <= lengthS2; j++ {
			//两个序列首部匹配
			if s1[i] == s2[j] {
				result[i][j] = result[i-1][j-1] + 1
				//1表示继续向左上[i-1][j-1]追踪
				trace[i][j] = 1
			} else if result[i-1][j] >= result[i][j-1] {
				result[i][j] = result[i-1][j]
				//2表示继续向上[i-1][j]追踪
				trace[i][j] = 2
			} else {
				result[i][j] = result[i][j-1]
				//3表示继续向左[i][j-1]追踪
				trace[i][j] = 3
			}
		}
	}

	return result[lengthS1][lengthS2]
}

func LCSTraceback(tailS1 int, tailS2 int, s1 string, trace [][]int) {
	if tailS1 == 0 || tailS2 == 0 {
		return
	}
	if trace[tailS1][tailS2] == 1 { //向左上追踪[i-1][j-1]
		LCSTraceback(tailS1-1, tailS2-1, s1, trace)
		fmt.Print(string(s1[tailS1]))
	} else if trace[tailS1][tailS2] == 2 { //继续向上追踪[i-1][j]
		LCSTraceback(tailS1-1, tailS2, s1, trace)
	} else {
		LCSTraceback(tailS1, tailS2-1, s1, trace)
	}
}

func main() {
	var s1, s2 string

	fmt.Println("请输入两个要比较的字符串：")
	fmt.Scanf("%s%s", &s1, &s2)
	sen1 := " " + s1
	sen2 := " " + s2
	fmt.Printf("s1: %v\ns2: %v\n", s1, s2)

	var resultMatrix [][]int
	var traceMatrix [][]int
	for i := 0; i < len(s1)+1; i++ {
		tmpArr := make([]int, len(s2)+1)
		tmp := make([]int, len(s2)+1)
		resultMatrix = append(resultMatrix, tmpArr)
		traceMatrix = append(traceMatrix, tmp)
	}

	fmt.Printf("The length of the longest common sequence is: %d \n", LCSLength(len(s1), len(s2), sen1, sen2, resultMatrix, traceMatrix))

	fmt.Println("The longest common sequece is: ")
	LCSTraceback(len(s1), len(s2), sen1, traceMatrix)
	fmt.Println()
}
