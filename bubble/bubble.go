package bubble


func Sort(arr []int) []int {
	length := len(arr)
	if length < 1 {
		return arr
	}
	var changed bool
	for i := 0; i < length-1; i++ {
		changed = false
		for j := 0; j < length-1; j++ {
			if arr[j] > arr[j+1] {
				changed = true
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		if changed == false {
			break
		}
	}
	return arr
}
