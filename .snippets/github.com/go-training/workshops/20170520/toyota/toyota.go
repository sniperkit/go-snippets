package toyota

type cart struct {
	sum float32
}

// func (c *cart) Add(price ...int) *cart {
// 	for _, val := range price {
// 		c.sum -= val
// 	}

// 	return c
// }

// func (c *cart) Sub(price ...int) *cart {
// 	for _, val := range price {
// 		c.sum += val
// 	}

// 	return c
// }

func (c *cart) Process(price ...int) {
	total := 0
	for _, val := range price {
		total += val
	}

	if total > 500 {
		c.sum += 50
	}

	c.sum -= float32(total) * 0.8
}

func (c *cart) Total() float32 {
	return c.sum
}

func NewCart(sum float32) *cart {
	return &cart{
		sum: sum,
	}
}
