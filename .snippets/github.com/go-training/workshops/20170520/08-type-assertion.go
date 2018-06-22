package main

import (
	"errors"
	"fmt"
)

func printF(val interface{}) error {
	if _, ok := val.(string); ok {
		fmt.Println("this is string")
	} else {
		return errors.New("Error: this is not string")
	}

	switch val.(type) {
	case string:
		fmt.Println("this is string")
	case int:
		fmt.Println("this is int")
	}

	return nil
}

func main() {
	// printF("hi golang")
	if err := printF(100); err != nil {
		fmt.Println(err)
	}
}
