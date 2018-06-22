package main


func main() {
}
func repeat_generator(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select{
			case <- done:
				return
				case takeStream <- <- valueStream:
			}
		}
	}()
	return takeStream
}


func int_generator(done <- chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers{
			select {
			case <-done:
				return
				case intStream <- i:
			}
		}
	}()
	return intStream
}
