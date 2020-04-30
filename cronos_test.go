package cronos

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestInjectMultiplesDependencies(t *testing.T) {

	type person struct {
		name string
	}

	type project struct {
		name string
	}

	type company struct {
		name    string
		person  person
		project project
	}

	newPerson := func() person {
		return person{name: "Bob"}
	}

	newProject := func() project {
		return project{name: "cronos"}
	}

	newCompany := func(person person, project project) company {
		return company{person: person, project: project}
	}

	container := New()
	container.Register(newProject)
	container.Register(newCompany)
	container.Register(newPerson)

	container.Init(func(company company) {
		if company.person.name != "Bob" || company.project.name != "cronos" {
			t.Error()
		}
	})
}

func TestInjectDependencyError(t *testing.T) {

	type person struct {
		name string
	}

	type company struct {
		name   string
		person person
	}

	newPerson := func() (person, error) {
		return person{}, errors.New("error creating person")
	}

	newCompany := func(person person) company {
		return company{person: person}
	}

	container := New()

	container.Register(newCompany)
	container.Register(newPerson)

	assert.Panics(t, func() { container.Init(func(company company) {}) }, "error creating person")
}

func TestInjectSingleton(t *testing.T) {

	type id struct{ number string }

	type employer struct{ id }

	type worker struct{ id }

	newID := func() id {
		return id{number: uuidGenerator()}
	}

	newEmployer := func(id id) employer {
		return employer{id}
	}

	newWorker := func(id id) worker {
		return worker{id}
	}

	container := New()
	container.Register(newID)
	container.Register(newEmployer)
	container.Register(newWorker)

	container.Init(func(employer employer, worker worker) {
		if employer.id != worker.id {
			t.Error()
		}
	})

}

func TestSpecifyingSingleton(t *testing.T) {

	type id struct{ number string }

	type employer struct{ id }

	type worker struct{ id }

	newID := func() id {
		return id{number: uuidGenerator()}
	}

	newEmployer := func(id id) employer {
		return employer{id}
	}

	newWorker := func(id id) worker {
		return worker{id}
	}

	container := New()
	container.Register(newID, Singleton(true))
	container.Register(newEmployer)
	container.Register(newWorker)

	container.Init(func(employer employer, worker worker) {
		if employer.id != worker.id {
			t.Error()
		}
	})
}

func TestInjectNotSingleton(t *testing.T) {

	type id struct{ number string }

	type employer struct{ id }

	type worker struct{ id }

	newID := func() id {
		return id{number: uuidGenerator()}
	}

	newEmployer := func(id id) employer {
		return employer{id}
	}

	newWorker := func(id id) worker {
		return worker{id}
	}

	container := New()
	container.Register(newID, Singleton(false))
	container.Register(newEmployer)
	container.Register(newWorker)

	container.Init(func(employer employer, worker worker) {
		if employer.id == worker.id {
			t.Error()
		}
	})
}

func TestImplementDefaultInterface(t *testing.T) {

	type subject struct {
		observable
	}

	newSubject := func(observable observable) subject {
		return subject{observable}
	}

	newObserver := func() observable {
		return observer{}
	}

	container := New()
	container.Register(newSubject)
	container.Register(newObserver, As(new(observable)))

	container.Init(func(subject subject) {
		if subject.observable.observe() != "observer" {
			t.Error()
		}
	})

}

func TestMultiplesImplementationsInterface(t *testing.T) {

	type person struct {
		observable observable
	}

	type subject struct {
		observable
	}

	newPerson := func(observable observable) person {
		return person{observable}
	}

	newSubject := func(observable observable) subject {
		return subject{observable}
	}

	newVisualizer := func() visualizer {
		return visualizer{}
	}

	newObserver := func() observer {
		return observer{}
	}

	container := New()
	container.Register(newSubject)
	container.Register(newObserver, As(new(observable)))
	container.Register(newVisualizer)
	container.Register(newPerson, Qualifier(new(visualizer), new(observable)))

	container.Init(func(person person, subject subject) {
		if subject.observable.observe() == person.observable.observe() {
			t.Error()
		}
	})
}

func TestInvalidConstructorManyReturns(t *testing.T) {

	type person struct{ name string }
	type dog struct{ name string }

	newPerson := func() (person, dog, error) {
		return person{}, dog{}, nil
	}

	container := New()

	assert.Panics(t, func() { container.Register(newPerson) })

}

func TestNoReturnConstructor(t *testing.T) {

	newPerson := func() {
	}

	container := New()

	assert.Panics(t, func() { container.Register(newPerson) })

}

func TestSecondReturnNotError(t *testing.T) {

	type person struct{ name string }
	type dog struct{ name string }

	newPerson := func() (person, dog) {
		return person{}, dog{}
	}

	container := New()

	assert.Panics(t, func() { container.Register(newPerson) })
}

func TestConstructorNotFunction(t *testing.T) {

	type person struct{ name string }

	container := New()

	assert.Panics(t, func() { container.Register(person{}) })
}
