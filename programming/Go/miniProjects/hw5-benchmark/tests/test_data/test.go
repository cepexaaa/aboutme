package testdata

import "time"

type ProProHuman struct {
	ProName string
}

type ProHuman struct {
	ProProHuman
	SuperName string
}

type Human struct {
	ProHuman
	Species string
	// Person
}

type Person struct {
	Human
	Name    string    `json:"name" validate:"required"`
	Age     int       `json:"age"`
	Email   string    `json:"email,omitempty"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

type MegaPerson struct {
	Person
	Id int
}

func (p Person) GetName() string {
	return p.Name
}

func (p *Person) SetName(name string) {
	p.Name = name
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}
