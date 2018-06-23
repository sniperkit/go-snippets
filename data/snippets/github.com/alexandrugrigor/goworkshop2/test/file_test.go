package test

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"fmt"
)

type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type persons []person

func (p person) String() string {
	return fmt.Sprintf("person{firstName='%s', lastName='%s'}", p.FirstName, p.LastName)
}

func TestReadFile(t *testing.T) {
	fileContent, err := ioutil.ReadFile("test_data.json")
	assert.NoError(t, err)

	var persons persons
	assert.NoError(t, json.Unmarshal(fileContent, &persons))
	fmt.Println(persons)
	assert.Len(t, persons, 2, )

	assert.Equal(t, "John", persons[0].FirstName)
	assert.Equal(t, "Doe", persons[0].LastName)

	assert.Equal(t, "Bob", persons[1].FirstName)
	assert.Equal(t, "Smith", persons[1].LastName)
}
