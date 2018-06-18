package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Dog struct {
	ID     int
	Name   string
	Breed  string
	BornAt time.Time
}

type JSONDog struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	BornAt int64  `json:"born_at"`
}

func (jd JSONDog) Dog() Dog {
	return Dog{
		jd.ID,
		jd.Name,
		jd.Breed,
		time.Unix(jd.BornAt, 0),
	}
}

func NewJSONDog(dog Dog) JSONDog {
	return JSONDog{
		dog.ID,
		dog.Name,
		dog.Breed,
		dog.BornAt.Unix(),
	}
}

func main() {
	dog := Dog{1, "bowser", "husky", time.Now()}
	b, err := json.Marshal(NewJSONDog(dog))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b = []byte(`{"id":1,"name":"bowser","breed":"husky","born_at":1480979203}`)
	var jsonDog JSONDog
	json.Unmarshal(b, &jsonDog)
	fmt.Println(jsonDog.Dog())
}
