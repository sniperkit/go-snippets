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
 *     Initial: 2017/11/08        Yang Chenglong
 */

package main

import "fmt"

type papapa []int

func main() {

	var head papapa
	head = append(head,0)

	for i:=0;i<20;i++{
		head.insert(i)
	}
	for i:=0;i<20;i++ {
		fmt.Println(head.deleteMax())
		fmt.Println(head)
	}
	a := []int{1,2,3,4}
	fmt.Println(a[1:3])

}

func(this *papapa)insert(key int){
	*this = append(*this,key)
	this.swim(len(*this)-1)
}

func(this *papapa)deleteMax()int{
	(*this)[1],(*this)[len(*this)-1] = (*this)[len(*this)-1],(*this)[1]
	res := (*this)[len(*this)-1]
	(*this) = (*this)[:len(*this)-1]
	this.sink(1)
	return res
}

func(this *papapa)swim(index int){
	if index/2 == 0{
		return
	}
	if (*this)[index] > (*this)[index/2]{
		(*this)[index],(*this)[index/2] = (*this)[index/2],(*this)[index]
		this.swim(index/2)
	}
	return
}

func(this *papapa)sink(index int){

	maxIndex := len(*this)-1

	if maxIndex>=2*index+1{

		//找出子节点最大的
		max := 2*index
		if (*this)[2*index+1]>(*this)[2*index]{
			max = 2*index+1
		}

		if (*this)[index] < (*this)[max]{
			(*this)[index],(*this)[max] = (*this)[max],(*this)[index]
			this.sink(max)
		}

	}else if maxIndex>=2*index{
		if (*this)[index] < (*this)[2*index]{
			(*this)[index],(*this)[2*index] = (*this)[2*index],(*this)[index]
			this.sink(2*index)
		}
	}

	return
}

