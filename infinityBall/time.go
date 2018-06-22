package main

type Time struct {
	Time float64
	prev float64
	Delta float64
}

func (time *Time) Set(secs float64) {
	time.prev = time.Time
	time.Time = secs
	time.Delta = time.Time - time.prev
}

