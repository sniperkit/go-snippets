package lexus

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

	if total > 1000 {
		c.sum += 200
	}

	c.sum -= float32(total) * 0.7
}

func (c *cart) Total() float32 {
	return c.sum
}

// NewCart ...
func NewCart(sum float32) *cart {
	return &cart{
		sum: sum,
	}
}
