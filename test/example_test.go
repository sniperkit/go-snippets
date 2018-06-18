package test

import "fmt"

func ExampleReverse() {
	fmt.Println(Reverse("The quick brown 狐 jumped over the lazy 犬"))
	// Output: 犬 yzal eht revo depmuj 狐 nworb kciuq ehT
}

func ExampleReverseWrong() {
	fmt.Println(Reverse("hello"))
	// Output: lleh
}
