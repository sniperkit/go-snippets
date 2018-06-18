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
	BornAt Time `json:"born_at"`
}

func (jd JSONDog) Dog() Dog {
	dog := Dog(jd.DogAlias)
	dog.BornAt = jd.BornAt.Time
	return dog
}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

func NewJSONDog(dog Dog) JSONDog {
	return JSONDog{
		DogAlias(dog),
		Time{dog.BornAt},
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
	err = json.Unmarshal(b, &dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(dog)
}
