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
 *     Initial: 2017/11/13        Yang Chenglong
 */

package main

import (
	"fmt"
)

var (
	p  [1005]int
	s  [1005][1005]int
	dp [1005][1005]int
)

func printv(i, j int) {
	if i == j {
		fmt.Printf("A%d", i)
		return
	}

	k := s[i][j]
	if i == k {
		fmt.Printf("A%d", i)
	} else if i < k-1 {
		fmt.Print("(")
		printv(i, k)
		fmt.Print(")")
	} else if i == k-1 {
		fmt.Printf("(A%dA%d)", i, k)
	}

	if k+1 == j {
		fmt.Printf("A%d", k+1)
	} else if k+1 < j-1 {
		fmt.Printf("(")
		printv(k+1, j)
		fmt.Printf(")")
	} else if k+1 == j-1 {
		fmt.Printf("(A%dA%d)", k+1, j)
	}
}

func main() {
	var n int
	fmt.Printf("input matrix number:")
	fmt.Scanf("%d",&n)
	for i := 0; i <= n; i++ {
		fmt.Scanf("%d",&p[i])
	}

	for lens := 1; lens < n; lens++ {
		for i := 1; i+lens <= n; i++ {
			dp[i][i+lens] = dp[i+1][i+lens] + p[i-1]*p[i]*p[i+lens]
			s[i][i+lens] = i
			for k := i + 1; k < i+lens; k++ {
				t := dp[i][k] + dp[k+1][i+lens] + p[i-1]*p[k]*p[i+lens]
				if dp[i][i+lens] > t {
					dp[i][i+lens] = t
					s[i][i+lens] = k
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if j < i {
				fmt.Printf("      ")
			} else {
				fmt.Printf("%6d",dp[i][j])
			}

		}
		fmt.Println()
	}

	fmt.Printf("answer: (")
	printv(1, n)
	fmt.Printf(")")
	fmt.Println()
}
