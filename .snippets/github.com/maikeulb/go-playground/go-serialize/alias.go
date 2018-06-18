package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type DogAlias Dog

type Dog struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Breed  string    `json:"breed"`
	BornAt time.Time `json:"-"`
}

func (d Dog) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONDog(d))
}

func (d *Dog) UnmarshalJSON(data []byte) error {
	var jd JSONDog
	if err := json.Unmarshal(data, &jd); err != nil {
		return err
	}
	*d = jd.Dog()
	return nil
}

type JSONDog struct {
	DogAlias
	BornAt int64 `json:"born_at"`
}

func (jd JSONDog) Dog() Dog {
	dog := Dog(jd.DogAlias)
	dog.BornAt = time.Unix(jd.BornAt, 0)
	return dog
}

func NewJSONDog(dog Dog) JSONDog {
	return JSONDog{
		DogAlias(dog),
		dog.BornAt.Unix(),
	}
}

func main() {
	dog := Dog{1, "bowser", "husky", time.Now()}
	b, err := json.Marshal(dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b = []byte(`{"id":1,"name":"bowser","breed":"husky","born_at":1480979203}`)
	dog = Dog{}
	json.Unmarshal(b, &dog)
	fmt.Println(dog)
}
