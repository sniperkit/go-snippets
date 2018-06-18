package selection

func Sort(arr []int) []int {
	length := len(arr)
	if length < 1 {
		return arr
	}
	for i := 0; i < length - 1; i++ {
		min := i
		for j := i + 1; j < length; j++ {
			if (arr[j] < arr[min]) {
				min = j
			}
		}
		if min != i {
			arr[i], arr[min] = arr[min], arr[i]
		}
	}

	return arr
}