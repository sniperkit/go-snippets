package main

import "fmt"

func main() {
	a := []string{"wangrui","weed","wangxinyi","haha","reg"}
	stringReverse(a)
	fmt.Println(a)
	b := sliceContains(a,"weed")
	fmt.Println(b)
}

func stringReverse(src []string) {
	if src == nil {
		panic(fmt.Errorf("the src can not be empty"))
	}
	count := len(src)
	mid := count/2
	for i := 0;i < mid; i ++ {
		tmp := src[i]
		src[i] = src[count - 1]
		src[count - 1] = tmp
		count --
	}
}
//判断是否包含
func sliceContains(src []string,value string)bool{
	isContain := false
	for _,srcValue := range src  {
		if(srcValue == value){
			isContain = true
			break
		}
	}
	return isContain
}

func mapContains(src map[string]interface{} ,key string) bool{
	if _, ok := src[key]; ok {
		return true
	}
	return false
}
