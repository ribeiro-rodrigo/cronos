package cronos

import (
	"reflect"
	"sort"
)

type key struct {
	typed reflect.Type
}

type Cronos struct {
	cache
}

// New - initializes the dependency injection container
func New() Cronos {
	return Cronos{
		cache{
			components:    map[key]component{},
			constructors:  map[key]constructor{},
			options:       OptionsList{},
			notSingletons: map[key]bool{},
		},
	}
}

func (cronos *Cronos) proccessOptions() {

	sort.Sort(cronos.cache.options)

	for i := 0; i < len(cronos.cache.options); i++ {
		op := cronos.cache.options[i]
		task := op.task
		task(op.key, cronos)
	}
}

func (cronos *Cronos) invokeConstructor(constructor constructor, args []reflect.Type) (returns []reflect.Value) {

	dependencies := make([]reflect.Value, len(args))

	for i := 0; i < len(dependencies); i++ {
		dependencie := cronos.Fetch(args[i])
		dependencies[i] = reflect.ValueOf(dependencie)
	}

	constructorValue := reflect.ValueOf(constructor)

	returns = constructorValue.Call(dependencies)
	return
}

func (cronos *Cronos) Fetch(typed reflect.Type) interface{} {
	key := key{typed}

	if object, found := cronos.cache.components[key]; found {
		return object
	}

	constructor := cronos.cache.constructors[key]
	args := cronos.getArgs(constructor)

	returns := cronos.invokeConstructor(constructor, args)

	if len(returns) == 2 && !returns[1].IsNil() {
		err := returns[1].Interface().(error)
		panic(err)
	}

	object := returns[0].Interface()

	isNotSingleton := cronos.cache.notSingletons[key]

	if !isNotSingleton {
		cronos.cache.components[key] = object
	}

	return object
}

func (cronos *Cronos) getArgs(constructor constructor) []reflect.Type {
	typec := reflect.TypeOf(constructor)

	deps := make([]reflect.Type, typec.NumIn())

	for i := 0; i < len(deps); i++ {
		deps[i] = typec.In(i)
	}

	return deps
}

type constructor interface{}
type component interface{}

type cache struct {
	components    map[key]component
	constructors  map[key]constructor
	notSingletons map[key]bool
	options       OptionsList
}
