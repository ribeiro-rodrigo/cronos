package cronos

import "reflect"

type key struct {
	typed reflect.Type
}

type Cronos struct {
	cache
}

func (cronos *Cronos) invokeConstructor(constructor constructor,
	args []reflect.Type) (returns []reflect.Value) {

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
	return nil
}

type constructor interface{}
type component interface{}

type cache struct {
	components    map[key]component
	constructors  map[key]constructor
	notSingletons map[key]bool
	options       OptionsList
}
