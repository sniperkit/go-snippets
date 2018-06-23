package car

type Car struct {
	Name  string
	Price int
	Color string
}

func (c *Car) GetName() string {
	return c.Name
}

func (c *Car) UpdateColor(color string) {
	c.Color = color
}

func (c *Car) GetColor() string {
	return c.Color
}

func New(name string, color string) *Car {

	return &Car{
		Name:  name,
		Color: color,
	}
}
