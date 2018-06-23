package main

//sector int
type sint int8
type sv2 [2]sint

type Overworld struct {
	all     map[Entity]*Colider
	sectors map[sv2]map[Entity]*Colider
}

func NewOverworld() *Overworld {
	var o Overworld
	o.all = make(map[Entity]*Colider)
	o.sectors = make(map[sv2]map[Entity]*Colider)
	return &o
}

func (o *Overworld) set(e Entity, x, y, r float32) {
	c, ok := o.all[e]
	if ok {
		for _, oldsector := range c.sectors() {
			delete(o.sectors[oldsector], e)
		}
	} else {
		c = new(Colider)
		o.all[e] = c
	}
	c.x, c.y, c.r = x, y, r
	for _, newSector := range c.sectors() {
		s, ok := o.sectors[newSector]
		if !ok {
			s = make(map[Entity]*Colider)
			o.sectors[newSector] = s
		}
		s[e] = c
	}
}

func (o *Overworld) remove(e Entity) {
	c, ok := o.all[e]
	if !ok {
		return
	}
	delete(o.all, e)
	for _, sector := range c.sectors() {
		delete(o.sectors[sector], e)
	}
}

func (o *Overworld) query(e Entity, x, y, r float32) []Entity {
	var entities []Entity

	c := Colider{x, y, r}

	for _, sector := range c.sectors() {
		for entity, col := range o.sectors[sector] {
			if entity == e {
				continue
			}
			dx := col.x - x
			dy := col.y - y
			totalR := col.r + r
			if dx*dx+dy*dy < totalR*totalR {
				alreadyPresent := false
				for _, other := range entities {
					if other == entity {
						alreadyPresent = true
					}
				}
				if !alreadyPresent {
					entities = append(entities, entity)
				}
			}
		}
	}
	return entities
}

type Colider struct {
	x float32
	y float32
	r float32
}

func (c *Colider) sectors() []sv2 {
	const sectorSize = 500
	sxmin := sint((c.x - c.r) / sectorSize)
	sxmax := sint((c.x+c.r)/sectorSize) + 1
	symin := sint((c.y - c.r) / sectorSize)
	symax := sint((c.y+c.r)/sectorSize) + 1

	result := make([]sv2, int(sxmax-sxmin)*int(symax-symin))
	i := 0
	for j := sxmin; j < sxmax; j++ {
		for k := symin; k < symax; k++ {
			result[i] = sv2{j, k}
			i++
		}
	}
	return result
}
