package capsule

import "github.com/cheekybits/genny/generic"

type Item generic.Type

type ItemCapsule struct {
	s []Item
}

func NewItemCapsule() *ItemCapsule {
	return &ItemCapsule{s: []Item{}}
}

func (c *ItemCapsule) Put(val Item) {
	c.s = append(c.s, val)
}

func (c *ItemCapsule) Get() Item {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}
