package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Class struct {
	ID   int
	Name string
}

func (c *Class) Depth() int {
	id, depth := c.ID, 2
	for depth > 0 && id%10 == 0 {
		depth--
		id /= 10
	}

	return depth
}

type DDC struct {
	Classes map[int]*Class
}

func LoadDDC(path string) (*DDC, error) {
	fin, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	ddc := &DDC{
		Classes: make(map[int]*Class),
	}

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("ddc: cannot parse class ID '%s'", parts[0])
		}
		if id < 0 || 1000 < id {
			return nil, fmt.Errorf("ddc: invalid ID '%d' not in [0; 1000[", id)
		}
		if _, ok := ddc.Classes[id]; ok {
			return nil, fmt.Errorf("ddc: duplicate ID '%d'", id)
		}

		ddc.Classes[id] = &Class{
			ID:   id,
			Name: parts[1],
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ddc, nil
}
