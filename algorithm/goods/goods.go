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
	"sort"
)

type Goods struct {
	weight    int     // the weight of goods重量
	value     int     // the value of goods价值
	ValPerWei float64 // value per weight权重
	load      float64 //物品装入背包的部分的系数（例如上图中物品全部装入则load为1，装入10/20则load为0.5）
}

type GoodSlice struct {
	G []Goods
}

func (g *GoodSlice)Less(i, j int) bool {
	if g.G[i].ValPerWei < g.G[j].ValPerWei {
		return false
	} else {
		return true
	}
}

func (g *GoodSlice) Len() int {
	return len(g.G)
}

func (g *GoodSlice) Swap(i, j int) {
	g.G[i], g.G[j] = g.G[j], g.G[i]
}

//贪心算法
func Greedy(g []Goods, good_num int, content int) {
	for i := 0; i < good_num; i++ {
		//如果背包足够装下整个物品
		if content > g[i].weight {
			content -= g[i].weight
			g[i].load = 1
		} else if content > 0 {                        // 如果背包不足以装下整个物品
			g[i].load = float64(content / g[i].weight) // 计算物品装入背包的部分
			content = 0                                // 背包容量置0
		}
	}
}

func main() {
	var goods_num, bagvol, total_value, total_weight int

	fmt.Println("Please input the volum of bag:")

	fmt.Scanf("%d", &bagvol)
	fmt.Println("Please input the number of goods:")
	fmt.Scanf("%d", &goods_num)

	Go := &GoodSlice{
		G:make([]Goods, goods_num+1),
	}

	fmt.Println("Please inuput the weight and value of goods:")

	for i := 0; i < goods_num; i++ {
		fmt.Scanf("%d%d", &Go.G[i].weight, &Go.G[i].value)
		Go.G[i].ValPerWei = float64(Go.G[i].value) / float64(Go.G[i].weight) //计算权重
		Go.G[i].load = 0                                        //load置0
	}

	sort.Sort(Go)

	Greedy(Go.G, goods_num, bagvol)
	fmt.Println("Info of loaded goods:")

	//output the info of goods
	for i := 0; i < goods_num; i++ {
		//如果检测到物品未被装入背包，则推出循环
		if Go.G[i].load == 0.0 {
			break
		}

		total_value += (Go.G[i].value * int(Go.G[i].load))   //装入背包的物品总价值
		total_weight += (Go.G[i].weight * int(Go.G[i].load)) //装入背包的物品总重量
		fmt.Println("weight: ", Go.G[i].weight, "  ", "value: ", Go.G[i].value, "  ",
			"\nthe value per weight of good: ", Go.G[i].ValPerWei, "  the part of goods: ", Go.G[i].load)
	}

	fmt.Println("the volum of bag is: ", bagvol)
	fmt.Println("the total weight is: ", total_weight)
	fmt.Println("the total value is: ", total_value)
}



