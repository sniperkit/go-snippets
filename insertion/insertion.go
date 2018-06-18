package insertion


func Sort(arr []int) []int {
	length := len(arr)
	if length < 1 {
		return arr
	}
	for i := 1; i < length; i++ {
		temp := arr[i]
		j := i - 1 //index
		for ; j >= 0 && arr[j] > temp; {
			arr[j + 1] = arr[j]
			j--
		}
		arr[j + 1] = temp
	}

	return arr
}