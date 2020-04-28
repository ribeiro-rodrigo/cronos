package cronos

import (
	"crypto/rand"
	"fmt"
	"log"
	"testing"
)

func uuidGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

type observer struct {
}

type observable interface {
	observe() string
}

func (o observer) observe() string {
	return "observer"
}

type visualizer struct{}

func (v visualizer) observe() string {
	return "visualizer"
}

func TestInject(t *testing.T) {

	type person struct {
		name string
	}

	type company struct {
		name   string
		person person
	}

	newPerson := func() person {
		return person{name: "Bob"}
	}

	newCompany := func(person person) company {
		return company{person: person}
	}

	container := New()

	container.Register(newCompany)
	container.Register(newPerson)

	container.Init(func(company company) {
		if company.person.name != "Bob" {
			t.Error()
		}
	})

}
